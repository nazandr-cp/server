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

// BuildPairs constructs a Merkle tree from address/amount pairs where each leaf
// is hashed as keccak256(account, amount). It returns the Merkle root and a map
// of proofs keyed by the account address bytes.
func BuildPairs(pairs []Pair) (root [32]byte, proofs map[[20]byte][][]byte) {
	proofs = make(map[[20]byte][][]byte, len(pairs))
	if len(pairs) == 0 {
		return root, proofs
	}

	for i := 1; i < len(pairs); i++ {
		prev := pairs[i-1]
		curr := pairs[i]
		if bytes.Compare(prev.Account.Bytes(), curr.Account.Bytes()) > 0 ||
			(bytes.Equal(prev.Account.Bytes(), curr.Account.Bytes()) && prev.Amount.Cmp(curr.Amount) > 0) {
			panic("unsorted recipients")
		}
	}

	level := make([][32]byte, len(pairs))
	for i, p := range pairs {
		var amt [32]byte
		if p.Amount != nil {
			p.Amount.FillBytes(amt[:])
		}

		level[i] = Keccak256(p.Account.Bytes(), amt[:])
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

	for i, p := range pairs {
		var addr [20]byte
		copy(addr[:], p.Account.Bytes())

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
