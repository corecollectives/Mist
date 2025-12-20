# Backup & Recovery

Backup and restore your Mist installation.

## What to Backup

### SQLite Database

Most important - contains all your data:
- Users and sessions
- Projects and applications
- Deployments history
- Environment variables
- Domains and configurations

**Location**: `/var/lib/mist/mist.db`

### Configuration Files

- `/etc/mist/config.yml`
- `/etc/mist/github-app.pem`
- `/etc/mist/traefik/` directory

### Optional

- Git repositories (`/var/lib/mist/repos/`)
- Build logs (`/var/lib/mist/logs/`)
- Docker volumes (if using databases)

## Manual Backup

### Database Backup

```bash
# Stop Mist to ensure consistency
sudo systemctl stop mist

# Backup database
sudo cp /var/lib/mist/mist.db /backup/mist_$(date +%Y%m%d).db

# Start Mist
sudo systemctl start mist
```

### Full Backup

```bash
# Backup everything
sudo tar -czf /backup/mist_full_$(date +%Y%m%d).tar.gz \
  /var/lib/mist \
  /etc/mist

# Copy to remote location
scp /backup/mist_full_*.tar.gz user@backup-server:/backups/
```

## Automated Backup Script

Create `/usr/local/bin/mist-backup.sh`:

```bash
#!/bin/bash

BACKUP_DIR="/var/backups/mist"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=7

mkdir -p $BACKUP_DIR

# Backup database
cp /var/lib/mist/mist.db $BACKUP_DIR/mist_$DATE.db

# Backup configuration
tar -czf $BACKUP_DIR/config_$DATE.tar.gz /etc/mist

# Remove old backups
find $BACKUP_DIR -name "mist_*.db" -mtime +$RETENTION_DAYS -delete
find $BACKUP_DIR -name "config_*.tar.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: mist_$DATE.db"
```

Make executable and schedule:

```bash
sudo chmod +x /usr/local/bin/mist-backup.sh

# Add to crontab (daily at 2 AM)
echo "0 2 * * * /usr/local/bin/mist-backup.sh" | sudo crontab -
```

## Backup to Cloud Storage

### AWS S3

```bash
# Install AWS CLI
sudo apt install awscli

# Configure credentials
aws configure

# Backup to S3
aws s3 cp /var/lib/mist/mist.db s3://your-bucket/mist-backups/mist_$(date +%Y%m%d).db
```

### Rsync to Remote Server

```bash
# Setup SSH keys first
rsync -avz /var/lib/mist/mist.db user@backup-server:/backups/mist/
```

## Restore from Backup

### Database Restore

```bash
# Stop Mist
sudo systemctl stop mist

# Restore database
sudo cp /backup/mist_20250115.db /var/lib/mist/mist.db

# Fix permissions
sudo chown mist:mist /var/lib/mist/mist.db
sudo chmod 600 /var/lib/mist/mist.db

# Start Mist
sudo systemctl start mist
```

### Full Restore

```bash
# Stop Mist
sudo systemctl stop mist

# Extract backup
sudo tar -xzf /backup/mist_full_20250115.tar.gz -C /

# Fix permissions
sudo chown -R mist:mist /var/lib/mist
sudo chown -R root:root /etc/mist

# Start Mist
sudo systemctl start mist
```

## Disaster Recovery

### Complete Server Loss

1. **Provision new server**
2. **Install Mist** (installation script)
3. **Stop Mist**: `sudo systemctl stop mist`
4. **Restore database** from backup
5. **Restore configurations** from backup
6. **Start Mist**: `sudo systemctl start mist`
7. **Verify** everything works
8. **Update DNS** if IP changed

### Database Corruption

```bash
# Check database integrity
sqlite3 /var/lib/mist/mist.db "PRAGMA integrity_check;"

# If corrupted, restore from last good backup
sudo systemctl stop mist
sudo cp /backup/mist_YYYYMMDD.db /var/lib/mist/mist.db
sudo systemctl start mist
```

## Backup Best Practices

### The 3-2-1 Rule

- **3** copies of data
- **2** different storage types
- **1** offsite copy

### Backup Checklist

- [ ] Automated daily backups configured
- [ ] Backups stored on separate disk/server
- [ ] Offsite/cloud backup configured
- [ ] Backup retention policy defined (7-30 days)
- [ ] Backup integrity tested regularly
- [ ] Restore procedure documented
- [ ] Recovery time objective (RTO) defined
- [ ] Recovery point objective (RPO) defined

### Test Your Backups

Regularly test restore procedures:

```bash
# Test restore on separate server
scp backup-server:/backups/mist_latest.db /tmp/test-restore.db
sqlite3 /tmp/test-restore.db "SELECT COUNT(*) FROM users;"
```

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Backup Features in Development</h4>
</div>

- **Built-in Backup UI** - Configure backups from dashboard
- **One-Click Restore** - Restore from UI
- **S3 Integration** - Direct cloud backup
- **Scheduled Backups** - Configurable schedules
- **Backup Verification** - Automatic integrity checks
- **Point-in-Time Recovery** - Restore to specific moment

## Getting Help

If you encounter issues with backup or recovery:

- [GitHub Issues](https://github.com/corecollectives/mist/issues)
- [GitHub Discussions](https://github.com/corecollectives/mist/discussions)
