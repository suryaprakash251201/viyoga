package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/viyoga/viyoga/internal/alerting"
	"github.com/viyoga/viyoga/internal/api"
	"github.com/viyoga/viyoga/internal/collector"
	"github.com/viyoga/viyoga/internal/config"
	"github.com/viyoga/viyoga/internal/dns"
	"github.com/viyoga/viyoga/internal/docker"
	"github.com/viyoga/viyoga/internal/hub"
	"github.com/viyoga/viyoga/internal/monitor"
	"github.com/viyoga/viyoga/internal/store"
)

const version = "0.2.0"

func main() {
	// CLI flags
	configPath := flag.String("config", "viyoga.yaml", "path to config file")
	showVersion := flag.Bool("version", false, "show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Viyoga v%s\n", version)
		os.Exit(0)
	}

	// Load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Setup logging
	setupLogging(cfg.Log.Level, cfg.Log.Format)

	log.Info().
		Str("version", version).
		Str("listen", cfg.ListenAddr()).
		Msg("Starting Viyoga server")

	// Initialize database
	db, err := store.New(cfg.DB.Path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()

	// Start metrics pruner
	if cfg.Metrics.HistoryEnabled {
		db.StartPruner(cfg.Metrics.RetentionHours, 1*time.Hour)
	}

	// Initialize WebSocket hub
	wsHub := hub.New()
	go wsHub.Run()

	// Initialize collectors
	cm := collector.NewManager(cfg.PollDuration())
	cm.Register(collector.NewCPUCollector())
	cm.Register(collector.NewMemoryCollector())
	cm.Register(collector.NewDiskCollector())
	cm.Register(collector.NewNetworkCollector())
	cm.Register(collector.NewSystemInfoCollector(version))

	// Start collectors
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cm.Start(ctx)

	// Initialize alerting engine
	alertEngine := alerting.NewEngine()

	// Bridge: forward collector snapshots to WebSocket hub + DB + Alerts
	go func() {
		for snap := range cm.Subscribe() {
			wsHub.Broadcast(snap)

			if cfg.Metrics.HistoryEnabled {
				if snap.CPU != nil {
					db.SaveMetrics("cpu", snap.CPU)
				}
				if snap.Memory != nil {
					db.SaveMetrics("memory", snap.Memory)
				}
			}

			// Evaluate alert thresholds
			if cfg.Alert.Enabled {
				if snap.CPU != nil {
					if cpu, ok := snap.CPU.(*collector.CPUMetrics); ok {
						alertEngine.Evaluate("cpu", cpu.UsagePercent)
					}
				}
				if snap.Memory != nil {
					if mem, ok := snap.Memory.(*collector.MemoryMetrics); ok {
						alertEngine.Evaluate("memory", mem.UsagePercent)
					}
				}
			}
		}
	}()

	// Initialize optional modules
	var dockerMgr *docker.Manager
	if cfg.Modules.ContainerManager {
		dockerMgr = docker.NewManager(cfg.Docker.Socket)
		if dockerMgr != nil {
			log.Info().Msg("Docker module enabled")
		}
	}

	var dnsClient *dns.TechnitiumClient
	if cfg.Modules.DNSGateway {
		dnsClient = dns.NewTechnitiumClient(cfg.DNS.APIURL, cfg.DNS.APIToken)
		if dnsClient != nil && dnsClient.IsAvailable() {
			log.Info().Msg("DNS gateway module enabled")
		}
	}

	var webMonitor *monitor.Manager
	if cfg.Modules.WebMonitor {
		webMonitor = monitor.NewManager()
		webMonitor.Start(ctx)
		log.Info().Msg("Web monitor module enabled")
	}

	// Initialize API server with all dependencies
	srv := api.NewServer(api.ServerDeps{
		Collector:   cm,
		Hub:         wsHub,
		Store:       db,
		Docker:      dockerMgr,
		DNS:         dnsClient,
		Monitor:     webMonitor,
		AlertEngine: alertEngine,
		FrontendDir: cfg.Server.FrontendDir,
	})

	// HTTP server
	httpSrv := &http.Server{
		Addr:         cfg.ListenAddr(),
		Handler:      srv.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Str("addr", cfg.ListenAddr()).Msg("HTTP server started")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Log enabled modules
	logModules(cfg)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down Viyoga server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}

	if webMonitor != nil {
		webMonitor.Stop()
	}

	cm.Stop()
	cancel()
	log.Info().Msg("Viyoga server stopped")
}

func logModules(cfg *config.Config) {
	modules := map[string]bool{
		"System Dashboard":   cfg.Modules.SystemDashboard,
		"Linux Management":   cfg.Modules.LinuxManagement,
		"Hardware Monitor":   cfg.Modules.HardwareMonitor,
		"Container Manager":  cfg.Modules.ContainerManager,
		"DNS Gateway":        cfg.Modules.DNSGateway,
		"Web Monitor":        cfg.Modules.WebMonitor,
		"Terminal":           cfg.Modules.Terminal,
		"Security":           cfg.Modules.Security,
	}
	for name, enabled := range modules {
		if enabled {
			log.Info().Str("module", name).Msg("Module enabled")
		}
	}
}

func setupLogging(level, format string) {
	switch format {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	default:
		// JSON format (default for production)
	}

	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
