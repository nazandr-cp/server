package merkletree

import (
	"bytes"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

// Recipient represents an address with earned amount used for deterministic tree.
type Recipient struct {
	Address     common.Address
	TotalEarned *big.Int
}

// BuildTree sorts recipients deterministically and builds the Merkle tree.
func BuildTree(recipients []Recipient) (root [32]byte, proofs map[[20]byte][][]byte) {
	sort.Slice(recipients, func(i, j int) bool {
		if recipients[i].Address == recipients[j].Address {
			return recipients[i].TotalEarned.Cmp(recipients[j].TotalEarned) < 0
		}
		return bytes.Compare(recipients[i].Address.Bytes(), recipients[j].Address.Bytes()) < 0
	})

	pairs := make([]Pair, len(recipients))
	for i, r := range recipients {
		pairs[i] = Pair{Account: r.Address, Amount: r.TotalEarned}
	}
	return BuildPairs(pairs)
}
