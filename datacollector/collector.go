package datacollector

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"go-server/contracts/collectionsvault"
	"go-server/contracts/epochmanager"
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
	epochManager           epochmanager.EpochManagerABI
	collectionsVaults      map[common.Address]collectionsvault.CollectionsVaultABI
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
	epochManagerContract, err := epochmanager.NewEpochManagerContract(epochManagerAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create EpochManager contract instance: %w", err)
	}

	// Create CollectionsVault contract instances
	vaults := make(map[common.Address]collectionsvault.CollectionsVaultABI)
	for _, vaultAddressStr := range collectionsVaultAddresses {
		vaultAddr := common.HexToAddress(vaultAddressStr)
		vaultContract, err := collectionsvault.NewCollectionsVaultContract(vaultAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create CollectionsVault contract instance for %s: %w", vaultAddressStr, err)
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
	epoch, err := d.epochManager.GetEpoch(callOpts, epochID)
	if err != nil {
		return nil, fmt.Errorf("failed to get epoch details for epoch %s: %w", epochID.String(), err)
	}

	log.Printf("Fetched epoch details - ID: %s, StartTime: %s, EndTime: %s, Status: %d, TotalYield: %s",
		epoch.Id.String(), epoch.StartTime.String(), epoch.EndTime.String(), epoch.Status, epoch.TotalYieldAvailableInEpoch.String())

	// Initialize aggregated data structure
	aggregatedData := &AggregatedEpochData{
		EpochID:                    epoch.Id,
		StartTime:                  epoch.StartTime,
		EndTime:                    epoch.EndTime,
		Status:                     uint8(epoch.Status),
		TotalYieldAvailableInEpoch: epoch.TotalYieldAvailableInEpoch,
		YieldPerVault:              make(map[string]*big.Int),
		VaultsInfo:                 make(map[string]VaultAssetInfo),
	}

	// Fetch vault-specific yield and asset information
	for vaultAddr, vaultContract := range d.collectionsVaults {
		vaultAddrHex := vaultAddr.Hex()

		// Get yield allocated by this vault for the epoch
		vaultYield, err := d.epochManager.GetVaultYieldForEpoch(callOpts, epochID, vaultAddr)
		if err != nil {
			log.Printf("Warning: failed to get vault yield for vault %s and epoch %s: %v", vaultAddrHex, epochID.String(), err)
			vaultYield = big.NewInt(0) // Default to 0 on error
		}

		aggregatedData.YieldPerVault[vaultAddrHex] = vaultYield
		log.Printf("Vault %s allocated yield: %s", vaultAddrHex, vaultYield.String())

		// Get asset information from the vault
		assetAddr, err := vaultContract.Asset(callOpts)
		if err != nil {
			log.Printf("Warning: failed to get asset address for vault %s: %v", vaultAddrHex, err)
			assetAddr = common.Address{} // Default to zero address on error
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
