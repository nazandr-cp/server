package config

import (
	"os"
	"strconv"
)

// Config holds environment configuration for the server.
type Config struct {
	RPCHTTPURL       string
	RPCWSURL         string
	PrivateKey       string
	SubgraphURL      string
	EpochManagerAddr string
	VaultAddr        string
	ChainID          int64
	HTTPPort         string
	AdminToken       string
}

// Load reads configuration from environment variables.
func Load() Config {
	cfg := Config{
		RPCHTTPURL:       os.Getenv("RPC_HTTP_URL"),
		RPCWSURL:         os.Getenv("RPC_WS_URL"),
		PrivateKey:       os.Getenv("PRIVATE_KEY"),
		SubgraphURL:      os.Getenv("SUBGRAPH_URL"),
		EpochManagerAddr: os.Getenv("EPOCH_MANAGER_ADDR"),
		VaultAddr:        os.Getenv("VAULT_ADDR"),
		HTTPPort:         envOrDefault("HTTP_PORT", "8080"),
		AdminToken:       os.Getenv("ADMIN_TOKEN"),
	}
	cid, _ := strconv.ParseInt(envOrDefault("CHAIN_ID", "1"), 10, 64)
	cfg.ChainID = cid
	return cfg
}

func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
