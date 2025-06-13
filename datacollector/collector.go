package datacollector

import (
	"context"
	"fmt"
	"go-server/contracts"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EpochYieldData represents yield data for a specific epoch
type EpochYieldData struct {
	EpochID             *big.Int
	TotalYieldAvailable *big.Int
	VaultYields         map[common.Address]*big.Int // Vault address -> yield allocated by this vault
}

// VaultAssetInfo represents information about a vault and its underlying asset
type VaultAssetInfo struct {
	VaultAddress common.Address
	AssetAddress common.Address
}

// AggregatedEpochData represents comprehensive data collected for an epoch
type AggregatedEpochData struct {
	EpochID                    *big.Int
	StartTime                  *big.Int
	EndTime                    *big.Int
	Status                     uint8 // From EpochManager.EpochStatus
	TotalYieldAvailableInEpoch *big.Int
	// Mapping of vault address to the amount of yield it has allocated for this epoch
	YieldPerVault map[string]*big.Int // string for common.Address.Hex() as map key
	// Information about each vault, like its underlying asset
	VaultsInfo map[string]VaultAssetInfo // Vault address (hex) -> asset info
}

// DataCollector aggregates data required for subsidy calculations
type DataCollector struct {
	ethClient              *ethclient.Client
	epochManager           *bind.BoundContract                    // Changed to *bind.BoundContract
	collectionsVaults      map[common.Address]*bind.BoundContract // Changed to *bind.BoundContract
	rpcUrl                 string
	epochManagerAddressStr string
}

// NewDataCollector creates a new DataCollector instance
func NewDataCollector(rpcUrl string, epochManagerAddress string, collectionsVaultAddresses []string) (*DataCollector, error) {
	// Initialize Ethereum client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Create EpochManager contract instance
	epochManagerAddr := common.HexToAddress(epochManagerAddress)
	iEpochManager := contracts.NewIEpochManager() // Use the generated constructor
	epochManagerContract := iEpochManager.Instance(client, epochManagerAddr)
	if epochManagerContract == nil {
		return nil, fmt.Errorf("failed to create EpochManager contract instance")
	}

	// Create CollectionsVault contract instances
	vaults := make(map[common.Address]*bind.BoundContract)
	iCollectionsVault := contracts.NewICollectionsVault() // Use the generated constructor
	for _, vaultAddressStr := range collectionsVaultAddresses {
		vaultAddr := common.HexToAddress(vaultAddressStr)
		vaultContract := iCollectionsVault.Instance(client, vaultAddr)
		if vaultContract == nil {
			return nil, fmt.Errorf("failed to create CollectionsVault contract instance for %s", vaultAddressStr)
		}
		vaults[vaultAddr] = vaultContract
	}

	return &DataCollector{
		ethClient:              client,
		epochManager:           epochManagerContract,
		collectionsVaults:      vaults,
		rpcUrl:                 rpcUrl,
		epochManagerAddressStr: epochManagerAddress,
	}, nil
}

// CollectDataForEpoch collects all relevant data for a given epoch
func (d *DataCollector) CollectDataForEpoch(epochID *big.Int) (*AggregatedEpochData, error) {
	log.Printf("Collecting data for epoch %s", epochID.String())

	// Create call options for read-only calls
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}

	// Fetch epoch details from EpochManager
	var outGetEpoch []interface{}
	err := d.epochManager.Call(callOpts, &outGetEpoch, "getEpoch", epochID)
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch details for epoch %s: %w", epochID.String(), err)
	}

	// Assuming getEpoch returns a struct where fields are ordered: Id, StartTime, EndTime, Status, TotalYieldAvailableInEpoch
	// This needs to match the actual struct definition in your Solidity contract and how abigen unpacks it.
	var epoch struct { // Define an anonymous struct to unpack into
		Id                         *big.Int
		StartTime                  *big.Int
		EndTime                    *big.Int
		Status                     uint8 // Assuming status is uint8 in contract or can be cast
		TotalYieldAvailableInEpoch *big.Int
		// Add other fields if your contract's Epoch struct has more that are returned
	}

	if len(outGetEpoch) > 0 {
		// If getEpoch returns the struct members as a flat list in `outGetEpoch[0]`
		// which itself is a slice or an array.
		epochData, ok := outGetEpoch[0].([]interface{}) // Or the specific type abigen unpacks to
		if ok && len(epochData) >= 5 {                  // Ensure we have enough elements for the assumed struct
			epoch.Id, _ = epochData[0].(*big.Int)
			epoch.StartTime, _ = epochData[1].(*big.Int)
			epoch.EndTime, _ = epochData[2].(*big.Int)
			if statusBigInt, ok := epochData[3].(*big.Int); ok { // Assuming status is returned as big.Int
				epoch.Status = uint8(statusBigInt.Uint64())
			}
			epoch.TotalYieldAvailableInEpoch, _ = epochData[4].(*big.Int)
		} else if !ok {
			// Fallback: if outGetEpoch itself is the flat list of members
			if len(outGetEpoch) >= 5 {
				epoch.Id, _ = outGetEpoch[0].(*big.Int)
				epoch.StartTime, _ = outGetEpoch[1].(*big.Int)
				epoch.EndTime, _ = outGetEpoch[2].(*big.Int)
				if statusBigInt, ok := outGetEpoch[3].(*big.Int); ok { // Assuming status is returned as big.Int
					epoch.Status = uint8(statusBigInt.Uint64())
				}
				epoch.TotalYieldAvailableInEpoch, _ = outGetEpoch[4].(*big.Int)
			} else {
				log.Printf("Error: getEpoch returned unexpected data structure or not enough elements in the primary array.")
				return nil, fmt.Errorf("getEpoch returned unexpected data structure or not enough elements in the primary array")
			}
		} else {
			log.Printf("Error: getEpoch returned unexpected data structure or not enough elements.")
			return nil, fmt.Errorf("getEpoch returned unexpected data structure or not enough elements")
		}
	} else {
		log.Printf("Error: getEpoch call returned no data.")
		return nil, fmt.Errorf("getEpoch call returned no data")
	}

	log.Printf("Fetched epoch details - ID: %s, StartTime: %s, EndTime: %s, Status: %d, TotalYield: %s",
		epoch.Id.String(), epoch.StartTime.String(), epoch.EndTime.String(), epoch.Status, epoch.TotalYieldAvailableInEpoch.String())

	// Initialize aggregated data structure
	aggregatedData := &AggregatedEpochData{
		EpochID:                    epoch.Id,
		StartTime:                  epoch.StartTime,
		EndTime:                    epoch.EndTime,
		Status:                     epoch.Status, // Already uint8
		TotalYieldAvailableInEpoch: epoch.TotalYieldAvailableInEpoch,
		YieldPerVault:              make(map[string]*big.Int),
		VaultsInfo:                 make(map[string]VaultAssetInfo),
	}

	// Fetch vault-specific yield and asset information
	for vaultAddr, vaultContractInstance := range d.collectionsVaults {
		vaultAddrHex := vaultAddr.Hex()

		// Get yield allocated by this vault for the epoch
		var outGetVaultYieldForEpoch []interface{}
		err := d.epochManager.Call(callOpts, &outGetVaultYieldForEpoch, "getVaultYieldForEpoch", epochID, vaultAddr)
		var vaultYield *big.Int
		if err != nil {
			log.Printf("Warning: failed to get vault yield for vault %s and epoch %s: %v", vaultAddrHex, epochID.String(), err)
			vaultYield = big.NewInt(0) // Default to 0 on error
		} else {
			if len(outGetVaultYieldForEpoch) > 0 {
				vaultYield, _ = outGetVaultYieldForEpoch[0].(*big.Int)
			} else {
				log.Printf("Warning: getVaultYieldForEpoch for vault %s returned no data", vaultAddrHex)
				vaultYield = big.NewInt(0)
			}
		}
		if vaultYield == nil { // Ensure vaultYield is not nil
			vaultYield = big.NewInt(0)
		}

		aggregatedData.YieldPerVault[vaultAddrHex] = vaultYield
		log.Printf("Vault %s allocated yield: %s", vaultAddrHex, vaultYield.String())

		// Get asset information from the vault
		var outAsset []interface{}
		err = vaultContractInstance.Call(callOpts, &outAsset, "asset")
		var assetAddr common.Address
		if err != nil {
			log.Printf("Warning: failed to get asset address for vault %s: %v", vaultAddrHex, err)
			assetAddr = common.Address{} // Default to zero address on error
		} else {
			if len(outAsset) > 0 {
				assetAddr, _ = outAsset[0].(common.Address)
			} else {
				log.Printf("Warning: asset call for vault %s returned no data", vaultAddrHex)
				assetAddr = common.Address{}
			}
		}

		aggregatedData.VaultsInfo[vaultAddrHex] = VaultAssetInfo{
			VaultAddress: vaultAddr,
			AssetAddress: assetAddr,
		}
		log.Printf("Vault %s asset address: %s", vaultAddrHex, assetAddr.Hex())
	}

	log.Printf("Successfully collected data for epoch %s with %d vaults", epochID.String(), len(d.collectionsVaults))
	return aggregatedData, nil
}

// Close closes the Ethereum client connection
func (d *DataCollector) Close() {
	if d.ethClient != nil {
		d.ethClient.Close()
	}
}
