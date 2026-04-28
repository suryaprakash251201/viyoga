package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/viyoga/viyoga/internal/linux"
)

func (s *Server) handleListServices(w http.ResponseWriter, r *http.Request) {
	services, err := s.linux.ListServices(r.Context())
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, services)
}

func (s *Server) handleGetService(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	detail, err := s.linux.GetService(r.Context(), name)
	if err != nil {
		s.respondError(w, http.StatusNotFound, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, detail)
}

func (s *Server) handleServiceAction(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	action := chi.URLParam(r, "action")

	var err error
	switch action {
	case "start":
		err = s.linux.StartService(r.Context(), name)
	case "stop":
		err = s.linux.StopService(r.Context(), name)
	case "restart":
		err = s.linux.RestartService(r.Context(), name)
	case "enable":
		err = s.linux.EnableService(r.Context(), name)
	case "disable":
		err = s.linux.DisableService(r.Context(), name)
	default:
		s.respondError(w, http.StatusBadRequest, "invalid action: "+action)
		return
	}

	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, map[string]string{"status": "ok", "action": action, "service": name})
}

func (s *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	filter := linux.LogFilter{
		Unit:     r.URL.Query().Get("unit"),
		Priority: r.URL.Query().Get("priority"),
		Since:    r.URL.Query().Get("since"),
		Grep:     r.URL.Query().Get("grep"),
		Lines:    200,
	}

	logs, err := s.logs.GetLogs(r.Context(), filter)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, logs)
}

func (s *Server) handleFirewallStatus(w http.ResponseWriter, r *http.Request) {
	status, err := s.firewall.GetStatus(r.Context())
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, status)
}

func (s *Server) handleFirewallAddRule(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Port   string `json:"port"`
		Proto  string `json:"proto"`
		Action string `json:"action"` // "allow" or "deny"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var err error
	switch body.Action {
	case "allow":
		err = s.firewall.AllowPort(r.Context(), body.Port, body.Proto)
	case "deny":
		err = s.firewall.DenyPort(r.Context(), body.Port, body.Proto)
	default:
		s.respondError(w, http.StatusBadRequest, "action must be 'allow' or 'deny'")
		return
	}

	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleFirewallDeleteRule(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RuleNumber int `json:"rule_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := s.firewall.DeleteRule(r.Context(), body.RuleNumber); err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	includeSystem := r.URL.Query().Get("system") == "true"
	users, err := s.users.ListUsers(r.Context(), includeSystem)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, users)
}

func (s *Server) handleListCronJobs(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	jobs, err := s.cron.ListCronJobs(r.Context(), user)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, jobs)
}
