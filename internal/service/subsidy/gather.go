package subsidy

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"go-server/internal/platform/graphql"
)

// gather collects total earned and previously claimed amounts for all recipients.
func (s *Service) gather(ctx context.Context, epoch uint64) ([]Recipient, error) {
	s.logger.Info("Gathering recipients for epoch", zap.Uint64("epoch", epoch))

	first := 1000
	skip := 0
	type subResp struct {
		AccountSubsidies []struct {
			Recipient     string `json:"recipient"`
			SecondsEarned string `json:"secondsEarned"`
		} `json:"accountSubsidies"`
	}

	recipientsMap := make(map[common.Address]*big.Int)
	query := `query GetSubsidies($first: Int!, $skip: Int!, $vault: String!) {
        accountSubsidies(first: $first, skip: $skip, where: { vault: $vault }) {
            recipient
            secondsEarned
        }
    }`

	client := graphql.NewClient(s.subgraphURL)
	for {
		vars := map[string]interface{}{
			"first": first,
			"skip":  skip,
			"vault": s.vaultAddr.Hex(),
		}
		var resp subResp
		if err := client.QueryWithVariables(ctx, query, vars, &resp); err != nil {
			return nil, fmt.Errorf("query accountSubsidies: %w", err)
		}
		if len(resp.AccountSubsidies) == 0 {
			break
		}
		for _, a := range resp.AccountSubsidies {
			addr := common.HexToAddress(a.Recipient)
			val, ok := new(big.Int).SetString(a.SecondsEarned, 10)
			if !ok {
				continue
			}
			if recipientsMap[addr] == nil {
				recipientsMap[addr] = new(big.Int)
			}
			recipientsMap[addr].Add(recipientsMap[addr], val)
		}
		if len(resp.AccountSubsidies) < first {
			break
		}
		skip += first
	}

	var recipients []Recipient
	for addr, total := range recipientsMap {
		claimed, err := s.claimedTotal(ctx, s.vaultAddr, addr)
		if err != nil {
			claimed = big.NewInt(0)
			s.logger.Warn("failed to get claimedTotals", zap.String("addr", addr.Hex()), zap.Error(err))
		}
		if total.Cmp(claimed) > 0 {
			recipients = append(recipients, Recipient{
				Address:     addr,
				TotalEarned: total,
				PrevClaimed: claimed,
			})
		}
	}

	s.logger.Info("Finished gathering recipients", zap.Int("count", len(recipients)))
	return recipients, nil
}

// claimedTotal calls the DebtSubsidizer contract to get previously claimed totals.
func (s *Service) claimedTotal(ctx context.Context, vault, recipient common.Address) (*big.Int, error) {
	var out []interface{}
	err := s.subsidizerContract.Call(&bind.CallOpts{Context: ctx}, &out, "claimedTotals", vault, recipient)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return big.NewInt(0), nil
	}
	return out[0].(*big.Int), nil
}
