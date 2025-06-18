package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"lend.fam/go-server/internal/gql"
)

// CollectionParticipantsResponse is collection participants.
type CollectionParticipantsResponse struct {
	CollectionAddress string                      `json:"collectionAddress"`
	TotalParticipants int                         `json:"totalParticipants"`
	Participants      []CollectionParticipantInfo `json:"participants"`
}

// CollectionParticipantInfo is individual participant info.
type CollectionParticipantInfo struct {
	UserAddress           string `json:"userAddress"`
	PrincipalShares       string `json:"principalShares"`
	PrincipalDeposited    string `json:"principalDeposited"`
	YieldAccrued          string `json:"yieldAccrued"`
	YieldClaimed          string `json:"yieldClaimed"`
	TotalSubsidies        string `json:"totalSubsidies"`
	TotalSubsidiesClaimed string `json:"totalSubsidiesClaimed"`
	AverageAPY            string `json:"averageAPY"`
	ParticipationStart    int64  `json:"participationStart"`
	LastUpdated           int64  `json:"lastUpdated"`
}

// CollectionMetricsResponse is collection metrics.
type CollectionMetricsResponse struct {
	CollectionAddress    string `json:"collectionAddress"`
	CollectionName       string `json:"collectionName"`
	CollectionSymbol     string `json:"collectionSymbol"`
	TotalParticipants    string `json:"totalParticipants"`
	TotalYieldGenerated  string `json:"totalYieldGenerated"`
	TotalSubsidies       string `json:"totalSubsidies"`
	TotalBorrowVolume    string `json:"totalBorrowVolume"`
	TotalNFTsDeposited   string `json:"totalNFTsDeposited"`
	YieldSharePercentage string `json:"yieldSharePercentage"`
	AverageAPY           string `json:"averageAPY"`
	IsActive             bool   `json:"isActive"`
	RegisteredTimestamp  int64  `json:"registeredTimestamp"`
	LastUpdated          int64  `json:"lastUpdated"`
}

// ActiveCollectionsResponse is active collections.
type ActiveCollectionsResponse struct {
	TotalCollections int                    `json:"totalCollections"`
	Collections      []ActiveCollectionInfo `json:"collections"`
}

// ActiveCollectionInfo is active collection info.
type ActiveCollectionInfo struct {
	CollectionAddress    string `json:"collectionAddress"`
	CollectionName       string `json:"collectionName"`
	CollectionSymbol     string `json:"collectionSymbol"`
	CollectionType       string `json:"collectionType"`
	TotalParticipants    string `json:"totalParticipants"`
	TotalYieldGenerated  string `json:"totalYieldGenerated"`
	TotalSubsidies       string `json:"totalSubsidies"`
	YieldSharePercentage string `json:"yieldSharePercentage"`
	IsActive             bool   `json:"isActive"`
}

