// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// IDebtSubsidizerSubsidy is an auto generated low-level Go binding around an user-defined struct.
type IDebtSubsidizerSubsidy struct {
	Account    common.Address
	Collection common.Address
	Vault      common.Address
	Amount     *big.Int
	Nonce      *big.Int
	Deadline   *big.Int
}

// IDebtSubsidizerVaultInfo is an auto generated low-level Go binding around an user-defined struct.
type IDebtSubsidizerVaultInfo struct {
	LendingManager common.Address
	CToken         common.Address
}

// IDebtSubsidizerWeightFunction is an auto generated low-level Go binding around an user-defined struct.
type IDebtSubsidizerWeightFunction struct {
	FnType uint8
	P1     *big.Int
	P2     *big.Int
}

// IDebtSubsidizerMetaData contains all meta data concerning the IDebtSubsidizer contract.
var IDebtSubsidizerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addVault\",\"inputs\":[{\"name\":\"vaultAddress_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lendingManagerAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collectionRewardBasis\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIDebtSubsidizer.RewardBasis\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isCollectionWhitelisted\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeCollection\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeVault\",\"inputs\":[{\"name\":\"vaultAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWeightFunction\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"weightFunction\",\"type\":\"tuple\",\"internalType\":\"structIDebtSubsidizer.WeightFunction\",\"components\":[{\"name\":\"fnType\",\"type\":\"uint8\",\"internalType\":\"enumIDebtSubsidizer.WeightFunctionType\"},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"subsidize\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subsidizes\",\"type\":\"tuple[]\",\"internalType\":\"structIDebtSubsidizer.Subsidy[]\",\"components\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"subsidySigner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCollectionPercentageShare\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newSharePercentageBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTrustedSigner\",\"inputs\":[{\"name\":\"newSigner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userNonce\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"nonce\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIDebtSubsidizer.VaultInfo\",\"components\":[{\"name\":\"lendingManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"cToken\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"whitelistCollection\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionType\",\"type\":\"uint8\",\"internalType\":\"enumIDebtSubsidizer.CollectionType\"},{\"name\":\"rewardBasis\",\"type\":\"uint8\",\"internalType\":\"enumIDebtSubsidizer.RewardBasis\"},{\"name\":\"sharePercentageBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CollectionYieldShareUpdated\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldSharePercentage\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"newSharePercentage\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DebtSubsidized\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewCollectionWhitelisted\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDebtSubsidizer.CollectionType\"},{\"name\":\"rewardBasis\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIDebtSubsidizer.RewardBasis\"},{\"name\":\"sharePercentage\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"weightFunction\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDebtSubsidizer.WeightFunction\",\"components\":[{\"name\":\"fnType\",\"type\":\"uint8\",\"internalType\":\"enumIDebtSubsidizer.WeightFunctionType\"},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TrustedSignerUpdated\",\"inputs\":[{\"name\":\"oldSigner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newSigner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"changedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultAdded\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"cTokenAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"lendingManagerAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultRemoved\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WeightFunctionConfigUpdated\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldWeightFunction\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDebtSubsidizer.WeightFunction\",\"components\":[{\"name\":\"fnType\",\"type\":\"uint8\",\"internalType\":\"enumIDebtSubsidizer.WeightFunctionType\"},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}]},{\"name\":\"newWeightFunction\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDebtSubsidizer.WeightFunction\",\"components\":[{\"name\":\"fnType\",\"type\":\"uint8\",\"internalType\":\"enumIDebtSubsidizer.WeightFunctionType\"},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WeightFunctionSet\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fn\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIDebtSubsidizer.WeightFunction\",\"components\":[{\"name\":\"fnType\",\"type\":\"uint8\",\"internalType\":\"enumIDebtSubsidizer.WeightFunctionType\"},{\"name\":\"p1\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"p2\",\"type\":\"int256\",\"internalType\":\"int256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhitelistCollectionRemoved\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotSetSignerToZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ClaimExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CollectionAlreadyExists\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CollectionAlreadyWhitelistedInVault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CollectionNotWhitelisted\",\"inputs\":[{\"name\":\"collection\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CollectionNotWhitelistedInVault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InsufficientYield\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCollectionInterface\",\"inputs\":[{\"name\":\"collectionAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"InvalidNonce\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSecondsColl\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidYieldSharePercentage\",\"inputs\":[{\"name\":\"totalSharePercentage\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidYieldSlice\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LendingManagerAssetMismatch\",\"inputs\":[{\"name\":\"vaultAsset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lmAsset\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"LendingManagerNotSetForVault\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"VaultAlreadyRegistered\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"VaultMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"VaultNotRegistered\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	ID:  "IDebtSubsidizer",
}

