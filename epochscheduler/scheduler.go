package epochscheduler

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"go-server/contracts/epochmanager"
)

type Scheduler struct {
	ethClient              *ethclient.Client
	epochManagerAddr       common.Address
	epochManager           epochmanager.EpochManagerABI
	privateKey             *ecdsa.PrivateKey
	ownerAddress           common.Address
	rpcUrl                 string
	epochManagerAddressStr string
	privateKeyStr          string
	pollingInterval        time.Duration
	epochDuration          time.Duration
	chainID                *big.Int
}

// NewScheduler creates a new scheduler instance
func NewScheduler(rpcUrl, epochManagerAddress, privateKeyStr string, pollingInterval time.Duration) (*Scheduler, error) {
	// Initialize Ethereum client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Get owner address from private key
	ownerAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Parse contract address
	contractAddr := common.HexToAddress(epochManagerAddress)

	// Create contract instance
	epochManagerContract, err := epochmanager.NewEpochManagerContract(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create EpochManager contract instance: %v", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	scheduler := &Scheduler{
		ethClient:              client,
		epochManagerAddr:       contractAddr,
		epochManager:           epochManagerContract,
		privateKey:             privateKey,
		ownerAddress:           ownerAddress,
		rpcUrl:                 rpcUrl,
		epochManagerAddressStr: epochManagerAddress,
		privateKeyStr:          privateKeyStr,
		pollingInterval:        pollingInterval,
		chainID:                chainID,
	}

	// Fetch epoch duration from contract
	epochDuration, err := scheduler.epochManager.EpochDuration(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch duration: %v", err)
	}
	scheduler.epochDuration = time.Duration(epochDuration.Int64()) * time.Second

	log.Printf("Scheduler initialized - Owner: %s, Contract: %s, Epoch Duration: %v",
		ownerAddress.Hex(), contractAddr.Hex(), scheduler.epochDuration)

	return scheduler, nil
}

// Run starts the main scheduler loop
func (s *Scheduler) Run(ctx context.Context) error {
	log.Printf("Starting EpochScheduler with polling interval: %v", s.pollingInterval)

	ticker := time.NewTicker(s.pollingInterval)
	defer ticker.Stop()

	// Initial check
	if err := s.checkAndProcessEpochTransition(); err != nil {
		log.Printf("Initial epoch check failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopping...")
			return ctx.Err()
		case <-ticker.C:
			if err := s.checkAndProcessEpochTransition(); err != nil {
				log.Printf("Epoch transition check failed: %v", err)
			}
		}
	}
}

// checkAndProcessEpochTransition checks if epoch transition is needed and processes it
func (s *Scheduler) checkAndProcessEpochTransition() error {
	// Get current epoch ID
	currentEpochId, err := s.epochManager.GetCurrentEpochId(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get current epoch ID: %v", err)
	}

	log.Printf("Current epoch ID: %s", currentEpochId.String())

	// If no epochs started yet (currentEpochId == 0)
	if currentEpochId.Cmp(big.NewInt(0)) == 0 {
		log.Println("No epochs started yet, attempting to start first epoch")
		return s.startFirstEpoch()
	}

	// Get current epoch details
	currentEpoch, err := s.epochManager.GetEpoch(&bind.CallOpts{}, currentEpochId)
	if err != nil {
		return fmt.Errorf("failed to get current epoch details: %v", err)
	}

	currentTime := time.Now().Unix()
	epochEndTime := currentEpoch.EndTime.Int64()

	log.Printf("Current time: %d, Epoch end time: %d, Status: %d",
		currentTime, epochEndTime, currentEpoch.Status)

	// Check if current epoch has ended and is still active
	if currentTime > epochEndTime && currentEpoch.Status == epochmanager.EpochStatusActive {
		log.Printf("Current epoch %s has ended, starting new epoch", currentEpochId.String())

		// Log that automated processing pipelines would be triggered here
		log.Printf("Would trigger automated processing pipelines for epoch %s", currentEpochId.String())

		return s.attemptToStartNewEpoch(currentEpoch.EndTime)
	}

	log.Printf("Epoch %s is still active (ends at %d)", currentEpochId.String(), epochEndTime)
	return nil
}

// startFirstEpoch starts the first epoch
func (s *Scheduler) startFirstEpoch() error {
	startTime := big.NewInt(time.Now().Unix())

	log.Printf("Starting first epoch with start time: %s", startTime.String())

	return s.executeStartNewEpoch(startTime)
}

// attemptToStartNewEpoch starts a new epoch after the previous one has ended
func (s *Scheduler) attemptToStartNewEpoch(previousEpochEndTime *big.Int) error {
	currentTime := time.Now().Unix()

	// Calculate next start time - use the later of previous epoch end time or current time
	var nextStartTime *big.Int
	if currentTime > previousEpochEndTime.Int64() {
		nextStartTime = big.NewInt(currentTime)
	} else {
		nextStartTime = previousEpochEndTime
	}

	log.Printf("Starting new epoch with start time: %s (previous epoch ended at: %s)",
		nextStartTime.String(), previousEpochEndTime.String())

	return s.executeStartNewEpoch(nextStartTime)
}

// executeStartNewEpoch executes the StartNewEpoch transaction
func (s *Scheduler) executeStartNewEpoch(startTime *big.Int) error {
	// Get current nonce
	nonce, err := s.ethClient.PendingNonceAt(context.Background(), s.ownerAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get gas price
	gasPrice, err := s.ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}

	// Create transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000) // Set a reasonable gas limit
	auth.GasPrice = gasPrice

	// Execute transaction
	tx, err := s.epochManager.StartNewEpoch(auth, startTime)
	if err != nil {
		return fmt.Errorf("failed to start new epoch: %v", err)
	}

	log.Printf("StartNewEpoch transaction submitted: %s", tx.Hash().Hex())

	// Wait for transaction confirmation (optional - could be done async)
	receipt, err := bind.WaitMined(context.Background(), s.ethClient, tx)
	if err != nil {
		log.Printf("Transaction may have failed, receipt error: %v", err)
		return err
	}

	if receipt.Status == 1 {
		log.Printf("New epoch started successfully! Transaction: %s", tx.Hash().Hex())
	} else {
		log.Printf("Transaction failed! Transaction: %s", tx.Hash().Hex())
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	return nil
}
