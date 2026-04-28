package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

// Store handles all database operations.
type Store struct {
	db *sql.DB
}

// MetricsRecord represents a stored metrics snapshot.
type MetricsRecord struct {
	ID         int64     `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	MetricType string    `json:"metric_type"`
	MetricData string    `json:"metric_data"`
}

// New creates a new Store instance and initializes the database.
func New(dbPath string) (*Store, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Enable WAL mode for better concurrent read/write performance
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("setting WAL mode: %w", err)
	}

	// Set busy timeout to avoid lock contention
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return nil, fmt.Errorf("setting busy timeout: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return s, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// migrate runs database schema migrations.
func (s *Store) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS metrics_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		metric_type TEXT NOT NULL,
		metric_data JSON NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_metrics_ts ON metrics_history(timestamp, metric_type);

	CREATE TABLE IF NOT EXISTS alert_rules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		metric_type TEXT NOT NULL,
		condition TEXT NOT NULL,
		threshold REAL NOT NULL,
		notify_channel TEXT DEFAULT 'webhook',
		notify_target TEXT,
		enabled BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS alert_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		rule_id INTEGER REFERENCES alert_rules(id),
		triggered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		value REAL,
		resolved_at DATETIME,
		acknowledged BOOLEAN DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT DEFAULT 'viewer',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS monitor_targets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		method TEXT DEFAULT 'GET',
		interval_seconds INTEGER DEFAULT 60,
		timeout_seconds INTEGER DEFAULT 10,
		expected_status INTEGER DEFAULT 200,
		enabled BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS monitor_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		target_id INTEGER REFERENCES monitor_targets(id),
		status_code INTEGER,
		response_time_ms INTEGER,
		is_up BOOLEAN,
		error_message TEXT,
		checked_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_monitor_results ON monitor_results(target_id, checked_at);
	`

	_, err := s.db.Exec(schema)
	return err
}

// SaveMetrics stores a metric snapshot in the database.
func (s *Store) SaveMetrics(metricType string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling metrics: %w", err)
	}

	_, err = s.db.Exec(
		"INSERT INTO metrics_history (metric_type, metric_data) VALUES (?, ?)",
		metricType, string(jsonData),
	)
	return err
}

// GetMetricsHistory returns historical metrics for a given type within a time range.
func (s *Store) GetMetricsHistory(metricType string, since time.Time, limit int) ([]MetricsRecord, error) {
	rows, err := s.db.Query(
		"SELECT id, timestamp, metric_type, metric_data FROM metrics_history WHERE metric_type = ? AND timestamp >= ? ORDER BY timestamp DESC LIMIT ?",
		metricType, since, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []MetricsRecord
	for rows.Next() {
		var r MetricsRecord
		if err := rows.Scan(&r.ID, &r.Timestamp, &r.MetricType, &r.MetricData); err != nil {
			continue
		}
		records = append(records, r)
	}
	return records, nil
}

// PruneOldMetrics removes metrics older than the retention period.
func (s *Store) PruneOldMetrics(retentionHours int) error {
	cutoff := time.Now().Add(-time.Duration(retentionHours) * time.Hour)
	result, err := s.db.Exec("DELETE FROM metrics_history WHERE timestamp < ?", cutoff)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		log.Info().Int64("pruned", rows).Int("retention_hours", retentionHours).Msg("Pruned old metrics")
	}
	return nil
}

// StartPruner runs a background goroutine that prunes old metrics periodically.
func (s *Store) StartPruner(retentionHours int, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.PruneOldMetrics(retentionHours); err != nil {
				log.Error().Err(err).Msg("Failed to prune metrics")
			}
		}
	}()
}
