package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/viyoga/viyoga/internal/collector"
	"github.com/viyoga/viyoga/internal/monitor"
)

func (s *Server) handleDNSStats(w http.ResponseWriter, r *http.Request) {
	if s.dns == nil {
		s.respondError(w, http.StatusServiceUnavailable, "DNS gateway not available")
		return
	}
	stats, err := s.dns.GetStats(r.Context())
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, stats)
}

func (s *Server) handleDNSQueryLog(w http.ResponseWriter, r *http.Request) {
	if s.dns == nil {
		s.respondError(w, http.StatusServiceUnavailable, "DNS gateway not available")
		return
	}
	entries, err := s.dns.GetQueryLog(r.Context(), 1, 100)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, entries)
}

func (s *Server) handleDNSBlockLists(w http.ResponseWriter, r *http.Request) {
	if s.dns == nil {
		s.respondError(w, http.StatusServiceUnavailable, "DNS gateway not available")
		return
	}
	lists, err := s.dns.GetBlockLists(r.Context())
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, lists)
}

func (s *Server) handleMonitorStatus(w http.ResponseWriter, r *http.Request) {
	if s.monitor == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Monitor not available")
		return
	}
	statuses := s.monitor.GetStatus()
	s.respondJSON(w, http.StatusOK, statuses)
}

func (s *Server) handleMonitorTargets(w http.ResponseWriter, r *http.Request) {
	if s.monitor == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Monitor not available")
		return
	}
	targets := s.monitor.GetTargets()
	s.respondJSON(w, http.StatusOK, targets)
}

func (s *Server) handleAddMonitorTarget(w http.ResponseWriter, r *http.Request) {
	if s.monitor == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Monitor not available")
		return
	}

	var target monitor.Target
	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		s.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	s.monitor.AddTarget(target)
	s.respondJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

func (s *Server) handleRemoveMonitorTarget(w http.ResponseWriter, r *http.Request) {
	if s.monitor == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Monitor not available")
		return
	}

	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		s.respondError(w, http.StatusBadRequest, "invalid target ID")
		return
	}

	s.monitor.RemoveTarget(id)
	s.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleProcessList(w http.ResponseWriter, r *http.Request) {
	pc := collector.NewProcessCollector()
	processes, err := pc.Collect(r.Context())
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, processes)
}

func (s *Server) handleAlertRules(w http.ResponseWriter, r *http.Request) {
	if s.alertEngine == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Alerting not available")
		return
	}
	rules := s.alertEngine.GetRules()
	s.respondJSON(w, http.StatusOK, rules)
}

func (s *Server) handleAlertEvents(w http.ResponseWriter, r *http.Request) {
	if s.alertEngine == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Alerting not available")
		return
	}
	events := s.alertEngine.GetEvents(100)
	s.respondJSON(w, http.StatusOK, events)
}
