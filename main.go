package main

import (
	"log"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"go-server/datacollector"
	"go-server/epochscheduler"
)

func main() {
	// Read configuration from environment variables
	rpcUrl := getEnvOrDefault("RPC_URL", "http://localhost:8545")
	epochManagerAddress := getEnvOrDefault("EPOCH_MANAGER_ADDRESS", "")
	privateKeyStr := getEnvOrDefault("PRIVATE_KEY", "")
	pollingIntervalStr := getEnvOrDefault("POLLING_INTERVAL", "30")
	collectionsVaultAddressesStr := getEnvOrDefault("COLLECTIONS_VAULT_ADDRESSES", "")

	// Validate required environment variables
	if epochManagerAddress == "" {
		log.Fatal("EPOCH_MANAGER_ADDRESS environment variable is required")
	}
	if privateKeyStr == "" {
		log.Fatal("PRIVATE_KEY environment variable is required")
	}

	// Parse polling interval
	pollingIntervalSec, err := strconv.Atoi(pollingIntervalStr)
	if err != nil {
		log.Fatalf("Invalid POLLING_INTERVAL: %v", err)
	}
	pollingInterval := time.Duration(pollingIntervalSec) * time.Second

	log.Printf("Starting EpochScheduler with configuration:")
	log.Printf("  RPC URL: %s", rpcUrl)
	log.Printf("  EpochManager Address: %s", epochManagerAddress)
	log.Printf("  Polling Interval: %v", pollingInterval)

	// Test DataCollector functionality if vault addresses are provided
	if collectionsVaultAddressesStr != "" {
		testDataCollector(rpcUrl, epochManagerAddress, collectionsVaultAddressesStr)
	}

	// Create scheduler instance
	scheduler, err := epochscheduler.NewScheduler(rpcUrl, epochManagerAddress, privateKeyStr, pollingInterval)
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	// Setup graceful shutdown
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start scheduler in a goroutine
	go func() {
		scheduler.Start()
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, stopping scheduler...")
	// cancel()

	// Give the scheduler some time to cleanup
	time.Sleep(2 * time.Second)
	log.Println("EpochScheduler stopped")
}

// testDataCollector demonstrates the DataCollector functionality
func testDataCollector(rpcUrl, epochManagerAddress, vaultAddressesStr string) {
	log.Println("Testing DataCollector functionality...")

	// Parse vault addresses
	vaultAddresses := strings.Split(vaultAddressesStr, ",")
	for i, addr := range vaultAddresses {
		vaultAddresses[i] = strings.TrimSpace(addr)
	}

	// Create DataCollector instance
	collector, err := datacollector.NewDataCollector(rpcUrl, epochManagerAddress, vaultAddresses)
	if err != nil {
		log.Printf("Failed to create DataCollector: %v", err)
		return
	}
	defer collector.Close()

	// Test with epoch ID 1
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
