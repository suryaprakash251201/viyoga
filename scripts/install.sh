#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────────────────
#  Viyoga — One-Line Installer for Ubuntu / Debian
#  
#  Usage:
#    curl -sSL https://raw.githubusercontent.com/suryaprakash251201/viyoga/main/scripts/install.sh | sudo bash
#
#  What it does:
#    1. Installs Go 1.22+ and Node.js 20+ (if missing)
#    2. Clones the Viyoga repository
#    3. Builds the Go backend binary
#    4. Builds the SvelteKit frontend (static)
#    5. Installs everything to /opt/viyoga
#    6. Creates a systemd service (viyoga.service)
#    7. Starts Viyoga on port 8080
# ─────────────────────────────────────────────────────────────────────────────

set -euo pipefail

# ── Colors ──
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

VIYOGA_VERSION="0.2.0"
INSTALL_DIR="/opt/viyoga"
CONFIG_DIR="/etc/viyoga"
DATA_DIR="/var/lib/viyoga"
LOG_DIR="/var/log/viyoga"
SERVICE_USER="viyoga"
REPO_URL="https://github.com/suryaprakash251201/viyoga.git"
BRANCH="main"
GO_VERSION="1.22.4"
NODE_VERSION="20"

# ── Helpers ──
info()    { echo -e "${CYAN}[INFO]${NC} $*"; }
success() { echo -e "${GREEN}[OK]${NC}   $*"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $*"; }
error()   { echo -e "${RED}[ERR]${NC}  $*"; exit 1; }

banner() {
    echo -e "${CYAN}"
    echo "  ╦  ╦╦╦ ╦╔═╗╔═╗╔═╗"
    echo "  ╚╗╔╝║╚╦╝║ ║║ ╦╠═╣"
    echo "   ╚╝ ╩ ╩ ╚═╝╚═╝╩ ╩"
    echo -e "${NC}"
    echo -e "  ${BOLD}Self-Hosted Server Dashboard${NC}"
    echo -e "  Version ${VIYOGA_VERSION}"
    echo ""
}

# ── Pre-flight checks ──
preflight() {
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root (use: sudo bash)"
    fi

    if ! grep -qiE 'ubuntu|debian' /etc/os-release 2>/dev/null; then
        warn "This installer is designed for Ubuntu/Debian. Continuing anyway..."
    fi

    info "Checking system requirements..."
    local ARCH
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)  GOARCH="amd64" ;;
        aarch64) GOARCH="arm64" ;;
        armv7l)  GOARCH="armv6l" ;;
        *)       error "Unsupported architecture: $ARCH" ;;
    esac
    success "Architecture: $ARCH ($GOARCH)"

    # Minimum 512MB RAM
    local TOTAL_MEM
    TOTAL_MEM=$(awk '/MemTotal/ {print int($2/1024)}' /proc/meminfo)
    if [[ $TOTAL_MEM -lt 512 ]]; then
        warn "Low memory detected (${TOTAL_MEM}MB). Viyoga recommends at least 512MB."
    fi
    success "Memory: ${TOTAL_MEM}MB"
}

# ── Install system dependencies ──
install_deps() {
    info "Installing system dependencies..."
    apt-get update -qq
    apt-get install -y -qq \
        curl wget git build-essential \
        ca-certificates gnupg lsb-release \
        sqlite3 ufw > /dev/null 2>&1
    success "System dependencies installed"
}

# ── Install Go ──
install_go() {
    if command -v go &>/dev/null; then
        local CURRENT_GO
        CURRENT_GO=$(go version | awk '{print $3}' | sed 's/go//')
        info "Go $CURRENT_GO already installed"
        # Check if version is sufficient (1.22+)
        local MAJOR MINOR
        MAJOR=$(echo "$CURRENT_GO" | cut -d. -f1)
        MINOR=$(echo "$CURRENT_GO" | cut -d. -f2)
        if [[ $MAJOR -ge 1 && $MINOR -ge 22 ]]; then
            success "Go version is sufficient ($CURRENT_GO >= 1.22)"
            return
        fi
        warn "Go version too old ($CURRENT_GO), installing $GO_VERSION..."
    fi

    info "Installing Go $GO_VERSION..."
    local GO_TAR="go${GO_VERSION}.linux-${GOARCH}.tar.gz"
    wget -q "https://go.dev/dl/${GO_TAR}" -O "/tmp/${GO_TAR}"
    rm -rf /usr/local/go
    tar -C /usr/local -xzf "/tmp/${GO_TAR}"
    rm -f "/tmp/${GO_TAR}"

    # Ensure Go is in PATH
    if ! grep -q '/usr/local/go/bin' /etc/profile.d/go.sh 2>/dev/null; then
        echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
    fi
    export PATH=$PATH:/usr/local/go/bin

    success "Go $(go version | awk '{print $3}') installed"
}

