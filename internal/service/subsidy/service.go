package subsidy

import (
	"bytes"
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

	config "lend.fam/go-server/configs"
	"lend.fam/go-server/contracts"
	"lend.fam/go-server/internal/gql"
	"lend.fam/go-server/internal/platform/graphql"
	"lend.fam/go-server/pkg/merkletree"
)

// Recipient is an account eligible for subsidy.
type Recipient struct {
	Address     common.Address
	TotalEarned *big.Int
	PrevClaimed *big.Int
}

// Service manages the subsidy distribution process.
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
			TotalEarned: r.TotalEarned,
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
	mtRecipients := make([]merkletree.Recipient, len(recipients))
	for i, r := range recipients {
		mtRecipients[i] = merkletree.Recipient{Address: r.Address, TotalEarned: r.TotalEarned}
	}
	root, proofs := merkletree.BuildTree(mtRecipients)

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

func (s *Service) CalculateUserEligibility(ctx context.Context, userAddr, vaultAddr, collectionAddr common.Address) (*gql.UserEpochEligibility, error) {
	query := gql.GetUserEpochEligibilityByUserAndEpochQuery()

	currentEpochQuery := gql.GetCurrentEpochQuery()
	var epochResult struct {
		Epochs []*gql.Epoch `json:"epochs"`
	}

	err := graphql.Query(ctx, s.subgraphURL, currentEpochQuery, &epochResult)
	if err != nil {
		return nil, fmt.Errorf("failed to get current epoch: %w", err)
	}

	if len(epochResult.Epochs) == 0 {
		return nil, fmt.Errorf("no active epoch found")
	}

	currentEpoch := epochResult.Epochs[0]

	client := graphql.NewClient(s.subgraphURL)
	variables := map[string]interface{}{
		"userId":  userAddr.Hex(),
		"epochId": currentEpoch.ID,
	}

	var result gql.UserEpochEligibilityResponse
	err = client.QueryWithVariables(ctx, query, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query user eligibility: %w", err)
	}

	for _, eligibility := range result.UserEpochEligibilities {
		if eligibility.Collection != nil &&
			eligibility.Collection.ID == collectionAddr.Hex() {
			return eligibility, nil
		}
	}

	return nil, fmt.Errorf("no eligibility found for user %s in collection %s", userAddr.Hex(), collectionAddr.Hex())
}

func (s *Service) GenerateMerkleProof(ctx context.Context, userAddr common.Address, amount *big.Int, recipients []Recipient) ([][32]byte, error) {
	mtRecipients := make([]merkletree.Recipient, len(recipients))
	for i, r := range recipients {
		mtRecipients[i] = merkletree.Recipient{Address: r.Address, TotalEarned: r.TotalEarned}
	}

	_, proofs := merkletree.BuildTree(mtRecipients)

	var addrKey [20]byte
	copy(addrKey[:], userAddr.Bytes())

	proofBytes, exists := proofs[addrKey]
	if !exists {
		return nil, fmt.Errorf("no proof found for user %s", userAddr.Hex())
	}

	var proof [][32]byte
	for _, b := range proofBytes {
		var tmp [32]byte
		copy(tmp[:], b)
		proof = append(proof, tmp)
	}

	return proof, nil
}

func (s *Service) ValidateClaimData(ctx context.Context, vaultAddr common.Address, claim contracts.IDebtSubsidizerClaimData) error {
	s.logger.Info("Validating claim data for user",
		zap.String("user", claim.Recipient.Hex()),
		zap.String("vault", vaultAddr.Hex()),
		zap.String("amount", claim.TotalEarned.String()))

	return nil
}

func (s *Service) ProcessBatchClaims(ctx context.Context, vaultAddresses []common.Address, claims []contracts.IDebtSubsidizerClaimData) error {
	if len(vaultAddresses) != len(claims) {
		return fmt.Errorf("vault addresses and claims length mismatch")
	}

	batchSize := s.config.SubsidyBatchSize
	for i := 0; i < len(claims); i += batchSize {
		end := i + batchSize
		if end > len(claims) {
			end = len(claims)
		}

		batchVaults := vaultAddresses[i:end]
		batchClaims := claims[i:end]

		err := s.processClaimBatch(ctx, batchVaults, batchClaims)
		if err != nil {
			return fmt.Errorf("failed to process claim batch at index %d: %w", i, err)
		}

		s.logger.Info("Successfully processed claim batch",
			zap.Int("start_index", i),
			zap.Int("end_index", end))
	}

	return nil
}

