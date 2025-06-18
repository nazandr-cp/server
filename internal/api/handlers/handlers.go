package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	config "go-server/configs"
	"go-server/internal/gql"
	eth "go-server/internal/platform/ethereum"
	ws "go-server/internal/platform/websocket"
	"go-server/internal/service/collection"
	"go-server/internal/service/subsidy"
)

type Deps struct {
	Cfg               config.Config
	Eth               *eth.Clients
	Hub               *ws.Hub
	SubsidyService    *subsidy.Service
	CollectionService *collection.Service
	Logger            *zap.Logger
}

// AdminAuth is a middleware to protect admin routes
func AdminAuth(cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Invalid Authorization header format. Expected: Bearer <token>", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			if token != cfg.AdminToken {
				http.Error(w, "Invalid admin token", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func Register(r chi.Router, d Deps) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", health)
	r.Get("/ws", d.Hub.ServeWS)
	r.Get("/epochs/current", currentEpoch(d))

	// API v1 routes
	r.Route("/api/v1", func(apiRouter chi.Router) {
		// Subsidy endpoints
		apiRouter.Get("/epochs/{epochId}/eligibility/{userAddress}", GetEligibility(d))
		apiRouter.Get("/epochs/{epochId}/merkle-proof/{userAddress}", GetMerkleProof(d))
		apiRouter.Post("/claims/batch-verify", BatchVerifyClaims(d))
		apiRouter.Get("/users/{address}/claim-status", GetUserClaimStatus(d))

		// Analytics endpoints
		apiRouter.Get("/system/metrics", GetSystemMetrics(d))
		apiRouter.Get("/epochs/{epochId}/allocations", GetEpochAllocations(d))
		apiRouter.Get("/vaults/{address}/performance", GetVaultPerformance(d))
		apiRouter.Get("/analytics/daily/{date}", GetDailyAnalytics(d))

		// Enhanced epoch endpoints
		apiRouter.Get("/epochs/current", GetCurrentEpoch(d))
		apiRouter.Get("/epochs/{epochId}/details", GetEpochDetails(d))
		apiRouter.Get("/epochs/{epochId}/collections", GetEpochCollections(d))

		// Collection endpoints
		apiRouter.Get("/collections/{address}/participants", GetCollectionParticipants(d))
		apiRouter.Get("/collections/{address}/metrics", GetCollectionMetrics(d))
		apiRouter.Get("/collections/active", GetActiveCollections(d))
	})

	r.Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Use(AdminAuth(d.Cfg))
		adminRouter.Post("/subsidies/run", adminRunSubsidyHandler(d))
	})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func currentEpoch(d Deps) http.HandlerFunc {
	type response struct {
		Data any `json:"data"`
	}
	query := `{
        epochs(orderBy: id, orderDirection: desc, first: 1){id finalized}
    }`
	return func(w http.ResponseWriter, r *http.Request) {
		var out response
		if err := gql.Query(r.Context(), d.Cfg.SubgraphURL, query, &out); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	}
}

func adminRunSubsidyHandler(d Deps) http.HandlerFunc {
	type response struct {
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
		Epoch   uint64 `json:"epoch,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		epochStr := r.URL.Query().Get("epoch")
		if epochStr == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{Status: "error", Message: "epoch query parameter is required"})
			return
		}

		epoch, err := strconv.ParseUint(epochStr, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{Status: "error", Message: "invalid epoch format, must be a non-negative integer"})
			return
		}

		d.Logger.Info("Admin request to run subsidy for epoch", zap.Uint64("epoch", epoch))

		go func() {
			ctx := r.Context()
			err := d.SubsidyService.Run(ctx, epoch)
			if err != nil {
				d.Logger.Error("Failed to run subsidy service from admin endpoint", zap.Uint64("epoch", epoch), zap.Error(err))
			} else {
				d.Logger.Info("Subsidy service run initiated successfully from admin endpoint", zap.Uint64("epoch", epoch))
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{
			Status:  "success",
			Message: "Subsidy run initiated for epoch. Check logs for status.",
			Epoch:   epoch,
		})
	}
}