# ── Install Node.js ──
install_node() {
    if command -v node &>/dev/null; then
        local CURRENT_NODE
        CURRENT_NODE=$(node -v | sed 's/v//')
        local NODE_MAJOR
        NODE_MAJOR=$(echo "$CURRENT_NODE" | cut -d. -f1)
        if [[ $NODE_MAJOR -ge $NODE_VERSION ]]; then
            success "Node.js $CURRENT_NODE already installed (>= $NODE_VERSION)"
            return
        fi
        warn "Node.js version too old ($CURRENT_NODE), installing v${NODE_VERSION}..."
    fi

    info "Installing Node.js v${NODE_VERSION}..."
    curl -fsSL "https://deb.nodesource.com/setup_${NODE_VERSION}.x" | bash - > /dev/null 2>&1
    apt-get install -y -qq nodejs > /dev/null 2>&1
    success "Node.js $(node -v) installed"
}

# ── Create viyoga system user ──
create_user() {
    if id "$SERVICE_USER" &>/dev/null; then
        success "User '$SERVICE_USER' already exists"
    else
        info "Creating system user '$SERVICE_USER'..."
        useradd --system --shell /usr/sbin/nologin --home-dir "$INSTALL_DIR" "$SERVICE_USER"
        success "User '$SERVICE_USER' created"
    fi

    # Add to docker group if docker is installed
    if getent group docker &>/dev/null; then
        usermod -aG docker "$SERVICE_USER" 2>/dev/null || true
        success "Added '$SERVICE_USER' to docker group"
    fi
}

# ── Clone & Build ──
build_viyoga() {
    local BUILD_DIR
    BUILD_DIR=$(mktemp -d)
    info "Cloning Viyoga from $REPO_URL ($BRANCH)..."
    git clone --depth 1 --branch "$BRANCH" "$REPO_URL" "$BUILD_DIR" > /dev/null 2>&1
    success "Repository cloned"

    # Build backend
    info "Building Go backend..."
    cd "$BUILD_DIR/apps/backend"
    export PATH=$PATH:/usr/local/go/bin
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -X main.Version=$VIYOGA_VERSION" -o viyoga ./cmd/viyoga/
    success "Backend binary built"

    # Build frontend
    info "Building SvelteKit frontend (this may take a minute)..."
    cd "$BUILD_DIR/apps/frontend"
    npm ci --silent 2>/dev/null
    npm run build 2>/dev/null
    success "Frontend built"

    # Install to /opt/viyoga
    info "Installing to $INSTALL_DIR..."
    mkdir -p "$INSTALL_DIR/frontend"
    cp "$BUILD_DIR/apps/backend/viyoga" "$INSTALL_DIR/viyoga"
    chmod +x "$INSTALL_DIR/viyoga"

    # Copy frontend build (SvelteKit adapter-auto output)
    if [[ -d "$BUILD_DIR/apps/frontend/build" ]]; then
        cp -r "$BUILD_DIR/apps/frontend/build/"* "$INSTALL_DIR/frontend/"
    elif [[ -d "$BUILD_DIR/apps/frontend/.svelte-kit/output" ]]; then
        cp -r "$BUILD_DIR/apps/frontend/.svelte-kit/output/"* "$INSTALL_DIR/frontend/"
    fi

    success "Files installed to $INSTALL_DIR"

    # Cleanup
    rm -rf "$BUILD_DIR"
}

# ── Configuration ──
setup_config() {
    mkdir -p "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR"

    if [[ ! -f "$CONFIG_DIR/viyoga.yaml" ]]; then
        info "Creating default configuration..."
        cat > "$CONFIG_DIR/viyoga.yaml" << 'YAML'
server:
  host: "0.0.0.0"
  port: 8080
  frontend_dir: "/opt/viyoga/frontend"

modules:
  system_dashboard: true
  linux_management: true
  hardware_monitor: true
  dns_gateway: false
  container_manager: true
  web_monitor: true
  terminal: true
  security: true

metrics:
  poll_interval: "2s"
  retention_hours: 24
  history_enabled: true

database:
  path: "/var/lib/viyoga/viyoga.db"

auth:
  enabled: false
  jwt_secret: ""
  session_hours: 24

docker:
  socket: "unix:///var/run/docker.sock"

dns:
  engine: "technitium"
  api_url: "http://localhost:5380"
  api_token: ""

alerting:
  enabled: false
  channels:
    - type: webhook
      url: ""

logging:
  level: "info"
  format: "console"
YAML
        success "Configuration written to $CONFIG_DIR/viyoga.yaml"
    else
        warn "Configuration already exists at $CONFIG_DIR/viyoga.yaml (keeping existing)"
    fi

    chown -R "$SERVICE_USER":"$SERVICE_USER" "$INSTALL_DIR" "$DATA_DIR" "$LOG_DIR"
    chown -R root:"$SERVICE_USER" "$CONFIG_DIR"
    chmod 750 "$CONFIG_DIR"
    chmod 640 "$CONFIG_DIR/viyoga.yaml"
}

