package view

import (
	"log/slog"
	"net/http"
	"time"

	"hack.moontide.ink/pingfisher/muffled/internal/sse"
)

func (v view) Events(w http.ResponseWriter, r *http.Request) {
	s, err := sse.NewSender(w)
	if err != nil {
		slog.Error("could not initialize sender", "error", err)
		http.Error(w, "Could not initialize sender", http.StatusInternalServerError)
		return
	}

	slog.Debug("client connected")

	ctx := r.Context()

	if latest, ok := v.nm.Latest(); ok {
		err := s.Send(latest)
		if err != nil {
			slog.Error("unable to send initial event", "error", err)
			return
		}
	}

	ec, unsubscribe := v.nm.Subscribe()
	defer unsubscribe()

	pt := time.NewTicker(10 * time.Second)
	defer pt.Stop()

	for {
		select {
		case event := <-ec:
			slog.Debug("received event", "event", event)

			err := s.Send(event)
			if err != nil {
				slog.Error("unable to send event", "error", err)
				return
			}
		case <-pt.C:
			if err := s.Ping(); err != nil {
				return
			}
		case <-ctx.Done():
			slog.Debug("client disconnected")
			return
		}
	}
}
