package view

import (
	"net/http"

	"hack.moontide.ink/pingfisher/muffled/internal/events"
	"hack.moontide.ink/pingfisher/muffled/internal/notify"
)

type view struct {
	nm     *notify.Manager[events.PlayingNowEvent]
	widget []byte
}

type View interface {
	Widget(w http.ResponseWriter, r *http.Request)
	Events(w http.ResponseWriter, r *http.Request)
}

func NewView(
	nm *notify.Manager[events.PlayingNowEvent],
	widget []byte,
) View {
	return view{
		nm,
		widget,
	}
}
