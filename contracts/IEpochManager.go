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

// IEpochManagerMetaData contains all meta data concerning the IEpochManager contract.
var IEpochManagerMetaData = bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"allocateVaultYield\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getCurrentEpochId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	ID:  "IEpochManager",
}

// IEpochManager is an auto generated Go binding around an Ethereum contract.
type IEpochManager struct {
	abi abi.ABI
}

// NewIEpochManager creates a new instance of IEpochManager.
func NewIEpochManager() *IEpochManager {
	parsed, err := IEpochManagerMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &IEpochManager{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *IEpochManager) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackAllocateVaultYield is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf05ca914.
//
// Solidity: function allocateVaultYield(address vault, uint256 amount) returns()
func (iEpochManager *IEpochManager) PackAllocateVaultYield(vault common.Address, amount *big.Int) []byte {
	enc, err := iEpochManager.abi.Pack("allocateVaultYield", vault, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackGetCurrentEpochId is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa29a839f.
//
// Solidity: function getCurrentEpochId() view returns(uint256)
func (iEpochManager *IEpochManager) PackGetCurrentEpochId() []byte {
	enc, err := iEpochManager.abi.Pack("getCurrentEpochId")
	if err != nil {
		panic(err)
	}
	return enc
}

// UnpackGetCurrentEpochId is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa29a839f.
//
// Solidity: function getCurrentEpochId() view returns(uint256)
func (iEpochManager *IEpochManager) UnpackGetCurrentEpochId(data []byte) (*big.Int, error) {
	out, err := iEpochManager.abi.Unpack("getCurrentEpochId", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, err
}