# ── Systemd Service ──
setup_systemd() {
    info "Creating systemd service..."
    cat > /etc/systemd/system/viyoga.service << EOF
[Unit]
Description=Viyoga Server Dashboard
Documentation=https://github.com/suryaprakash251201/viyoga
After=network-online.target docker.service
Wants=network-online.target

[Service]
Type=simple
User=${SERVICE_USER}
Group=${SERVICE_USER}
ExecStart=${INSTALL_DIR}/viyoga --config ${CONFIG_DIR}/viyoga.yaml
WorkingDirectory=${INSTALL_DIR}
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=${DATA_DIR} ${LOG_DIR}
PrivateTmp=true

# Environment
Environment=VIYOGA_LOG_DIR=${LOG_DIR}

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=viyoga

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable viyoga.service
    success "Systemd service created and enabled"
}

# ── UFW Firewall Rule ──
setup_firewall() {
    if command -v ufw &>/dev/null && ufw status | grep -q "active"; then
        info "Adding UFW rule for port 8080..."
        ufw allow 8080/tcp comment "Viyoga Dashboard" > /dev/null 2>&1
        success "Firewall rule added (port 8080)"
    fi
}

# ── Start ──
start_viyoga() {
    info "Starting Viyoga..."
    systemctl start viyoga.service
    sleep 2

    if systemctl is-active --quiet viyoga.service; then
        success "Viyoga is running!"
    else
        warn "Viyoga may have failed to start. Check: journalctl -u viyoga -f"
    fi
}

# ── Summary ──
print_summary() {
    local IP
    IP=$(hostname -I 2>/dev/null | awk '{print $1}' || echo "localhost")

    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                                                              ║${NC}"
    echo -e "${GREEN}║   ${BOLD}✅ Viyoga installed successfully!${NC}${GREEN}                          ║${NC}"
    echo -e "${GREEN}║                                                              ║${NC}"
    echo -e "${GREEN}╠══════════════════════════════════════════════════════════════╣${NC}"
    echo -e "${GREEN}║${NC}                                                              ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}   🌐 Dashboard:  ${BOLD}http://${IP}:8080${NC}                       ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}   📁 Install:    ${INSTALL_DIR}                              ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}   ⚙️  Config:     ${CONFIG_DIR}/viyoga.yaml               ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}   💾 Data:       ${DATA_DIR}                          ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}   📋 Logs:       journalctl -u viyoga -f                    ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}                                                              ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}   ${CYAN}Commands:${NC}                                                  ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}     sudo systemctl status viyoga                             ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}     sudo systemctl restart viyoga                            ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}     sudo systemctl stop viyoga                               ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}     sudo nano ${CONFIG_DIR}/viyoga.yaml                    ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}                                                              ${GREEN}║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# ── Uninstall function (called with --uninstall flag) ──
uninstall_viyoga() {
    banner
    warn "Uninstalling Viyoga..."

    systemctl stop viyoga.service 2>/dev/null || true
    systemctl disable viyoga.service 2>/dev/null || true
    rm -f /etc/systemd/system/viyoga.service
    systemctl daemon-reload

    rm -rf "$INSTALL_DIR"
    # Keep config and data by default
    echo ""
    read -rp "Remove configuration ($CONFIG_DIR)? [y/N]: " REMOVE_CONFIG
    if [[ "$REMOVE_CONFIG" =~ ^[Yy]$ ]]; then
        rm -rf "$CONFIG_DIR"
        success "Configuration removed"
    fi

    read -rp "Remove data ($DATA_DIR)? [y/N]: " REMOVE_DATA
    if [[ "$REMOVE_DATA" =~ ^[Yy]$ ]]; then
        rm -rf "$DATA_DIR"
        success "Data removed"
    fi

    userdel "$SERVICE_USER" 2>/dev/null || true
    success "Viyoga uninstalled"
}

# ── Main ──
main() {
    banner

    # Handle flags
    if [[ "${1:-}" == "--uninstall" ]]; then
        uninstall_viyoga
        exit 0
    fi

    if [[ "${1:-}" == "--version" ]]; then
        echo "Viyoga Installer v${VIYOGA_VERSION}"
        exit 0
    fi

    echo -e "${BOLD}This will install Viyoga v${VIYOGA_VERSION} to ${INSTALL_DIR}${NC}"
    echo ""

    preflight
    install_deps
    install_go
    install_node
    create_user
    build_viyoga
    setup_config
    setup_systemd
    setup_firewall
    start_viyoga
    print_summary
}

main "$@"
