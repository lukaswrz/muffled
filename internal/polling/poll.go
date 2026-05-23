package polling

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"hack.moontide.ink/pingfisher/muffled/internal/events"
	"hack.moontide.ink/pingfisher/muffled/internal/listenbrainz"
	"hack.moontide.ink/pingfisher/muffled/internal/notify"
)

func Poll(cm *notify.Manager[events.PlayingNowEvent], lb *listenbrainz.Client, interval int, username string) error {
	first := true
	old := events.PlayingNowEvent{}

	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer t.Stop()

	for {
		slog.Debug("requesting current track")
		response, err := lb.GetPlayingNow(username)
		if err != nil {
			return fmt.Errorf("unable to request playing now event: %w", err)
		}

		new, err := events.MapPlayingNowEvent(response)
		if err != nil {
			return fmt.Errorf("unable to map response to event: %w", err)
		}

		slog.Debug("received track", "title", new.Title, "artist", new.Artist, "release", new.Release)

		if first || !reflect.DeepEqual(old, new) {
			first = false
			cm.Broadcast(new)
		}

		old = new

		<-t.C
	}
}
