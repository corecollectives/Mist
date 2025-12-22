# CLI Tool <span class="coming-soon-badge">Coming Soon</span>

Command-line interface for deploying and managing applications.

<div class="coming-soon-banner">
  <h4>🚧 Feature in Development</h4>
  <p>CLI tool is a high-priority feature for Phase 5 (Developer Experience).</p>
</div>

## Planned Commands

### Authentication

```bash
# Login to Mist instance
mist login https://mist.example.com

# Generate API token
mist token create --name "my-laptop"

# Logout
mist logout
```

### Deployment

```bash
# Deploy from current directory
mist deploy

# Deploy specific app
mist deploy --app my-app

# Deploy with environment variables
mist deploy --env NODE_ENV=production
```

### Logs

```bash
# Stream application logs
mist logs my-app

# Follow logs in real-time
mist logs my-app --follow

# Filter by time
mist logs my-app --since 1h
```

### Container Management

```bash
# List running containers
mist ps

# Restart application
mist restart my-app

# Stop application
mist stop my-app

# Start application
mist start my-app
```

### Environment Variables

```bash
# List variables
mist env list my-app

# Set variable
mist env set my-app KEY=value

# Import from file
mist env import my-app .env

# Export to file
mist env export my-app > .env
```

### Database Management

```bash
# List databases
mist db list

# Create database
mist db create postgres --name mydb

# Backup database
mist db backup mydb

# Restore database
mist db restore mydb backup.sql
```

### Projects

```bash
# List projects
mist projects

# Create project
mist project create "My Project"

# List apps in project
mist apps my-project
```

### One-Off Commands

```bash
# Run command in container
mist run my-app npm run migrate

# Open shell in container
mist shell my-app

# Execute SQL query
mist db exec mydb "SELECT * FROM users;"
```

## Installation

```bash
# Linux/macOS
curl -sSL https://get.mist.sh | sh

# Homebrew
brew install mist

# NPM
npm install -g @mist/cli

# Go
go install github.com/corecollectives/mist/cli@latest
```

## Configuration

```yaml
# ~/.mist/config.yml
default:
  endpoint: https://mist.example.com
  token: your-api-token

instances:
  production:
    endpoint: https://prod.mist.example.com
    token: prod-token
  staging:
    endpoint: https://staging.mist.example.com
    token: staging-token
```

## Expected Release

Target: Q3 2025

[View implementation plan →](https://github.com/corecollectives/mist/blob/main/roadmap.md#cli-tool)
