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

	"go-server/contracts"
	"go-server/internal/service/subsidy"
)

const (
	EpochStatusActive = 0
)

type Scheduler struct {
	ethClient              *ethclient.Client
	epochManagerAddr       common.Address
	epochManager           *bind.BoundContract
	privateKey             *ecdsa.PrivateKey
	ownerAddress           common.Address
	rpcUrl                 string
	epochManagerAddressStr string
	privateKeyStr          string
	pollingInterval        time.Duration
	epochDuration          time.Duration // This will be fetched from contract
	chainID                *big.Int
	subsidyService         *subsidy.Service
}

// NewScheduler creates a new scheduler instance
func NewScheduler(rpcUrl, epochManagerAddress, privateKeyStr string, pollingInterval time.Duration, subsidyService *subsidy.Service) (*Scheduler, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	pkey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	ownerAddr := crypto.PubkeyToAddress(pkey.PublicKey)
	contractAddr := common.HexToAddress(epochManagerAddress)

	iEpochManager := contracts.NewIEpochManager()
	epochManagerBinding := iEpochManager.Instance(client, contractAddr)
	if epochManagerBinding == nil {
		return nil, fmt.Errorf("failed to create EpochManager contract binding instance")
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	s := &Scheduler{
		ethClient:              client,
		epochManagerAddr:       contractAddr,
		epochManager:           epochManagerBinding,
		privateKey:             pkey,
		ownerAddress:           ownerAddr,
		rpcUrl:                 rpcUrl,
		epochManagerAddressStr: epochManagerAddress,
		privateKeyStr:          privateKeyStr,
		pollingInterval:        pollingInterval,
		chainID:                chainID,
		subsidyService:         subsidyService,
	}

	// Fetch epoch duration from the contract upon initialization
	epochDur, err := s.GetEpochDurationFromContract()
	if err != nil {
		log.Printf("Warning: failed to get epoch duration from contract: %v. Using default or zero.", err)
		s.epochDuration = 0
	} else {
		s.epochDuration = epochDur
		log.Printf("Epoch duration set from contract: %v", s.epochDuration)
	}

	return s, nil
}

// GetEpochDurationFromContract fetches the epoch duration from the smart contract.
func (s *Scheduler) GetEpochDurationFromContract() (time.Duration, error) {
	var out []interface{}
	err := s.epochManager.Call(&bind.CallOpts{Context: context.Background()}, &out, "epochDuration")
	if err != nil {
		return 0, fmt.Errorf("failed to call epochDuration on contract: %w", err)
	}
	if len(out) == 0 {
		return 0, fmt.Errorf("epochDuration call returned no data")
	}
	epochDurationBig, ok := out[0].(*big.Int)
	if !ok {
		return 0, fmt.Errorf("failed to cast epoch duration to *big.Int, got %T", out[0])
	}
	return time.Duration(epochDurationBig.Int64()) * time.Second, nil
}

func (s *Scheduler) Start() {
	log.Println("Starting epoch scheduler...")
	ticker := time.NewTicker(s.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndStartNewEpoch()
		}
	}
}

