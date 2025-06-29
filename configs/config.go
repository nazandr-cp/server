package config

import (
	"encoding/json"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// DeployedContracts holds the addresses of deployed smart contracts.
type DeployedContracts struct {
	Asset              string `json:"asset"`
	Nft                string `json:"nft"`
	CToken             string `json:"cToken"`
	LendingManager     string `json:"lendingManager"`
	CollectionRegistry string `json:"collectionRegistry"`
	CollectionsVault   string `json:"collectionsVault"`
	EpochManager       string `json:"epochManager"`
	DebtSubsidizer     string `json:"debtSubsidizer"`
}

// Config holds environment configuration for the server.
type Config struct {
	RPCHTTPURL               string
	RPCWSURL                 string
	PrivateKey               string
	SubgraphURL              string
	ChainID                  int64
	HTTPPort                 string
	AdminToken               string
	EnableMerkleDistribution bool
	BatchProcessingSize      int
	AnalyticsRetentionDays   int
	Subsidy                  SubsidyConfig
	Contracts                DeployedContracts
	EpochManagerAddr         string
	VaultAddr                string
}

// SubsidyConfig holds configuration for subsidy payments.
type SubsidyConfig struct {
	SubsidizerAddr        string
	SubsidyTokenAddr      string
	SubsidyMinAmount      *big.Int
	SubsidyBatchSize      int
	SubsidyPayoutGasTip   *big.Int
	MerkleTreeDepth       int
	MaxClaimsPerBatch     int
	ClaimExpirationBlocks int64
}

// Load reads configuration from environment variables and contract addresses from a file.
func Load() Config {
	// Load contract addresses
	contracts, err := loadContracts("/artifacts/deployed-contracts.json")
	if err != nil {
		log.Fatalf("Failed to load contract addresses: %v", err)
	}

	cfg := Config{
		RPCHTTPURL:               os.Getenv("RPC_HTTP_URL"),
		RPCWSURL:                 os.Getenv("RPC_WS_URL"),
		PrivateKey:               os.Getenv("PRIVATE_KEY"),
		SubgraphURL:              os.Getenv("SUBGRAPH_URL"),
		HTTPPort:                 envOrDefault("HTTP_PORT", "8080"),
		AdminToken:               os.Getenv("ADMIN_TOKEN"),
		EnableMerkleDistribution: parseBool(envOrDefault("ENABLE_MERKLE_DISTRIBUTION", "true")),
		BatchProcessingSize:      parseInt(envOrDefault("BATCH_PROCESSING_SIZE", "100")),
		AnalyticsRetentionDays:   parseInt(envOrDefault("ANALYTICS_RETENTION_DAYS", "90")),
		Contracts:                contracts,
		EpochManagerAddr:         os.Getenv("EPOCH_MANAGER_ADDR"),
		VaultAddr:                os.Getenv("VAULT_ADDR"),
	}

	cid, _ := strconv.ParseInt(envOrDefault("CHAIN_ID", "1"), 10, 64)
	cfg.ChainID = cid

	// Subsidy configuration
	cfg.Subsidy = SubsidyConfig{
		SubsidizerAddr:        contracts.DebtSubsidizer,
		SubsidyTokenAddr:      contracts.Asset,
		SubsidyMinAmount:      parseBigInt(envOrDefault("SUBSIDY_MIN_AMOUNT", "0")),
		SubsidyBatchSize:      parseInt(envOrDefault("SUBSIDY_BATCH_SIZE", "100")),
		SubsidyPayoutGasTip:   parseBigInt(envOrDefault("SUBSIDY_PAYOUT_GASTIP", "0")),
		MerkleTreeDepth:       parseInt(envOrDefault("MERKLE_TREE_DEPTH", "20")),
		MaxClaimsPerBatch:     parseInt(envOrDefault("MAX_CLAIMS_PER_BATCH", "50")),
		ClaimExpirationBlocks: int64(parseInt(envOrDefault("CLAIM_EXPIRATION_BLOCKS", "100800"))), // ~2 weeks at 12s blocks
	}

	return cfg
}

func loadContracts(filePath string) (DeployedContracts, error) {
	var contracts DeployedContracts
	file, err := os.ReadFile(filePath)
	if err != nil {
		return contracts, err
	}
	err = json.Unmarshal(file, &contracts)
	return contracts, err
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Error parsing int: %v. Using default 0.", err)
		return 0
	}
	return i
}

func parseBigInt(s string) *big.Int {
	if strings.HasPrefix(s, "0x") {
		s = strings.TrimPrefix(s, "0x")
	}
	n, ok := new(big.Int).SetString(s, 10)
	if !ok {
		log.Printf("Error parsing big.Int: %v. Using default 0.", s)
		return big.NewInt(0)
	}
	return n
}

func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Printf("Error parsing bool: %v. Using default false.", err)
		return false
	}
	return b
}

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