// IDebtSubsidizer is an auto generated Go binding around an Ethereum contract.
type IDebtSubsidizer struct {
	abi abi.ABI
}

// NewIDebtSubsidizer creates a new instance of IDebtSubsidizer.
func NewIDebtSubsidizer() *IDebtSubsidizer {
	parsed, err := IDebtSubsidizerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IDebtSubsidizer{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IDebtSubsidizer) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAddVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xec3a7823.
//
// Solidity: function addVault(address vaultAddress_, address lendingManagerAddress_) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackAddVault(vaultAddress common.Address, lendingManagerAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("addVault", vaultAddress, lendingManagerAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCollectionRewardBasis is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe6aeb97b.
//
// Solidity: function collectionRewardBasis(address vaultAddress, address collectionAddress) view returns(uint8)
func (iDebtSubsidizer *IDebtSubsidizer) PackCollectionRewardBasis(vaultAddress common.Address, collectionAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("collectionRewardBasis", vaultAddress, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackCollectionRewardBasis is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe6aeb97b.
//
// Solidity: function collectionRewardBasis(address vaultAddress, address collectionAddress) view returns(uint8)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionRewardBasis(data []byte) (uint8, error) {
	out, err := iDebtSubsidizer.abi.Unpack("collectionRewardBasis", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, err
}

// PackIsCollectionWhitelisted is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x41afb808.
//
// Solidity: function isCollectionWhitelisted(address vaultAddress, address collectionAddress) view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) PackIsCollectionWhitelisted(vaultAddress common.Address, collectionAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("isCollectionWhitelisted", vaultAddress, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackIsCollectionWhitelisted is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x41afb808.
//
// Solidity: function isCollectionWhitelisted(address vaultAddress, address collectionAddress) view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackIsCollectionWhitelisted(data []byte) (bool, error) {
	out, err := iDebtSubsidizer.abi.Unpack("isCollectionWhitelisted", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackPause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x8456cb59.
//
// Solidity: function pause() returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackPause() []byte {
	enc, err := iDebtSubsidizer.abi.Pack("pause")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackPaused is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) PackPaused() []byte {
	enc, err := iDebtSubsidizer.abi.Pack("paused")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackPaused is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackPaused(data []byte) (bool, error) {
	out, err := iDebtSubsidizer.abi.Unpack("paused", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, err
}

// PackRemoveCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9f1bcbad.
//
// Solidity: function removeCollection(address vaultAddress, address collectionAddress) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackRemoveCollection(vaultAddress common.Address, collectionAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("removeCollection", vaultAddress, collectionAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackRemoveVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xceb68c23.
//
// Solidity: function removeVault(address vaultAddress_) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackRemoveVault(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("removeVault", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSetWeightFunction is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x075d01f0.
//
// Solidity: function setWeightFunction(address vaultAddress, address collectionAddress, (uint8,int256,int256) weightFunction) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackSetWeightFunction(vaultAddress common.Address, collectionAddress common.Address, weightFunction IDebtSubsidizerWeightFunction) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("setWeightFunction", vaultAddress, collectionAddress, weightFunction)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSubsidize is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x64559408.
//
// Solidity: function subsidize(address vaultAddress, (address,address,address,uint256,uint256,uint256)[] subsidizes, bytes signature) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackSubsidize(vaultAddress common.Address, subsidizes []IDebtSubsidizerSubsidy, signature []byte) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("subsidize", vaultAddress, subsidizes, signature)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackSubsidySigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x804096e0.
//
// Solidity: function subsidySigner() view returns(address)
func (iDebtSubsidizer *IDebtSubsidizer) PackSubsidySigner() []byte {
	enc, err := iDebtSubsidizer.abi.Pack("subsidySigner")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackSubsidySigner is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x804096e0.
//
// Solidity: function subsidySigner() view returns(address)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackSubsidySigner(data []byte) (common.Address, error) {
	out, err := iDebtSubsidizer.abi.Unpack("subsidySigner", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, err
}

// PackUnpause is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackUnpause() []byte {
	enc, err := iDebtSubsidizer.abi.Pack("unpause")
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateCollectionPercentageShare is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x930e203a.
//
// Solidity: function updateCollectionPercentageShare(address vaultAddress, address collectionAddress, uint16 newSharePercentageBps) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackUpdateCollectionPercentageShare(vaultAddress common.Address, collectionAddress common.Address, newSharePercentageBps uint16) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("updateCollectionPercentageShare", vaultAddress, collectionAddress, newSharePercentageBps)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUpdateTrustedSigner is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7530625f.
//
// Solidity: function updateTrustedSigner(address newSigner) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackUpdateTrustedSigner(newSigner common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("updateTrustedSigner", newSigner)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackUserNonce is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f6cdb0d.
//
// Solidity: function userNonce(address vaultAddress, address userAddress) view returns(uint64 nonce)
func (iDebtSubsidizer *IDebtSubsidizer) PackUserNonce(vaultAddress common.Address, userAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("userNonce", vaultAddress, userAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackUserNonce is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2f6cdb0d.
//
// Solidity: function userNonce(address vaultAddress, address userAddress) view returns(uint64 nonce)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackUserNonce(data []byte) (uint64, error) {
	out, err := iDebtSubsidizer.abi.Unpack("userNonce", data)
	if err != nil {
		return *new(uint64), err
	}
	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	return out0, err
}

// PackVault is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf815c03d.
//
// Solidity: function vault(address vaultAddress) view returns((address,address))
func (iDebtSubsidizer *IDebtSubsidizer) PackVault(vaultAddress common.Address) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("vault", vaultAddress)
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackVault is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf815c03d.
//
// Solidity: function vault(address vaultAddress) view returns((address,address))
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVault(data []byte) (IDebtSubsidizerVaultInfo, error) {
	out, err := iDebtSubsidizer.abi.Unpack("vault", data)
	if err != nil {
		return *new(IDebtSubsidizerVaultInfo), err
	}
	out0 := *abi.ConvertType(out[0], new(IDebtSubsidizerVaultInfo)).(*IDebtSubsidizerVaultInfo)
	return out0, err
}

// PackWhitelistCollection is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x120c030f.
//
// Solidity: function whitelistCollection(address vaultAddress, address collectionAddress, uint8 collectionType, uint8 rewardBasis, uint16 sharePercentageBps) returns()
func (iDebtSubsidizer *IDebtSubsidizer) PackWhitelistCollection(vaultAddress common.Address, collectionAddress common.Address, collectionType uint8, rewardBasis uint8, sharePercentageBps uint16) []byte {
	enc, err := iDebtSubsidizer.abi.Pack("whitelistCollection", vaultAddress, collectionAddress, collectionType, rewardBasis, sharePercentageBps)
	if err != nil {
		panic(err)
	}
	return enc
}

// IDebtSubsidizerCollectionYieldShareUpdated represents a CollectionYieldShareUpdated event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionYieldShareUpdated struct {
	VaultAddress       common.Address
	CollectionAddress  common.Address
	OldSharePercentage uint16
	NewSharePercentage uint16
	Raw                *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerCollectionYieldShareUpdatedEventName = "CollectionYieldShareUpdated"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerCollectionYieldShareUpdated) ContractEventName() string {
	return IDebtSubsidizerCollectionYieldShareUpdatedEventName
}

// UnpackCollectionYieldShareUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event CollectionYieldShareUpdated(address indexed vaultAddress, address indexed collectionAddress, uint16 oldSharePercentage, uint16 newSharePercentage)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionYieldShareUpdatedEvent(log *types.Log) (*IDebtSubsidizerCollectionYieldShareUpdated, error) {
	event := "CollectionYieldShareUpdated"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerCollectionYieldShareUpdated)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IDebtSubsidizerDebtSubsidized represents a DebtSubsidized event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerDebtSubsidized struct {
	VaultAddress      common.Address
	User              common.Address
	CollectionAddress common.Address
	Amount            *big.Int
	Raw               *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerDebtSubsidizedEventName = "DebtSubsidized"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerDebtSubsidized) ContractEventName() string {
	return IDebtSubsidizerDebtSubsidizedEventName
}

// UnpackDebtSubsidizedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event DebtSubsidized(address indexed vaultAddress, address indexed user, address indexed collectionAddress, uint256 amount)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackDebtSubsidizedEvent(log *types.Log) (*IDebtSubsidizerDebtSubsidized, error) {
	event := "DebtSubsidized"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerDebtSubsidized)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IDebtSubsidizerNewCollectionWhitelisted represents a NewCollectionWhitelisted event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerNewCollectionWhitelisted struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
	CollectionType    uint8
	RewardBasis       uint8
	SharePercentage   uint16
	WeightFunction    IDebtSubsidizerWeightFunction
	Raw               *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerNewCollectionWhitelistedEventName = "NewCollectionWhitelisted"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerNewCollectionWhitelisted) ContractEventName() string {
	return IDebtSubsidizerNewCollectionWhitelistedEventName
}

// UnpackNewCollectionWhitelistedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event NewCollectionWhitelisted(address indexed vaultAddress, address indexed collectionAddress, uint8 collectionType, uint8 rewardBasis, uint16 sharePercentage, (uint8,int256,int256) weightFunction)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackNewCollectionWhitelistedEvent(log *types.Log) (*IDebtSubsidizerNewCollectionWhitelisted, error) {
	event := "NewCollectionWhitelisted"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerNewCollectionWhitelisted)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IDebtSubsidizerTrustedSignerUpdated represents a TrustedSignerUpdated event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerTrustedSignerUpdated struct {
	OldSigner common.Address
	NewSigner common.Address
	ChangedBy common.Address
	Raw       *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerTrustedSignerUpdatedEventName = "TrustedSignerUpdated"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerTrustedSignerUpdated) ContractEventName() string {
	return IDebtSubsidizerTrustedSignerUpdatedEventName
}

// UnpackTrustedSignerUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event TrustedSignerUpdated(address oldSigner, address newSigner, address indexed changedBy)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackTrustedSignerUpdatedEvent(log *types.Log) (*IDebtSubsidizerTrustedSignerUpdated, error) {
	event := "TrustedSignerUpdated"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerTrustedSignerUpdated)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IDebtSubsidizerVaultAdded represents a VaultAdded event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultAdded struct {
	VaultAddress          common.Address
	CTokenAddress         common.Address
	LendingManagerAddress common.Address
	Raw                   *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerVaultAddedEventName = "VaultAdded"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerVaultAdded) ContractEventName() string {
	return IDebtSubsidizerVaultAddedEventName
}

// UnpackVaultAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VaultAdded(address indexed vaultAddress, address indexed cTokenAddress, address indexed lendingManagerAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultAddedEvent(log *types.Log) (*IDebtSubsidizerVaultAdded, error) {
	event := "VaultAdded"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerVaultAdded)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IDebtSubsidizerVaultRemoved represents a VaultRemoved event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultRemoved struct {
	VaultAddress common.Address
	Raw          *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerVaultRemovedEventName = "VaultRemoved"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerVaultRemoved) ContractEventName() string {
	return IDebtSubsidizerVaultRemovedEventName
}

// UnpackVaultRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event VaultRemoved(address indexed vaultAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultRemovedEvent(log *types.Log) (*IDebtSubsidizerVaultRemoved, error) {
	event := "VaultRemoved"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerVaultRemoved)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IDebtSubsidizerWeightFunctionConfigUpdated represents a WeightFunctionConfigUpdated event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerWeightFunctionConfigUpdated struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
	OldWeightFunction IDebtSubsidizerWeightFunction
	NewWeightFunction IDebtSubsidizerWeightFunction
	Raw               *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerWeightFunctionConfigUpdatedEventName = "WeightFunctionConfigUpdated"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerWeightFunctionConfigUpdated) ContractEventName() string {
	return IDebtSubsidizerWeightFunctionConfigUpdatedEventName
}

// UnpackWeightFunctionConfigUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WeightFunctionConfigUpdated(address indexed vaultAddress, address indexed collectionAddress, (uint8,int256,int256) oldWeightFunction, (uint8,int256,int256) newWeightFunction)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackWeightFunctionConfigUpdatedEvent(log *types.Log) (*IDebtSubsidizerWeightFunctionConfigUpdated, error) {
	event := "WeightFunctionConfigUpdated"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerWeightFunctionConfigUpdated)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IDebtSubsidizerWeightFunctionSet represents a WeightFunctionSet event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerWeightFunctionSet struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
	Fn                IDebtSubsidizerWeightFunction
	Raw               *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerWeightFunctionSetEventName = "WeightFunctionSet"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerWeightFunctionSet) ContractEventName() string {
	return IDebtSubsidizerWeightFunctionSetEventName
}

// UnpackWeightFunctionSetEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WeightFunctionSet(address indexed vaultAddress, address indexed collectionAddress, (uint8,int256,int256) fn)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackWeightFunctionSetEvent(log *types.Log) (*IDebtSubsidizerWeightFunctionSet, error) {
	event := "WeightFunctionSet"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerWeightFunctionSet)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// IDebtSubsidizerWhitelistCollectionRemoved represents a WhitelistCollectionRemoved event raised by the IDebtSubsidizer contract.
type IDebtSubsidizerWhitelistCollectionRemoved struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
	Raw               *types.Log // Blockchain specific contextual infos
}

const IDebtSubsidizerWhitelistCollectionRemovedEventName = "WhitelistCollectionRemoved"

// ContractEventName returns the user-defined event name.
func (IDebtSubsidizerWhitelistCollectionRemoved) ContractEventName() string {
	return IDebtSubsidizerWhitelistCollectionRemovedEventName
}

// UnpackWhitelistCollectionRemovedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event WhitelistCollectionRemoved(address indexed vaultAddress, address indexed collectionAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackWhitelistCollectionRemovedEvent(log *types.Log) (*IDebtSubsidizerWhitelistCollectionRemoved, error) {
	event := "WhitelistCollectionRemoved"
	if log.Topics[0] != iDebtSubsidizer.abi.Events[event].ID {
		return nil, errors.New("event signature mismatch")
	}
	out := new(IDebtSubsidizerWhitelistCollectionRemoved)
	if len(log.Data) > 0 {
		if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range iDebtSubsidizer.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (iDebtSubsidizer *IDebtSubsidizer) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["AddressZero"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackAddressZeroError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["ArrayLengthMismatch"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackArrayLengthMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CannotSetSignerToZeroAddress"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCannotSetSignerToZeroAddressError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["ClaimExpired"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackClaimExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CollectionAlreadyExists"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCollectionAlreadyExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CollectionAlreadyWhitelistedInVault"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCollectionAlreadyWhitelistedInVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CollectionNotWhitelisted"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCollectionNotWhitelistedError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["CollectionNotWhitelistedInVault"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackCollectionNotWhitelistedInVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InsufficientYield"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInsufficientYieldError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidCollectionInterface"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidCollectionInterfaceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidNonce"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidNonceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidSecondsColl"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidSecondsCollError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidYieldSharePercentage"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidYieldSharePercentageError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["InvalidYieldSlice"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackInvalidYieldSliceError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["LendingManagerAssetMismatch"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackLendingManagerAssetMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["LendingManagerNotSetForVault"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackLendingManagerNotSetForVaultError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["VaultAlreadyRegistered"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackVaultAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["VaultMismatch"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackVaultMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], iDebtSubsidizer.abi.Errors["VaultNotRegistered"].ID.Bytes()[:4]) {
		return iDebtSubsidizer.UnpackVaultNotRegisteredError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// IDebtSubsidizerAddressZero represents a AddressZero error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerAddressZero struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AddressZero()
func IDebtSubsidizerAddressZeroErrorID() common.Hash {
	return common.HexToHash("0x9fabe1c19979afc45ec7efec1bde2c38021c590a0ce42965cf55b3f518197f02")
}

// UnpackAddressZeroError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AddressZero()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackAddressZeroError(raw []byte) (*IDebtSubsidizerAddressZero, error) {
	out := new(IDebtSubsidizerAddressZero)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "AddressZero", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerArrayLengthMismatch represents a ArrayLengthMismatch error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerArrayLengthMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ArrayLengthMismatch()
func IDebtSubsidizerArrayLengthMismatchErrorID() common.Hash {
	return common.HexToHash("0xa24a13a6c9c749fdebc1ced0c54b040f90ec2bad4921a2449a09961f99596abe")
}

// UnpackArrayLengthMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ArrayLengthMismatch()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackArrayLengthMismatchError(raw []byte) (*IDebtSubsidizerArrayLengthMismatch, error) {
	out := new(IDebtSubsidizerArrayLengthMismatch)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "ArrayLengthMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCannotSetSignerToZeroAddress represents a CannotSetSignerToZeroAddress error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCannotSetSignerToZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CannotSetSignerToZeroAddress()
func IDebtSubsidizerCannotSetSignerToZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xc4aaac8bb50f06fadf7bb35d9bd7eb2b4ca9c665ea3aaa56621dc042c7b14b3f")
}

// UnpackCannotSetSignerToZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CannotSetSignerToZeroAddress()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCannotSetSignerToZeroAddressError(raw []byte) (*IDebtSubsidizerCannotSetSignerToZeroAddress, error) {
	out := new(IDebtSubsidizerCannotSetSignerToZeroAddress)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CannotSetSignerToZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerClaimExpired represents a ClaimExpired error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerClaimExpired struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ClaimExpired()
func IDebtSubsidizerClaimExpiredErrorID() common.Hash {
	return common.HexToHash("0x82a49d9e1a771843d39e8826b2cc5ec620f1a84fb3845ddd134da6fe9b0b747c")
}

// UnpackClaimExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ClaimExpired()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackClaimExpiredError(raw []byte) (*IDebtSubsidizerClaimExpired, error) {
	out := new(IDebtSubsidizerClaimExpired)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "ClaimExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCollectionAlreadyExists represents a CollectionAlreadyExists error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionAlreadyExists struct {
	Collection common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionAlreadyExists(address collection)
func IDebtSubsidizerCollectionAlreadyExistsErrorID() common.Hash {
	return common.HexToHash("0xf53cf25dd2089419f7b729be10af57f42fce4446592b3e7b63db5e9a7dea1d22")
}

// UnpackCollectionAlreadyExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionAlreadyExists(address collection)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionAlreadyExistsError(raw []byte) (*IDebtSubsidizerCollectionAlreadyExists, error) {
	out := new(IDebtSubsidizerCollectionAlreadyExists)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CollectionAlreadyExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCollectionAlreadyWhitelistedInVault represents a CollectionAlreadyWhitelistedInVault error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionAlreadyWhitelistedInVault struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionAlreadyWhitelistedInVault(address vaultAddress, address collectionAddress)
func IDebtSubsidizerCollectionAlreadyWhitelistedInVaultErrorID() common.Hash {
	return common.HexToHash("0x5646f04d83aa0b86b63d2636695294c9fabe78eec635ce30317c8d3f2a7445ee")
}

// UnpackCollectionAlreadyWhitelistedInVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionAlreadyWhitelistedInVault(address vaultAddress, address collectionAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionAlreadyWhitelistedInVaultError(raw []byte) (*IDebtSubsidizerCollectionAlreadyWhitelistedInVault, error) {
	out := new(IDebtSubsidizerCollectionAlreadyWhitelistedInVault)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CollectionAlreadyWhitelistedInVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCollectionNotWhitelisted represents a CollectionNotWhitelisted error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionNotWhitelisted struct {
	Collection common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionNotWhitelisted(address collection)
func IDebtSubsidizerCollectionNotWhitelistedErrorID() common.Hash {
	return common.HexToHash("0xd76f0d44a12a3e65ee33e72c9198ae0d715d644e6916f257dfab83e4fe49142d")
}

// UnpackCollectionNotWhitelistedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionNotWhitelisted(address collection)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionNotWhitelistedError(raw []byte) (*IDebtSubsidizerCollectionNotWhitelisted, error) {
	out := new(IDebtSubsidizerCollectionNotWhitelisted)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CollectionNotWhitelisted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerCollectionNotWhitelistedInVault represents a CollectionNotWhitelistedInVault error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerCollectionNotWhitelistedInVault struct {
	VaultAddress      common.Address
	CollectionAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error CollectionNotWhitelistedInVault(address vaultAddress, address collectionAddress)
func IDebtSubsidizerCollectionNotWhitelistedInVaultErrorID() common.Hash {
	return common.HexToHash("0x063b14e113f34053abd31ac480d3cf89b7e139870716916231508db61ef4537f")
}

// UnpackCollectionNotWhitelistedInVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error CollectionNotWhitelistedInVault(address vaultAddress, address collectionAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackCollectionNotWhitelistedInVaultError(raw []byte) (*IDebtSubsidizerCollectionNotWhitelistedInVault, error) {
	out := new(IDebtSubsidizerCollectionNotWhitelistedInVault)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "CollectionNotWhitelistedInVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInsufficientYield represents a InsufficientYield error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInsufficientYield struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InsufficientYield()
func IDebtSubsidizerInsufficientYieldErrorID() common.Hash {
	return common.HexToHash("0xfa3505cfeaa6d565a857daf3cb6754e23c09dd9555dc3027e00c0110bf42262e")
}

// UnpackInsufficientYieldError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InsufficientYield()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInsufficientYieldError(raw []byte) (*IDebtSubsidizerInsufficientYield, error) {
	out := new(IDebtSubsidizerInsufficientYield)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InsufficientYield", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidCollectionInterface represents a InvalidCollectionInterface error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidCollectionInterface struct {
	CollectionAddress common.Address
	InterfaceId       [4]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidCollectionInterface(address collectionAddress, bytes4 interfaceId)
func IDebtSubsidizerInvalidCollectionInterfaceErrorID() common.Hash {
	return common.HexToHash("0x77ec721aca2ae0444bced9a0d7a5f5188a2be2a2efd5a50c57eebbebdf070e76")
}

// UnpackInvalidCollectionInterfaceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidCollectionInterface(address collectionAddress, bytes4 interfaceId)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidCollectionInterfaceError(raw []byte) (*IDebtSubsidizerInvalidCollectionInterface, error) {
	out := new(IDebtSubsidizerInvalidCollectionInterface)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidCollectionInterface", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidNonce represents a InvalidNonce error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidNonce struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidNonce()
func IDebtSubsidizerInvalidNonceErrorID() common.Hash {
	return common.HexToHash("0x756688fec2871909d72599c334b663ffcc94654c438569966c7fd3ab3a351f34")
}

// UnpackInvalidNonceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidNonce()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidNonceError(raw []byte) (*IDebtSubsidizerInvalidNonce, error) {
	out := new(IDebtSubsidizerInvalidNonce)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidNonce", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidSecondsColl represents a InvalidSecondsColl error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidSecondsColl struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSecondsColl()
func IDebtSubsidizerInvalidSecondsCollErrorID() common.Hash {
	return common.HexToHash("0x8d50ad91ce4aa4f6a33fbdf3db5deb828c1f770a89b6a63e157eb9d7ae851df9")
}

// UnpackInvalidSecondsCollError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSecondsColl()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidSecondsCollError(raw []byte) (*IDebtSubsidizerInvalidSecondsColl, error) {
	out := new(IDebtSubsidizerInvalidSecondsColl)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidSecondsColl", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidSignature represents a InvalidSignature error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidSignature struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature()
func IDebtSubsidizerInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x8baa579fce362245063d36f11747a89dd489c54795634fc673cc0e0db51fedc5")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidSignatureError(raw []byte) (*IDebtSubsidizerInvalidSignature, error) {
	out := new(IDebtSubsidizerInvalidSignature)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidYieldSharePercentage represents a InvalidYieldSharePercentage error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidYieldSharePercentage struct {
	TotalSharePercentage *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidYieldSharePercentage(uint256 totalSharePercentage)
func IDebtSubsidizerInvalidYieldSharePercentageErrorID() common.Hash {
	return common.HexToHash("0xbb366faa5bb0b720039d1f82c07bbb07db4fa7da184ae4181de9d51a39d893d4")
}

// UnpackInvalidYieldSharePercentageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidYieldSharePercentage(uint256 totalSharePercentage)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidYieldSharePercentageError(raw []byte) (*IDebtSubsidizerInvalidYieldSharePercentage, error) {
	out := new(IDebtSubsidizerInvalidYieldSharePercentage)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidYieldSharePercentage", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerInvalidYieldSlice represents a InvalidYieldSlice error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerInvalidYieldSlice struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidYieldSlice()
func IDebtSubsidizerInvalidYieldSliceErrorID() common.Hash {
	return common.HexToHash("0x1d18dd9358b8bac99c04d0a1149ec9fca50265639c372e50e9468ac996d78ee1")
}

// UnpackInvalidYieldSliceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidYieldSlice()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackInvalidYieldSliceError(raw []byte) (*IDebtSubsidizerInvalidYieldSlice, error) {
	out := new(IDebtSubsidizerInvalidYieldSlice)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "InvalidYieldSlice", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerLendingManagerAssetMismatch represents a LendingManagerAssetMismatch error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerLendingManagerAssetMismatch struct {
	VaultAsset common.Address
	LmAsset    common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerAssetMismatch(address vaultAsset, address lmAsset)
func IDebtSubsidizerLendingManagerAssetMismatchErrorID() common.Hash {
	return common.HexToHash("0x42a981d74553ce955a75472e22ff8184e6a5612f06aec91a457e357383d769f0")
}

// UnpackLendingManagerAssetMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerAssetMismatch(address vaultAsset, address lmAsset)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackLendingManagerAssetMismatchError(raw []byte) (*IDebtSubsidizerLendingManagerAssetMismatch, error) {
	out := new(IDebtSubsidizerLendingManagerAssetMismatch)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "LendingManagerAssetMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerLendingManagerNotSetForVault represents a LendingManagerNotSetForVault error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerLendingManagerNotSetForVault struct {
	VaultAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error LendingManagerNotSetForVault(address vaultAddress)
func IDebtSubsidizerLendingManagerNotSetForVaultErrorID() common.Hash {
	return common.HexToHash("0x1930e66911b345fdaa39a639549dc16c2105f0f4bc3ae1c1e39867527469205e")
}

// UnpackLendingManagerNotSetForVaultError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error LendingManagerNotSetForVault(address vaultAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackLendingManagerNotSetForVaultError(raw []byte) (*IDebtSubsidizerLendingManagerNotSetForVault, error) {
	out := new(IDebtSubsidizerLendingManagerNotSetForVault)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "LendingManagerNotSetForVault", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerVaultAlreadyRegistered represents a VaultAlreadyRegistered error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultAlreadyRegistered struct {
	VaultAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error VaultAlreadyRegistered(address vaultAddress)
func IDebtSubsidizerVaultAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x38bfcc166bcba4b5532c5572994cadc9ff436e7c322bb1d41b75caf05248d852")
}

// UnpackVaultAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error VaultAlreadyRegistered(address vaultAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultAlreadyRegisteredError(raw []byte) (*IDebtSubsidizerVaultAlreadyRegistered, error) {
	out := new(IDebtSubsidizerVaultAlreadyRegistered)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "VaultAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerVaultMismatch represents a VaultMismatch error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultMismatch struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error VaultMismatch()
func IDebtSubsidizerVaultMismatchErrorID() common.Hash {
	return common.HexToHash("0xc1faacc51457ff2b9dc86cf152d0b03efa014271f93627a43af5c22f0bb7c5ec")
}

// UnpackVaultMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error VaultMismatch()
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultMismatchError(raw []byte) (*IDebtSubsidizerVaultMismatch, error) {
	out := new(IDebtSubsidizerVaultMismatch)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "VaultMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// IDebtSubsidizerVaultNotRegistered represents a VaultNotRegistered error raised by the IDebtSubsidizer contract.
type IDebtSubsidizerVaultNotRegistered struct {
	VaultAddress common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error VaultNotRegistered(address vaultAddress)
func IDebtSubsidizerVaultNotRegisteredErrorID() common.Hash {
	return common.HexToHash("0x299f3425f9f512f17e69ec652b8e51ec669fa2236d04369d83965289ca6a5afb")
}

// UnpackVaultNotRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error VaultNotRegistered(address vaultAddress)
func (iDebtSubsidizer *IDebtSubsidizer) UnpackVaultNotRegisteredError(raw []byte) (*IDebtSubsidizerVaultNotRegistered, error) {
	out := new(IDebtSubsidizerVaultNotRegistered)
	if err := iDebtSubsidizer.abi.UnpackIntoInterface(out, "VaultNotRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}
