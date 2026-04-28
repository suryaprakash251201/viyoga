package api

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/viyoga/viyoga/internal/alerting"
	"github.com/viyoga/viyoga/internal/collector"
	"github.com/viyoga/viyoga/internal/dns"
	"github.com/viyoga/viyoga/internal/docker"
	"github.com/viyoga/viyoga/internal/hub"
	"github.com/viyoga/viyoga/internal/linux"
	"github.com/viyoga/viyoga/internal/monitor"
	"github.com/viyoga/viyoga/internal/store"
)

// Server holds all API dependencies.
type Server struct {
	router      chi.Router
	collector   *collector.Manager
	hub         *hub.Hub
	store       *store.Store
	upgrader    websocket.Upgrader
	linux       *linux.SystemdManager
	logs        *linux.LogManager
	firewall    *linux.FirewallManager
	users       *linux.UserManager
	cron        *linux.CronManager
	docker      *docker.Manager
	dns         *dns.TechnitiumClient
	monitor     *monitor.Manager
	alertEngine *alerting.Engine
	frontendDir string
}

// ServerDeps holds optional dependencies for the API server.
type ServerDeps struct {
	Collector   *collector.Manager
	Hub         *hub.Hub
	Store       *store.Store
	Docker      *docker.Manager
	DNS         *dns.TechnitiumClient
	Monitor     *monitor.Manager
	AlertEngine *alerting.Engine
	FrontendDir string
}

// NewServer creates a new API server with all dependencies.
func NewServer(deps ServerDeps) *Server {
	srv := &Server{
		collector:   deps.Collector,
		hub:         deps.Hub,
		store:       deps.Store,
		docker:      deps.Docker,
		frontendDir: deps.FrontendDir,
		dns:         deps.DNS,
		monitor:     deps.Monitor,
		alertEngine: deps.AlertEngine,
		linux:       linux.NewSystemdManager(),
		logs:        linux.NewLogManager(),
		firewall:    linux.NewFirewallManager(),
		users:       linux.NewUserManager(),
		cron:        linux.NewCronManager(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	srv.setupRoutes()
	return srv
}

// Router returns the configured chi router.
func (s *Server) Router() chi.Router {
	return s.router
}

func (s *Server) setupRoutes() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(requestLogger)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Health check
		r.Get("/health", s.handleHealth)

		// Metrics
		r.Route("/metrics", func(r chi.Router) {
			r.Get("/current", s.handleCurrentMetrics)
			r.Get("/cpu", s.handleCPUMetrics)
			r.Get("/memory", s.handleMemoryMetrics)
			r.Get("/disk", s.handleDiskMetrics)
			r.Get("/network", s.handleNetworkMetrics)
			r.Get("/history/{type}", s.handleMetricsHistory)
		})

		// System info
		r.Get("/system", s.handleSystemInfo)

		// Linux management (Phase 2)
		r.Route("/linux", func(r chi.Router) {
			r.Get("/services", s.handleListServices)
			r.Get("/services/{name}", s.handleGetService)
			r.Post("/services/{name}/{action}", s.handleServiceAction)
			r.Get("/logs", s.handleGetLogs)
			r.Get("/firewall", s.handleFirewallStatus)
			r.Post("/firewall/rules", s.handleFirewallAddRule)
			r.Delete("/firewall/rules", s.handleFirewallDeleteRule)
			r.Get("/users", s.handleListUsers)
			r.Get("/cron", s.handleListCronJobs)
		})

		// Hardware / Processes (Phase 3)
		r.Route("/hardware", func(r chi.Router) {
			r.Get("/processes", s.handleProcessList)
			r.Get("/alerts/rules", s.handleAlertRules)
			r.Get("/alerts/events", s.handleAlertEvents)
		})

		// Docker (Phase 4)
		r.Route("/docker", func(r chi.Router) {
			r.Get("/containers", s.handleListContainers)
			r.Post("/containers/{id}/{action}", s.handleContainerAction)
			r.Get("/containers/{id}/logs", s.handleContainerLogs)
			r.Get("/containers/{id}/inspect", s.handleContainerInspect)
			r.Get("/images", s.handleListImages)
			r.Delete("/images/{id}", s.handleRemoveImage)
			r.Post("/prune", s.handlePrune)
		})

		// DNS Gateway (Phase 5)
		r.Route("/dns", func(r chi.Router) {
			r.Get("/stats", s.handleDNSStats)
			r.Get("/querylog", s.handleDNSQueryLog)
			r.Get("/blocklists", s.handleDNSBlockLists)
		})

		// Web Monitor (Phase 6)
		r.Route("/monitor", func(r chi.Router) {
			r.Get("/status", s.handleMonitorStatus)
			r.Get("/targets", s.handleMonitorTargets)
			r.Post("/targets", s.handleAddMonitorTarget)
			r.Delete("/targets/{id}", s.handleRemoveMonitorTarget)
		})
	})

	// WebSocket endpoint
	r.Get("/ws/metrics", s.handleWebSocket)

	// Serve frontend static files (SPA mode)
	if s.frontendDir != "" {
		s.serveFrontend(r)
	}

	s.router = r
}

// serveFrontend serves the SvelteKit static build with SPA fallback.
func (s *Server) serveFrontend(r chi.Router) {
	absPath, err := filepath.Abs(s.frontendDir)
	if err != nil {
		log.Warn().Str("dir", s.frontendDir).Err(err).Msg("Invalid frontend directory")
		return
	}

	// Check if directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Warn().Str("dir", absPath).Msg("Frontend directory not found, skipping static serving")
		return
	}

	fsys := os.DirFS(absPath)
	fileServer := http.FileServer(http.FS(fsys))

	log.Info().Str("dir", absPath).Msg("Serving frontend static files")

	// Catch-all handler: serve static files, fallback to index.html for SPA
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try to serve the exact file
		if f, err := fs.Stat(fsys, path); err == nil && !f.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Try path with index.html appended (for directory-style routes)
		idxPath := path + "/index.html"
		if _, err := fs.Stat(fsys, idxPath); err == nil {
			r.URL.Path = "/" + idxPath
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for all non-file routes
		indexFile, err := fs.ReadFile(fsys, "index.html")
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(indexFile)
	})
}

// requestLogger is a custom middleware for structured request logging.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		log.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.Status()).
			Dur("latency", time.Since(start)).
			Msg("Request")
	})
}
