package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the full application configuration.
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Modules ModulesConfig `yaml:"modules"`
	Metrics MetricsConfig `yaml:"metrics"`
	DB      DBConfig      `yaml:"database"`
	Auth    AuthConfig    `yaml:"auth"`
	Docker  DockerConfig  `yaml:"docker"`
	DNS     DNSConfig     `yaml:"dns"`
	Alert   AlertConfig   `yaml:"alerting"`
	Log     LogConfig     `yaml:"logging"`
}

type ServerConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	FrontendDir string `yaml:"frontend_dir"`
}

type ModulesConfig struct {
	SystemDashboard  bool `yaml:"system_dashboard"`
	LinuxManagement  bool `yaml:"linux_management"`
	HardwareMonitor  bool `yaml:"hardware_monitor"`
	DNSGateway       bool `yaml:"dns_gateway"`
	ContainerManager bool `yaml:"container_manager"`
	WebMonitor       bool `yaml:"web_monitor"`
	Terminal         bool `yaml:"terminal"`
	Security         bool `yaml:"security"`
}

type MetricsConfig struct {
	PollInterval   string `yaml:"poll_interval"`
	RetentionHours int    `yaml:"retention_hours"`
	HistoryEnabled bool   `yaml:"history_enabled"`
}

type DBConfig struct {
	Path string `yaml:"path"`
}

type AuthConfig struct {
	Enabled      bool   `yaml:"enabled"`
	JWTSecret    string `yaml:"jwt_secret"`
	SessionHours int    `yaml:"session_hours"`
}

type DockerConfig struct {
	Socket string `yaml:"socket"`
}

type DNSConfig struct {
	Engine   string `yaml:"engine"`
	APIURL   string `yaml:"api_url"`
	APIToken string `yaml:"api_token"`
}

type AlertConfig struct {
	Enabled  bool            `yaml:"enabled"`
	Channels []AlertChannel  `yaml:"channels"`
}

type AlertChannel struct {
	Type     string `yaml:"type"`
	URL      string `yaml:"url,omitempty"`
	BotToken string `yaml:"bot_token,omitempty"`
	ChatID   string `yaml:"chat_id,omitempty"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Host:        "0.0.0.0",
			Port:        8080,
			FrontendDir: "./frontend/build",
		},
		Modules: ModulesConfig{
			SystemDashboard:  true,
			LinuxManagement:  true,
			HardwareMonitor:  true,
			DNSGateway:       false,
			ContainerManager: true,
			WebMonitor:       true,
			Terminal:         true,
			Security:         true,
		},
		Metrics: MetricsConfig{
			PollInterval:   "2s",
			RetentionHours: 24,
			HistoryEnabled: true,
		},
		DB: DBConfig{
			Path: "./data/viyoga.db",
		},
		Auth: AuthConfig{
			Enabled:      false,
			SessionHours: 24,
		},
		Docker: DockerConfig{
			Socket: "unix:///var/run/docker.sock",
		},
		DNS: DNSConfig{
			Engine: "technitium",
			APIURL: "http://localhost:5380",
		},
		Alert: AlertConfig{
			Enabled: false,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "console",
		},
	}
}

// Load reads config from a YAML file and merges with defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Use defaults if no config file
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

// PollDuration returns the poll interval as a time.Duration.
func (c *Config) PollDuration() time.Duration {
	d, err := time.ParseDuration(c.Metrics.PollInterval)
	if err != nil {
		return 2 * time.Second
	}
	return d
}

// ListenAddr returns the formatted listen address.
func (c *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
