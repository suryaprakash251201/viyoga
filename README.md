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
- 🔧 **Hardware Monitor** — Sensors, process tree, threshold alerts
- 🛡️ **DNS Gateway** — Network-level ad/tracker blocking (Pi-hole style)
- 🐳 **Container Manager** — Docker lifecycle management
- 📡 **Web Monitor** — HTTP endpoint health checks, uptime tracking
- 💻 **Terminal** — In-browser shell (xterm.js)
- 🔒 **Security** — Fail2Ban, UFW, SSL/TLS, RBAC

## Quick Start

```bash
# Clone the repository
git clone https://github.com/viyoga/viyoga.git
cd viyoga

# Start the Go backend
cd apps/backend
go run ./cmd/viyoga/

# In another terminal, start the frontend
cd apps/frontend
npm install
npm run dev
```

Open `http://localhost:5173` in your browser.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go (Chi router, gorilla/websocket) |
| Frontend | SvelteKit 5 + DaisyUI + ApexCharts |
| Database | SQLite (embedded, zero-config) |
| Metrics | gopsutil (/proc, /sys parsing) |
| Deployment | Docker Compose |

## Project Structure

```
viyoga/
├── apps/
│   ├── backend/     # Go API server + metric collectors
│   └── frontend/    # SvelteKit dashboard UI
├── deploy/          # Docker & deployment configs
├── configs/         # Default configuration templates
├── docs/            # Documentation
└── scripts/         # Build & dev scripts
```

## Configuration

Copy the example config:
```bash
cp configs/viyoga.example.yaml viyoga.yaml
```

## Design Principles

- 🪶 Total idle memory < 150MB RAM
- 🏠 100% self-hosted, no cloud dependency
- 📦 Modular — enable/disable features via config
- 📱 Mobile-responsive, dark mode by default
- ⚡ Single binary backend

## License

MIT
