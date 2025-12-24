#!/bin/bash
set -Eeuo pipefail

# Update script for Mist - Self-updating with safety mechanisms
# This script MUST be bulletproof as it updates itself

LOG_FILE="/tmp/mist-update.log"
: > "$LOG_FILE"

REAL_USER="${SUDO_USER:-$USER}"
REAL_HOME="$(getent passwd "$REAL_USER" | cut -d: -f6)"

REPO="https://github.com/corecollectives/mist"
BRANCH="release"
APP_NAME="mist"
INSTALL_DIR="/opt/mist"
GO_BACKEND_DIR="server"
GO_BINARY_NAME="mist"
DB_FILE="/var/lib/mist/mist.db"
BACKUP_DIR="/var/lib/mist/backups"
LOCK_FILE="/var/lib/mist/update.lock"

SPINNER_PID=""
SUDO_KEEPALIVE_PID=""
BACKUP_COMMIT=""
BACKUP_DB=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[INFO]${NC} $1" | tee -a "$LOG_FILE"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

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

    if bash -c "$cmd" >>"$LOG_FILE" 2>&1; then
        kill "$SPINNER_PID" >/dev/null 2>&1 || true
        wait "$SPINNER_PID" 2>/dev/null || true
        printf "\r\033[K✔ Done\n"
        return 0
    else
        kill "$SPINNER_PID" >/dev/null 2>&1 || true
        wait "$SPINNER_PID" 2>/dev/null || true
        printf "\r\033[K✘ Failed\n"
        return 1
    fi
}

cleanup() {
    kill "$SPINNER_PID" >/dev/null 2>&1 || true
    kill "$SUDO_KEEPALIVE_PID" >/dev/null 2>&1 || true
    
    # Remove lock file
    rm -f "$LOCK_FILE"
}
trap cleanup EXIT

rollback() {
    error "Update failed! Attempting rollback..."
    
    # Rollback git
    if [ -n "$BACKUP_COMMIT" ]; then
        warn "Rolling back to commit: $BACKUP_COMMIT"
        cd "$INSTALL_DIR"
        git reset --hard "$BACKUP_COMMIT" >>"$LOG_FILE" 2>&1 || true
        
        # Rebuild from rollback
        cd "$INSTALL_DIR/$GO_BACKEND_DIR"
        go mod tidy >>"$LOG_FILE" 2>&1 || true
        go build -o "$GO_BINARY_NAME" >>"$LOG_FILE" 2>&1 || true
        
        # Restart service
        sudo systemctl restart "$APP_NAME" >>"$LOG_FILE" 2>&1 || true
    fi
    
    # Rollback database if backup exists
    if [ -n "$BACKUP_DB" ] && [ -f "$BACKUP_DB" ]; then
        warn "Rolling back database to backup"
        cp "$BACKUP_DB" "$DB_FILE" || true
    fi
    
    error "Rollback completed. Check logs at: $LOG_FILE"
    cleanup
    exit 1
}
trap rollback ERR

# ---------------- Pre-flight checks ----------------

log "Starting Mist update process..."

# Check if running as root or with sudo
if [ "$EUID" -eq 0 ] || [ -n "${SUDO_USER:-}" ]; then
    log "Running with elevated privileges"
else
    error "This script requires sudo privileges"
    exit 1
fi

# Verify sudo access
echo "🔐 Verifying sudo access..."
sudo -v

# Keep sudo alive
(
    while true; do
        sleep 60
        sudo -n true || exit
    done
) 2>/dev/null &
SUDO_KEEPALIVE_PID=$!

# Check for lock file (prevent concurrent updates)
if [ -f "$LOCK_FILE" ]; then
    error "Another update is already in progress!"
    error "If this is incorrect, remove: $LOCK_FILE"
    exit 1
fi

# Create lock file
echo "$$" > "$LOCK_FILE"
log "Update lock acquired"

# Check if Mist is installed
if [ ! -d "$INSTALL_DIR/.git" ]; then
    error "Mist is not installed at $INSTALL_DIR"
    error "Please run install.sh first"
    exit 1
