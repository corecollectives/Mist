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

# Prompt for Let's Encrypt email
echo "📧 Let's Encrypt Configuration"
read -p "Enter your email for Let's Encrypt certificates (default: admin@example.com): " LETSENCRYPT_EMAIL
LETSENCRYPT_EMAIL=${LETSENCRYPT_EMAIL:-admin@example.com}

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
    echo "❌ Unsupported Linux distro. Please install git, curl, and build tools manually."
    exit 1
fi

echo "📦 Installing dependencies..."
eval $PKG_INSTALL

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

# ===============================
# ✅ ADDED: Create Traefik network
# ===============================
echo "🌐 Ensuring traefik-net Docker network exists..."
if ! docker network inspect traefik-net >/dev/null 2>&1; then
    docker network create traefik-net
fi

# ===============================
# Configure Traefik static config
# ===============================
echo "⚙️ Configuring Traefik..."
cd $INSTALL_DIR
sed -i "s/email: admin@example.com/email: $LETSENCRYPT_EMAIL/" traefik-static.yml

# ===============================
# ✅ ADDED: Start Traefik
# ===============================
echo "🚦 Starting Traefik..."
docker compose -f traefik-compose.yml up -d

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

# -------------------------------
# Setup Traefik configuration directory
# -------------------------------
echo "🚦 Setting up Traefik configuration directory..."
sudo mkdir -p /var/lib/mist/traefik
sudo chown $USER:$USER /var/lib/mist/traefik

# -------------------------------
# Open firewall port
# -------------------------------
echo "🌐 Checking firewall rules..."
if command -v ufw &>/dev/null; then
    sudo ufw allow $PORT/tcp
    sudo ufw reload
elif command -v firewall-cmd &>/dev/null; then
    sudo firewall-cmd --permanent --add-port=${PORT}/tcp
    sudo firewall-cmd --reload
elif command -v iptables &>/dev/null; then
    sudo iptables -C INPUT -p tcp --dport $PORT -j ACCEPT 2>/dev/null || \
        sudo iptables -A INPUT -p tcp --dport $PORT -j ACCEPT
    if command -v netfilter-persistent &>/dev/null; then
        sudo netfilter-persistent save
    fi
else
    echo "⚠️ No recognized firewall found. Ensure port $PORT is open manually if needed."
fi

# -------------------------------
# Create or update systemd service
# -------------------------------
echo "🛠️ Creating/Updating systemd service..."
sudo bash -c "cat > $SERVICE_FILE" <<EOL
[Unit]
Description=$APP_NAME Service
After=network.target

[Service]
Type=simple
WorkingDirectory=$INSTALL_DIR/$GO_BACKEND_DIR
ExecStart=$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME
Restart=always
RestartSec=5
User=$USER
Environment=PORT=$PORT

[Install]
WantedBy=multi-user.target
EOL

sudo systemctl daemon-reload
sudo systemctl enable $APP_NAME
sudo systemctl restart $APP_NAME

echo "✅ $APP_NAME updated and running on port $PORT!"
echo "📍 Installation path: $INSTALL_DIR"
echo "🧩 Service: systemctl status $APP_NAME"
echo "🌐 Access the application at http://<your-server-ip>:$PORT"