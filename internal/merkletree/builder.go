package merkletree

import (
	"bytes"
	"encoding/binary"
)

// Build constructs a Merkle tree from leaves and returns the root hash and proofs.
// Proofs are keyed by the account address bytes.
func Build(leaves []Leaf) (root [32]byte, proofs map[[20]byte][][]byte) {
	proofs = make(map[[20]byte][][]byte, len(leaves))
	if len(leaves) == 0 {
		return root, proofs
	}

	level := make([][32]byte, len(leaves))
	for i, l := range leaves {
		var idx [32]byte
		binary.BigEndian.PutUint32(idx[28:], l.Index)

		var amt [32]byte
		if l.Amount != nil {
			l.Amount.FillBytes(amt[:])
		}

		level[i] = Keccak256(idx[:], l.Account.Bytes(), amt[:])
	}

	levels := [][][32]byte{level}
	for len(level) > 1 {
		next := make([][32]byte, (len(level)+1)/2)
		for i := 0; i < len(level); i += 2 {
			if i+1 < len(level) {
				left := level[i]
				right := level[i+1]
				if bytes.Compare(left[:], right[:]) > 0 {
					left, right = right, left
				}
				next[i/2] = Keccak256(left[:], right[:])
			} else {
				next[i/2] = level[i]
			}
		}
		level = next
		levels = append(levels, level)
	}
	root = level[0]

	// Build proofs
	for i, l := range leaves {
		var addr [20]byte
		copy(addr[:], l.Account.Bytes())

		idx := i
		for lvl := 0; lvl < len(levels)-1; lvl++ {
			sib := idx ^ 1
			if sib < len(levels[lvl]) {
				h := make([]byte, 32)
				copy(h, levels[lvl][sib][:])
				proofs[addr] = append(proofs[addr], h)
			}
			idx /= 2
		}
	}

	return root, proofs
}
