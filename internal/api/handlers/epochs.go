package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"go-server/internal/gql"
)

// CurrentEpochResponse is current epoch information.
type CurrentEpochResponse struct {
	EpochID                   string `json:"epochId"`
	EpochNumber               string `json:"epochNumber"`
	Status                    string `json:"status"`
	StartTimestamp            int64  `json:"startTimestamp"`
	EndTimestamp              int64  `json:"endTimestamp"`
	TotalYieldAvailable       string `json:"totalYieldAvailable"`
	TotalYieldAllocated       string `json:"totalYieldAllocated"`
	TotalSubsidiesDistributed string `json:"totalSubsidiesDistributed"`
	TotalEligibleUsers        string `json:"totalEligibleUsers"`
	ParticipantCount          string `json:"participantCount"`
	ProcessingTimeMs          string `json:"processingTimeMs"`
}

// EpochDetailsResponse is detailed epoch information.
type EpochDetailsResponse struct {
	EpochID                       string                `json:"epochId"`
	EpochNumber                   string                `json:"epochNumber"`
	Status                        string                `json:"status"`
	StartTimestamp                int64                 `json:"startTimestamp"`
	EndTimestamp                  int64                 `json:"endTimestamp"`
	ProcessingStartedTimestamp    int64                 `json:"processingStartedTimestamp"`
	ProcessingCompletedTimestamp  int64                 `json:"processingCompletedTimestamp"`
	TotalYieldAvailable           string                `json:"totalYieldAvailable"`
	TotalYieldAllocated           string                `json:"totalYieldAllocated"`
	TotalYieldDistributed         string                `json:"totalYieldDistributed"`
	RemainingYield                string                `json:"remainingYield"`
	TotalSubsidiesDistributed     string                `json:"totalSubsidiesDistributed"`
	TotalEligibleUsers            string                `json:"totalEligibleUsers"`
	TotalParticipatingCollections string                `json:"totalParticipatingCollections"`
	ParticipantCount              string                `json:"participantCount"`
	ProcessingTimeMs              string                `json:"processingTimeMs"`
	EstimatedProcessingTime       string                `json:"estimatedProcessingTime"`
	ProcessingGasUsed             string                `json:"processingGasUsed"`
	ProcessingTransactionCount    string                `json:"processingTransactionCount"`
	VaultAllocations              []VaultAllocationInfo `json:"vaultAllocations"`
}

// EpochCollectionsResponse is collections participating in an epoch.
type EpochCollectionsResponse struct {
	EpochID     string                `json:"epochId"`
	Collections []EpochCollectionInfo `json:"collections"`
}

// EpochCollectionInfo is collection info for an epoch.
type EpochCollectionInfo struct {
	CollectionAddress string `json:"collectionAddress"`
	CollectionName    string `json:"collectionName"`
	CollectionSymbol  string `json:"collectionSymbol"`
	TotalParticipants string `json:"totalParticipants"`
	TotalSubsidies    string `json:"totalSubsidies"`
	YieldShare        string `json:"yieldShare"`
	IsActive          bool   `json:"isActive"`
}