fi

# Check if service is running
if ! sudo systemctl is-active --quiet "$APP_NAME"; then
    warn "Mist service is not running!"
    warn "Continuing anyway, but this is unusual..."
fi

# Check available disk space (need at least 500MB)
AVAILABLE_SPACE=$(df "$INSTALL_DIR" | tail -1 | awk '{print $4}')
if [ "$AVAILABLE_SPACE" -lt 500000 ]; then
    error "Insufficient disk space! Need at least 500MB free"
    exit 1
fi
log "Disk space check passed"

# Check if Go is available
if ! command -v go >/dev/null 2>&1; then
    error "Go is not installed or not in PATH"
    exit 1
fi
log "Go installation verified"

# Check if Docker is available
if ! command -v docker >/dev/null 2>&1; then
    error "Docker is not installed or not in PATH"
    exit 1
fi
log "Docker installation verified"

# ---------------- Create backup directory ----------------

sudo mkdir -p "$BACKUP_DIR"
sudo chown "$REAL_USER:$REAL_USER" "$BACKUP_DIR"
log "Backup directory ready"

# ---------------- Backup database ----------------

TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_DB="$BACKUP_DIR/mist-$TIMESTAMP.db"

if [ -f "$DB_FILE" ]; then
    log "Creating database backup..."
    if cp "$DB_FILE" "$BACKUP_DB"; then
        log "Database backed up to: $BACKUP_DB"
    else
        warn "Failed to backup database, but continuing..."
        BACKUP_DB=""
    fi
fi

# ---------------- Fetch latest from release branch ----------------

run_step "Fetching latest updates from $BRANCH branch" "
    cd '$INSTALL_DIR' &&
    git fetch origin '$BRANCH'
"

# Check if updates are available
cd "$INSTALL_DIR"
LOCAL_COMMIT=$(git rev-parse HEAD)
REMOTE_COMMIT=$(git rev-parse origin/$BRANCH)
BACKUP_COMMIT="$LOCAL_COMMIT"

if [ "$LOCAL_COMMIT" = "$REMOTE_COMMIT" ]; then
    log "Already up to date"
    exit 0
fi

log "Updates available"
log "Current commit: ${LOCAL_COMMIT:0:7}"
log "Latest commit:  ${REMOTE_COMMIT:0:7}"

# Show what's changing
log "Changes in this update:"
git log --oneline "$LOCAL_COMMIT..$REMOTE_COMMIT" | head -10 | tee -a "$LOG_FILE"

# ---------------- Create git backup tag ----------------

run_step "Creating backup tag" "
    cd '$INSTALL_DIR' &&
    git tag -f backup-$TIMESTAMP
"

# ---------------- Update repository ----------------

if ! run_step "Updating repository to latest version" "
    cd '$INSTALL_DIR' &&
    git reset --hard origin/'$BRANCH'
"; then
    error "Failed to update repository"
    rollback
fi

# Verify the update
NEW_COMMIT=$(git rev-parse HEAD)
if [ "$NEW_COMMIT" != "$REMOTE_COMMIT" ]; then
    error "Repository update verification failed!"
    rollback
fi
log "Repository update verified"

# ---------------- Run database migrations ----------------

log "Running database migrations (if any)..."
cd "$INSTALL_DIR/$GO_BACKEND_DIR"

# Build migration runner first (if it exists)
if [ -d "db/migrations" ]; then
    log "Migration files detected"
    # Migrations will run automatically on service start
fi

# ---------------- Rebuild backend ----------------

if ! run_step "Rebuilding backend binary" "
    cd '$INSTALL_DIR/$GO_BACKEND_DIR' &&
    go mod tidy &&
    go build -o '$GO_BINARY_NAME'
"; then
    error "Failed to rebuild backend"
    rollback
fi

# Verify binary was created and is executable
if [ ! -f "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME" ]; then
    error "Backend binary was not created!"
    rollback
fi

