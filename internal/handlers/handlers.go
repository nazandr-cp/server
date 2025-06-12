package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go-server/internal/config"
	"go-server/internal/eth"
	"go-server/internal/gql"
	"go-server/internal/ws"
)

type Deps struct {
	Cfg config.Config
	Eth *eth.Clients
	Hub *ws.Hub
}

func Register(r chi.Router, d Deps) {
	r.Get("/healthz", healthz)
	r.Get("/ws", d.Hub.ServeWS)
	r.Get("/epochs/current", currentEpoch(d))
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
		json.NewEncoder(w).Encode(out)
	}
}
