#!/bin/bash
set -e

REPO="https://github.com/corecollectives/mist"
BRANCH="main"
APP_NAME="mist"
INSTALL_DIR="/opt/mist"
GO_BACKEND_DIR="server"
VITE_FRONTEND_DIR="dash"
GO_BINARY_NAME="mist"
PORT=8080
MIST_FILE="/var/lib/mist/mist.db"
SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"
TRAEFIK_COMPOSE="traefik-compose.yml"

echo "🔍 Detecting package manager..."
if command -v apt >/dev/null; then
    PKG_INSTALL="sudo apt update && sudo apt install -y git curl build-essential wget unzip docker.io docker-compose"
elif command -v dnf >/dev/null; then
    PKG_INSTALL="sudo dnf install -y git curl gcc make wget unzip docker docker-compose"
elif command -v yum >/dev/null; then
    PKG_INSTALL="sudo yum install -y git curl gcc make wget unzip docker docker-compose"
elif command -v pacman >/dev/null; then
    PKG_INSTALL="sudo pacman -Sy --noconfirm git curl base-devel wget unzip docker docker-compose"
else
    echo "❌ Unsupported Linux distro. Please install git, curl, docker, and build tools manually."
    exit 1
fi

echo "📦 Installing dependencies..."
eval $PKG_INSTALL

# Start and enable Docker
sudo systemctl enable docker --now
sudo usermod -aG docker $USER

# -------------------------------
# Install Go
# -------------------------------
if ! command -v go &>/dev/null; then
    echo "🐹 Installing Go..."
    GO_URL="https://go.dev/dl/go1.22.11.linux-amd64.tar.gz"
    wget -q $GO_URL -O /tmp/go.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
fi

# -------------------------------
# Install Bun
# -------------------------------
if ! command -v bun &>/dev/null; then
    echo "🥖 Installing Bun..."
    curl -fsSL https://bun.sh/install | bash
    export PATH=$HOME/.bun/bin:$PATH
    echo 'export PATH=$HOME/.bun/bin:$PATH' >> ~/.bashrc
fi

# -------------------------------
# Clone or update Mist repo
# -------------------------------
if [ -d "$INSTALL_DIR/.git" ]; then
    echo "🔄 Updating existing Mist installation..."
    cd $INSTALL_DIR
    git fetch origin $BRANCH
    git reset --hard origin/$BRANCH
else
    echo "📥 Cloning Mist repository..."
    sudo mkdir -p $INSTALL_DIR
    sudo chown $USER:$USER $INSTALL_DIR
    git clone -b $BRANCH --single-branch $REPO $INSTALL_DIR
fi

# -------------------------------
# Setup Traefik (before building Mist)
# -------------------------------
echo "🌐 Setting up Traefik reverse proxy..."
cd $INSTALL_DIR

# Create traefik network
docker network create traefik-net 2>/dev/null || echo "✅ traefik-net already exists"

# Start Traefik with docker-compose
if [ -f "$TRAEFIK_COMPOSE" ]; then
    echo "🚀 Starting Traefik..."
    docker compose -f $TRAEFIK_COMPOSE up -d
    echo "✅ Traefik running on ports 80/443"
else
    echo "⚠️ $TRAEFIK_COMPOSE not found, skipping Traefik setup"
fi

# -------------------------------
# Build frontend
# -------------------------------
echo "🧱 Building frontend..."
cd $INSTALL_DIR/$VITE_FRONTEND_DIR
bun install
bun run build
cd ..

if [ ! -d "$GO_BACKEND_DIR/static" ]; then
    mkdir -p "$GO_BACKEND_DIR/static"
fi
rm -rf "$INSTALL_DIR/$GO_BACKEND_DIR/static/*"
cp -r "$VITE_FRONTEND_DIR/dist/"* "$INSTALL_DIR/$GO_BACKEND_DIR/static/"

# -------------------------------
# Build backend
# -------------------------------
echo "⚙️ Building backend..."
cd "$INSTALL_DIR/$GO_BACKEND_DIR"
go mod tidy
go build -o "$GO_BINARY_NAME"

# -------------------------------
# Setup database file
# -------------------------------
echo "🗃️ Ensuring Mist database file exists..."
sudo mkdir -p $(dirname $MIST_FILE)
sudo touch $MIST_FILE
sudo chown $USER:$USER $MIST_FILE
sudo chmod 666 $MIST_FILE  # SQLite needs full access

# -------------------------------
# Create mist data directories
# -------------------------------
sudo mkdir -p /var/lib/mist/uploads/avatar /var/lib/mist/projects /var/lib/mist/logs
sudo chown -R $USER:$USER /var/lib/mist
sudo chmod -R 775 /var/lib/mist

# -------------------------------
# Open firewall ports (Traefik + Mist)
# -------------------------------
echo "🌐 Checking firewall rules..."
if command -v ufw &>/dev/null; then
    sudo ufw allow 80/tcp
    sudo ufw allow 443/tcp
    sudo ufw allow $PORT/tcp
    sudo ufw reload
elif command -v firewall-cmd &>/dev/null; then
    sudo firewall-cmd --permanent --add-port=80/tcp
    sudo firewall-cmd --permanent --add-port=443/tcp
    sudo firewall-cmd --permanent --add-port=${PORT}/tcp
    sudo firewall-cmd --reload
else
    echo "⚠️ No recognized firewall found. Ensure ports 80, 443, $PORT are open."
fi

# -------------------------------
# Create or update systemd service
# -------------------------------
echo "🛠️ Creating/Updating systemd service..."
sudo bash -c "cat > $SERVICE_FILE" <<EOL
[Unit]
Description=$APP_NAME Service
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
WorkingDirectory=$INSTALL_DIR/$GO_BACKEND_DIR
ExecStart=$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME
Restart=always
RestartSec=5
User=$USER
Group=$USER
Environment=PORT=$PORT
Environment=DB_PATH=$MIST_FILE

[Install]
WantedBy=multi-user.target
EOL

sudo systemctl daemon-reload
sudo systemctl enable $APP_NAME
sudo systemctl restart $APP_NAME

echo "✅ $APP_NAME updated and running!"
echo "📍 Installation path: $INSTALL_DIR"
echo "🌐 Traefik: http://localhost (ports 80/443)"
echo "🧩 Mist service: systemctl status $APP_NAME"
echo "🐳 Docker networks: docker network ls | grep traefik-net"
echo "🔍 Traefik dashboard: http://localhost:8080/dashboard/ (if enabled)"
