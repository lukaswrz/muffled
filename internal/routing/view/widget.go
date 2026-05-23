package view

import (
	"log/slog"
	"net/http"
)

func (v view) Widget(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(v.widget)
	if err != nil {
		slog.Error("unable to return widget", "error", err)
		return
	}
}
