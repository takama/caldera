package metrics

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// MetricFunc defines func for regular update of the metric.
type MetricFunc func() error

// Monitor is used for regular updates of the metrics.
type Monitor struct {
	funcs []MetricFunc

	period   time.Duration
	quitOnce sync.Once
	quit     chan bool

	log *zap.Logger
}

// NewMonitor returns a new instance of Monitor.
// nolint: gomnd
func NewMonitor(log *zap.Logger, funcs ...MetricFunc) *Monitor {
	return &Monitor{
		funcs:  funcs,
		period: 50 * time.Second,
		quit:   make(chan bool),
		log:    log,
	}
}

// Run real-time update of the metrics.
func (m *Monitor) Run() {
	for {
		select {
		case <-time.After(m.period):
			m.update()
		case <-m.quit:
			return
		}
	}
}

// Shutdown process graceful shutdown for the monitor.
func (m *Monitor) Shutdown(ctx context.Context) error {
	m.quitOnce.Do(func() {
		close(m.quit)
	})

	return nil
}

// update all metrics via funcs.
func (m *Monitor) update() {
	for _, f := range m.funcs {
		if err := f(); err != nil {
			m.log.Error("сannot update metrics", zap.Error(err))

			return
		}
	}
}
