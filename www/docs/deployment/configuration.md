# Configuration

Configure Mist for your environment.

## Environment Variables

Mist can be configured via environment variables or config file.

### Database

```bash
# SQLite database path (default: ./mist.db)
MIST_DB_PATH=/var/lib/mist/mist.db
```

### Server

```bash
# Server port (default: 3000)
MIST_PORT=3000

# Server host (default: 0.0.0.0)
MIST_HOST=0.0.0.0
```

### GitHub Integration

```bash
# GitHub App credentials
GITHUB_APP_ID=your-app-id
GITHUB_APP_PRIVATE_KEY=/path/to/private-key.pem
GITHUB_WEBHOOK_SECRET=your-webhook-secret
```

### Security

```bash
# JWT secret for token signing
JWT_SECRET=your-random-secret-key

# Session timeout in hours (default: 168 = 7 days)
SESSION_TIMEOUT=168
```

## Configuration File

Create `/etc/mist/config.yml`:

```yaml
server:
  port: 3000
  host: "0.0.0.0"

database:
  path: "/var/lib/mist/mist.db"

github:
  app_id: "123456"
  private_key_path: "/etc/mist/github-app.pem"
  webhook_secret: "your-secret"

security:
  jwt_secret: "your-random-secret"
  session_timeout: 168
```

## Docker Configuration

Mist uses the Docker socket at `/var/run/docker.sock`.

Ensure the user running Mist has Docker permissions:

```bash
sudo usermod -aG docker mist-user
```

## Traefik Configuration

See [Traefik Setup](./traefik) for detailed Traefik configuration.

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Configuration Features</h4>
</div>

- **SMTP Configuration** - Email settings
- **Backup Settings** - Automatic backup configuration
- **Resource Limits** - Global resource limits
- **Logging Configuration** - Log levels and destinations
