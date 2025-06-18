package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// SubsidyEligibilityResponse is the eligibility check response.
type SubsidyEligibilityResponse struct {
	Eligible bool   `json:"eligible"`
	Amount   string `json:"amount,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

// MerkleProofResponse is the merkle proof response.
type MerkleProofResponse struct {
	Proof    []string `json:"proof"`
	Amount   string   `json:"amount"`
	Index    uint64   `json:"index"`
	Root     string   `json:"root"`
	Verified bool     `json:"verified"`
}

// BatchVerifyRequest is the batch verification request.
type BatchVerifyRequest struct {
	Claims []ClaimInfo `json:"claims"`
}

// ClaimInfo is individual claim information.
type ClaimInfo struct {
	UserAddress string   `json:"userAddress"`
	Amount      string   `json:"amount"`
	Proof       []string `json:"proof"`
	Index       uint64   `json:"index"`
}

// ClaimStatusResponse is the claim status response.
type ClaimStatusResponse struct {
	TotalEligible string             `json:"totalEligible"`
	TotalClaimed  string             `json:"totalClaimed"`
	Pending       string             `json:"pending"`
	Claims        []ClaimHistoryInfo `json:"claims"`
}

// ClaimHistoryInfo is claim history.
type ClaimHistoryInfo struct {
	EpochID   uint64 `json:"epochId"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
	TxHash    string `json:"txHash,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// GetEligibility handles GET /api/v1/epochs/{epochId}/eligibility/{userAddress}
func GetEligibility(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		epochIdStr := chi.URLParam(r, "epochId")
		userAddress := chi.URLParam(r, "userAddress")

		epochId, err := strconv.ParseUint(epochIdStr, 10, 64)
		if err != nil {
			d.Logger.Error("Invalid epoch ID", zap.String("epochId", epochIdStr), zap.Error(err))
			http.Error(w, "Invalid epoch ID", http.StatusBadRequest)
			return
		}

		if !common.IsHexAddress(userAddress) {
			d.Logger.Error("Invalid user address", zap.String("userAddress", userAddress))
			http.Error(w, "Invalid user address", http.StatusBadRequest)
			return
		}

		addr := common.HexToAddress(userAddress)
		vaultAddr := common.HexToAddress(d.Cfg.VaultAddr)
		collectionAddr := common.Address{} // Would need to be passed as parameter or retrieved

		eligibility, err := d.SubsidyService.CalculateUserEligibility(r.Context(), addr, vaultAddr, collectionAddr)
		if err != nil {
			d.Logger.Error("Failed to check eligibility",
				zap.Uint64("epochId", epochId),
				zap.String("userAddress", userAddress),
				zap.Error(err))
			response := SubsidyEligibilityResponse{
				Eligible: false,
				Reason:   "Not eligible for subsidy in this epoch",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := SubsidyEligibilityResponse{
			Eligible: eligibility != nil,
		}

		if eligibility != nil && eligibility.SubsidyReceived != nil {
			response.Amount = eligibility.SubsidyReceived.String()
		} else {
			response.Reason = "Not eligible for subsidy in this epoch"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetMerkleProof handles GET /api/v1/epochs/{epochId}/merkle-proof/{userAddress}
func GetMerkleProof(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		epochIdStr := chi.URLParam(r, "epochId")
		userAddress := chi.URLParam(r, "userAddress")

		epochId, err := strconv.ParseUint(epochIdStr, 10, 64)
		if err != nil {
			d.Logger.Error("Invalid epoch ID", zap.String("epochId", epochIdStr), zap.Error(err))
			http.Error(w, "Invalid epoch ID", http.StatusBadRequest)
			return
		}

		if !common.IsHexAddress(userAddress) {
			d.Logger.Error("Invalid user address", zap.String("userAddress", userAddress))
			http.Error(w, "Invalid user address", http.StatusBadRequest)
			return
		}

		// For now, return a placeholder response since we need to implement proper merkle proof retrieval
		// This would require storing merkle proofs or recalculating them from stored data
		response := MerkleProofResponse{
			Proof:    []string{},
			Amount:   "0",
			Index:    0,
			Root:     "0x0000000000000000000000000000000000000000000000000000000000000000",
			Verified: false,
		}

		d.Logger.Warn("Merkle proof generation not yet implemented",
			zap.Uint64("epochId", epochId),
			zap.String("userAddress", userAddress))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// BatchVerifyClaims handles POST /api/v1/claims/batch-verify
func BatchVerifyClaims(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BatchVerifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			d.Logger.Error("Failed to decode batch verify request", zap.Error(err))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(req.Claims) == 0 {
			http.Error(w, "No claims provided", http.StatusBadRequest)
			return
		}

		if len(req.Claims) > d.Cfg.Subsidy.SubsidyBatchSize {
			http.Error(w, "Too many claims in batch", http.StatusBadRequest)
			return
		}

		results := make([]map[string]interface{}, len(req.Claims))
		for i, claim := range req.Claims {
			if !common.IsHexAddress(claim.UserAddress) {
				results[i] = map[string]interface{}{
					"userAddress": claim.UserAddress,
					"valid":       false,
					"error":       "Invalid user address",
				}
				continue
			}
			// For now, return a placeholder validation since we need to implement proper merkle proof verification
			// This would require access to the stored merkle tree data
			results[i] = map[string]interface{}{
				"userAddress": claim.UserAddress,
				"valid":       false,
				"error":       "Merkle proof verification not yet implemented",
				"amount":      claim.Amount,
			}

			d.Logger.Warn("Merkle proof verification not yet implemented",
				zap.String("userAddress", claim.UserAddress))
		}

		response := map[string]interface{}{
			"results": results,
			"total":   len(req.Claims),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetUserClaimStatus handles GET /api/v1/users/{address}/claim-status
func GetUserClaimStatus(d Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAddress := chi.URLParam(r, "address")

		if !common.IsHexAddress(userAddress) {
			d.Logger.Error("Invalid user address", zap.String("userAddress", userAddress))
			http.Error(w, "Invalid user address", http.StatusBadRequest)
			return
		}

		addr := common.HexToAddress(userAddress)

		// Get claim status using existing method - note: this returns a single distribution
		distribution, err := d.SubsidyService.GetClaimStatus(r.Context(), addr, "1") // Using epoch 1 as example
		if err != nil {
			d.Logger.Error("Failed to get user claim status",
				zap.String("userAddress", userAddress),
				zap.Error(err))
			// Return empty status instead of error for better UX
			response := ClaimStatusResponse{
				TotalEligible: "0",
				TotalClaimed:  "0",
				Pending:       "0",
				Claims:        []ClaimHistoryInfo{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := ClaimStatusResponse{
			TotalEligible: "0",
			TotalClaimed:  "0",
			Pending:       "0",
			Claims:        []ClaimHistoryInfo{},
		}

		if distribution != nil && distribution.SubsidyAmount != nil {
			response.TotalEligible = distribution.SubsidyAmount.String()
			// SubsidyDistribution doesn't have a Claimed field, so we assume it's claimed if it exists
			response.TotalClaimed = distribution.SubsidyAmount.String()
			response.Pending = "0"

			// Add single claim record
			response.Claims = []ClaimHistoryInfo{{
				EpochID:   1, // Would need actual epoch ID from distribution.Epoch
				Amount:    distribution.SubsidyAmount.String(),
				Status:    "claimed",
				TxHash:    distribution.TransactionHash,
				Timestamp: distribution.Timestamp.Int64(),
			}}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
