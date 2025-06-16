package merkletree

import "golang.org/x/crypto/sha3"

// Keccak256 returns Keccak-256 hash of the concatenated data chunks.
func Keccak256(data ...[]byte) [32]byte {
	h := sha3.NewLegacyKeccak256()
	for _, b := range data {
		h.Write(b)
	}
	var out [32]byte
	sum := h.Sum(nil)
	copy(out[:], sum)
	return out
}
