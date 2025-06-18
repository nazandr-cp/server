package handlers

import (
	"encoding/json"
	"fmt"
	"math/big"
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

		// Check user's eligibility by querying subsidy distributions for this epoch
		distributions, err := d.SubsidyService.GetUserSubsidyDistributions(r.Context(), addr, strconv.FormatUint(epochId, 10))
		if err != nil {
			d.Logger.Error("Failed to check eligibility",
				zap.Uint64("epochId", epochId),
				zap.String("userAddress", userAddress),
				zap.Error(err))
			response := SubsidyEligibilityResponse{
				Eligible: false,
				Reason:   "Failed to check eligibility",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		response := SubsidyEligibilityResponse{
			Eligible: len(distributions) > 0,
		}

		if len(distributions) > 0 {
			// Calculate total eligible amount
			totalAmount := big.NewInt(0)
			for _, dist := range distributions {
				if dist.SubsidyAmount != nil {
					totalAmount.Add(totalAmount, dist.SubsidyAmount)
				}
			}
			response.Amount = totalAmount.String()
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

		addr := common.HexToAddress(userAddress)
		vaultID := d.Cfg.VaultAddr

		// Get merkle distribution for this epoch and vault
		merkleDistribution, err := d.SubsidyService.GetMerkleDistribution(r.Context(), epochIdStr, vaultID)
		if err != nil {
			d.Logger.Error("Failed to get merkle distribution",
				zap.Uint64("epochId", epochId),
				zap.String("vaultId", vaultID),
				zap.Error(err))
			response := MerkleProofResponse{
				Proof:    []string{},
				Amount:   "0",
				Index:    0,
				Root:     "0x0000000000000000000000000000000000000000000000000000000000000000",
				Verified: false,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Get user's subsidy distributions for this epoch
		distributions, err := d.SubsidyService.GetUserSubsidyDistributions(r.Context(), addr, epochIdStr)
		if err != nil || len(distributions) == 0 {
			d.Logger.Error("No distributions found for user",
				zap.Uint64("epochId", epochId),
				zap.String("userAddress", userAddress),
				zap.Error(err))
			response := MerkleProofResponse{
				Proof:    []string{},
				Amount:   "0",
				Index:    0,
				Root:     merkleDistribution.MerkleRoot,
				Verified: false,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Calculate total amount for this user
		totalAmount := big.NewInt(0)
		for _, dist := range distributions {
			if dist.SubsidyAmount != nil {
				totalAmount.Add(totalAmount, dist.SubsidyAmount)
			}
		}

		// Generate merkle proof for this user and amount
		// Note: This requires the original recipients list from the epoch processing
		// For now, we'll return basic info and indicate that proof generation needs the full recipient set
		response := MerkleProofResponse{
			Proof:    []string{}, // Would need full recipient list to regenerate proof
			Amount:   totalAmount.String(),
			Index:    0, // Would need to find user's index in the original merkle tree
			Root:     merkleDistribution.MerkleRoot,
			Verified: false, // Set to false since we can't verify without the full proof
		}

		d.Logger.Info("Merkle proof requested but requires full recipient data from epoch processing",
			zap.Uint64("epochId", epochId),
			zap.String("userAddress", userAddress),
			zap.String("amount", totalAmount.String()))

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

			// Parse amount
			amount := big.NewInt(0)
			if _, ok := amount.SetString(claim.Amount, 10); !ok {
				results[i] = map[string]interface{}{
					"userAddress": claim.UserAddress,
					"valid":       false,
					"error":       "Invalid amount format",
				}
				continue
			}

			// Verify merkle proof
			userAddr := common.HexToAddress(claim.UserAddress)

			// For verification, we need to get the merkle root for this claim
			// This would typically be stored with the claim or retrieved from the epoch data
			// For now, we'll attempt to find the latest merkle root from vault distributions

			// Try to get current epoch merkle distribution
			vaultID := d.Cfg.VaultAddr

			// Get latest merkle distributions to find a root to verify against
			merkleDistributions, err := d.SubsidyService.GetLatestMerkleDistributions(r.Context(), vaultID)
			if err != nil || len(merkleDistributions) == 0 {
				results[i] = map[string]interface{}{
					"userAddress": claim.UserAddress,
					"valid":       false,
					"error":       "No merkle root found for verification",
					"amount":      claim.Amount,
				}
				continue
			}

			// Use the most recent merkle root
			merkleRoot := merkleDistributions[0].MerkleRoot

			// Verify the proof
			verified, err := d.SubsidyService.VerifyMerkleProof(claim.Proof, merkleRoot, userAddr, amount)
			if err != nil {
				results[i] = map[string]interface{}{
					"userAddress": claim.UserAddress,
					"valid":       false,
					"error":       fmt.Sprintf("Verification error: %v", err),
					"amount":      claim.Amount,
				}
				continue
			}

			results[i] = map[string]interface{}{
				"userAddress": claim.UserAddress,
				"valid":       verified,
				"amount":      claim.Amount,
				"merkleRoot":  merkleRoot,
			}

			if verified {
				d.Logger.Info("Merkle proof verified successfully",
					zap.String("userAddress", claim.UserAddress),
					zap.String("amount", claim.Amount))
			} else {
				d.Logger.Warn("Merkle proof verification failed",
					zap.String("userAddress", claim.UserAddress),
					zap.String("amount", claim.Amount))
			}
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

		// Get all historical claims for this user
		claims, err := d.SubsidyService.GetUserClaimHistory(r.Context(), addr)
		if err != nil {
			d.Logger.Error("Failed to get user claim history",
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

		// Calculate totals and build claim history
		totalEligible := big.NewInt(0)
		totalClaimed := big.NewInt(0)
		var claimHistory []ClaimHistoryInfo

		for _, claim := range claims {
			if claim.SubsidyAmount != nil {
				totalEligible.Add(totalEligible, claim.SubsidyAmount)
				totalClaimed.Add(totalClaimed, claim.SubsidyAmount)

				epochID := uint64(0)
				if claim.Epoch != nil && claim.Epoch.EpochNumber != nil {
					epochID = claim.Epoch.EpochNumber.Uint64()
				}

				timestamp := int64(0)
				if claim.Timestamp != nil {
					timestamp = claim.Timestamp.Int64()
				}

				claimHistory = append(claimHistory, ClaimHistoryInfo{
					EpochID:   epochID,
					Amount:    claim.SubsidyAmount.String(),
					Status:    "claimed",
					TxHash:    claim.TransactionHash,
					Timestamp: timestamp,
				})
			}
		}

		// For now, pending is 0 since we're only looking at historical claims
		// In a more advanced implementation, we could check for unclaimed eligible amounts
		pending := big.NewInt(0)

		response := ClaimStatusResponse{
			TotalEligible: totalEligible.String(),
			TotalClaimed:  totalClaimed.String(),
			Pending:       pending.String(),
			Claims:        claimHistory,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
