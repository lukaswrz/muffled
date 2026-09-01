package polling

import (
	"fmt"
	"log/slog"
	"reflect"

	"hack.moontide.ink/lukas/muffled/internal/events"
	"hack.moontide.ink/lukas/muffled/internal/listenbrainz"
	"hack.moontide.ink/lukas/muffled/internal/notify"
)

type refresher struct {
	old      *events.PlayingNowEvent
	username string
	lb       *listenbrainz.Client
	nm       *notify.Manager[events.PlayingNowEvent]
}

func (r *refresher) refresh() error {
	response, err := r.lb.GetPlayingNow(r.username)
	if err != nil {
		return fmt.Errorf("unable to request playing now event: %w", err)
	}

	new, err := events.MapPlayingNowEvent(response)
	if err != nil {
		return fmt.Errorf("unable to map response to event: %w", err)
	}

	slog.Debug("received track", "title", new.Title, "artist", new.Artist, "release", new.Release)

	if r.old == nil || !reflect.DeepEqual(*r.old, new) {
		r.old = &new
		r.nm.Broadcast(new)
	}

	return nil
}
