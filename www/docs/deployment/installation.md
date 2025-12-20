# Installation

Install Mist on your server with a single command.

## Prerequisites

### System Requirements

- **OS**: Ubuntu 20.04+, Debian 11+, or any Linux with systemd
- **RAM**: Minimum 1GB, recommended 2GB+
- **Disk**: Minimum 10GB free space
- **Docker**: Version 20.10 or later
- **Ports**: 80, 443, 3000 must be available

### Docker Installation

If Docker is not installed:

```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Verify installation
docker --version
```

## Quick Install

### One-Line Installation

```bash
curl -sSL https://raw.githubusercontent.com/corecollectives/mist/main/install.sh | bash
```

This script will:
1. Download the latest Mist binary
2. Install to `/usr/local/bin/mist`
3. Create data directory at `/var/lib/mist`
4. Set up Traefik reverse proxy
5. Create systemd service for auto-start
6. Configure firewall (UFW if installed)
7. Start Mist service

### What Gets Installed

```bash
/usr/local/bin/mist              # Mist binary
/var/lib/mist/                   # Data directory
  ├── mist.db                    # SQLite database
  ├── repos/                     # Git repositories
  └── builds/                    # Build artifacts
/etc/systemd/system/mist.service # Systemd service
/etc/mist/                       # Configuration files
  └── traefik/                   # Traefik configs
```

## Manual Installation

### 1. Download Binary

```bash
# Find latest release
LATEST=$(curl -s https://api.github.com/repos/corecollectives/mist/releases/latest | grep tag_name | cut -d '"' -f 4)

# Download for your architecture
# Linux AMD64
wget https://github.com/corecollectives/mist/releases/download/${LATEST}/mist-linux-amd64

# Linux ARM64
wget https://github.com/corecollectives/mist/releases/download/${LATEST}/mist-linux-arm64

# Make executable
chmod +x mist-linux-*
sudo mv mist-linux-* /usr/local/bin/mist
```

### 2. Create Data Directory

```bash
sudo mkdir -p /var/lib/mist
sudo chown $USER:$USER /var/lib/mist
```

### 3. Set Up Traefik

Create `docker-compose.yml` for Traefik:

```yaml
version: '3'

services:
  traefik:
    image: traefik:v2.10
    container_name: traefik
    restart: unless-stopped
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
```

Start Traefik:

```bash
docker-compose up -d
```

### 4. Create Systemd Service

Create `/etc/systemd/system/mist.service`:

```ini
[Unit]
Description=Mist PaaS
After=network.target docker.service

[Service]
Type=simple
User=root
WorkingDirectory=/var/lib/mist
ExecStart=/usr/local/bin/mist
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable mist
sudo systemctl start mist
```

### 5. Configure Firewall

```bash
# UFW
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 3000/tcp

# firewalld
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --permanent --add-port=3000/tcp
sudo firewall-cmd --reload
```

## Verify Installation

### Check Service Status

```bash
sudo systemctl status mist
```

### Check Logs

```bash
sudo journalctl -u mist -f
```

### Access Dashboard

Navigate to `http://your-server-ip:3000`

## Post-Installation

### 1. Create Admin Account

Visit the web interface and create your admin account.

### 2. Configure GitHub Integration

Follow the [GitHub App Setup Guide](./github-app).

### 3. Set Up SSL (Coming Soon)

Currently requires manual Traefik configuration. Automatic Let's Encrypt integration coming soon.

## Uninstallation

```bash
# Stop service
sudo systemctl stop mist
sudo systemctl disable mist

# Remove files
sudo rm /usr/local/bin/mist
sudo rm -rf /var/lib/mist
sudo rm /etc/systemd/system/mist.service

# Stop Traefik
docker-compose down

# Remove data (optional)
docker volume prune
```

## Troubleshooting

### Mist Won't Start

Check logs:
```bash
sudo journalctl -u mist -n 50
```

### Port Already in Use

Check what's using port 3000:
```bash
sudo lsof -i :3000
```

### Docker Socket Permission Denied

Add user to docker group:
```bash
sudo usermod -aG docker $USER
newgrp docker
```

## Next Steps

- [Configuration](./configuration)
- [GitHub App Setup](./github-app)
- [Security Best Practices](./security)
