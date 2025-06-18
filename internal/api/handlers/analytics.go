package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"lend.fam/go-server/internal/gql"
)

// SystemMetricsResponse is system-wide metrics.
type SystemMetricsResponse struct {
	TotalValueLocked          string `json:"totalValueLocked"`
	TotalYieldDistributed     string `json:"totalYieldDistributed"`
	TotalSubsidiesDistributed string `json:"totalSubsidiesDistributed"`
	TotalUsers                string `json:"totalUsers"`
	TotalCollections          string `json:"totalCollections"`
	TotalVaults               string `json:"totalVaults"`
	SystemUtilizationRate     string `json:"systemUtilizationRate"`
	AverageAPY                string `json:"averageAPY"`
	LastUpdatedTimestamp      int64  `json:"lastUpdatedTimestamp"`
}

// EpochAllocationsResponse is epoch allocations.
type EpochAllocationsResponse struct {
	EpochID               uint64                `json:"epochId"`
	TotalYieldAllocated   string                `json:"totalYieldAllocated"`
	TotalYieldDistributed string                `json:"totalYieldDistributed"`
	VaultAllocations      []VaultAllocationInfo `json:"vaultAllocations"`
}

// VaultAllocationInfo is allocation info for a vault.
type VaultAllocationInfo struct {
	VaultAddress          string `json:"vaultAddress"`
	YieldAllocated        string `json:"yieldAllocated"`
	SubsidiesDistributed  string `json:"subsidiesDistributed"`
	ParticipantCount      string `json:"participantCount"`
	AverageSubsidyPerUser string `json:"averageSubsidyPerUser"`
	UtilizationRate       string `json:"utilizationRate"`
}

// VaultPerformanceResponse is vault performance metrics.
type VaultPerformanceResponse struct {
	VaultAddress            string `json:"vaultAddress"`
	TotalShares             string `json:"totalShares"`
	TotalDeposits           string `json:"totalDeposits"`
	TotalCTokens            string `json:"totalCTokens"`
	TotalPrincipalDeposited string `json:"totalPrincipalDeposited"`
	GlobalDepositIndex      string `json:"globalDepositIndex"`
	CurrentAPY              string `json:"currentAPY"`
	TotalParticipants       string `json:"totalParticipants"`
	TotalCollections        string `json:"totalCollections"`
	LastUpdatedTimestamp    int64  `json:"lastUpdatedTimestamp"`
}

// DailyAnalyticsResponse is daily analytics data.
type DailyAnalyticsResponse struct {
	Date                      string `json:"date"`
	DailyVolumeUSD            string `json:"dailyVolumeUSD"`
	DailyTransactionCount     string `json:"dailyTransactionCount"`
	DailyActiveUsers          string `json:"dailyActiveUsers"`
	DailyYieldDistributed     string `json:"dailyYieldDistributed"`
	DailySubsidiesDistributed string `json:"dailySubsidiesDistributed"`
	AverageAPY                string `json:"averageAPY"`
	UtilizationRate           string `json:"utilizationRate"`
	Timestamp                 int64  `json:"timestamp"`
}

