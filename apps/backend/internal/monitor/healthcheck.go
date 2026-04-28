package monitor

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Target represents an HTTP endpoint to monitor.
type Target struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	Method          string `json:"method"`
	IntervalSeconds int    `json:"interval_seconds"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
	ExpectedStatus  int    `json:"expected_status"`
	Enabled         bool   `json:"enabled"`
}

// Result records a health check outcome.
type Result struct {
	TargetID       int       `json:"target_id"`
	TargetName     string    `json:"target_name"`
	URL            string    `json:"url"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int64    `json:"response_time_ms"`
	IsUp           bool      `json:"is_up"`
	Error          string    `json:"error,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
	CertExpiry     *time.Time `json:"cert_expiry,omitempty"`
}

// TargetStatus combines a target with its latest results.
type TargetStatus struct {
	Target        Target   `json:"target"`
	LastResult    *Result  `json:"last_result,omitempty"`
	UptimePercent float64  `json:"uptime_percent"`
	AvgResponse   int64    `json:"avg_response_ms"`
	History       []Result `json:"history,omitempty"`
}

// Manager runs periodic health checks on configured targets.
type Manager struct {
	targets []Target
	results map[int][]Result // target_id -> results
	mu      sync.RWMutex
	cancel  context.CancelFunc
	client  *http.Client
}

// NewManager creates a web monitor manager.
func NewManager() *Manager {
	return &Manager{
		results: make(map[int][]Result),
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
		},
	}
}

// SetTargets configures the monitoring targets.
func (m *Manager) SetTargets(targets []Target) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.targets = targets
}

// GetTargets returns all targets.
func (m *Manager) GetTargets() []Target {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]Target{}, m.targets...)
}

// GetStatus returns the status of all targets.
func (m *Manager) GetStatus() []TargetStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var statuses []TargetStatus
	for _, target := range m.targets {
		status := TargetStatus{Target: target}
		if results, ok := m.results[target.ID]; ok && len(results) > 0 {
			last := results[len(results)-1]
			status.LastResult = &last

			// Calculate uptime & avg response
			upCount := 0
			var totalResp int64
			for _, r := range results {
				if r.IsUp {
					upCount++
				}
				totalResp += r.ResponseTimeMs
			}
			if len(results) > 0 {
				status.UptimePercent = float64(upCount) / float64(len(results)) * 100
				status.AvgResponse = totalResp / int64(len(results))
			}

			// Last 20 results
			start := len(results) - 20
			if start < 0 {
				start = 0
			}
			status.History = results[start:]
		}
		statuses = append(statuses, status)
	}
	return statuses
}

// Start begins periodic health checks.
func (m *Manager) Start(ctx context.Context) {
	ctx, m.cancel = context.WithCancel(ctx)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		// Run immediately
		m.checkAll(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.checkAll(ctx)
			}
		}
	}()

	log.Info().Msg("Web monitor started")
}

// Stop halts health checks.
func (m *Manager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
}

// CheckTarget performs a single health check.
func (m *Manager) CheckTarget(ctx context.Context, target Target) Result {
	result := Result{
		TargetID:   target.ID,
		TargetName: target.Name,
		URL:        target.URL,
		CheckedAt:  time.Now(),
	}

	timeout := time.Duration(target.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	method := target.Method
	if method == "" {
		method = "GET"
	}

	req, err := http.NewRequestWithContext(reqCtx, method, target.URL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	req.Header.Set("User-Agent", "Viyoga-Monitor/1.0")

	start := time.Now()
	resp, err := m.client.Do(req)
	result.ResponseTimeMs = time.Since(start).Milliseconds()

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// Check SSL cert expiry
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		expiry := resp.TLS.PeerCertificates[0].NotAfter
		result.CertExpiry = &expiry
	}

	expectedStatus := target.ExpectedStatus
	if expectedStatus == 0 {
		expectedStatus = 200
	}
	result.IsUp = resp.StatusCode == expectedStatus

	return result
}

func (m *Manager) checkAll(ctx context.Context) {
	m.mu.RLock()
	targets := append([]Target{}, m.targets...)
	m.mu.RUnlock()

	for _, target := range targets {
		if !target.Enabled {
			continue
		}

		result := m.CheckTarget(ctx, target)

		m.mu.Lock()
		m.results[target.ID] = append(m.results[target.ID], result)
		// Keep only last 1440 results (24h at 1/min)
		if len(m.results[target.ID]) > 1440 {
			m.results[target.ID] = m.results[target.ID][len(m.results[target.ID])-720:]
		}
		m.mu.Unlock()
	}
}

// AddTarget dynamically adds a monitoring target.
func (m *Manager) AddTarget(target Target) {
	m.mu.Lock()
	defer m.mu.Unlock()
	target.ID = len(m.targets) + 1
	m.targets = append(m.targets, target)
}

// RemoveTarget removes a monitoring target by ID.
func (m *Manager) RemoveTarget(id int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, t := range m.targets {
		if t.ID == id {
			m.targets = append(m.targets[:i], m.targets[i+1:]...)
			delete(m.results, id)
			return
		}
	}
}
