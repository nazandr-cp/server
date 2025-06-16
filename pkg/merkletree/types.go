package merkletree

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Leaf represents a merkle tree leaf.
type Leaf struct {
	Index   uint32
	Account common.Address
	Amount  *big.Int
}

// Pair represents an address/amount pair used for simple Merkle trees where
// the leaf hash is `keccak256(account, amount)`.
type Pair struct {
	Account common.Address
	Amount  *big.Int
}