func (s *Service) processClaimBatch(ctx context.Context, vaultAddresses []common.Address, claims []contracts.IDebtSubsidizerClaimData) error {
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

	receipt, err := bind.WaitMined(ctx, s.ethClient, tx)
	if err != nil {
		return fmt.Errorf("failed to mine claimAllSubsidies transaction: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("claimAllSubsidies transaction failed: %s", receipt.TxHash.Hex())
	}

	return nil
}

func (s *Service) GetClaimStatus(ctx context.Context, userAddr common.Address, epochID string) (*gql.SubsidyDistribution, error) {
	query := gql.GetSubsidyDistributionsByEpochQuery()

	variables := map[string]interface{}{
		"epochId": epochID,
		"first":   1000,
		"skip":    0,
	}

	client := graphql.NewClient(s.subgraphURL)
	var result gql.SubsidyDistributionResponse
	err := client.QueryWithVariables(ctx, query, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query subsidy distributions: %w", err)
	}

	for _, distribution := range result.SubsidyDistributions {
		if distribution.User != nil && distribution.User.ID == userAddr.Hex() {
			return distribution, nil
		}
	}

	return nil, fmt.Errorf("no claim found for user %s in epoch %s", userAddr.Hex(), epochID)
}

// GetUserSubsidyDistributions gets all subsidy distributions for a user in an epoch
func (s *Service) GetUserSubsidyDistributions(ctx context.Context, userAddr common.Address, epochID string) ([]*gql.SubsidyDistribution, error) {
	query := gql.GetSubsidyDistributionsByEpochQuery()

	variables := map[string]interface{}{
		"epochId": epochID,
		"first":   1000,
		"skip":    0,
	}

	client := graphql.NewClient(s.subgraphURL)
	var result gql.SubsidyDistributionResponse
	err := client.QueryWithVariables(ctx, query, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query subsidy distributions: %w", err)
	}

	var userDistributions []*gql.SubsidyDistribution
	for _, distribution := range result.SubsidyDistributions {
		if distribution.User != nil && distribution.User.ID == userAddr.Hex() {
			userDistributions = append(userDistributions, distribution)
		}
	}

	return userDistributions, nil
}

// GetUserClaimHistory gets all historical claims for a user across all epochs
func (s *Service) GetUserClaimHistory(ctx context.Context, userAddr common.Address) ([]*gql.SubsidyDistribution, error) {
	query := gql.GetSubsidyDistributionsQuery()

	variables := map[string]interface{}{
		"first": 1000,
		"skip":  0,
		"where": map[string]interface{}{
			"user": userAddr.Hex(),
		},
	}

	client := graphql.NewClient(s.subgraphURL)
	var result gql.SubsidyDistributionResponse
	err := client.QueryWithVariables(ctx, query, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query user claim history: %w", err)
	}

	return result.SubsidyDistributions, nil
}

// GetAllSubsidyDistributions gets all historical subsidy distributions
func (s *Service) GetAllSubsidyDistributions(ctx context.Context) ([]*gql.SubsidyDistribution, error) {
	query := gql.GetSubsidyDistributionsQuery()

	variables := map[string]interface{}{
		"first": 1000, // Fetch a reasonable number, consider pagination for very large datasets
		"skip":  0,
	}

	client := graphql.NewClient(s.subgraphURL)
	var result gql.SubsidyDistributionResponse
	err := client.QueryWithVariables(ctx, query, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query all subsidy distributions: %w", err)
	}

	return result.SubsidyDistributions, nil
}

// GetMerkleDistribution gets merkle distribution data for an epoch and vault
func (s *Service) GetMerkleDistribution(ctx context.Context, epochID, vaultID string) (*gql.MerkleDistribution, error) {
	query := gql.GetMerkleDistributionByEpochAndVaultQuery()

	variables := map[string]interface{}{
		"epochId": epochID,
		"vaultId": vaultID,
	}

	client := graphql.NewClient(s.subgraphURL)
	var result gql.MerkleDistributionResponse
	err := client.QueryWithVariables(ctx, query, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle distribution: %w", err)
	}

	if len(result.MerkleDistributions) == 0 {
		return nil, fmt.Errorf("no merkle distribution found for epoch %s and vault %s", epochID, vaultID)
	}

	return result.MerkleDistributions[0], nil
}

// GetLatestMerkleDistributions gets the latest merkle distributions for a vault
func (s *Service) GetLatestMerkleDistributions(ctx context.Context, vaultID string) ([]*gql.MerkleDistribution, error) {
	query := gql.GetMerkleDistributionsQuery()

	variables := map[string]interface{}{
		"first": 10,
		"skip":  0,
		"where": map[string]interface{}{
			"vault": vaultID,
		},
	}

	client := graphql.NewClient(s.subgraphURL)
	var result gql.MerkleDistributionResponse
	err := client.QueryWithVariables(ctx, query, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest merkle distributions: %w", err)
	}

	return result.MerkleDistributions, nil
}

// VerifyMerkleProof verifies a merkle proof against a root and leaf data
func (s *Service) VerifyMerkleProof(proof []string, root string, userAddr common.Address, amount *big.Int) (bool, error) {
	// Convert proof strings to bytes32
	proofBytes := make([][32]byte, len(proof))
	for i, p := range proof {
		proofHash := common.HexToHash(p)
		copy(proofBytes[i][:], proofHash.Bytes())
	}

	// Create leaf hash
	var amt [32]byte
	if amount != nil {
		amount.FillBytes(amt[:])
	}
	leafHash := merkletree.Keccak256(userAddr.Bytes(), amt[:])

	// Verify proof
	rootHash := common.HexToHash(root)
	return s.verifyProof(proofBytes, rootHash, leafHash), nil
}

// verifyProof verifies a merkle proof
func (s *Service) verifyProof(proof [][32]byte, root [32]byte, leaf [32]byte) bool {
	computedHash := leaf
	for _, proofElement := range proof {
		if bytes.Compare(computedHash[:], proofElement[:]) < 0 {
			computedHash = merkletree.Keccak256(computedHash[:], proofElement[:])
		} else {
			computedHash = merkletree.Keccak256(proofElement[:], computedHash[:])
		}
	}
	return bytes.Equal(computedHash[:], root[:])
}
