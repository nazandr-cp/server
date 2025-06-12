package eth

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"go-server/contracts/collectionsvault"
	"go-server/contracts/epochmanager"
	"go-server/internal/config"
)

// Clients bundles Ethereum related clients and contracts.
type Clients struct {
	RPC    *ethclient.Client
	WS     *ethclient.Client
	EM     epochmanager.EpochManagerABI
	VA     collectionsvault.CollectionsVaultABI
	TxOpts *bind.TransactOpts

	events chan types.Log
}

// New creates Clients based on configuration.
func New(ctx context.Context, cfg config.Config) (*Clients, error) {
	rpc, err := ethclient.DialContext(ctx, cfg.RPCHTTPURL)
	if err != nil {
		return nil, err
	}

	var ws *ethclient.Client
	if cfg.RPCWSURL != "" {
		ws, err = ethclient.DialContext(ctx, cfg.RPCWSURL)
		if err != nil {
			rpc.Close()
			return nil, err
		}
	}

	pk, err := crypto.HexToECDSA(trim0x(cfg.PrivateKey))
	if err != nil {
		rpc.Close()
		if ws != nil {
			ws.Close()
		}
		return nil, err
	}

	chainID := big.NewInt(cfg.ChainID)
	txOpts, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		rpc.Close()
		if ws != nil {
			ws.Close()
		}
		return nil, err
	}

	emAddr := common.HexToAddress(cfg.EpochManagerAddr)
	em, _ := epochmanager.NewEpochManagerContract(emAddr, rpc)

	vaAddr := common.HexToAddress(cfg.VaultAddr)
	va, _ := collectionsvault.NewCollectionsVaultContract(vaAddr, rpc)

	c := &Clients{
		RPC:    rpc,
		WS:     ws,
		EM:     em,
		VA:     va,
		TxOpts: txOpts,
		events: make(chan types.Log, 16),
	}

	if ws != nil {
		go c.subscribeLogs(ctx, emAddr)
	}

	return c, nil
}

func (c *Clients) subscribeLogs(ctx context.Context, addr common.Address) {
	q := ethereum.FilterQuery{Addresses: []common.Address{addr}}
	logsCh := make(chan types.Log)
	sub, err := c.WS.SubscribeFilterLogs(ctx, q, logsCh)
	if err != nil {
		log.Printf("ws subscribe error: %v", err)
		return
	}
	for {
		select {
		case l := <-logsCh:
			c.events <- l
		case err := <-sub.Err():
			log.Printf("subscription error: %v", err)
			return
		case <-ctx.Done():
			sub.Unsubscribe()
			return
		}
	}
}

// Events returns a read-only channel of blockchain logs.
func (c *Clients) Events() <-chan types.Log { return c.events }

// Close closes all clients.
func (c *Clients) Close() {
	if c.RPC != nil {
		c.RPC.Close()
	}
	if c.WS != nil {
		c.WS.Close()
	}
	close(c.events)
}

func trim0x(s string) string {
	if len(s) > 1 && s[:2] == "0x" {
		return s[2:]
	}
	return s
}
