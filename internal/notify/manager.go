package notify

import (
	"sync"

	"github.com/google/uuid"
)

type Manager[T any] struct {
	bufsize int

	muClients sync.RWMutex
	clients   map[uuid.UUID]chan T

	muState   sync.RWMutex
	latest    T
	hasLatest bool
}

func NewManager[T any](bufsize int) *Manager[T] {
	return &Manager[T]{
		clients: make(map[uuid.UUID]chan T),
		bufsize: bufsize,
	}
}

func (m *Manager[T]) Subscribe() (<-chan T, func()) {
	id := uuid.New()
	c := make(chan T, m.bufsize)

	m.connect(id, c)

	unsubscribe := func() {
		m.disconnect(id)
	}

	return c, unsubscribe
}

func (m *Manager[T]) Broadcast(message T) {
	m.update(message)

	m.muClients.RLock()

	clients := make([]chan T, 0, len(m.clients))
	for _, c := range m.clients {
		clients = append(clients, c)
	}

	m.muClients.RUnlock()

	for _, c := range clients {
		select {
		case c <- message:
		default:
		}
	}
}

func (m *Manager[T]) Latest() (T, bool) {
	m.muState.RLock()
	defer m.muState.RUnlock()

	if !m.hasLatest {
		var zero T
		return zero, false
	}

	return m.latest, true
}

func (m *Manager[T]) connect(id uuid.UUID, ch chan T) {
	m.muClients.Lock()
	defer m.muClients.Unlock()

	m.clients[id] = ch
}

func (m *Manager[T]) disconnect(id uuid.UUID) {
	m.muClients.Lock()
	defer m.muClients.Unlock()

	ch, ok := m.clients[id]
	if !ok {
		return
	}

	delete(m.clients, id)
	close(ch)
}

func (m *Manager[T]) update(message T) {
	m.muState.Lock()
	defer m.muState.Unlock()

	m.latest = message
	m.hasLatest = true
}
