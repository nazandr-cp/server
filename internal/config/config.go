package config

import (
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Config holds environment configuration for the server.
type Config struct {
	Environment         string
	RPCHTTPURL          string
	RPCWSURL            string
	PrivateKey          string
	SubgraphURL         string
	TestSubgraphURL     string
	EpochManagerAddr    string
	VaultAddr           string
	ChainID             int64
	HTTPPort            string
	AdminToken          string
	HealthCheckInterval string
	LogLevel            string
	Subsidy             SubsidyConfig
}

// SubsidyConfig holds configuration for subsidy payments.
type SubsidyConfig struct {
	SubsidizerAddr      string
	SubsidyTokenAddr    string
	SubsidyMinAmount    *big.Int
	SubsidyBatchSize    int
	SubsidyPayoutGasTip *big.Int
}

// Load reads configuration from environment variables.
func Load() Config {
	cfg := Config{
		Environment:         envOrDefault("ENVIRONMENT", "development"),
		RPCHTTPURL:          os.Getenv("RPC_HTTP_URL"),
		RPCWSURL:            os.Getenv("RPC_WS_URL"),
		PrivateKey:          os.Getenv("PRIVATE_KEY"),
		SubgraphURL:         os.Getenv("SUBGRAPH_URL"),
		TestSubgraphURL:     os.Getenv("TEST_SUBGRAPH_URL"),
		EpochManagerAddr:    os.Getenv("EPOCH_MANAGER_ADDR"),
		VaultAddr:           os.Getenv("VAULT_ADDR"),
		HTTPPort:            envOrDefault("HTTP_PORT", "8080"),
		AdminToken:          os.Getenv("ADMIN_TOKEN"),
		HealthCheckInterval: envOrDefault("HEALTH_CHECK_INTERVAL", "30s"),
		LogLevel:            envOrDefault("LOG_LEVEL", "info"),
	}

	cid, _ := strconv.ParseInt(envOrDefault("CHAIN_ID", "1"), 10, 64)
	cfg.ChainID = cid

	// Environment-specific configuration
	if cfg.Environment == "test" {
		if cfg.TestSubgraphURL != "" {
			cfg.SubgraphURL = cfg.TestSubgraphURL
		}
	}

	// Subsidy configuration
	cfg.Subsidy = SubsidyConfig{
		SubsidizerAddr:      os.Getenv("SUBSIDIZER_ADDR"),
		SubsidyTokenAddr:    os.Getenv("SUBSIDY_TOKEN_ADDR"),
		SubsidyMinAmount:    parseBigInt(envOrDefault("SUBSIDY_MIN_AMOUNT", "0")),
		SubsidyBatchSize:    parseInt(envOrDefault("SUBSIDY_BATCH_SIZE", "100")),
		SubsidyPayoutGasTip: parseBigInt(envOrDefault("SUBSIDY_PAYOUT_GASTIP", "0")),
	}

	cfg.Validate()

	return cfg
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

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

// Validate ensures that required configuration fields are set and applies environment-specific defaults.
func (c *Config) Validate() {
	if c.SubgraphURL == "" {
		log.Fatal("SUBGRAPH_URL is not set")
	}
	if c.RPCHTTPURL == "" {
		log.Fatal("RPC_HTTP_URL is not set")
	}
	if c.RPCWSURL == "" {
		log.Fatal("RPC_WS_URL is not set")
	}
	if c.PrivateKey == "" {
		log.Fatal("PRIVATE_KEY is not set")
	}
	if c.EpochManagerAddr == "" {
		log.Fatal("EPOCH_MANAGER_ADDR is not set")
	}
	if c.VaultAddr == "" {
		log.Fatal("VAULT_ADDR is not set")
	}

	// Test environment specific validations
	if c.Environment == "test" {
		if c.TestSubgraphURL == "" {
			log.Println("Warning: TEST_SUBGRAPH_URL is not set in test environment. Using default SUBGRAPH_URL.")
		}
	}
}