// GetCollectionParticipants handles GET /api/v1/collections/{address}/participants
func GetCollectionParticipants(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collectionAddress := chi.URLParam(r, "address")

		if !common.IsHexAddress(collectionAddress) {
			d.Logger.Error("Invalid collection address", zap.String("collectionAddress", collectionAddress))
			http.Error(w, "Invalid collection address", http.StatusBadRequest)
			return
		}

		// Query collection participations directly from subgraph since CollectionService is not in Deps yet
		query := `{
			collectionParticipations(where: {collection: "` + collectionAddress + `"}, first: 1000) {
				id
				principalShares
				principalDeposited
				yieldAccrued
				yieldClaimed
				totalSubsidies
				totalSubsidiesClaimed
				averageAPY
				createdAtTimestamp
				updatedAtTimestamp
			}
		}`

		var result struct {
			CollectionParticipations []*gql.CollectionParticipation `json:"collectionParticipations"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query collection participants",
				zap.String("collectionAddress", collectionAddress),
				zap.Error(err))
			http.Error(w, "Failed to get collection participants", http.StatusInternalServerError)
			return
		}

		response := CollectionParticipantsResponse{
			CollectionAddress: collectionAddress,
			TotalParticipants: len(result.CollectionParticipations),
			Participants:      make([]CollectionParticipantInfo, len(result.CollectionParticipations)),
		}

		for i, participant := range result.CollectionParticipations {
			participantInfo := CollectionParticipantInfo{
				UserAddress:           participant.ID, // Using participation ID as user identifier
				PrincipalShares:       "0",
				PrincipalDeposited:    "0",
				YieldAccrued:          "0",
				YieldClaimed:          "0",
				TotalSubsidies:        "0",
				TotalSubsidiesClaimed: "0",
				AverageAPY:            "0",
				ParticipationStart:    0,
				LastUpdated:           0,
			}

			if participant.PrincipalShares != nil {
				participantInfo.PrincipalShares = participant.PrincipalShares.String()
			}
			if participant.PrincipalDeposited != nil {
				participantInfo.PrincipalDeposited = participant.PrincipalDeposited.String()
			}
			if participant.YieldAccrued != nil {
				participantInfo.YieldAccrued = participant.YieldAccrued.String()
			}
			if participant.YieldClaimed != nil {
				participantInfo.YieldClaimed = participant.YieldClaimed.String()
			}
			if participant.TotalSubsidies != nil {
				participantInfo.TotalSubsidies = participant.TotalSubsidies.String()
			}
			if participant.TotalSubsidiesClaimed != nil {
				participantInfo.TotalSubsidiesClaimed = participant.TotalSubsidiesClaimed.String()
			}
			if participant.AverageAPY != nil {
				participantInfo.AverageAPY = participant.AverageAPY.String()
			}
			if participant.CreatedAtTimestamp != nil {
				participantInfo.ParticipationStart = participant.CreatedAtTimestamp.Int64()
			}
			if participant.UpdatedAtTimestamp != nil {
				participantInfo.LastUpdated = participant.UpdatedAtTimestamp.Int64()
			}

			response.Participants[i] = participantInfo
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetCollectionMetrics handles GET /api/v1/collections/{address}/metrics
func GetCollectionMetrics(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collectionAddress := chi.URLParam(r, "address")

		if !common.IsHexAddress(collectionAddress) {
			d.Logger.Error("Invalid collection address", zap.String("collectionAddress", collectionAddress))
			http.Error(w, "Invalid collection address", http.StatusBadRequest)
			return
		}

		// Query collection details and participations to calculate metrics
		query := `{
			collections(where: {contractAddress: "` + collectionAddress + `"}) {
				id
				contractAddress
				name
				symbol
				totalYieldGenerated
				totalSubsidiesReceived
				totalBorrowVolume
				totalNFTsDeposited
				yieldSharePercentage
				isActive
				registeredAtTimestamp
				updatedAtTimestamp
				participations {
					id
					totalParticipants
					yieldAccrued
					totalSubsidies
					averageAPY
				}
			}
		}`

		var result struct {
			Collections []*gql.Collection `json:"collections"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query collection details",
				zap.String("collectionAddress", collectionAddress),
				zap.Error(err))
			http.Error(w, "Failed to get collection details", http.StatusInternalServerError)
			return
		}

		response := CollectionMetricsResponse{
			CollectionAddress:    collectionAddress,
			CollectionName:       "Unknown",
			CollectionSymbol:     "UNKNOWN",
			TotalParticipants:    "0",
			TotalYieldGenerated:  "0",
			TotalSubsidies:       "0",
			TotalBorrowVolume:    "0",
			TotalNFTsDeposited:   "0",
			YieldSharePercentage: "0",
			AverageAPY:           "0",
			IsActive:             false,
			RegisteredTimestamp:  0,
			LastUpdated:          0,
		}

		if len(result.Collections) > 0 {
			collection := result.Collections[0]
			response.CollectionName = collection.Name
			response.CollectionSymbol = collection.Symbol
			response.IsActive = collection.IsActive

			if collection.TotalYieldGenerated != nil {
				response.TotalYieldGenerated = collection.TotalYieldGenerated.String()
			}
			if collection.TotalSubsidiesReceived != nil {
				response.TotalSubsidies = collection.TotalSubsidiesReceived.String()
			}
			if collection.TotalBorrowVolume != nil {
				response.TotalBorrowVolume = collection.TotalBorrowVolume.String()
			}
			if collection.TotalNFTsDeposited != nil {
				response.TotalNFTsDeposited = collection.TotalNFTsDeposited.String()
			}
			if collection.YieldSharePercentage != nil {
				response.YieldSharePercentage = collection.YieldSharePercentage.String()
			}
			if collection.RegisteredAtTimestamp != nil {
				response.RegisteredTimestamp = collection.RegisteredAtTimestamp.Int64()
			}
			if collection.UpdatedAtTimestamp != nil {
				response.LastUpdated = collection.UpdatedAtTimestamp.Int64()
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetActiveCollections handles GET /api/v1/collections/active
func GetActiveCollections(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `{
			collections(where: {isActive: true}, first: 1000) {
				id
				contractAddress
				name
				symbol
				collectionType
				totalYieldGenerated
				totalSubsidiesReceived
				yieldSharePercentage
				isActive
				participations {
					id
					totalParticipants
				}
			}
		}`

		var result struct {
			Collections []*gql.Collection `json:"collections"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query active collections", zap.Error(err))
			http.Error(w, "Failed to get active collections", http.StatusInternalServerError)
			return
		}

		response := ActiveCollectionsResponse{
			TotalCollections: len(result.Collections),
			Collections:      make([]ActiveCollectionInfo, len(result.Collections)),
		}

		for i, collection := range result.Collections {
			collectionInfo := ActiveCollectionInfo{
				CollectionAddress:    collection.ContractAddress,
				CollectionName:       collection.Name,
				CollectionSymbol:     collection.Symbol,
				CollectionType:       string(collection.CollectionType),
				TotalParticipants:    "0",
				TotalYieldGenerated:  "0",
				TotalSubsidies:       "0",
				YieldSharePercentage: "0",
				IsActive:             collection.IsActive,
			}

			if collection.TotalYieldGenerated != nil {
				collectionInfo.TotalYieldGenerated = collection.TotalYieldGenerated.String()
			}
			if collection.TotalSubsidiesReceived != nil {
				collectionInfo.TotalSubsidies = collection.TotalSubsidiesReceived.String()
			}
			if collection.YieldSharePercentage != nil {
				collectionInfo.YieldSharePercentage = collection.YieldSharePercentage.String()
			}

			// Calculate total participants from participations
			totalParticipants := int64(0)
			for _, participation := range collection.Participations {
				if participation.TotalParticipants != nil {
					totalParticipants += participation.TotalParticipants.Int64()
				}
			}
			collectionInfo.TotalParticipants = strconv.FormatInt(totalParticipants, 10)

			response.Collections[i] = collectionInfo
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// CollectionDepositsResponse is collection deposits information.
type CollectionDepositsResponse struct {
	CollectionAddress string                  `json:"collectionAddress"`
	TotalDeposits     string                  `json:"totalDeposits"`
	TotalNFTs         string                  `json:"totalNFTs"`
	Deposits          []CollectionDepositInfo `json:"deposits"`
}

// CollectionDepositInfo is individual deposit info.
type CollectionDepositInfo struct {
	UserAddress     string `json:"userAddress"`
	NFTTokenId      string `json:"nftTokenId"`
	DepositAmount   string `json:"depositAmount"`
	Timestamp       int64  `json:"timestamp"`
	TransactionHash string `json:"transactionHash"`
	BlockNumber     string `json:"blockNumber"`
}

// GetCollectionDeposits handles GET /api/collections/{collection}/deposits
func GetCollectionDeposits(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collectionAddress := chi.URLParam(r, "collection")

		if !common.IsHexAddress(collectionAddress) {
			d.Logger.Error("Invalid collection address", zap.String("collectionAddress", collectionAddress))
			http.Error(w, "Invalid collection address", http.StatusBadRequest)
			return
		}

		query := `{
			collectionParticipations(where: {collection: "` + collectionAddress + `"}, first: 1000) {
				id
				principalDeposited
				collection {
					id
					contractAddress
					name
					totalNFTsDeposited
				}
				vault {
					id
				}
				createdAtTimestamp
				updatedAtTimestamp
			}
		}`

		var result struct {
			CollectionParticipations []*gql.CollectionParticipation `json:"collectionParticipations"`
		}

		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &result); err != nil {
			d.Logger.Error("Failed to query collection deposits",
				zap.String("collectionAddress", collectionAddress),
				zap.Error(err))
			http.Error(w, "Failed to get collection deposits", http.StatusInternalServerError)
			return
		}

		response := CollectionDepositsResponse{
			CollectionAddress: collectionAddress,
			TotalDeposits:     "0",
			TotalNFTs:         "0",
			Deposits:          []CollectionDepositInfo{},
		}

		totalDeposits := int64(0)
		totalNFTs := int64(0)

		for _, participation := range result.CollectionParticipations {
			if participation.PrincipalDeposited != nil {
				totalDeposits += participation.PrincipalDeposited.Int64()
			}

			depositInfo := CollectionDepositInfo{
				UserAddress:     participation.ID,
				NFTTokenId:      "N/A",
				DepositAmount:   "0",
				Timestamp:       0,
				TransactionHash: "N/A",
				BlockNumber:     "N/A",
			}

			if participation.PrincipalDeposited != nil {
				depositInfo.DepositAmount = participation.PrincipalDeposited.String()
			}
			if participation.CreatedAtTimestamp != nil {
				depositInfo.Timestamp = participation.CreatedAtTimestamp.Int64()
			}

			response.Deposits = append(response.Deposits, depositInfo)
		}

		if len(result.CollectionParticipations) > 0 && result.CollectionParticipations[0].Collection != nil {
			if result.CollectionParticipations[0].Collection.TotalNFTsDeposited != nil {
				totalNFTs = result.CollectionParticipations[0].Collection.TotalNFTsDeposited.Int64()
			}
		}

		response.TotalDeposits = strconv.FormatInt(totalDeposits, 10)
		response.TotalNFTs = strconv.FormatInt(totalNFTs, 10)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
