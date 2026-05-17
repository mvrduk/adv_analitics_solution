package stats

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
)

var logger *zap.Logger

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

func (m *Manager) AppendRows(rows Rows) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range rows {
		m.unsafeAppend(k, v)

	}
}

func (m *Manager) unsafeAppend(k Key, v Value) {
	current := m.rows[k]
	current = current.Assign(v)

	m.rows[k] = current
}

func (m *Manager) start() {
	logger.Info("Stat loop started")

	go m.loop()
}

func (m *Manager) loop() {
	for {
		select {
		case <-time.After(m.flushInterval):
			m.startInserting()
		case <-m.ctx.Done():
			m.startInserting()
			return
		}
	}
}

func (m *Manager) startInserting() {
	logger.Info("Start stats inserting")
	rows := m.withdraw()
	if len(rows) == 0 {
		logger.Warn("No stats rows, skipping")
		return
	}

	if err := m.writer.Insert(rows); err != nil {
		logger.Error("Failed to write stats", zap.Error(err))
		logger.Warn("Return stats rows to map: ", zap.Int("rows_count", len(rows)))

		m.AppendRows(rows)
		return
	}

	logger.Info("Stats rows successfully written:", zap.Int("len rows", len(rows)))
}

func (m *Manager) withdraw() Rows {

	m.mu.Lock()
	defer m.mu.Unlock()

	rows := m.rows
	m.rows = newRows()

	return rows

}

func newRows() Rows {
	return make(Rows, defaultCapacity)
}
