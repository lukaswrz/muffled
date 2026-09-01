package polling

import (
	"log/slog"
	"time"

	"hack.moontide.ink/lukas/muffled/internal/events"
	"hack.moontide.ink/lukas/muffled/internal/listenbrainz"
	"hack.moontide.ink/lukas/muffled/internal/notify"
)

func Poll(nm *notify.Manager[events.PlayingNowEvent], lb *listenbrainz.Client, interval int, username string) error {
	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer t.Stop()

	r := refresher{
		username: username,
		lb:       lb,
		nm:       nm,
	}

	c := 0
	for {
		for c == 0 {
			c = <-nm.ClientsC
		}

		slog.Debug("refreshing...")
		if err := r.refresh(); err != nil {
			slog.Error("refresh failed", "error", err)
		}

		select {
		case c = <-nm.ClientsC:
		case <-t.C:
		}
	}
}
