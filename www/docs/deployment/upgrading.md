# Upgrading Mist

Keep your Mist installation up to date.

## Check Current Version

```bash
mist --version
```

## Upgrade Methods

### Automatic Upgrade (Recommended)

<div class="coming-soon-banner">
  <h4>🚧 Coming Soon</h4>
  <p>One-click upgrade from dashboard is planned.</p>
</div>

### Manual Upgrade

```bash
# 1. Backup database
sudo systemctl stop mist
cp /var/lib/mist/mist.db /var/lib/mist/mist.db.backup

# 2. Download latest version
LATEST=$(curl -s https://api.github.com/repos/corecollectives/mist/releases/latest | grep tag_name | cut -d '"' -f 4)
wget https://github.com/corecollectives/mist/releases/download/${LATEST}/mist-linux-amd64

# 3. Replace binary
chmod +x mist-linux-amd64
sudo mv mist-linux-amd64 /usr/local/bin/mist

# 4. Restart service
sudo systemctl start mist
```

### Via Installation Script

```bash
# Re-run installation script (preserves data)
curl -sSL https://raw.githubusercontent.com/corecollectives/mist/main/install.sh | bash
```

## Database Migrations

Mist automatically runs database migrations on startup. Check logs:

```bash
sudo journalctl -u mist -n 50
```

## Rollback

If upgrade fails:

```bash
# Stop Mist
sudo systemctl stop mist

# Restore previous version
sudo cp /backup/mist-v0.1.0 /usr/local/bin/mist

# Restore database if needed
sudo cp /var/lib/mist/mist.db.backup /var/lib/mist/mist.db

# Start Mist
sudo systemctl start mist
```

## Breaking Changes

Check [CHANGELOG.md](https://github.com/corecollectives/mist/blob/main/CHANGELOG.md) for breaking changes before upgrading.

## Coming Soon

- In-app update notifications
- One-click upgrade button
- Automatic updates (opt-in)
- Rollback functionality
