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

	return &Service{
		subgraphURL:        subgraphURL,
		subsidizerContract: subsidizerContract,
		tokenContract:      tokenContract,
		signer:             privateKey,
		ownerAddress:       ownerAddress,
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

func (s *Service) processPayments(ctx context.Context, recipients []Recipient) error {
	if len(recipients) == 0 {
		return nil
	}

	s.logger.Info("Processing payment batch", zap.Int("count", len(recipients)))

	var (
		addresses []common.Address
		amounts   []*big.Int
	)

	for _, r := range recipients {
		addresses = append(addresses, r.Address)
		amounts = append(amounts, r.Amount)
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

	var vaultAddresses []common.Address
	var claims []contracts.IDebtSubsidizerClaimData

	for i := range addresses {
		vaultAddresses = append(vaultAddresses, common.HexToAddress("0x0000000000000000000000000000000000000000"))
		claims = append(claims, contracts.IDebtSubsidizerClaimData{
			Recipient:   addresses[i],
			TotalEarned: amounts[i],
			MerkleProof: [][32]byte{},
		})
	}

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

	s.logger.Info("Total recipients to process", zap.Int("count", len(recipients)))

	batchSize := s.config.SubsidyBatchSize
	for i := 0; i < len(recipients); i += batchSize {
		end := i + batchSize
		if end > len(recipients) {
			end = len(recipients)
		}
		batch := recipients[i:end]

		err := s.processPayments(ctx, batch)
		if err != nil {
			s.logger.Error("Failed to process payment batch", zap.Error(err), zap.Int("start_index", i), zap.Int("end_index", end))
			return fmt.Errorf("failed to process batch: %w", err)
		}
		s.logger.Info("Successfully processed batch", zap.Int("start_index", i), zap.Int("end_index", end))
	}

	s.logger.Info("Subsidy distribution completed for epoch", zap.Uint64("epoch", epoch))
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
