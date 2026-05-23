package sse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type Sender struct {
	mu      sync.Mutex
	w       http.ResponseWriter
	flusher http.Flusher
}

func NewSender(w http.ResponseWriter) (*Sender, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("flushing unsupported")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	return &Sender{
		w:       w,
		flusher: flusher,
	}, nil
}

func (c *Sender) Send(payload any) error {
	j, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	_, err = fmt.Fprintf(c.w, "data: %s\n\n", j)
	if err != nil {
		return fmt.Errorf("unable to send data via SSE: %w", err)
	}
	c.flusher.Flush()
	return nil
}

func (c *Sender) Ping() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := fmt.Fprint(c.w, ": ping\n\n")
	if err != nil {
		return fmt.Errorf("unable to ping via SSE: %w", err)
	}
	c.flusher.Flush()
	return nil
}
