package merkletree

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"math/rand"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func sortPairs(p []Pair) {
	sort.Slice(p, func(i, j int) bool {
		cmp := bytes.Compare(p[i].Account.Bytes(), p[j].Account.Bytes())
		if cmp == 0 {
			return p[i].Amount.Cmp(p[j].Amount) < 0
		}
		return cmp < 0
	})
}

func TestBuildPairsGolden(t *testing.T) {
	pairs := []Pair{
		{Account: common.HexToAddress("0x0000000000000000000000000000000000000001"), Amount: big.NewInt(10)},
		{Account: common.HexToAddress("0x0000000000000000000000000000000000000002"), Amount: big.NewInt(20)},
		{Account: common.HexToAddress("0x0000000000000000000000000000000000000003"), Amount: big.NewInt(30)},
		{Account: common.HexToAddress("0x0000000000000000000000000000000000000004"), Amount: big.NewInt(40)},
	}
	sortPairs(pairs)
	root, _ := BuildPairs(pairs)
	exp := "f9c15a18fa344d61d1ecae51a79db08cd99ec9f64d704dc02ea0648935c86696"
	if hex.EncodeToString(root[:]) != exp {
		t.Fatalf("unexpected root: %s", hex.EncodeToString(root[:]))
	}
}

func FuzzDeterministicRoot(f *testing.F) {
	base := []Pair{
		{Account: common.HexToAddress("0x0000000000000000000000000000000000000001"), Amount: big.NewInt(10)},
		{Account: common.HexToAddress("0x0000000000000000000000000000000000000002"), Amount: big.NewInt(20)},
		{Account: common.HexToAddress("0x0000000000000000000000000000000000000003"), Amount: big.NewInt(30)},
		{Account: common.HexToAddress("0x0000000000000000000000000000000000000004"), Amount: big.NewInt(40)},
	}
	sortPairs(base)
	exp, _ := BuildPairs(base)
	f.Add(int64(1))
	f.Add(int64(2))
	f.Fuzz(func(t *testing.T, seed int64) {
		r := rand.New(rand.NewSource(seed))
		pairs := append([]Pair{}, base...)
		r.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
		sortPairs(pairs)
		root, _ := BuildPairs(pairs)
		if root != exp {
			t.Fatalf("root changed: %x vs %x", root, exp)
		}
	})
}
