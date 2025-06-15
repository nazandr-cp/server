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
