# Applications

Applications are the core deployable units in Mist. Each application represents a containerized service that can be a web application, background service, or managed database.

## Application Types

Mist supports three distinct application types:

### Web Applications (`web`)

Web applications are services exposed via HTTP/HTTPS through Traefik reverse proxy:

- **Public access** via custom domains
- **Auto-generated domains** (if wildcard DNS configured)
- **SSL certificates** managed automatically
- **Port configuration** (default: 3000)
- **Git-based deployments** from GitHub repositories

**Use cases**: Websites, web apps, REST APIs, GraphQL servers

### Service Applications (`service`)

Background services that run without external HTTP access:

- **Internal only** - Not exposed to public internet
- **Port 3000** (internal, not routable)
- **Git-based deployments** from GitHub repositories
- **Same deployment flow** as web apps

**Use cases**: Workers, queue processors, scheduled tasks, internal services

### Database Applications (`database`)

Managed database services using pre-configured Docker templates:

- **Template-based** deployment (PostgreSQL, MySQL, Redis, MongoDB)
- **No Git repository** required
- **Pre-configured** CPU and memory limits from template
- **Version control** via Docker image tags
- **Environment variables** for configuration (passwords, databases, etc.)

**Use cases**: PostgreSQL, MySQL, Redis, MongoDB, or any other database service

## Creating an Application

### Via Dashboard

1. Navigate to your project
2. Click **"New Application"**
3. Select application type:
   - **Web** - For public-facing applications
   - **Service** - For background workers
   - **Database** - For managed database services
4. Fill in the required fields based on type

### Web/Service Configuration

- **Name**: Unique identifier within project
- **Description**: Optional application description
- **Port**: Port your app listens on (default: 3000 for web)
- **Environment Variables**: Optional key-value pairs

### Database Configuration

- **Name**: Database instance name
- **Description**: Optional description
- **Template**: Select from available templates:
  - PostgreSQL
  - MySQL
  - Redis
  - MongoDB
- **Environment Variables**: Database credentials and settings

::: tip Auto-Generated Domains
For web applications, Mist automatically generates a subdomain if wildcard DNS is configured (format: `{project-name}-{app-name}.{domain}`).
:::

## Resource Configuration

### CPU Limits

Set CPU core limits for your application:

- **Default**: No limit (uses available CPU)
- **Range**: 0.1 to N cores (e.g., 0.5, 1.0, 2.0)
- **Database defaults**: Applied automatically from templates

```javascript
// API example
{
  "appId": 1,
  "cpuLimit": 1.0  // 1 CPU core
}
```

### Memory Limits

Set memory limits in megabytes:

- **Default**: No limit (uses available memory)
- **Range**: 128MB to system maximum
- **Database defaults**: Applied automatically from templates

```javascript
// API example
{
  "appId": 1,
  "memoryLimit": 512  // 512 MB
}
```

### Restart Policies

Control container restart behavior:

- **`no`**: Never restart automatically
- **`always`**: Always restart if stopped
- **`on-failure`**: Restart only on failure
- **`unless-stopped`**: Restart unless manually stopped (default)

```javascript
// API example
{
  "appId": 1,
  "restartPolicy": "unless-stopped"
}
```

## Container Configuration

### Dockerfile Support

Mist supports custom Dockerfiles in your repository:

1. **Default location**: `./Dockerfile` in repository root
2. **Custom path**: Specify via `dockerfilePath` setting
3. **Root directory**: Set `rootDirectory` for builds in subdirectories

```dockerfile
# Example Dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --production
COPY . .
EXPOSE 3000
CMD ["node", "server.js"]
```

::: tip Multi-Stage Builds
Use multi-stage Dockerfiles to keep production images small and secure.
:::

### Health Checks

Configure health checks for automatic container monitoring:

- `healthcheck_enabled` - Enable/disable health checks
- `healthcheck_path` - HTTP endpoint to check (e.g., `/health`)
- `healthcheck_interval` - Seconds between checks
- `healthcheck_timeout` - Seconds before timeout
- `healthcheck_retries` - Failed attempts before marking unhealthy

