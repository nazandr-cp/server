package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"lend.fam/go-server/internal/gql"
)

// GetUserVaultBalance handles GET /api/users/{address}/vault-balance
func GetUserVaultBalance(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAddress := chi.URLParam(r, "address")

		if !common.IsHexAddress(userAddress) {
			d.Logger.Error("Invalid user address", zap.String("userAddress", userAddress))
			http.Error(w, "Invalid user address", http.StatusBadRequest)
			return
		}

		query := `{
			accounts(where: {id: "` + userAddress + `"}) {
				id
				vaultBalance
			}
		}`

		var result struct {
			Accounts []*gql.Account `json:"accounts"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query user vault balance",
				zap.String("userAddress", userAddress),
				zap.Error(err))
			http.Error(w, "Failed to get user vault balance", http.StatusInternalServerError)
			return
		}

		vaultBalance := "0"
		if len(result.Accounts) > 0 && result.Accounts[0].VaultBalance != nil {
			vaultBalance = result.Accounts[0].VaultBalance.String()
		}

		response := map[string]string{
			"userAddress":  userAddress,
			"vaultBalance": vaultBalance,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
