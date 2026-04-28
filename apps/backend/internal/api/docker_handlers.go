package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (s *Server) handleListContainers(w http.ResponseWriter, r *http.Request) {
	if s.docker == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Docker not available")
		return
	}
	all := r.URL.Query().Get("all") == "true"
	containers, err := s.docker.ListContainers(r.Context(), all)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, containers)
}

func (s *Server) handleContainerAction(w http.ResponseWriter, r *http.Request) {
	if s.docker == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Docker not available")
		return
	}
	id := chi.URLParam(r, "id")
	action := chi.URLParam(r, "action")

	var err error
	switch action {
	case "start":
		err = s.docker.StartContainer(r.Context(), id)
	case "stop":
		err = s.docker.StopContainer(r.Context(), id)
	case "restart":
		err = s.docker.RestartContainer(r.Context(), id)
	case "remove":
		err = s.docker.RemoveContainer(r.Context(), id)
	default:
		s.respondError(w, http.StatusBadRequest, "invalid action: "+action)
		return
	}

	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, map[string]string{"status": "ok", "action": action})
}

func (s *Server) handleContainerLogs(w http.ResponseWriter, r *http.Request) {
	if s.docker == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Docker not available")
		return
	}
	id := chi.URLParam(r, "id")
	lines := 100
	if l, err := strconv.Atoi(r.URL.Query().Get("lines")); err == nil && l > 0 {
		lines = l
	}

	logs, err := s.docker.GetContainerLogs(r.Context(), id, lines)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, map[string]string{"logs": logs})
}

func (s *Server) handleContainerInspect(w http.ResponseWriter, r *http.Request) {
	if s.docker == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Docker not available")
		return
	}
	id := chi.URLParam(r, "id")
	info, err := s.docker.InspectContainer(r.Context(), id)
	if err != nil {
		s.respondError(w, http.StatusNotFound, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, info)
}

func (s *Server) handleListImages(w http.ResponseWriter, r *http.Request) {
	if s.docker == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Docker not available")
		return
	}
	images, err := s.docker.ListImages(r.Context())
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, images)
}

func (s *Server) handleRemoveImage(w http.ResponseWriter, r *http.Request) {
	if s.docker == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Docker not available")
		return
	}
	id := chi.URLParam(r, "id")
	if err := s.docker.RemoveImage(r.Context(), id); err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handlePrune(w http.ResponseWriter, r *http.Request) {
	if s.docker == nil {
		s.respondError(w, http.StatusServiceUnavailable, "Docker not available")
		return
	}

	var body struct {
		Type string `json:"type"` // "containers" or "images"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var spaceReclaimed uint64
	var err error
	switch body.Type {
	case "containers":
		spaceReclaimed, err = s.docker.PruneContainers(r.Context())
	case "images":
		spaceReclaimed, err = s.docker.PruneImages(r.Context())
	default:
		s.respondError(w, http.StatusBadRequest, "type must be 'containers' or 'images'")
		return
	}

	if err != nil {
		s.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":          "ok",
		"space_reclaimed": spaceReclaimed,
	})
}
