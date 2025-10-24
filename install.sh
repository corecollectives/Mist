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

if [ -d "$INSTALL_DIR/.git" ]; then
    echo "Repository exists. Pulling latest changes from branch $BRANCH..."
    cd $INSTALL_DIR
    git fetch origin $BRANCH
    git reset --hard origin/$BRANCH
else
    echo "Cloning repository from branch $BRANCH..."
    git clone -b $BRANCH --single-branch $REPO $INSTALL_DIR
fi

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

echo "Opening port $PORT in firewall if applicable..."

if command -v ufw &>/dev/null; then
    echo "Configuring UFW..."
    sudo ufw allow $PORT/tcp
    sudo ufw reload
elif command -v firewall-cmd &>/dev/null; then
    echo "Configuring firewalld..."
    sudo firewall-cmd --permanent --add-port=${PORT}/tcp
    sudo firewall-cmd --reload
elif command -v iptables &>/dev/null; then
    echo "Configuring iptables..."
    sudo iptables -C INPUT -p tcp --dport $PORT -j ACCEPT 2>/dev/null || \
        sudo iptables -A INPUT -p tcp --dport $PORT -j ACCEPT
    # Save iptables rules for persistence
    if command -v netfilter-persistent &>/dev/null; then
        sudo netfilter-persistent save
    elif command -v service &>/dev/null; then
        sudo service iptables save || true
    fi
else
    echo "No recognized firewall found. Make sure port $PORT is open manually if needed."
fi


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
sudo systemctl restart $APP_NAME

echo "Installation complete! $APP_NAME is running on port $PORT"
