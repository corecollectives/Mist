# Over-The-Air (OTA) Updates Implementation

This document describes the OTA update system implementation for Mist PAAS.

## Overview

The OTA update system allows administrators to check for, download, and install updates directly from the web interface. Updates are fetched from GitHub releases and built from source on the server.

## Architecture

### Backend Components

1. **Database (server/db/migrations/024_Create_System_updates.sql)**
   - `system_info` table: Stores system version and other metadata
   - `update_history` table: Tracks update attempts and their status

2. **Models (server/models/system.go)**
   - `SystemInfo`: Current version information
   - `UpdateHistory`: Update attempt records
   - `GithubRelease`: GitHub release data structure
   - Helper functions for version management

3. **API Handlers (server/api/handlers/system/systemHandler.go)**
   - `GET /api/system/version` - Get current version
   - `GET /api/system/health` - System health check
   - `GET /api/system/updates/check` - Check for updates from GitHub
   - `POST /api/system/updates/trigger` - Trigger update process
   - `GET /api/system/updates/history` - Get update history
   - `GET /api/system/updates/status` - Get current update status

4. **Update Script (scripts/update.sh)**
   - Clones repository from GitHub
   - Builds frontend (Bun)
   - Builds backend (Go)
   - Creates backups before updating
   - Stops service, replaces binaries, restarts service
   - Auto-rollback on failure

### Frontend Components

1. **Types (dash/src/types/system.ts)**
   - TypeScript interfaces for all system update related data

2. **Service (dash/src/services/system.service.ts)**
   - API client for system endpoints

3. **UI Component (dash/src/components/common/system-updates.tsx)**
   - System version display
   - Update checker with GitHub integration
   - Changelog viewer (markdown support)
   - Update trigger button
   - Update history with status badges
   - Real-time update progress

4. **Settings Integration (dash/src/pages/Settings.tsx)**
   - System Updates section (admin-only)
   - Positioned at the top of settings page

## How It Works

### Update Flow

1. **Check for Updates**
   - User clicks "Check for Updates"
   - System fetches latest release from GitHub API
   - Compares current version with latest release
   - Shows update notification if available

2. **View Changelog**
   - Release notes are displayed in markdown format
   - Expandable/collapsible changelog viewer
   - Shows new features, bug fixes, breaking changes

3. **Trigger Update**
   - User clicks "Update Now"
   - Confirmation dialog appears
   - Backend creates update history record
   - Update script runs in background:
     ```
     ├─ Backup current installation
     ├─ Clone repo from GitHub
     ├─ Build frontend (bun)
     ├─ Build backend (go)
     ├─ Stop mist service
     ├─ Replace binary & static files
     ├─ Start mist service
     └─ Health check
     ```

4. **Auto-reload**
   - Frontend polls for completion
   - Page reloads automatically after 10 seconds
   - User sees new version

### Rollback Mechanism

If update fails at any step:
- Service is stopped
- Backup is restored
- Service is restarted
- Error is logged in update_history
- User is notified

## Security

- **Admin Only**: Only users with role="admin" can access update features
- **Confirmation Required**: User must confirm before updating
- **Audit Logging**: All update attempts are logged with user ID
- **Source Verification**: Updates only from official GitHub repository

## Installation Requirements

The update system requires these tools to be installed:
- **Git**: For cloning repository
- **Go**: For building backend
- **Bun**: For building frontend
- **systemd**: For service management

These are installed automatically by the main `install.sh` script.

## Usage

### For End Users (Admins)

1. Navigate to **Settings** in the web interface
2. At the top, you'll see the **System Version** section
3. Click **"Check for Updates"**
4. If an update is available:
   - View the changelog
   - Click **"Update Now"**
   - Confirm the update
   - Wait for the system to restart (30-60 seconds)

### For Developers

**Creating a New Release:**

1. Update version in code if needed
2. Create a new tag:
   ```bash
   git tag -a v1.1.0 -m "Release v1.1.0"
   git push origin v1.1.0
   ```

3. Create a GitHub release:
   - Go to GitHub repository
   - Click "Releases" → "Draft a new release"
   - Select the tag
   - Write release notes (markdown supported)
   - Publish release

4. Users can now update to this version

## Environment Variables

- `SERVER_IP`: (Optional) Override auto-detected server IP
- `PORT`: (Default: 8080) API server port

## File Locations

- **Install Directory**: `/opt/mist`
- **Database**: `/var/lib/mist/mist.db`
- **Backup Directory**: `/opt/mist-backup`
- **Update Script**: `/opt/mist/scripts/update.sh`
- **Service File**: `/etc/systemd/system/mist.service`

## Troubleshooting

### Update Fails

Check the update history in Settings to see the error message.

Common issues:
- **Git not found**: Install git
- **Go not found**: Install Go 1.22+
- **Bun not found**: Install Bun
- **Permission denied**: Ensure user has sudo access
- **Network error**: Check internet connection

### Manual Rollback

If automatic rollback fails:

```bash
sudo systemctl stop mist
cd /opt/mist/server
cp /opt/mist-backup/mist ./mist
rm -rf static
cp -r /opt/mist-backup/static ./
sudo systemctl start mist
```

### Check Service Status

```bash
sudo systemctl status mist
```

### View Logs

```bash
sudo journalctl -u mist -f
```

## Future Enhancements

Possible future improvements:
- [ ] Scheduled updates (maintenance windows)
- [ ] Auto-update option (with opt-in)
- [ ] Update channels (stable, beta, nightly)
- [ ] Delta updates (only changed files)
- [ ] Blue-green deployment (zero downtime)
- [ ] Email notifications on update completion
- [ ] Backup management (keep last N backups)
- [ ] Database migration rollback support
- [ ] Update pre-checks (disk space, etc.)

## API Reference

### Check for Updates

```http
GET /api/system/updates/check
```

Response:
```json
{
  "success": true,
  "data": {
    "hasUpdate": true,
    "currentVersion": "1.0.0",
    "latestVersion": "1.1.0",
    "release": {
      "tag_name": "v1.1.0",
      "name": "Release 1.1.0",
      "body": "## Changes\n\n- Feature A\n- Bug fix B",
      "published_at": "2024-01-15T10:00:00Z",
      "html_url": "https://github.com/..."
    }
  }
}
```

### Trigger Update

```http
POST /api/system/updates/trigger
Content-Type: application/json

{
  "version": "1.1.0",
  "branch": "main"
}
```

Response:
```json
{
  "success": true,
  "data": {
    "updateId": 1,
    "message": "Update started. The system will restart shortly."
  }
}
```

### Get Update History

```http
GET /api/system/updates/history
```

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "fromVersion": "1.0.0",
      "toVersion": "1.1.0",
      "status": "success",
      "startedAt": "2024-01-15T10:00:00Z",
      "completedAt": "2024-01-15T10:02:00Z",
      "initiatedBy": 1
    }
  ]
}
```

## Testing

To test the update system:

1. **Development Environment:**
   ```bash
   # Set a low version in database
   sqlite3 /var/lib/mist/mist.db "UPDATE system_info SET value='0.9.0' WHERE key='version'"
   
   # Create a test release on GitHub
   # Then use the UI to check and install updates
   ```

2. **Production Environment:**
   - Always test updates in a staging environment first
   - Create a backup before major updates
   - Monitor logs during update process

## License

Part of the Mist PAAS platform.
