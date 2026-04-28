package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Timestamp time.Time `json:"timestamp"`
	WSClients int       `json:"ws_clients,omitempty"`
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: status >= 200 && status < 300,
		Data:    data,
		Meta:    &Meta{Timestamp: time.Now(), WSClients: s.hub.ClientCount()},
	})
}

func (s *Server) respondError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{Success: false, Error: msg, Meta: &Meta{Timestamp: time.Now()}})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, http.StatusOK, map[string]interface{}{"status": "healthy", "version": "0.1.0", "ws_clients": s.hub.ClientCount()})
}

func (s *Server) handleCurrentMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.collector.Latest()
	if snap == nil {
		s.respondError(w, http.StatusServiceUnavailable, "metrics not yet available")
		return
	}
	s.respondJSON(w, http.StatusOK, snap)
}

func (s *Server) handleCPUMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.collector.Latest()
	if snap == nil || snap.CPU == nil {
		s.respondError(w, http.StatusServiceUnavailable, "CPU metrics not available")
		return
	}
	s.respondJSON(w, http.StatusOK, snap.CPU)
}

func (s *Server) handleMemoryMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.collector.Latest()
	if snap == nil || snap.Memory == nil {
		s.respondError(w, http.StatusServiceUnavailable, "memory metrics not available")
		return
	}
	s.respondJSON(w, http.StatusOK, snap.Memory)
}

func (s *Server) handleDiskMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.collector.Latest()
	if snap == nil || snap.Disk == nil {
		s.respondError(w, http.StatusServiceUnavailable, "disk metrics not available")
		return
	}
	s.respondJSON(w, http.StatusOK, snap.Disk)
}

func (s *Server) handleNetworkMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.collector.Latest()
	if snap == nil || snap.Network == nil {
		s.respondError(w, http.StatusServiceUnavailable, "network metrics not available")
		return
	}
	s.respondJSON(w, http.StatusOK, snap.Network)
}

func (s *Server) handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	snap := s.collector.Latest()
	if snap == nil || snap.SystemInfo == nil {
		s.respondError(w, http.StatusServiceUnavailable, "system info not available")
		return
	}
	s.respondJSON(w, http.StatusOK, snap.SystemInfo)
}

func (s *Server) handleMetricsHistory(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	if metricType == "" {
		s.respondError(w, http.StatusBadRequest, "metric type required")
		return
	}
	limit := 100
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 1000 {
		limit = l
	}
	since := time.Now().Add(-1 * time.Hour)
	if d, err := time.ParseDuration(r.URL.Query().Get("since")); err == nil {
		since = time.Now().Add(-d)
	}
	records, err := s.store.GetMetricsHistory(metricType, since, limit)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "failed to fetch history")
		return
	}
	s.respondJSON(w, http.StatusOK, records)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	s.hub.RegisterClient(conn)
}
