package stats

import (
	"context"
	"sync"
	"time"
)

const (
	defaultCapacity = 1000
)

type (
	Writer interface {
		Insert(rows Rows) error
	}

	Manager struct {
		writer        Writer
		flushInterval time.Duration
		ctx           context.Context
		cancel        context.CancelFunc
		mu            sync.RWMutex
		rows          Rows
	}
)

func NewManager(w Writer, flushInterval time.Duration) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		writer:        w,
		flushInterval: flushInterval,
		ctx:           ctx,
		cancel:        cancel,
		rows:          nil,
	}
}

func (m *Manager) Append(k Key, v Value) {
	m.mu.Lock()
	defer m.mu.Unlock()

	current := m.rows[k]

	current = current.Assign(v)
	m.rows[k] = current

}

func newRows() Rows {
	return make(Rows, defaultCapacity)
}