```javascript
// API example
{
  "appId": 1,
  "healthcheck_enabled": true,
  "healthcheck_path": "/health",
  "healthcheck_interval": 30,
  "healthcheck_timeout": 5,
  "healthcheck_retries": 3
}
```

::: warning Health Check Endpoint Required
Your application must implement the health check endpoint. It should return a 200 status code when healthy.
:::

## Git Integration

### Connecting a Repository

For web and service applications:

1. **Install GitHub App** on your GitHub account/organization
2. **Grant access** to repositories in Mist settings
3. **Select repository** when creating application
4. **Choose branch** to deploy from

### Configuration Settings

- **Git Repository**: Full HTTPS URL (e.g., `https://github.com/user/repo`)
- **Git Branch**: Branch name to deploy (e.g., `main`, `develop`)
- **Root Directory**: Subdirectory for monorepos (optional)
- **Dockerfile Path**: Custom Dockerfile location (optional)

```javascript
// API example
{
  "appId": 1,
  "gitRepository": "https://github.com/myorg/myapp",
  "gitBranch": "main",
  "rootDirectory": "services/api",
  "dockerfilePath": "docker/Dockerfile.prod"
}
```

## Environment Variables

Applications can have unlimited environment variables for configuration. See the [Environment Variables guide](./environment-variables) for details.

### Adding Variables

1. Go to **Environment Variables** tab
2. Click **"Add Variable"** or use **"Bulk Paste"**
3. Enter key-value pairs
4. Variables are available immediately (restart may be needed)

### Build-time vs Runtime

- **Build-time**: Available during Docker image build
- **Runtime**: Available when container runs
- **Both**: All variables are available in both phases

::: warning Secrets Management
Never commit secrets to Git. Always use environment variables for sensitive data like API keys, passwords, and tokens.
:::

## Deployment

### Manual Deployment

Trigger deployment from dashboard:

1. Click **"Deploy"** button
2. Latest commit is fetched automatically
3. Deployment added to queue
4. Build starts when worker is available

For databases, deployment uses the configured Docker image version from the template.

### Automatic Deployment

With GitHub App configured:

- **Push events**: Auto-deploy on push to configured branch
- **Webhooks**: GitHub triggers deployment automatically
- **Status updates**: Real-time deployment progress

[Learn more about deployments →](./deployments)

## Deployment Strategies

Choose how new versions are deployed:

- **`rolling`**: Stop old container, start new container (default)
- Future: Blue-green, canary deployments

```javascript
{
  "deploymentStrategy": "rolling"
}
```

## Domains

Add custom domains to web applications:

1. Go to **Domains** tab
2. Click **"Add Domain"**
3. Enter domain name (e.g., `app.example.com`)
4. Configure DNS A record to point to your server
5. Wait for DNS verification
6. SSL certificate issued automatically via Let's Encrypt

[Learn more about domains →](./domains)

## Container Controls

### Start/Stop/Restart

Control container lifecycle from dashboard or API:

- **Stop**: Gracefully stops container
- **Start**: Starts stopped container
- **Restart**: Restarts running container

### Application Status

Container states:

- **`stopped`**: Container is not running
- **`running`**: Container is active and healthy
- **`error`**: Container failed to start or crashed
- **`building`**: Docker image is being built
- **`deploying`**: Container is being deployed

## Monitoring

### Real-time Logs

View live container logs via WebSocket:

- **Real-time streaming**: See logs as they're generated
- **Last 100 lines**: Initial connection shows tail
- **Auto-reconnect**: Maintains connection with ping/pong
- **Timestamps**: Each log line includes timestamp

[View logs tab in dashboard or use WebSocket API](/api/websockets)

### Deployment History

Track all deployments:

- **Deployment number**: Sequential numbering
- **Commit info**: Hash and message (for Git apps)
- **Status**: pending → building → deploying → success/failed
- **Duration**: Build and deployment time
- **Build logs**: Complete build output
- **Active deployment**: Currently running version

### Container Statistics

Monitor resource usage in real-time:

- CPU usage percentage
- Memory consumption
- Network I/O
- Container uptime