// GetCurrentEpoch handles GET /api/v1/epochs/current
func GetCurrentEpoch(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `{
			epochs(first: 1, orderBy: epochNumber, orderDirection: desc, where: {status: ACTIVE}) {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				totalYieldAvailable
				totalYieldAllocated
				totalSubsidiesDistributed
				totalEligibleUsers
				participantCount
				processingTimeMs
			}
		}`

		var result struct {
			Epochs []*gql.Epoch `json:"epochs"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query current epoch", zap.Error(err))
			http.Error(w, "Failed to get current epoch", http.StatusInternalServerError)
			return
		}

		if len(result.Epochs) == 0 {
			d.Logger.Info("No active epoch found")
			http.Error(w, "No active epoch found", http.StatusNotFound)
			return
		}

		epoch := result.Epochs[0]
		response := CurrentEpochResponse{
			EpochID:                   epoch.ID,
			Status:                    string(epoch.Status),
			TotalYieldAvailable:       "0",
			TotalYieldAllocated:       "0",
			TotalSubsidiesDistributed: "0",
			TotalEligibleUsers:        "0",
			ParticipantCount:          "0",
			ProcessingTimeMs:          "0",
		}

		if epoch.EpochNumber != nil {
			response.EpochNumber = epoch.EpochNumber.String()
		}
		if epoch.StartTimestamp != nil {
			response.StartTimestamp = epoch.StartTimestamp.Int64()
		}
		if epoch.EndTimestamp != nil {
			response.EndTimestamp = epoch.EndTimestamp.Int64()
		}
		if epoch.TotalYieldAvailable != nil {
			response.TotalYieldAvailable = epoch.TotalYieldAvailable.String()
		}
		if epoch.TotalYieldAllocated != nil {
			response.TotalYieldAllocated = epoch.TotalYieldAllocated.String()
		}
		if epoch.TotalSubsidiesDistributed != nil {
			response.TotalSubsidiesDistributed = epoch.TotalSubsidiesDistributed.String()
		}
		if epoch.TotalEligibleUsers != nil {
			response.TotalEligibleUsers = epoch.TotalEligibleUsers.String()
		}
		if epoch.ParticipantCount != nil {
			response.ParticipantCount = epoch.ParticipantCount.String()
		}
		if epoch.ProcessingTimeMs != nil {
			response.ProcessingTimeMs = epoch.ProcessingTimeMs.String()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetEpochDetails handles GET /api/v1/epochs/{epochId}/details
func GetEpochDetails(d Deps) http.HandlerFunc {
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
				status
				startTimestamp
				endTimestamp
				processingStartedTimestamp
				processingCompletedTimestamp
				totalYieldAvailable
				totalYieldAllocated
				totalYieldDistributed
				remainingYield
				totalSubsidiesDistributed
				totalEligibleUsers
				totalParticipatingCollections
				participantCount
				processingTimeMs
				estimatedProcessingTime
				processingGasUsed
				processingTransactionCount
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
			d.Logger.Error("Failed to query epoch details",
				zap.Uint64("epochId", epochId),
				zap.Error(err))
			http.Error(w, "Failed to get epoch details", http.StatusInternalServerError)
			return
		}

		if len(result.Epochs) == 0 {
			d.Logger.Info("Epoch not found", zap.Uint64("epochId", epochId))
			http.Error(w, "Epoch not found", http.StatusNotFound)
			return
		}

		epoch := result.Epochs[0]
		response := EpochDetailsResponse{
			EpochID:                       epoch.ID,
			Status:                        string(epoch.Status),
			TotalYieldAvailable:           "0",
			TotalYieldAllocated:           "0",
			TotalYieldDistributed:         "0",
			RemainingYield:                "0",
			TotalSubsidiesDistributed:     "0",
			TotalEligibleUsers:            "0",
			TotalParticipatingCollections: "0",
			ParticipantCount:              "0",
			ProcessingTimeMs:              "0",
			EstimatedProcessingTime:       "0",
			ProcessingGasUsed:             "0",
			ProcessingTransactionCount:    "0",
			VaultAllocations:              []VaultAllocationInfo{},
		}

		if epoch.EpochNumber != nil {
			response.EpochNumber = epoch.EpochNumber.String()
		}
		if epoch.StartTimestamp != nil {
			response.StartTimestamp = epoch.StartTimestamp.Int64()
		}
		if epoch.EndTimestamp != nil {
			response.EndTimestamp = epoch.EndTimestamp.Int64()
		}
		if epoch.ProcessingStartedTimestamp != nil {
			response.ProcessingStartedTimestamp = epoch.ProcessingStartedTimestamp.Int64()
		}
		if epoch.ProcessingCompletedTimestamp != nil {
			response.ProcessingCompletedTimestamp = epoch.ProcessingCompletedTimestamp.Int64()
		}
		if epoch.TotalYieldAvailable != nil {
			response.TotalYieldAvailable = epoch.TotalYieldAvailable.String()
		}
		if epoch.TotalYieldAllocated != nil {
			response.TotalYieldAllocated = epoch.TotalYieldAllocated.String()
		}
		if epoch.TotalYieldDistributed != nil {
			response.TotalYieldDistributed = epoch.TotalYieldDistributed.String()
		}
		if epoch.RemainingYield != nil {
			response.RemainingYield = epoch.RemainingYield.String()
		}
		if epoch.TotalSubsidiesDistributed != nil {
			response.TotalSubsidiesDistributed = epoch.TotalSubsidiesDistributed.String()
		}
		if epoch.TotalEligibleUsers != nil {
			response.TotalEligibleUsers = epoch.TotalEligibleUsers.String()
		}
		if epoch.TotalParticipatingCollections != nil {
			response.TotalParticipatingCollections = epoch.TotalParticipatingCollections.String()
		}
		if epoch.ParticipantCount != nil {
			response.ParticipantCount = epoch.ParticipantCount.String()
		}
		if epoch.ProcessingTimeMs != nil {
			response.ProcessingTimeMs = epoch.ProcessingTimeMs.String()
		}
		if epoch.EstimatedProcessingTime != nil {
			response.EstimatedProcessingTime = epoch.EstimatedProcessingTime.String()
		}
		if epoch.ProcessingGasUsed != nil {
			response.ProcessingGasUsed = epoch.ProcessingGasUsed.String()
		}
		if epoch.ProcessingTransactionCount != nil {
			response.ProcessingTransactionCount = epoch.ProcessingTransactionCount.String()
		}

		// Add vault allocations
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetEpochCollections handles GET /api/v1/epochs/{epochId}/collections
func GetEpochCollections(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		epochIdStr := chi.URLParam(r, "epochId")

		epochId, err := strconv.ParseUint(epochIdStr, 10, 64)
		if err != nil {
			d.Logger.Error("Invalid epoch ID", zap.String("epochId", epochIdStr), zap.Error(err))
			http.Error(w, "Invalid epoch ID", http.StatusBadRequest)
			return
		}

		query := `{
			userEpochEligibilities(where: {epoch: "` + epochIdStr + `"}, first: 1000) {
				id
				collection {
					id
					contractAddress
					name
					symbol
					isActive
				}
				subsidyReceived
				yieldShare
			}
		}`

		var result struct {
			UserEpochEligibilities []*gql.UserEpochEligibility `json:"userEpochEligibilities"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query epoch collections",
				zap.Uint64("epochId", epochId),
				zap.Error(err))
			http.Error(w, "Failed to get epoch collections", http.StatusInternalServerError)
			return
		}

		// Aggregate collection data
		collectionMap := make(map[string]*EpochCollectionInfo)

		for _, eligibility := range result.UserEpochEligibilities {
			if eligibility.Collection == nil {
				continue
			}

			collectionAddr := eligibility.Collection.ContractAddress
			if _, exists := collectionMap[collectionAddr]; !exists {
				collectionMap[collectionAddr] = &EpochCollectionInfo{
					CollectionAddress: collectionAddr,
					CollectionName:    eligibility.Collection.Name,
					CollectionSymbol:  eligibility.Collection.Symbol,
					TotalParticipants: "0",
					TotalSubsidies:    "0",
					YieldShare:        "0",
					IsActive:          eligibility.Collection.IsActive,
				}
			}

			// Increment participant count
			currentParticipants, _ := strconv.ParseInt(collectionMap[collectionAddr].TotalParticipants, 10, 64)
			collectionMap[collectionAddr].TotalParticipants = strconv.FormatInt(currentParticipants+1, 10)

			// Add subsidies
			if eligibility.SubsidyReceived != nil {
				currentSubsidies, _ := strconv.ParseInt(collectionMap[collectionAddr].TotalSubsidies, 10, 64)
				subsidyAmount := eligibility.SubsidyReceived.Int64()
				collectionMap[collectionAddr].TotalSubsidies = strconv.FormatInt(currentSubsidies+subsidyAmount, 10)
			}

			// Add yield share
			if eligibility.YieldShare != nil {
				currentYield, _ := strconv.ParseInt(collectionMap[collectionAddr].YieldShare, 10, 64)
				yieldAmount := eligibility.YieldShare.Int64()
				collectionMap[collectionAddr].YieldShare = strconv.FormatInt(currentYield+yieldAmount, 10)
			}
		}

		// Convert map to slice
		collections := make([]EpochCollectionInfo, 0, len(collectionMap))
		for _, info := range collectionMap {
			collections = append(collections, *info)
		}

		response := EpochCollectionsResponse{
			EpochID:     epochIdStr,
			Collections: collections,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
