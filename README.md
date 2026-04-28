# Viyoga

**Self-Hosted Ubuntu Server Dashboard**

A lightweight, modular, open-source Linux server operations dashboard for Ubuntu Server — serving solo developers, homelab admins, and small teams.

![Status](https://img.shields.io/badge/status-alpha-orange)
![License](https://img.shields.io/badge/license-MIT-blue)
![Go](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)
![Svelte](https://img.shields.io/badge/svelte-5-FF3E00?logo=svelte)

## Features

- 📊 **System Dashboard** — Real-time CPU, RAM, Disk I/O, Network stats
- 🖥️ **Linux Management** — Services (systemd), users, logs, cron, firewall (UFW)
- 🔧 **Hardware Monitor** — Process tree, threshold alerts, webhook/Telegram notifications
- 🛡️ **DNS Gateway** — Network-level ad/tracker blocking via Technitium DNS
- 🐳 **Container Manager** — Docker lifecycle management (start/stop/restart/remove/logs)
- 📡 **Web Monitor** — HTTP endpoint health checks, SSL expiry, uptime tracking
- 💻 **Terminal** — In-browser web terminal with real API commands
- 🔒 **Security** — Security checklist, user audit, firewall overview

---

## ⚡ One-Line Install (Ubuntu / Debian)

```bash
curl -sSL https://raw.githubusercontent.com/suryaprakash251201/viyoga/main/scripts/install.sh | sudo bash
```

This single command will:
1. Install Go 1.22+ and Node.js 20+ (if missing)
2. Clone the repository & build from source
3. Install the binary to `/opt/viyoga`
4. Create a systemd service (`viyoga.service`)
5. Start the dashboard on **port 8080**

After installation, open `http://<your-server-ip>:8080` in your browser.

### Manage the Service

```bash
# Check status
sudo systemctl status viyoga

# View live logs
sudo journalctl -u viyoga -f

# Restart after config changes
sudo systemctl restart viyoga

# Stop the service
sudo systemctl stop viyoga

# Edit configuration
sudo nano /etc/viyoga/viyoga.yaml
```

### Uninstall

```bash
curl -sSL https://raw.githubusercontent.com/suryaprakash251201/viyoga/main/scripts/install.sh | sudo bash -s -- --uninstall
```

---

## 🛠️ Development Setup

```bash
# Clone the repository
git clone https://github.com/suryaprakash251201/viyoga.git
cd viyoga

# Start the Go backend
cd apps/backend
go run ./cmd/viyoga/ --config ../../configs/viyoga.example.yaml

# In another terminal, start the frontend
cd apps/frontend
npm install
npm run dev
```

Open `http://localhost:5173` in your browser.

## Tech Stack

| Layer | Technology |
|-------|-----------:|
| Backend | Go (Chi router, gorilla/websocket) |
| Frontend | SvelteKit 5 + DaisyUI + ApexCharts |
| Database | SQLite (embedded, zero-config) |
| Metrics | gopsutil (/proc, /sys parsing) |
| Deployment | systemd, Docker Compose |

## Project Structure

```
viyoga/
├── apps/
│   ├── backend/         # Go API server + metric collectors
│   │   ├── cmd/viyoga/  # Main entry point
│   │   └── internal/    # Modules: collector, linux, docker, dns, monitor, alerting
│   └── frontend/        # SvelteKit dashboard UI
│       └── src/routes/  # Pages: dashboard, system, hardware, docker, dns, monitor, terminal, security
├── scripts/             # install.sh (one-line installer)
├── deploy/              # Docker & deployment configs
├── configs/             # Default configuration templates
└── docs/                # Documentation
```

## Configuration

The config file is located at `/etc/viyoga/viyoga.yaml` (production) or `configs/viyoga.example.yaml` (development).

Key sections:
```yaml
server:
  host: "0.0.0.0"
  port: 8080

modules:
  system_dashboard: true    # CPU, RAM, Disk, Network gauges
  linux_management: true    # systemd, journalctl, UFW, users, cron
  hardware_monitor: true    # Process list, threshold alerts
  dns_gateway: false        # Requires Technitium DNS
  container_manager: true   # Docker container/image management
  web_monitor: true         # HTTP health checks
  terminal: true            # Web terminal
  security: true            # Security center

docker:
  socket: "unix:///var/run/docker.sock"

alerting:
  enabled: false
  channels:
    - type: webhook
      url: "https://hooks.slack.com/..."
    - type: telegram
      bot_token: ""
      chat_id: ""
```

## API Endpoints

| Module | Endpoints |
|--------|-----------|
| System | `GET /api/v1/system`, `GET /api/v1/metrics/current`, `WS /api/v1/ws` |
| Linux | `GET/POST /api/v1/linux/services`, `/logs`, `/firewall`, `/users`, `/cron` |
| Hardware | `GET /api/v1/hardware/processes`, `/alerts/rules`, `/alerts/events` |
| Docker | `GET/POST/DELETE /api/v1/docker/containers`, `/images`, `/prune` |
| DNS | `GET /api/v1/dns/stats`, `/querylog`, `/blocklists` |
| Monitor | `GET/POST/DELETE /api/v1/monitor/targets`, `/status` |

## Design Principles

- 🪶 Total idle memory < 150MB RAM
- 🏠 100% self-hosted, no cloud dependency
- 📦 Modular — enable/disable features via config
- 📱 Mobile-responsive, light & dark mode
- ⚡ Single binary backend
- 🔐 Runs as unprivileged systemd service

## Requirements

- **OS:** Ubuntu 20.04+ / Debian 11+
- **RAM:** 512MB minimum
- **CPU:** Any (x86_64, ARM64, ARMv7)
- **Go:** 1.22+ (auto-installed)
- **Node.js:** 20+ (auto-installed, build-time only)

## License

MIT
