#!/bin/bash
set -Eeuo pipefail

LOG_FILE="/tmp/mist-install.log"
: > "$LOG_FILE"

REAL_USER="${SUDO_USER:-$USER}"
REAL_HOME="$(getent passwd "$REAL_USER" | cut -d: -f6)"

REPO="https://github.com/corecollectives/mist"
BRANCH="release"
APP_NAME="mist"
INSTALL_DIR="/opt/mist"
GO_BACKEND_DIR="server"
GO_BINARY_NAME="mist"
PORT=8080
MIST_FILE="/var/lib/mist/mist.db"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"

SPINNER_PID=""
SUDO_KEEPALIVE_PID=""

spinner() {
    local i=0
    local chars='|/-\'
    while :; do
        i=$(( (i + 1) % 4 ))
        printf "\r⏳ %c" "${chars:$i:1}"
        sleep 0.1
    done
}

run_step() {
    local msg="$1"
    local cmd="$2"

    printf "\n▶ %s\n" "$msg"
    spinner &
    SPINNER_PID=$!

    bash -c "$cmd" >>"$LOG_FILE" 2>&1 || true

    kill "$SPINNER_PID" >/dev/null 2>&1 || true
    wait "$SPINNER_PID" 2>/dev/null || true
    printf "\r\033[K✔ Done\n"
}

cleanup() {
    kill "$SPINNER_PID" >/dev/null 2>&1 || true
    kill "$SUDO_KEEPALIVE_PID" >/dev/null 2>&1 || true
}
trap cleanup EXIT

fail() {
    cleanup
    echo
    tail -20 "$LOG_FILE"
    echo "❌ Installation failed"
    exit 1
}
trap fail ERR

echo "📧 Let's Encrypt configuration"
read -rp "Email (default: admin@example.com): " LETSENCRYPT_EMAIL
LETSENCRYPT_EMAIL="${LETSENCRYPT_EMAIL:-admin@example.com}"

echo "🔐 Verifying sudo access..."
sudo -v

(
    while true; do
        sleep 60
        sudo -n true || exit
    done
) 2>/dev/null &
SUDO_KEEPALIVE_PID=$!

if command -v apt >/dev/null; then
    PKG_INSTALL="sudo apt update && sudo apt install -y git curl build-essential wget unzip"
elif command -v dnf >/dev/null; then
    PKG_INSTALL="sudo dnf install -y git curl gcc make wget unzip"
elif command -v yum >/dev/null; then
    PKG_INSTALL="sudo yum install -y git curl gcc make wget unzip"
elif command -v pacman >/dev/null; then
    PKG_INSTALL="sudo pacman -Sy --noconfirm git curl base-devel wget unzip"
else
    exit 1
fi

run_step "Installing system dependencies" "$PKG_INSTALL"

command -v docker >/dev/null || exit 1

if ! command -v go >/dev/null; then
    run_step "Installing Go" "
        wget -q https://go.dev/dl/go1.22.11.linux-amd64.tar.gz -O /tmp/go.tar.gz &&
        sudo rm -rf /usr/local/go &&
        sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    "
    grep -q '/usr/local/go/bin' "$REAL_HOME/.bashrc" || \
        echo 'export PATH=$PATH:/usr/local/go/bin' >>"$REAL_HOME/.bashrc"
    export PATH="$PATH:/usr/local/go/bin"
fi

if [ -d "$INSTALL_DIR/.git" ]; then
    run_step "Updating Mist ($BRANCH)" "
        cd '$INSTALL_DIR' &&
        git fetch origin '$BRANCH' &&
        git reset --hard origin/'$BRANCH'
    "
else
    run_step "Cloning Mist repository" "
        sudo mkdir -p '$INSTALL_DIR' &&
        sudo chown '$REAL_USER:$REAL_USER' '$INSTALL_DIR' &&
        git clone -b '$BRANCH' --single-branch '$REPO' '$INSTALL_DIR'
    "
fi

run_step "Ensuring Traefik Docker network" "
    docker network inspect traefik-net >/dev/null 2>&1 ||
    docker network create traefik-net
"

[ -f "$INSTALL_DIR/traefik-static.yml" ] || exit 1

run_step "Configuring Traefik" "
    grep -q '$LETSENCRYPT_EMAIL' '$INSTALL_DIR/traefik-static.yml' ||
    sed -i \"s|^\\s*email:.*|  email: $LETSENCRYPT_EMAIL|\" '$INSTALL_DIR/traefik-static.yml'
"

run_step "Starting Traefik" "
    docker compose -f '$INSTALL_DIR/traefik-compose.yml' up -d
"

[ -d "$INSTALL_DIR/$GO_BACKEND_DIR/static" ] || exit 1

run_step "Building backend" "
    cd '$INSTALL_DIR/$GO_BACKEND_DIR' &&
    go mod tidy &&
    go build -o '$GO_BINARY_NAME'
"

run_step "Preparing data directories" "
    sudo mkdir -p /var/lib/mist/traefik &&
    sudo mkdir -p $(dirname "$MIST_FILE") &&
    sudo touch '$MIST_FILE' &&
    sudo chown -R '$REAL_USER:$REAL_USER' /var/lib/mist
"

run_step "Configuring firewall" "
    if command -v ufw >/dev/null; then
        sudo ufw allow $PORT/tcp || true
        sudo ufw reload || true
    elif command -v firewall-cmd >/dev/null; then
        sudo firewall-cmd --permanent --add-port=${PORT}/tcp || true
        sudo firewall-cmd --reload || true
    fi
"

run_step "Installing systemd service" "
    sudo tee '$SERVICE_FILE' >/dev/null <<EOF
[Unit]
Description=$APP_NAME Service
After=network.target docker.service
Requires=docker.service

[Service]
WorkingDirectory=$INSTALL_DIR/$GO_BACKEND_DIR
ExecStart=$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME
Restart=always
RestartSec=5
User=$REAL_USER
Environment=PORT=$PORT
StandardOutput=journal
StandardError=journal
SyslogIdentifier=mist

[Install]
WantedBy=multi-user.target
EOF
"

run_step "Starting $APP_NAME service" "
    sudo systemctl daemon-reload &&
    sudo systemctl enable '$APP_NAME' || true &&
    sudo systemctl restart '$APP_NAME'
"

SERVER_IP="$(curl -fsSL https://api.ipify.org || hostname -I | awk '{print $1}')"
URL="http://$SERVER_IP:$PORT"

echo
echo "╔════════════════════════════════════════════╗"
echo "║ 🎉 Mist is now running                      ║"
echo "║                                            ║"
echo "║ 👉 Open this in your browser:               ║"
echo "║                                            ║"
printf "║   %-40s ║\n" "$URL"
echo "║                                            ║"
echo "║ (Ctrl + Click the link above if supported) ║"
echo "╚════════════════════════════════════════════╝"
echo
echo "📄 Logs: $LOG_FILE"
