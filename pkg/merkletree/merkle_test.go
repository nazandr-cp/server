package merkletree

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestBuildTreeGolden(t *testing.T) {
	recipients := []Recipient{
		{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), TotalEarned: big.NewInt(100)},
		{Address: common.HexToAddress("0x2222222222222222222222222222222222222222"), TotalEarned: big.NewInt(50)},
		{Address: common.HexToAddress("0x3333333333333333333333333333333333333333"), TotalEarned: big.NewInt(75)},
	}
	root, _ := BuildTree(recipients)
	got := common.Bytes2Hex(root[:])
	want := "1a4419cabd8afad85626382b2d212a6e6a1d49d5412ce08a76f8e8b18aab9907"
	if got != want {
		t.Fatalf("unexpected root %s", got)
	}
}

func FuzzBuildTreeDeterministic(f *testing.F) {
	base := []Recipient{
		{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), TotalEarned: big.NewInt(100)},
		{Address: common.HexToAddress("0x2222222222222222222222222222222222222222"), TotalEarned: big.NewInt(50)},
		{Address: common.HexToAddress("0x3333333333333333333333333333333333333333"), TotalEarned: big.NewInt(75)},
	}
	root, _ := BuildTree(base)
	want := root
	f.Fuzz(func(t *testing.T, seed uint64) {
		rs := make([]Recipient, len(base))
		copy(rs, base)
		rand.Seed(int64(seed))
		rand.Shuffle(len(rs), func(i, j int) { rs[i], rs[j] = rs[j], rs[i] })
		r2, _ := BuildTree(rs)
		if r2 != want {
			t.Fatalf("root mismatch")
		}
	})
}
