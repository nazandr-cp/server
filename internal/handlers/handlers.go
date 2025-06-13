package handlers

import (
	"encoding/json"
	"net/http"

	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"go-server/internal/config"
	"go-server/internal/eth"
	"go-server/internal/gql"
	"go-server/internal/ws"
	"go-server/services/subsidy"
)

type Deps struct {
	Cfg            config.Config
	Eth            *eth.Clients
	Hub            *ws.Hub
	SubsidyService *subsidy.Service
	Logger         *zap.Logger
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

	r.Get("/healthz", healthz)
	r.Get("/ws", d.Hub.ServeWS)
	r.Get("/epochs/current", currentEpoch(d))

	r.Route("/admin", func(adminRouter chi.Router) {
		adminRouter.Use(AdminAuth(d.Cfg))
		adminRouter.Post("/subsidies/run", adminRunSubsidyHandler(d))
	})
}

func healthz(w http.ResponseWriter, r *http.Request) {
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
