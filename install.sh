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

TRAEFIK_DIR="/var/lib/mist/traefik"
TRAEFIK_COMPOSE_FILE="$INSTALL_DIR/traefik-compose.yml"
TRAEFIK_NETWORK="traefik-net"

echo "🔍 Detecting package manager..."
if command -v apt >/dev/null; then
    PKG_INSTALL="sudo apt update && sudo apt install -y git curl build-essential wget unzip"
elif command -v dnf >/dev/null; then
    PKG_INSTALL="sudo dnf install -y git curl gcc make wget unzip"
elif command -v yum >/dev/null; then
    PKG_INSTALL="sudo yum install -y git curl gcc make wget unzip"
elif command -v pacman >/dev/null; then
    PKG_INSTALL="sudo pacman -Sy --noconfirm git curl base-devel wget unzip"
else
    echo "❌ Unsupported Linux distro."
    exit 1
fi

echo "📦 Installing dependencies..."
eval $PKG_INSTALL

# -------------------------------
# Install Docker
# -------------------------------
if ! command -v docker &>/dev/null; then
    echo "🐳 Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    sudo systemctl enable docker
    sudo systemctl start docker
    sudo usermod -aG docker $USER
fi

# -------------------------------
# Install Docker Compose plugin
# -------------------------------
if ! docker compose version &>/dev/null; then
    echo "🧩 Installing Docker Compose plugin..."
    sudo mkdir -p /usr/local/lib/docker/cli-plugins
    sudo curl -SL https://github.com/docker/compose/releases/download/v2.29.2/docker-compose-linux-x86_64 \
        -o /usr/local/lib/docker/cli-plugins/docker-compose
    sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
fi

# -------------------------------
# Install Go
# -------------------------------
if ! command -v go &>/dev/null; then
    echo "🐹 Installing Go..."
    wget -q https://go.dev/dl/go1.22.11.linux-amd64.tar.gz -O /tmp/go.tar.gz
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
# Clone or update Mist
# -------------------------------
if [ -d "$INSTALL_DIR/.git" ]; then
    echo "🔄 Updating Mist..."
    cd $INSTALL_DIR
    git fetch origin $BRANCH
    git reset --hard origin/$BRANCH
else
    echo "📥 Cloning Mist..."
    sudo mkdir -p $INSTALL_DIR
    sudo chown $USER:$USER $INSTALL_DIR
    git clone -b $BRANCH --single-branch $REPO $INSTALL_DIR
fi

# -------------------------------
# Build frontend
# -------------------------------
echo "🧱 Building frontend..."
cd $INSTALL_DIR/$VITE_FRONTEND_DIR
bun install
bun run build

mkdir -p "$INSTALL_DIR/$GO_BACKEND_DIR/static"
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
# Setup Mist database
# -------------------------------
echo "🗃️ Ensuring database exists..."
sudo mkdir -p $(dirname $MIST_FILE)
sudo touch $MIST_FILE
sudo chown $USER:$USER $MIST_FILE

# -------------------------------
# Create Traefik network
# -------------------------------
echo "🌐 Creating Traefik network..."
if ! docker network inspect $TRAEFIK_NETWORK >/dev/null 2>&1; then
    docker network create $TRAEFIK_NETWORK
fi

# -------------------------------
# Start Traefik
# -------------------------------
echo "🚦 Starting Traefik..."
sudo mkdir -p $TRAEFIK_DIR

cd $INSTALL_DIR
docker compose -f "$TRAEFIK_COMPOSE_FILE" up -d

# -------------------------------
# Open firewall
# -------------------------------
echo "🌐 Configuring firewall..."
if command -v ufw &>/dev/null; then
    sudo ufw allow 80/tcp
    sudo ufw allow 443/tcp
    sudo ufw allow $PORT/tcp
    sudo ufw reload
fi

# -------------------------------
# Create systemd service
# -------------------------------
echo "🛠️ Creating systemd service..."
sudo bash -c "cat > $SERVICE_FILE" <<EOL
[Unit]
Description=Mist Backend
After=network.target docker.service
Requires=docker.service

[Service]
WorkingDirectory=$INSTALL_DIR/$GO_BACKEND_DIR
ExecStart=$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME
Restart=always
User=$USER
Environment=PORT=$PORT

[Install]
WantedBy=multi-user.target
EOL

sudo systemctl daemon-reload
sudo systemctl enable mist
sudo systemctl restart mist

echo ""
echo "✅ Mist installed successfully"
echo "🌍 Backend: http://localhost:$PORT"
echo "🔐 Traefik running with automatic HTTPS"
echo ""
