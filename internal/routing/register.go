package routing

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"hack.moontide.ink/pingfisher/muffled/internal/assets"
	"hack.moontide.ink/pingfisher/muffled/internal/events"
	"hack.moontide.ink/pingfisher/muffled/internal/notify"
	"hack.moontide.ink/pingfisher/muffled/internal/routing/view"
)

func Register(r *chi.Mux, nm *notify.Manager[events.PlayingNowEvent], widget []byte) error {
	v := view.NewView(nm, widget)

	r.Handle("/static/*", http.FileServer(http.FS(assets.FS)))

	r.Get("/", v.Widget)
	r.Get("/events", v.Events)

	return nil
}