func (s *Scheduler) checkAndStartNewEpoch() {
	log.Println("Checking for new epoch...")

	var outGetCurrentEpochId []interface{}
	err := s.epochManager.Call(&bind.CallOpts{Context: context.Background()}, &outGetCurrentEpochId, "getCurrentEpochId")
	if err != nil {
		log.Printf("Error getting current epoch ID: %v", err)
		return
	}
	var currentEpochId *big.Int
	if len(outGetCurrentEpochId) > 0 {
		currentEpochId, _ = outGetCurrentEpochId[0].(*big.Int)
	}
	if currentEpochId == nil {
		log.Printf("Error: currentEpochId is nil after calling getCurrentEpochId")
		return
	}
	log.Printf("Current epoch ID: %s", currentEpochId.String())

	var outGetEpoch []interface{}
	err = s.epochManager.Call(&bind.CallOpts{Context: context.Background()}, &outGetEpoch, "getEpoch", currentEpochId)
	if err != nil {
		log.Printf("Error getting current epoch data for epoch %s: %v", currentEpochId.String(), err)
		return
	}

	var epochEndTime *big.Int
	var epochStatus *big.Int

	if len(outGetEpoch) > 0 {
		// Assuming getEpoch returns a struct/tuple where elements are directly in outGetEpoch[0]
		// or outGetEpoch itself is the list of elements if the function returns multiple values.
		// The exact parsing depends on the Solidity function signature and abigen's behavior.
		// If getEpoch returns (uint256 id, uint256 startTime, uint256 endTime, uint256 status, ...)
		// and abigen unpacks this into a []interface{} where each element is a field.

		// Try to assert outGetEpoch[0] to []interface{} first, as structs are often wrapped.
		epochDataSlice, isSlice := outGetEpoch[0].([]interface{})
		if isSlice {
			if len(epochDataSlice) >= 4 { // id, startTime, endTime, status
				if et, ok := epochDataSlice[2].(*big.Int); ok { // endTime is 3rd field (index 2)
					epochEndTime = et
				}
				if st, ok := epochDataSlice[3].(*big.Int); ok { // status is 4th field (index 3)
					epochStatus = st
				}
			} else {
				log.Printf("Error: getEpoch returned slice with not enough elements: got %d, expected at least 4", len(epochDataSlice))
				return
			}
		} else if len(outGetEpoch) >= 4 { // Fallback: if outGetEpoch itself is the flat list of members
			if et, ok := outGetEpoch[2].(*big.Int); ok { // endTime is 3rd field (index 2)
				epochEndTime = et
			}
			if st, ok := outGetEpoch[3].(*big.Int); ok { // status is 4th field (index 3)
				epochStatus = st
			}
		} else {
			log.Printf("Error: getEpoch returned unexpected data structure or not enough elements. Len(outGetEpoch): %d", len(outGetEpoch))
			return
		}
	} else {
		log.Printf("Error: getEpoch call returned no data.")
		return
	}

	if epochEndTime == nil || epochStatus == nil {
		log.Printf("Error parsing epoch data from contract call (endTime or status is nil)")
		return
	}

	currentTime := big.NewInt(time.Now().Unix())
	log.Printf("Current time: %s, Epoch end time: %s, Epoch status: %s", currentTime.String(), epochEndTime.String(), epochStatus.String())

	if currentTime.Cmp(epochEndTime) > 0 && epochStatus.Cmp(big.NewInt(EpochStatusActive)) == 0 {
		log.Printf("Epoch %s has ended, attempting to start new epoch...", currentEpochId.String())
		s.startNewEpoch() // This attempts to start the *next* epoch

		// After attempting to start a new epoch, trigger subsidy for the epoch that just finalized.
		if s.subsidyService != nil {
			log.Printf("Triggering subsidy payments for finalized epoch ID: %s", currentEpochId.String())
			// Use a background context for the subsidy run, or a more specific one if available/needed.
			// The subsidy operation should not block the scheduler's main loop indefinitely.
			// Errors from subsidyService.Run should be logged but not halt the epoch scheduler.
			go func(epochIDToSubsidize *big.Int) {
				if epochIDToSubsidize == nil {
					log.Printf("Error: epochIDToSubsidize is nil, cannot run subsidy.")
					return
				}
				err := s.subsidyService.Run(context.Background(), epochIDToSubsidize.Uint64())
				if err != nil {
					log.Printf("Error running subsidy service for epoch %s: %v", epochIDToSubsidize.String(), err)
				} else {
					log.Printf("Subsidy service successfully processed epoch %s", epochIDToSubsidize.String())
				}
			}(currentEpochId) // Pass currentEpochId to the goroutine to avoid race conditions if currentEpochId is modified later.
		} else {
			log.Println("Subsidy service is not configured, skipping subsidy payment.")
		}
	} else {
		log.Println("Epoch has not ended or is not active.")
	}
}

func (s *Scheduler) startNewEpoch() {
	log.Println("Attempting to start a new epoch...")
	nonce, err := s.ethClient.PendingNonceAt(context.Background(), s.ownerAddress)
	if err != nil {
		log.Printf("Failed to get nonce: %v", err)
		return
	}

	suggestedGasPrice, err := s.ethClient.SuggestGasPrice(context.Background()) // Renamed variable
	if err != nil {
		log.Printf("Failed to suggest gas price: %v", err)
		return
	}

	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, s.chainID)
	if err != nil {
		log.Printf("Failed to create transactor: %v", err)
		return
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = suggestedGasPrice

	currentUnixTime := big.NewInt(time.Now().Unix())
	tx, err := s.epochManager.Transact(auth, "startNewEpoch", currentUnixTime)
	if err != nil {
		log.Printf("Failed to send StartNewEpoch transaction: %v", err)
		return
	}
	log.Printf("StartNewEpoch transaction sent: %s", tx.Hash().Hex())
}
