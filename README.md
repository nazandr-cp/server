# Go Rewards Server

This repository implements a lightweight orchestration layer for the
`EpochManager` and `CollectionsVault` contracts. The server exposes a small
REST API, consumes data from a The Graph subgraph and pushes blockchain events
to connected WebSocket clients.

```
┌────────────┐   GraphQL HTTP        ┌──────────────┐
│  Subgraph  │ <──────────────────── │ Go-Server    │
└────────────┘                        │ • chi router │
       ▲  ETH JSON-RPC WebSocket      │ • scheduler  │
       │                              │ • cache      │
       ▼                              └─────▲────────┘
┌────────────┐  signed TX / logs            │
│  Ethereum  │ ─────────────────────────────┘
└────────────┘
```

## Getting started

```
RPC_HTTP_URL=<http-url>
RPC_WS_URL=<ws-url>
PRIVATE_KEY=<hex>
SUBGRAPH_URL=<the-graph-url>
EPOCH_MANAGER_ADDR=<0x..>
VAULT_ADDR=<0x..>
CHAIN_ID=1
HTTP_PORT=8080
```

Run the server:

```
go run ./cmd/server
```

Available endpoints:

- `GET /healthz` – simple ping
- `GET /epochs/current` – last epoch info from the subgraph
- `GET /ws` – WebSocket stream of contract logs

This codebase is intentionally minimal and only relies on
[go-chi/chi](https://github.com/go-chi/chi) and
[go-ethereum](https://github.com/ethereum/go-ethereum).

## Docker

You can build a container image and run the server with Docker Compose:

```sh
docker compose up --build
```

The compose file exposes port `8080` and relies on the following environment variables, which you can adjust in `docker-compose.yml`:

- `RPC_HTTP_URL`
- `RPC_WS_URL`
- `PRIVATE_KEY`
- `SUBGRAPH_URL`
- `EPOCH_MANAGER_ADDR`
- `VAULT_ADDR`
- `CHAIN_ID`
- `HTTP_PORT`