if [ ! -x "$INSTALL_DIR/$GO_BACKEND_DIR/$GO_BINARY_NAME" ]; then
    error "Backend binary is not executable!"
    rollback
fi
log "Backend binary verified"

# ---------------- Rebuild CLI ----------------

if [ -d "$INSTALL_DIR/cli" ]; then
    if ! run_step "Rebuilding CLI tool" "
        cd '$INSTALL_DIR/cli' &&
        go mod tidy &&
        go build -o mist-cli
    "; then
        warn "Failed to rebuild CLI, but continuing..."
    else
        if ! run_step "Installing CLI tool" "
            sudo cp '$INSTALL_DIR/cli/mist-cli' /usr/local/bin/mist-cli &&
            sudo chmod +x /usr/local/bin/mist-cli
        "; then
            warn "Failed to install CLI, but continuing..."
        fi
    fi
fi

# ---------------- Restart service ----------------

log "Stopping Mist service..."
if ! sudo systemctl stop "$APP_NAME" >>"$LOG_FILE" 2>&1; then
    warn "Failed to stop service gracefully, forcing..."
    sudo systemctl kill "$APP_NAME" >>"$LOG_FILE" 2>&1 || true
    sleep 2
fi

run_step "Reloading systemd and starting service" "
    sudo systemctl daemon-reload &&
    sudo systemctl start '$APP_NAME'
"

# ---------------- Health check ----------------

log "Waiting for service to start..."
RETRY_COUNT=0
MAX_RETRIES=30

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if sudo systemctl is-active --quiet "$APP_NAME"; then
        log "Service is active"
        break
    fi
    sleep 1
    RETRY_COUNT=$((RETRY_COUNT + 1))
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    error "Service failed to start within 30 seconds!"
    error "Rolling back..."
    rollback
fi

# Wait a bit more and check if service is still running
sleep 5

if ! sudo systemctl is-active --quiet "$APP_NAME"; then
    error "Service started but crashed immediately!"
    error "Check logs: sudo journalctl -u $APP_NAME -n 50"
    rollback
fi

log "Service health check passed"

# Try to connect to the service
log "Checking HTTP endpoint..."
if curl -f -s -o /dev/null -w "%{http_code}" "http://localhost:8080/api/health" | grep -q "200"; then
    log "HTTP health check passed"
else
    warn "HTTP health check failed, but service is running"
    warn "This might be normal if the service is still initializing"
fi

# ---------------- Update Traefik ----------------

if [ -f "$INSTALL_DIR/traefik-compose.yml" ]; then
    if ! run_step "Updating Traefik" "
        docker compose -f '$INSTALL_DIR/traefik-compose.yml' up -d
    "; then
        warn "Failed to update Traefik, but Mist update succeeded"
    fi
fi

# ---------------- Cleanup old backups ----------------

log "Cleaning up old backups (keeping last 5)..."
cd "$BACKUP_DIR"
ls -t mist-*.db 2>/dev/null | tail -n +6 | xargs -r rm -f

# Keep last 10 backup tags
cd "$INSTALL_DIR"
git tag | grep "^backup-" | sort -r | tail -n +11 | xargs -r git tag -d

# ---------------- Success ----------------

log ""
log "╔════════════════════════════════════════════╗"
log "║ 🎉 Mist has been updated successfully      ║"
log "║                                            ║"
log "║ From: ${LOCAL_COMMIT:0:7}                             ║"
log "║ To:   ${REMOTE_COMMIT:0:7}                             ║"
log "╚════════════════════════════════════════════╝"
log ""
log "📄 Update logs: $LOG_FILE"
log "💾 Database backup: $BACKUP_DB"
log "🔄 Service status: sudo systemctl status $APP_NAME"
log "📋 Service logs: sudo journalctl -u $APP_NAME -n 50"
log ""
log "If you experience issues, you can rollback with:"
log "  cd $INSTALL_DIR && git reset --hard backup-$TIMESTAMP"
log "  sudo systemctl restart $APP_NAME"