// GetSystemMetrics handles GET /api/v1/system/metrics
func GetSystemMetrics(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `{
			systemStates(first: 1, orderBy: lastUpdatedTimestamp, orderDirection: desc) {
				id
				totalVaults
				totalCollections
				totalUsers
				totalValueLocked
				totalYieldDistributed
				totalSubsidiesDistributed
				systemUtilizationRate
				averageAPY
				lastUpdatedTimestamp
			}
		}`

		var result struct {
			SystemStates []*gql.SystemState `json:"systemStates"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query system metrics", zap.Error(err))
			http.Error(w, "Failed to get system metrics", http.StatusInternalServerError)
			return
		}

		response := SystemMetricsResponse{
			TotalValueLocked:          "0",
			TotalYieldDistributed:     "0",
			TotalSubsidiesDistributed: "0",
			TotalUsers:                "0",
			TotalCollections:          "0",
			TotalVaults:               "0",
			SystemUtilizationRate:     "0",
			AverageAPY:                "0",
			LastUpdatedTimestamp:      time.Now().Unix(),
		}

		if len(result.SystemStates) > 0 {
			state := result.SystemStates[0]
			if state.TotalValueLocked != nil {
				response.TotalValueLocked = state.TotalValueLocked.String()
			}
			if state.TotalYieldDistributed != nil {
				response.TotalYieldDistributed = state.TotalYieldDistributed.String()
			}
			if state.TotalSubsidiesDistributed != nil {
				response.TotalSubsidiesDistributed = state.TotalSubsidiesDistributed.String()
			}
			if state.TotalUsers != nil {
				response.TotalUsers = state.TotalUsers.String()
			}
			if state.TotalCollections != nil {
				response.TotalCollections = state.TotalCollections.String()
			}
			if state.TotalVaults != nil {
				response.TotalVaults = state.TotalVaults.String()
			}
			if state.SystemUtilizationRate != nil {
				response.SystemUtilizationRate = state.SystemUtilizationRate.String()
			}
			if state.AverageAPY != nil {
				response.AverageAPY = state.AverageAPY.String()
			}
			if state.LastUpdatedTimestamp != nil {
				response.LastUpdatedTimestamp = state.LastUpdatedTimestamp.Int64()
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetEpochAllocations handles GET /api/v1/epochs/{epochId}/allocations
func GetEpochAllocations(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		epochIdStr := chi.URLParam(r, "epochId")

		epochId, err := strconv.ParseUint(epochIdStr, 10, 64)
		if err != nil {
			d.Logger.Error("Invalid epoch ID", zap.String("epochId", epochIdStr), zap.Error(err))
			http.Error(w, "Invalid epoch ID", http.StatusBadRequest)
			return
		}

		query := `{
			epochs(where: {epochNumber: "` + epochIdStr + `"}) {
				id
				epochNumber
				totalYieldAllocated
				totalYieldDistributed
				vaultAllocations {
					id
					vault {
						id
					}
					yieldAllocated
					subsidiesDistributed
					participantCount
					averageSubsidyPerUser
					utilizationRate
				}
			}
		}`

		var result struct {
			Epochs []*gql.Epoch `json:"epochs"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query epoch allocations",
				zap.Uint64("epochId", epochId),
				zap.Error(err))
			http.Error(w, "Failed to get epoch allocations", http.StatusInternalServerError)
			return
		}

		response := EpochAllocationsResponse{
			EpochID:               epochId,
			TotalYieldAllocated:   "0",
			TotalYieldDistributed: "0",
			VaultAllocations:      []VaultAllocationInfo{},
		}

		if len(result.Epochs) > 0 {
			epoch := result.Epochs[0]
			if epoch.TotalYieldAllocated != nil {
				response.TotalYieldAllocated = epoch.TotalYieldAllocated.String()
			}
			if epoch.TotalYieldDistributed != nil {
				response.TotalYieldDistributed = epoch.TotalYieldDistributed.String()
			}

			for _, allocation := range epoch.VaultAllocations {
				vaultInfo := VaultAllocationInfo{
					VaultAddress: allocation.Vault.ID,
				}
				if allocation.YieldAllocated != nil {
					vaultInfo.YieldAllocated = allocation.YieldAllocated.String()
				}
				if allocation.SubsidiesDistributed != nil {
					vaultInfo.SubsidiesDistributed = allocation.SubsidiesDistributed.String()
				}
				if allocation.ParticipantCount != nil {
					vaultInfo.ParticipantCount = allocation.ParticipantCount.String()
				}
				if allocation.AverageSubsidyPerUser != nil {
					vaultInfo.AverageSubsidyPerUser = allocation.AverageSubsidyPerUser.String()
				}
				if allocation.UtilizationRate != nil {
					vaultInfo.UtilizationRate = allocation.UtilizationRate.String()
				}
				response.VaultAllocations = append(response.VaultAllocations, vaultInfo)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetVaultPerformance handles GET /api/v1/vaults/{address}/performance
func GetVaultPerformance(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vaultAddress := chi.URLParam(r, "address")

		if !common.IsHexAddress(vaultAddress) {
			d.Logger.Error("Invalid vault address", zap.String("vaultAddress", vaultAddress))
			http.Error(w, "Invalid vault address", http.StatusBadRequest)
			return
		}

		query := `{
			collectionsVaults(where: {id: "` + vaultAddress + `"}) {
				id
				totalShares
				totalDeposits
				totalCTokens
				totalPrincipalDeposited
				globalDepositIndex
				collectionParticipations {
					id
					totalParticipants
				}
				updatedAtTimestamp
			}
		}`

		var result struct {
			CollectionsVaults []*gql.CollectionsVault `json:"collectionsVaults"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query vault performance",
				zap.String("vaultAddress", vaultAddress),
				zap.Error(err))
			http.Error(w, "Failed to get vault performance", http.StatusInternalServerError)
			return
		}

		response := VaultPerformanceResponse{
			VaultAddress:            vaultAddress,
			TotalShares:             "0",
			TotalDeposits:           "0",
			TotalCTokens:            "0",
			TotalPrincipalDeposited: "0",
			GlobalDepositIndex:      "0",
			CurrentAPY:              "0",
			TotalParticipants:       "0",
			TotalCollections:        "0",
			LastUpdatedTimestamp:    time.Now().Unix(),
		}

		if len(result.CollectionsVaults) > 0 {
			vault := result.CollectionsVaults[0]
			if vault.TotalShares != nil {
				response.TotalShares = vault.TotalShares.String()
			}
			if vault.TotalDeposits != nil {
				response.TotalDeposits = vault.TotalDeposits.String()
			}
			if vault.TotalCTokens != nil {
				response.TotalCTokens = vault.TotalCTokens.String()
			}
			if vault.TotalPrincipalDeposited != nil {
				response.TotalPrincipalDeposited = vault.TotalPrincipalDeposited.String()
			}
			if vault.GlobalDepositIndex != nil {
				response.GlobalDepositIndex = vault.GlobalDepositIndex.String()
			}
			if vault.UpdatedAtTimestamp != nil {
				response.LastUpdatedTimestamp = vault.UpdatedAtTimestamp.Int64()
			}

			response.TotalCollections = strconv.Itoa(len(vault.CollectionParticipations))

			totalParticipants := int64(0)
			for _, participation := range vault.CollectionParticipations {
				if participation.TotalParticipants != nil {
					totalParticipants += participation.TotalParticipants.Int64()
				}
			}
			response.TotalParticipants = strconv.FormatInt(totalParticipants, 10)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetDailyAnalytics handles GET /api/v1/analytics/daily/{date}
func GetDailyAnalytics(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := chi.URLParam(r, "date")

		// Validate date format (YYYY-MM-DD)
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			d.Logger.Error("Invalid date format", zap.String("date", date), zap.Error(err))
			http.Error(w, "Invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		query := `{
			dailyMetrics(where: {date: "` + date + `"}) {
				id
				date
				dailyVolumeUSD
				dailyTransactionCount
				dailyActiveUsers
				dailyYieldDistributed
				dailySubsidiesDistributed
				averageAPY
				utilizationRate
				timestamp
			}
		}`

		var result struct {
			DailyMetrics []*gql.DailyMetrics `json:"dailyMetrics"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query daily analytics",
				zap.String("date", date),
				zap.Error(err))
			http.Error(w, "Failed to get daily analytics", http.StatusInternalServerError)
			return
		}

		response := DailyAnalyticsResponse{
			Date:                      date,
			DailyVolumeUSD:            "0",
			DailyTransactionCount:     "0",
			DailyActiveUsers:          "0",
			DailyYieldDistributed:     "0",
			DailySubsidiesDistributed: "0",
			AverageAPY:                "0",
			UtilizationRate:           "0",
			Timestamp:                 time.Now().Unix(),
		}

		if len(result.DailyMetrics) > 0 {
			metrics := result.DailyMetrics[0]
			if metrics.DailyVolumeUSD != nil {
				response.DailyVolumeUSD = metrics.DailyVolumeUSD.String()
			}
			if metrics.DailyTransactionCount != nil {
				response.DailyTransactionCount = metrics.DailyTransactionCount.String()
			}
			if metrics.DailyActiveUsers != nil {
				response.DailyActiveUsers = metrics.DailyActiveUsers.String()
			}
			if metrics.DailyYieldDistributed != nil {
				response.DailyYieldDistributed = metrics.DailyYieldDistributed.String()
			}
			if metrics.DailySubsidiesDistributed != nil {
				response.DailySubsidiesDistributed = metrics.DailySubsidiesDistributed.String()
			}
			if metrics.AverageAPY != nil {
				response.AverageAPY = metrics.AverageAPY.String()
			}
			if metrics.UtilizationRate != nil {
				response.UtilizationRate = metrics.UtilizationRate.String()
			}
			if metrics.Timestamp != nil {
				response.Timestamp = metrics.Timestamp.Int64()
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
