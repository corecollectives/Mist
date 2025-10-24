#!/bin/bash
set -e

REPO="https://github.com/corecollectives/mist"
BRANCH="scripts"
APP_NAME="mist"
INSTALL_DIR="/opt/mist"
GO_BACKEND_DIR="server"
VITE_FRONTEND_DIR="dash"
GO_BINARY_NAME="mist"
PORT=8080
MIST_FILE="/var/lib/mist/mist.db"

echo "Detecting package manager..."
if command -v apt >/dev/null; then
    PKG_INSTALL="sudo apt update && sudo apt install -y git curl build-essential wget unzip"
elif command -v dnf >/dev/null; then
    PKG_INSTALL="sudo dnf install -y git curl gcc make wget unzip"
elif command -v yum >/dev/null; then
    PKG_INSTALL="sudo yum install -y git curl gcc make wget unzip"
elif command -v pacman >/dev/null; then
    PKG_INSTALL="sudo pacman -Sy --noconfirm git curl base-devel wget unzip"
else
    echo "Unsupported Linux distro. Please install git, curl, build tools manually."
    exit 1
fi

echo "Installing dependencies..."
eval $PKG_INSTALL

if ! command -v go &>/dev/null; then
    echo "Go not found, installing..."
    GO_URL="https://go.dev/dl/go1.22.11.linux-amd64.tar.gz"
    wget $GO_URL -O /tmp/go.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
fi

if ! command -v bun &>/dev/null; then
    echo "Installing Bun..."
    curl -fsSL https://bun.sh/install | bash
    export PATH=$HOME/.bun/bin:$PATH
    echo 'export PATH=$HOME/.bun/bin:$PATH' >> ~/.bashrc
fi

echo "Cloning repository..."
sudo rm -rf $INSTALL_DIR
git clone -b $BRANCH --single-branch $REPO $INSTALL_DIR

echo "Building frontend..."
cd $INSTALL_DIR/$VITE_FRONTEND_DIR
bun install
bun run build
cd ..

if [ ! -d "$GO_BACKEND_DIR/static" ]; then
    mkdir -p "$GO_BACKEND_DIR/static"
fi
rm -rf "$GO_BACKEND_DIR/static/*"
cp -r "$VITE_FRONTEND_DIR/dist/"* "$GO_BACKEND_DIR/static/"

echo "Building backend..."
cd $GO_BACKEND_DIR
go mod tidy
go build -o $GO_BINARY_NAME
cd ..

echo "Creating $MIST_FILE..."
sudo mkdir -p $(dirname $MIST_FILE)
sudo touch $MIST_FILE
sudo chown $USER:$USER $MIST_FILE

SERVICE_FILE="/etc/systemd/system/$APP_NAME.service"

echo "Creating systemd service..."
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
sudo systemctl start $APP_NAME

echo "Installation complete! $APP_NAME is running on port $PORT"
