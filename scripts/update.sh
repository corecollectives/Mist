#!/bin/bash
set -e

# Mist Self-Update Script
# This script updates Mist to a new version from GitHub

REPO="https://github.com/corecollectives/mist"
INSTALL_DIR="/opt/mist"
GO_BACKEND_DIR="server"
VITE_FRONTEND_DIR="dash"
GO_BINARY_NAME="mist"
SERVICE_NAME="mist"
TEMP_DIR="/tmp/mist-update-$$"
BACKUP_DIR="/opt/mist-backup"
DB_FILE="/var/lib/mist/mist.db"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Log function
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if running as appropriate user
if [ "$EUID" -eq 0 ]; then
    error "Do not run this script as root. Run as the user that owns the mist service."
    exit 1
fi

# Parse arguments
VERSION=${1:-"latest"}
BRANCH=${2:-"main"}

log "Starting Mist update to version: $VERSION"

# Step 1: Create backup
log "Creating backup of current installation..."
if [ -d "$BACKUP_DIR" ]; then
    rm -rf "$BACKUP_DIR"
fi
mkdir -p "$BACKUP_DIR"

# Backup binary and static files (not the entire directory)
if [ -f "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME" ]; then
    cp "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME" "$BACKUP_DIR/"
    success "Backed up binary"
fi

if [ -d "$INSTALL_DIR/$GO_BACKEND_DIR/static" ]; then
    cp -r "$INSTALL_DIR/$GO_BACKEND_DIR/static" "$BACKUP_DIR/"
    success "Backed up static files"
fi

# Backup database (just in case)
if [ -f "$DB_FILE" ]; then
    cp "$DB_FILE" "$BACKUP_DIR/mist.db.backup"
    success "Backed up database"
fi

# Step 2: Clone repository to temp location
log "Downloading latest code from GitHub..."
rm -rf "$TEMP_DIR"
git clone -b "$BRANCH" --single-branch --depth 1 "$REPO" "$TEMP_DIR"

if [ $? -ne 0 ]; then
    error "Failed to clone repository"
    exit 1
fi

success "Code downloaded successfully"

# Step 3: Build frontend
log "Building frontend..."
cd "$TEMP_DIR/$VITE_FRONTEND_DIR"

if ! command -v bun &>/dev/null; then
    error "Bun is not installed. Cannot build frontend."
    exit 1
fi

bun install
if [ $? -ne 0 ]; then
    error "Failed to install frontend dependencies"
    exit 1
fi

bun run build
if [ $? -ne 0 ]; then
    error "Failed to build frontend"
    exit 1
fi

success "Frontend built successfully"

# Step 4: Build backend
log "Building backend..."
cd "$TEMP_DIR/$GO_BACKEND_DIR"

if ! command -v go &>/dev/null; then
    error "Go is not installed. Cannot build backend."
    exit 1
fi

go mod tidy
if [ $? -ne 0 ]; then
    error "Failed to tidy Go modules"
    exit 1
fi

go build -o "$GO_BINARY_NAME"
if [ $? -ne 0 ]; then
    error "Failed to build backend"
    exit 1
fi

success "Backend built successfully"

# Step 5: Stop the service
log "Stopping Mist service..."
sudo systemctl stop "$SERVICE_NAME"

if [ $? -ne 0 ]; then
    warning "Failed to stop service, continuing anyway..."
fi

sleep 2

# Step 6: Replace binary
log "Replacing binary..."
cp "$TEMP_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME" "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME"

if [ $? -ne 0 ]; then
    error "Failed to replace binary"
    # Restore from backup
    log "Restoring from backup..."
    if [ -f "$BACKUP_DIR/$GO_BINARY_NAME" ]; then
        cp "$BACKUP_DIR/$GO_BINARY_NAME" "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME"
    fi
    sudo systemctl start "$SERVICE_NAME"
    exit 1
fi

success "Binary replaced"

# Step 7: Replace frontend static files
log "Replacing frontend files..."
rm -rf "$INSTALL_DIR/$GO_BACKEND_DIR/static"
mkdir -p "$INSTALL_DIR/$GO_BACKEND_DIR/static"
cp -r "$TEMP_DIR/$VITE_FRONTEND_DIR/dist/"* "$INSTALL_DIR/$GO_BACKEND_DIR/static/"

if [ $? -ne 0 ]; then
    error "Failed to replace frontend files"
    # Restore from backup
    log "Restoring from backup..."
    if [ -d "$BACKUP_DIR/static" ]; then
        rm -rf "$INSTALL_DIR/$GO_BACKEND_DIR/static"
        cp -r "$BACKUP_DIR/static" "$INSTALL_DIR/$GO_BACKEND_DIR/"
    fi
    if [ -f "$BACKUP_DIR/$GO_BINARY_NAME" ]; then
        cp "$BACKUP_DIR/$GO_BINARY_NAME" "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME"
    fi
    sudo systemctl start "$SERVICE_NAME"
    exit 1
fi

success "Frontend files replaced"

# Step 8: Start the service
log "Starting Mist service..."
sudo systemctl start "$SERVICE_NAME"

if [ $? -ne 0 ]; then
    error "Failed to start service"
    # Restore from backup
    log "Restoring from backup..."
    if [ -f "$BACKUP_DIR/$GO_BINARY_NAME" ]; then
        cp "$BACKUP_DIR/$GO_BINARY_NAME" "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME"
    fi
    if [ -d "$BACKUP_DIR/static" ]; then
        rm -rf "$INSTALL_DIR/$GO_BACKEND_DIR/static"
        cp -r "$BACKUP_DIR/static" "$INSTALL_DIR/$GO_BACKEND_DIR/"
    fi
    sudo systemctl start "$SERVICE_NAME"
    exit 1
fi

# Step 9: Wait and check service health
log "Checking service health..."
sleep 5

if ! systemctl is-active --quiet "$SERVICE_NAME"; then
    error "Service is not running after update"
    # Restore from backup
    log "Restoring from backup..."
    sudo systemctl stop "$SERVICE_NAME"
    if [ -f "$BACKUP_DIR/$GO_BINARY_NAME" ]; then
        cp "$BACKUP_DIR/$GO_BINARY_NAME" "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME"
    fi
    if [ -d "$BACKUP_DIR/static" ]; then
        rm -rf "$INSTALL_DIR/$GO_BACKEND_DIR/static"
        cp -r "$BACKUP_DIR/static" "$INSTALL_DIR/$GO_BACKEND_DIR/"
    fi
    sudo systemctl start "$SERVICE_NAME"
    exit 1
fi

success "Service is running"

# Step 10: Cleanup
log "Cleaning up..."
rm -rf "$TEMP_DIR"
success "Cleanup complete"

# Keep backup for potential rollback
log "Backup kept at: $BACKUP_DIR"

success "Update completed successfully!"
log "Mist has been updated to $VERSION"
