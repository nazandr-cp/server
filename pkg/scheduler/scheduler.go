package scheduler

import (
	"context"
	"time"

	eth "lend.fam/go-server/internal/platform/ethereum"
	ws "lend.fam/go-server/internal/platform/websocket"
)

// Run starts the main loop that reacts to blockchain events and ticks.
func Run(ctx context.Context, ethc *eth.Clients, dur time.Duration, hub *ws.Hub) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// TODO: orchestration logic (start/finalize epochs, allocate yield)
		case ev := <-ethc.Events():
			hub.Broadcast("epochs", ev)
		case <-ctx.Done():
			return
		}
	}
}