## Application Settings

Update application configuration:

- **General**: Name, description
- **Git**: Repository URL, branch, paths
- **Resources**: CPU, memory limits
- **Container**: Restart policy, health checks
- **Status**: Control running state

::: tip Audit Logging
All configuration changes are logged in the audit log with user information and timestamps.
:::

## Rollback

Rollback to a previous successful deployment:

1. Go to **Deployments** tab
2. Find successful deployment
3. Click **"Rollback"** button
4. Previous version is redeployed

::: warning Database State
Rolling back code doesn't rollback database migrations. Plan rollback strategies for schema changes.
:::

## Service Templates

Database applications use pre-configured service templates stored in the database:

### Available Templates

- **PostgreSQL**: Latest PostgreSQL with recommended 1 CPU / 512MB RAM
- **MySQL**: Latest MySQL with custom configuration
- **Redis**: Latest Redis for caching
- **MongoDB**: Latest MongoDB for document storage

### Template Properties

- `docker_image`: Docker image name
- `docker_image_version`: Image tag (e.g., `latest`, `14-alpine`)
- `default_port`: Container port
- `recommended_cpu`: Suggested CPU limit
- `recommended_memory`: Suggested memory limit (MB)

## Best Practices

### Application Structure

Organize your repository:

```
my-app/
├── src/               # Source code
├── Dockerfile         # Production Dockerfile
├── .dockerignore      # Exclude unnecessary files
├── package.json       # Dependencies
├── .env.example       # Document required env vars
└── README.md          # Setup instructions
```

### Resource Planning

- **Start small**: Begin with default limits
- **Monitor usage**: Check metrics before scaling
- **Database memory**: Databases need sufficient RAM
- **CPU limits**: Set realistic limits for consistent performance

### Environment Variables

- Document all required variables in README
- Use descriptive key names (e.g., `DATABASE_URL` not `DB`)
- Separate development and production configs
- Rotate secrets regularly

### Build Optimization

Use `.dockerignore` to exclude files:

```
node_modules
.git
.env
*.log
*.md
.vscode
.idea
```

Use multi-stage builds:

```dockerfile
# Build stage
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Production stage
FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
CMD ["node", "dist/server.js"]
```

### Health Check Implementation

Add health check endpoint to your application:

```javascript
// Express.js example
app.get('/health', (req, res) => {
  // Check database connection
  // Check external dependencies
  res.status(200).json({ 
    status: 'healthy',
    timestamp: new Date().toISOString()
  });
});
```

```python
# Flask example
@app.route('/health')
def health():
    return {'status': 'healthy'}, 200
```

## Troubleshooting

### Build Failures

Check deployment logs for errors:

- **Missing dependencies**: Check package.json/requirements.txt
- **Build command failed**: Verify build scripts
- **Out of disk space**: Clean up old images
- **Permission denied**: Check file permissions in Dockerfile

### Container Won't Start

Common issues:

- **Wrong start command**: Container exits immediately
- **Port mismatch**: Application port doesn't match config
- **Missing env vars**: Required variables not set
- **Syntax errors**: Check application logs for errors

### Application Not Accessible

Checklist:

- ✅ Container status is `running` (green indicator)
- ✅ Domain has A record pointing to server
- ✅ DNS has propagated (check with `dig` or `nslookup`)
- ✅ Firewall allows ports 80/443
- ✅ Application is listening on correct port
- ✅ No errors in container logs

### Performance Issues

Investigate:

- Check CPU/memory usage in metrics
- Review container logs for errors
- Increase resource limits if needed
- Check database query performance
- Review application profiling data

## API Reference

See the [Applications API documentation](/api/applications) for programmatic access to:

- Create applications
- Update configuration  
- Control containers
- Get status and logs
- Manage resources

## Related Topics

- [Deployments](./deployments) - Deploy your applications
- [Environment Variables](./environment-variables) - Configure applications
- [Domains](./domains) - Add custom domains
- [Git Integration](./git-integration) - Connect GitHub repositories
- [Logs](./logs) - Monitor application logs
