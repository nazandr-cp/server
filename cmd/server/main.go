package main

import (
	"context"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	config "go-server/configs"
	"go-server/internal/api/handlers"
	eth "go-server/internal/platform/ethereum"
	ws "go-server/internal/platform/websocket"
	"go-server/internal/service/datacollector"
	epochscheduler "go-server/internal/service/epoch"
	"go-server/internal/service/subsidy"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Validate required environment variables from loaded config
	if cfg.EpochManagerAddr == "" {
		logger.Fatal("EPOCH_MANAGER_ADDR environment variable is required")
	}
	if cfg.PrivateKey == "" {
		logger.Fatal("PRIVATE_KEY environment variable is required")
	}
	if cfg.RPCHTTPURL == "" {
		logger.Fatal("RPC_HTTP_URL environment variable is required")
	}

	// Parse polling interval
	pollingIntervalSec, err := strconv.Atoi(getEnvOrDefault("POLLING_INTERVAL", "30")) // Keep using getEnvOrDefault for this one if not in cfg
	if err != nil {
		logger.Fatal("Invalid POLLING_INTERVAL", zap.Error(err))
	}
	pollingInterval := time.Duration(pollingIntervalSec) * time.Second

	logger.Info("Starting application with configuration:",
		zap.String("rpcUrl", cfg.RPCHTTPURL),
		zap.String("epochManagerAddress", cfg.EpochManagerAddr),
		zap.Duration("pollingInterval", pollingInterval),
		zap.String("httpPort", cfg.HTTPPort),
		zap.String("subgraphUrl", cfg.SubgraphURL),
	)

	// Initialize Ethereum clients
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ethClients, err := eth.New(ctx, cfg)
	if err != nil {
		logger.Fatal("Failed to initialize Ethereum clients", zap.Error(err))
	}
	defer ethClients.Close()

	// Initialize Subsidy Service
	subsidyService, err := subsidy.NewService(cfg, logger)
	if err != nil {
		logger.Warn("Failed to initialize Subsidy Service. Continuing without it.", zap.Error(err))
	} else {
		logger.Info("Subsidy Service initialized successfully")
		defer subsidyService.Close()
	}

	// Test DataCollector functionality if vault addresses are provided
	collectionsVaultAddressesStr := getEnvOrDefault("COLLECTIONS_VAULT_ADDRESSES", "")
	if collectionsVaultAddressesStr != "" {
		testDataCollector(cfg.RPCHTTPURL, cfg.EpochManagerAddr, collectionsVaultAddressesStr)
	}

	// Create scheduler instance
	scheduler, err := epochscheduler.NewScheduler(cfg.RPCHTTPURL, cfg.EpochManagerAddr, cfg.PrivateKey, pollingInterval, subsidyService)
	if err != nil {
		logger.Fatal("Failed to create scheduler", zap.Error(err))
	}

	// Initialize WebSocket Hub
	hub := ws.NewHub()

	// Initialize HTTP router and handlers
	router := chi.NewRouter()
	deps := handlers.Deps{
		Cfg:            cfg,
		Eth:            ethClients,
		Hub:            hub,
		SubsidyService: subsidyService,
		Logger:         logger,
	}
	handlers.Register(router, deps)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		logger.Info("Starting HTTP server", zap.String("port", cfg.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Could not listen on", zap.String("port", cfg.HTTPPort), zap.Error(err))
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		scheduler.Start()
	}()

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutdown signal received, stopping services...")
	cancel() // Signal cancellation to all components using the context

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
	} else {
		logger.Info("HTTP server stopped")
	}

	// Give the scheduler and other services some time to cleanup based on context cancellation
	time.Sleep(1 * time.Second)
	logger.Info("Application stopped")
}

func testDataCollector(rpcUrl, epochManagerAddress, vaultAddressesStr string) {
	log.Println("Testing DataCollector functionality...")

	vaultAddresses := strings.Split(vaultAddressesStr, ",")
	for i, addr := range vaultAddresses {
		vaultAddresses[i] = strings.TrimSpace(addr)
	}

	collector, err := datacollector.NewDataCollector(rpcUrl, epochManagerAddress, vaultAddresses)
	if err != nil {
		log.Printf("Failed to create DataCollector: %v", err)
		return
	}
	defer collector.Close()

	epochID := big.NewInt(1)
	data, err := collector.CollectDataForEpoch(epochID)
	if err != nil {
		log.Printf("Failed to collect data for epoch %s: %v", epochID.String(), err)
		return
	}

	log.Printf("DataCollector test completed successfully for epoch %s", epochID.String())
	log.Printf("  Epoch Start Time: %s", data.StartTime.String())
	log.Printf("  Epoch End Time: %s", data.EndTime.String())
	log.Printf("  Total Yield Available: %s", data.TotalYieldAvailableInEpoch.String())
	log.Printf("  Number of vaults: %d", len(data.VaultsInfo))
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
