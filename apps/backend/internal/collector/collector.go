package collector

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Collector defines the interface for all metric collectors.
type Collector interface {
	// Name returns the unique name of this collector.
	Name() string
	// Collect gathers current metrics and returns them.
	Collect(ctx context.Context) (interface{}, error)
}

// Snapshot represents a point-in-time collection of all metrics.
type Snapshot struct {
	Timestamp  time.Time              `json:"timestamp"`
	CPU        interface{}            `json:"cpu,omitempty"`
	Memory     interface{}            `json:"memory,omitempty"`
	Disk       interface{}            `json:"disk,omitempty"`
	Network    interface{}            `json:"network,omitempty"`
	SystemInfo interface{}            `json:"system_info,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// Manager orchestrates all collectors and produces metric snapshots.
type Manager struct {
	collectors map[string]Collector
	interval   time.Duration
	mu         sync.RWMutex
	latest     *Snapshot
	outCh      chan *Snapshot
	stopCh     chan struct{}
}

// NewManager creates a new collector manager.
func NewManager(interval time.Duration) *Manager {
	return &Manager{
		collectors: make(map[string]Collector),
		interval:   interval,
		outCh:      make(chan *Snapshot, 10),
		stopCh:     make(chan struct{}),
	}
}

// Register adds a collector to the manager.
func (m *Manager) Register(c Collector) {
	m.collectors[c.Name()] = c
	log.Info().Str("collector", c.Name()).Msg("Registered collector")
}

// Subscribe returns a channel that receives metric snapshots.
func (m *Manager) Subscribe() <-chan *Snapshot {
	return m.outCh
}

// Latest returns the most recently collected snapshot.
func (m *Manager) Latest() *Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.latest
}

// Start begins the collection loop in a goroutine.
func (m *Manager) Start(ctx context.Context) {
	go m.run(ctx)
}

// Stop signals the collection loop to stop.
func (m *Manager) Stop() {
	close(m.stopCh)
}

func (m *Manager) run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Collect immediately on start
	m.collect(ctx)

	for {
		select {
		case <-ticker.C:
			m.collect(ctx)
		case <-m.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) collect(ctx context.Context) {
	snap := &Snapshot{
		Timestamp: time.Now(),
		Extra:     make(map[string]interface{}),
	}

	for name, c := range m.collectors {
		data, err := c.Collect(ctx)
		if err != nil {
			log.Warn().Err(err).Str("collector", name).Msg("Collection failed")
			continue
		}

		switch name {
		case "cpu":
			snap.CPU = data
		case "memory":
			snap.Memory = data
		case "disk":
			snap.Disk = data
		case "network":
			snap.Network = data
		case "system_info":
			snap.SystemInfo = data
		default:
			snap.Extra[name] = data
		}
	}

	m.mu.Lock()
	m.latest = snap
	m.mu.Unlock()

	// Non-blocking send to channel
	select {
	case m.outCh <- snap:
	default:
		// Drop if no one is listening (buffer full)
	}
}
