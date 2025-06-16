package subsidy

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	config "go-server/configs"
	"go-server/contracts"
	gql "go-server/internal/platform/graphql"
	"go-server/pkg/merkletree"
)

// Recipient represents an account eligible for subsidy
type Recipient struct {
	Address common.Address
	Amount  *big.Int
}

// Service manages the subsidy distribution process
type Service struct {
	subgraphURL        string
	subsidizerContract *bind.BoundContract
	tokenContract      *bind.BoundContract
	signer             *ecdsa.PrivateKey
	ownerAddress       common.Address
	vaultAddr          common.Address
	config             config.SubsidyConfig
	logger             *zap.Logger
	ethClient          *ethclient.Client
	chainID            *big.Int
}

// NewService creates a new Subsidy Service instance
func NewService(cfg config.Config, logger *zap.Logger) (*Service, error) {
	ethClient, err := ethclient.Dial(cfg.RPCHTTPURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	ownerAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Subsidizer Contract
	subsidizerAddr := common.HexToAddress(cfg.Subsidy.SubsidizerAddr)
	iDebtSubsidizer := contracts.NewIDebtSubsidizer()
	subsidizerContract := iDebtSubsidizer.Instance(ethClient, subsidizerAddr)
	if subsidizerContract == nil {
		return nil, fmt.Errorf("failed to create IDebtSubsidizer contract instance")
	}

	// Token Contract (IERC20)
	tokenAddr := common.HexToAddress(cfg.Subsidy.SubsidyTokenAddr)
	iERC20 := contracts.NewIERC20()
	tokenContract := iERC20.Instance(ethClient, tokenAddr)
	if tokenContract == nil {
		return nil, fmt.Errorf("failed to create IERC20 contract instance")
	}

	// Subgraph Client
	subgraphURL := cfg.SubgraphURL

	vaultAddr := common.HexToAddress(cfg.VaultAddr)

	return &Service{
		subgraphURL:        subgraphURL,
		subsidizerContract: subsidizerContract,
		tokenContract:      tokenContract,
		signer:             privateKey,
		ownerAddress:       ownerAddress,
		vaultAddr:          vaultAddr,
		config:             cfg.Subsidy,
		logger:             logger,
		ethClient:          ethClient,
		chainID:            chainID,
	}, nil
}

func (s *Service) gather(ctx context.Context, epoch uint64) ([]Recipient, error) {
	s.logger.Info("Gathering recipients for epoch", zap.Uint64("epoch", epoch))

	var recipients []Recipient
	skip := 0
	first := 1000

	for {
		query := fmt.Sprintf(`
		query GetAccounts($skip: Int!, $first: Int!, $minAmount: BigInt!) {
			accounts(first: $first, skip: $skip, where: {subsidyEarned_gt: "%s"}) {
				id
				subsidyEarned
			}
		}
	`, s.config.SubsidyMinAmount.String())

		var result struct {
			Accounts []struct {
				ID            string `json:"id"`
				SubsidyEarned string `json:"subsidyEarned"`
			} `json:"accounts"`
		}

		err := gql.Query(ctx, s.subgraphURL, query, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to query subgraph: %w", err)
		}

		if len(result.Accounts) == 0 {
			break
		}

		for _, acc := range result.Accounts {
			address := common.HexToAddress(acc.ID)
			amount, ok := new(big.Int).SetString(acc.SubsidyEarned, 10)
			if !ok {
				s.logger.Warn("Failed to parse subsidyEarned amount", zap.String("amount", acc.SubsidyEarned))
				continue
			}
			recipients = append(recipients, Recipient{Address: address, Amount: amount})
		}

		if len(result.Accounts) < first {
			break
		}
		skip += first
	}

	s.logger.Info("Finished gathering recipients", zap.Int("count", len(recipients)))
	return recipients, nil
}

func (s *Service) processPayments(ctx context.Context, root [32]byte, proofs map[[20]byte][][]byte, recipients []Recipient) error {
	if len(recipients) == 0 {
		return nil
	}

	s.logger.Info("Processing payment batch", zap.Int("count", len(recipients)))

	vaultAddress := s.vaultAddr

	var claims []contracts.IDebtSubsidizerClaimData
	var vaultAddresses []common.Address

	for _, r := range recipients {
		var addrKey [20]byte
		copy(addrKey[:], r.Address.Bytes())

		proofBytes := proofs[addrKey]
		var merkleProof [][32]byte
		for _, b := range proofBytes {
			var tmp [32]byte
			copy(tmp[:], b)
			merkleProof = append(merkleProof, tmp)
		}

		claims = append(claims, contracts.IDebtSubsidizerClaimData{
			Recipient:   r.Address,
			TotalEarned: r.Amount,
			MerkleProof: merkleProof,
		})
		vaultAddresses = append(vaultAddresses, vaultAddress)
	}

	nonce, err := s.ethClient.PendingNonceAt(ctx, s.ownerAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasTipCap := s.config.SubsidyPayoutGasTip
	if gasTipCap.Cmp(big.NewInt(0)) == 0 {
		suggestedGasTip, err := s.ethClient.SuggestGasTipCap(ctx)
		if err != nil {
			s.logger.Warn("Failed to suggest gas tip cap, using default 0", zap.Error(err))
			gasTipCap = big.NewInt(0)
		} else {
			gasTipCap = suggestedGasTip
		}
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.signer, s.chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(0)
	auth.GasTipCap = gasTipCap

	tx, err := s.subsidizerContract.Transact(auth, "claimAllSubsidies", vaultAddresses, claims)
	if err != nil {
		return fmt.Errorf("failed to send claimAllSubsidies transaction: %w", err)
	}

	s.logger.Info("ClaimAllSubsidies transaction sent", zap.String("tx_hash", tx.Hash().Hex()))

	receipt, err := bind.WaitMined(ctx, s.ethClient, tx)
	if err != nil {
		return fmt.Errorf("failed to mine claimAllSubsidies transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("claimAllSubsidies transaction failed: %s", receipt.TxHash.Hex())
	}

	s.logger.Info("ClaimAllSubsidies transaction successful", zap.String("tx_hash", receipt.TxHash.Hex()))
	return nil
}

func (s *Service) Run(ctx context.Context, epoch uint64) error {
	s.logger.Info("Starting subsidy distribution for epoch", zap.Uint64("epoch", epoch))

	recipients, err := s.gather(ctx, epoch)
	if err != nil {
		return fmt.Errorf("failed to gather recipients: %w", err)
	}

	if len(recipients) == 0 {
		s.logger.Info("No recipients found for subsidy in this epoch.")
		return nil
	}

	// Build Merkle tree for all recipients
	pairs := make([]merkletree.Pair, len(recipients))
	for i, r := range recipients {
		pairs[i] = merkletree.Pair{Account: r.Address, Amount: r.Amount}
	}
	root, proofs := merkletree.BuildPairs(pairs)

	// Update Merkle root on the contract
	if err := s.updateMerkleRoot(ctx, root); err != nil {
		return fmt.Errorf("failed to update merkle root: %w", err)
	}

	s.logger.Info("Total recipients to process", zap.Int("count", len(recipients)))

	batchSize := s.config.SubsidyBatchSize
	for i := 0; i < len(recipients); i += batchSize {
		end := i + batchSize
		if end > len(recipients) {
			end = len(recipients)
		}
		batch := recipients[i:end]

		err := s.processPayments(ctx, root, proofs, batch)
		if err != nil {
			s.logger.Error("Failed to process payment batch", zap.Error(err), zap.Int("start_index", i), zap.Int("end_index", end))
			return fmt.Errorf("failed to process batch: %w", err)
		}
		s.logger.Info("Successfully processed batch", zap.Int("start_index", i), zap.Int("end_index", end))
	}

	s.logger.Info("Subsidy distribution completed for epoch", zap.Uint64("epoch", epoch))
	return nil
}

// updateMerkleRoot submits a transaction to update the Merkle root for the
// configured vault.
func (s *Service) updateMerkleRoot(ctx context.Context, root [32]byte) error {
	nonce, err := s.ethClient.PendingNonceAt(ctx, s.ownerAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasTipCap := s.config.SubsidyPayoutGasTip
	if gasTipCap.Cmp(big.NewInt(0)) == 0 {
		suggestedGasTip, err := s.ethClient.SuggestGasTipCap(ctx)
		if err != nil {
			s.logger.Warn("Failed to suggest gas tip cap, using default 0", zap.Error(err))
			gasTipCap = big.NewInt(0)
		} else {
			gasTipCap = suggestedGasTip
		}
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.signer, s.chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(0)
	auth.GasTipCap = gasTipCap

	vaultAddr := s.vaultAddr
	tx, err := s.subsidizerContract.Transact(auth, "updateMerkleRoot", vaultAddr, root)
	if err != nil {
		return fmt.Errorf("failed to send updateMerkleRoot transaction: %w", err)
	}

	s.logger.Info("updateMerkleRoot transaction sent", zap.String("tx_hash", tx.Hash().Hex()))
	receipt, err := bind.WaitMined(ctx, s.ethClient, tx)
	if err != nil {
		return fmt.Errorf("failed to mine updateMerkleRoot transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("updateMerkleRoot transaction failed: %s", receipt.TxHash.Hex())
	}

	s.logger.Info("Merkle root updated", zap.String("tx_hash", receipt.TxHash.Hex()))
	return nil
}

func (s *Service) Close() {
	if s.ethClient != nil {
		s.ethClient.Close()
	}
	if s.logger != nil {
		s.logger.Sync()
	}
}
