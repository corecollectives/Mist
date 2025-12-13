# Mist App Type System - Complete Implementation Guide

## Overview

Mist now supports three distinct application types, each with different deployment behaviors, networking configurations, and UI presentations:

1. **Web Applications** - HTTP servers that need external access
2. **Background Services** - Workers/bots with no external access
3. **Database/Pre-made Services** - One-click deployments from Docker images

---

## App Type Characteristics

### Web Applications (`app_type = 'web'`)

**Purpose:** User-facing applications that need HTTP access from the internet

**Characteristics:**
- Built from Git repositories
- Requires port configuration
- Gets Traefik routing with domains
- External port mapping if no domains configured
- Shows git information in UI
- Supports custom build/start commands

**Docker Networking:**
- **With domains:** Connected to `traefik-net`, routed via Traefik (no direct port mapping)
- **Without domains:** Direct port mapping `HOST_PORT:CONTAINER_PORT`

**Example Use Cases:**
- React/Vue/Angular frontends
- Next.js applications
- REST APIs
- Express/FastAPI servers

---

### Background Services (`app_type = 'service'`)

**Purpose:** Background workers that don't need external HTTP access

**Characteristics:**
- Built from Git repositories
- Port configured but not exposed externally
- No Traefik routing
- Shows git information in UI
- Can communicate with other containers on internal network

**Docker Networking:**
- Connected to `traefik-net` (internal only)
- No external port mapping
- No Traefik labels

**Example Use Cases:**
- Discord bots
- Queue workers (Celery, Bull)
- Cron jobs
- Data processors

---

### Database/Pre-made Services (`app_type = 'database'`)

**Purpose:** One-click deployments of popular databases and services

**Characteristics:**
- Deployed from official Docker images (no git repo)
- Uses service templates
- Port configured but not exposed externally
- Environment variables auto-configured
- No git information in UI
- Resource limits from template recommendations

**Docker Networking:**
- Connected to `traefik-net` (internal only)
- No external port mapping
- No Traefik labels

**Example Use Cases:**
- PostgreSQL
- Redis
- MySQL/MariaDB
- MongoDB
- RabbitMQ
- MinIO

---

## Database Schema Changes

### Apps Table (`apps`)

```sql
ALTER TABLE apps ADD COLUMN app_type VARCHAR(20) DEFAULT 'web';
ALTER TABLE apps ADD COLUMN template_name VARCHAR(100);
ALTER TABLE apps ADD COLUMN cpu_limit DECIMAL(4,2);
ALTER TABLE apps ADD COLUMN memory_limit INTEGER;
ALTER TABLE apps ADD COLUMN restart_policy VARCHAR(20) DEFAULT 'unless-stopped';
```

**New Columns:**
- `app_type` - Enum: 'web', 'service', 'database'
- `template_name` - References service template (for database apps)
- `cpu_limit` - CPU cores limit (e.g., 0.5, 2.0)
- `memory_limit` - Memory in MB
- `restart_policy` - Docker restart policy

### Service Templates Table (`service_templates`)

```sql
CREATE TABLE service_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    category VARCHAR(50) NOT NULL,
    description TEXT,
    icon_url TEXT,
    docker_image VARCHAR(200) NOT NULL,
    docker_image_version VARCHAR(50),
    default_port INTEGER NOT NULL,
    default_env_vars TEXT,      -- JSON
    required_env_vars TEXT,     -- JSON
    default_volume_path TEXT,
    volume_required BOOLEAN DEFAULT 0,
    recommended_cpu DECIMAL(4,2),
    recommended_memory INTEGER,
    min_memory INTEGER,
    healthcheck_command TEXT,
    healthcheck_interval INTEGER DEFAULT 30,
    admin_ui_image TEXT,
    admin_ui_port INTEGER,
    setup_instructions TEXT,
    is_active BOOLEAN DEFAULT 1,
    is_featured BOOLEAN DEFAULT 0,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Pre-populated Templates:**
- PostgreSQL (postgres:16-alpine)
- Redis (redis:7-alpine)
- MySQL (mysql:8)
- MariaDB (mariadb:11)
- MongoDB (mongo:7)
- RabbitMQ (rabbitmq:3-management)
- MinIO (minio/minio:latest)

---

## Backend Implementation

### File: `server/models/app.go`

```go
type AppType string

const (
    AppTypeWeb      AppType = "web"
    AppTypeService  AppType = "service"
    AppTypeDatabase AppType = "database"
)

type App struct {
    // ... existing fields
    AppType       AppType   `db:"app_type" json:"appType"`
    TemplateName  *string   `db:"template_name" json:"templateName"`
    CPULimit      *float64  `db:"cpu_limit" json:"cpuLimit"`
    MemoryLimit   *int      `db:"memory_limit" json:"memoryLimit"`
    RestartPolicy string    `db:"restart_policy" json:"restartPolicy"`
}
```

**NULL-Safe Git Operations:**
- `GetAppRepoInfo()` - Uses `sql.NullString` for git_repository
- `GetAppRepoAndBranch()` - Returns error if database app

---

### File: `server/api/handlers/applications/create.go`

**Request Structure:**
```go
type CreateAppRequest struct {
    Name         string             `json:"name"`
    Description  string             `json:"description"`
    ProjectID    int64              `json:"projectId"`
    AppType      string             `json:"appType"`
    TemplateName *string            `json:"templateName"`
    Port         *int               `json:"port"`
    EnvVars      map[string]string  `json:"envVars"`
}
```

**Logic:**
1. Validates app type ('web', 'service', 'database')
2. For web apps: Sets port (default 3000)
3. For service apps: Sets internal port (not exposed)
4. For database apps:
   - Requires template name
   - Fetches template configuration
   - Sets port from template
   - Applies resource recommendations
5. Creates environment variables after app creation

---

### File: `server/docker/deployer.go`

**Deployment Branch Logic:**

```go
if app.AppType == models.AppTypeDatabase {
    // Pull Docker image (no git clone/build)
    template := models.GetServiceTemplateByName(*app.TemplateName)
    imageName := template.DockerImage + ":" + template.DockerImageVersion
    PullDockerImage(imageName, logfile)
    imageTag = imageName
} else {
    // Git clone and build (web/service apps)
    BuildImage(imageTag, appContextPath, envVars, logfile)
}
```

**Environment Variables:**
- For database apps: Merges template defaults with user-defined vars
- User vars override template defaults

---

### File: `server/docker/build.go`

**Container Networking by Type:**

```go
func RunContainer(app *models.App, imageTag, containerName string, 
                  domains []string, Port int, envVars map[string]string) error {
    
    switch app.AppType {
    case models.AppTypeWeb:
        if len(domains) > 0 {
            // Use Traefik routing
            --network traefik-net
            -l traefik.enable=true
            -l traefik.http.routers.<name>.rule=Host(`domain.com`)
        } else {
            // Direct port mapping
            -p PORT:PORT
        }
    
    case models.AppTypeService:
        // Internal network only
        --network traefik-net
        // No port mapping, no Traefik labels
    
    case models.AppTypeDatabase:
        // Internal network only
        --network traefik-net
        // No port mapping, no Traefik labels
    }
}
```

**Resource Limits:**
```bash
--cpus <cpu_limit>      # If app.CPULimit is set
-m <memory_limit>m      # If app.MemoryLimit is set
--restart <policy>      # app.RestartPolicy (default: unless-stopped)
```

---

### File: `server/queue/handleWork.go`

**Skip Git Operations for Database Apps:**
```go
if app.AppType == models.AppTypeDatabase {
    logger.Info("Database app - skipping git clone")
    // Skip to deployment
} else {
    // Clone git repository
    github.CloneRepo(...)
}
```

---

### File: `server/api/handlers/deployments/AddDeployHandler.go`

**Commit Hash Handling:**
```go
if app.AppType == models.AppTypeDatabase {
    // Use template version as "commit hash"
    template := models.GetServiceTemplateByName(*app.TemplateName)
    deployment.CommitHash = *template.DockerImageVersion
    deployment.CommitMessage = "Database deployment"
} else {
    // Fetch actual git commit
    commit := github.GetLatestCommit(...)
    deployment.CommitHash = commit.SHA
    deployment.CommitMessage = commit.Message
}
```

---

## Frontend Implementation

### File: `dash/src/types/app.ts`

```typescript
export type AppType = 'web' | 'service' | 'database';

export type App = {
  id: number;
  appType: AppType;
  templateName: string | null;
  cpuLimit: number | null;
  memoryLimit: number | null;
  restartPolicy: RestartPolicy;
  // ... other fields
};

export type CreateAppRequest = {
  projectId: number;
  appType: AppType;
  templateName?: string;
  envVars?: Record<string, string>;
  // ... other fields
};
```

---

### File: `dash/src/features/projects/components/CreateAppModal.tsx`

**Three-Step Flow:**
1. **App Type Selection** - User chooses: Web, Service, or Database
2. **Type-Specific Form:**
   - Web: Name, Description, Port
   - Service: Name, Description (no port field)
   - Database: Template selector, Name, Password generator
3. **Submit** - Creates app with appropriate configuration

---

### File: `dash/src/features/projects/components/DatabaseForm.tsx`

**Features:**
- Loads templates from API
- Displays categorized template cards
- Shows port, RAM, CPU recommendations
- Password generator
- Maps passwords to correct env var names:
  - PostgreSQL: `POSTGRES_PASSWORD`
  - MySQL: `MYSQL_ROOT_PASSWORD`
  - MariaDB: `MARIADB_ROOT_PASSWORD`
  - MongoDB: `MONGO_INITDB_ROOT_PASSWORD`, `MONGO_INITDB_ROOT_USERNAME`
  - Redis: `REDIS_PASSWORD`

---

### File: `dash/src/features/applications/AppPage.tsx`

**UI Changes:**
- Git tab hidden for database apps
- Dynamic tab columns (5 for database, 6 for web/service)
- Fetches latest commit only for non-database apps

---

### File: `dash/src/components/applications/app-info.tsx`

**Conditional Rendering:**

**For Database Apps:**
- Shows: Status, Deployment Strategy, Template Name, Port
- Hides: Git repo, branch, commit, build/start commands

**For Web/Service Apps:**
- Shows: Status, Git repo, branch, commit, port, build/start commands
- Hides: Template name

---

### File: `dash/src/components/applications/app-settings.tsx`

**Database App Restrictions:**
- Git repository fields: Disabled with info alert
- Git branch: Disabled
- Root directory: Disabled
- Shows app type badge
- Resource limits editable

---

### File: `dash/src/components/deployments/deployment-list.tsx`

**Deployment Display:**

**For Database Apps:**
```
Version: 16-alpine – Database deployment
```

**For Web/Service Apps:**
```
a7b3c2d – Fix authentication bug
```

---

## Docker Networking Summary

### Network Architecture

All containers connect to `traefik-net` Docker network for inter-container communication.

| App Type | External Port | Traefik Routing | Internal Network |
|----------|--------------|-----------------|------------------|
| Web (with domains) | No | Yes | traefik-net |
| Web (no domains) | Yes (`-p`) | No | host + traefik-net |
| Service | No | No | traefik-net |
| Database | No | No | traefik-net |

### Example Configurations

**Web App with Domain:**
```bash
docker run -d \
  --name app-12345 \
  --network traefik-net \
  -l traefik.enable=true \
  -l traefik.http.routers.app-12345.rule=Host(`myapp.com`) \
  -l traefik.http.routers.app-12345.entrypoints=web \
  -l traefik.http.services.app-12345.loadbalancer.server.port=3000 \
  myapp:latest
```

**Web App without Domain:**
```bash
docker run -d \
  --name app-12345 \
  -p 3000:3000 \
  myapp:latest
```

**Service App:**
```bash
docker run -d \
  --name app-67890 \
  --network traefik-net \
  mybot:latest
```

**Database App:**
```bash
docker run -d \
  --name app-54321 \
  --network traefik-net \
  --restart unless-stopped \
  --cpus 1.0 \
  -m 512m \
  -e POSTGRES_PASSWORD=securepass123 \
  postgres:16-alpine
```

---

## Testing Checklist

### Web App Testing
- [ ] Create web app with port 3000
- [ ] Add domain and verify Traefik routing
- [ ] Verify external access via domain
- [ ] Remove domain and verify port mapping works
- [ ] Check git tab shows commit info
- [ ] Verify build and deployment succeeds

### Service App Testing
- [ ] Create service app (Discord bot)
- [ ] Verify NO external port mapping (`docker ps`)
- [ ] Verify connected to traefik-net
- [ ] Check git tab shows commit info
- [ ] Verify deployment succeeds
- [ ] Test internal communication with web app

### Database App Testing
- [ ] Create PostgreSQL database
- [ ] Verify password created as env var
- [ ] Verify NO external port mapping
- [ ] Verify connected to traefik-net
- [ ] Check git tab is hidden
- [ ] Verify deployment shows "Version: 16-alpine"
- [ ] Test connection from web app (internal network)
- [ ] Verify resource limits applied

### Cross-App Communication Testing
- [ ] Web app connects to database (internal network)
- [ ] Service app connects to database (internal network)
- [ ] Web app receives external traffic
- [ ] Database NOT accessible from external network

---

## Key Files Modified

### Backend (Go)
1. `server/db/migrations/005_Create_App.sql` - Added columns
2. `server/db/migrations/018_Create_Service_templates.sql` - New table
3. `server/models/app.go` - AppType enum, NULL handling
4. `server/models/serviceTemplate.go` - Template CRUD
5. `server/api/handlers/applications/create.go` - Type-aware creation + env vars
6. `server/api/handlers/applications/update.go` - Resource limits
7. `server/api/handlers/applications/getLatestCommit.go` - Skip database apps
8. `server/api/handlers/deployments/AddDeployHandler.go` - Version as commit
9. `server/api/handlers/templates/list.go` - Template API
10. `server/docker/deployer.go` - Branch logic, env var merging
11. `server/docker/build.go` - Network configuration by type
12. `server/queue/handleWork.go` - Skip git for database apps

### Frontend (React/TypeScript)
1. `dash/src/types/app.ts` - Type definitions
2. `dash/src/api/endpoints/templates.ts` - Template API client
3. `dash/src/features/projects/components/AppTypeSelection.tsx` - Type picker
4. `dash/src/features/projects/components/WebAppForm.tsx` - Web form
5. `dash/src/features/projects/components/ServiceForm.tsx` - Service form
6. `dash/src/features/projects/components/DatabaseForm.tsx` - Database form
7. `dash/src/features/projects/components/CreateAppModal.tsx` - Modal wrapper
8. `dash/src/features/projects/ProjectPage.tsx` - Uses new modal
9. `dash/src/features/applications/AppPage.tsx` - Hide git tab
10. `dash/src/components/applications/app-info.tsx` - Conditional rendering
11. `dash/src/components/applications/app-settings.tsx` - Disable fields
12. `dash/src/components/deployments/deployment-list.tsx` - Show version

---

## Common Issues & Solutions

### Issue: Database exposed on host port
**Symptom:** `docker ps` shows `0.0.0.0:6379->6379/tcp`  
**Cause:** Wrong network configuration in `build.go`  
**Solution:** Fixed in `RunContainer()` - database apps use `--network traefik-net` without `-p`

### Issue: NULL git repository causes SQL errors
**Symptom:** "sql: Scan error on column index 5"  
**Cause:** Attempting to scan NULL into string  
**Solution:** Use `sql.NullString` in `GetAppRepoInfo()`

### Issue: Git operations fail for database apps
**Symptom:** "failed to clone repository"  
**Cause:** Database apps don't have git repos  
**Solution:** Skip git operations in `handleWork.go` and `AddDeployHandler.go`

### Issue: Wrong commit hash for database deployments
**Symptom:** Deployment shows empty commit  
**Cause:** No git repo to fetch commit from  
**Solution:** Use template version as commit hash in `AddDeployHandler.go`

---

## Future Enhancements

1. **Volume Management:** Persistent storage for databases
2. **Backup System:** Automated database backups
3. **Admin UIs:** Deploy phpMyAdmin, pgAdmin alongside databases
4. **Health Checks:** Template-specific health check commands
5. **Multi-Container Templates:** Redis + Redis Commander
6. **Custom Templates:** User-defined service templates
7. **Resource Monitoring:** Real-time CPU/memory usage
8. **Auto-Scaling:** Scale containers based on metrics

---

## Deployment Flow Diagram

```
User Creates App
       |
       v
  App Type?
       |
   +---+---+
   |   |   |
  Web Svc DB
   |   |   |
   v   v   v
[Git][Git][Template]
   |   |   |
   v   v   v
[Build][Build][Pull]
   |   |   |
   +---+---+
       |
       v
   Run Container
       |
   +---+---+
   |   |   |
  Web Svc DB
   |   |   |
   v   v   v
[Traefik][Internal][Internal]
[or -p]
```

---

## Database Connection Examples

### Web App Connecting to PostgreSQL

**Environment Variables in Web App:**
```
DATABASE_URL=postgresql://postgres:securepass123@app-54321:5432/mydb
```

**Connection String Format:**
```
protocol://user:password@container_name:port/database
```

**Key Point:** Use container name as hostname (Docker DNS resolution on traefik-net)

### Service App Connecting to Redis

**Environment Variables in Service:**
```
REDIS_HOST=app-67890
REDIS_PORT=6379
REDIS_PASSWORD=redispass456
```

**Node.js Example:**
```javascript
const redis = new Redis({
  host: process.env.REDIS_HOST,
  port: process.env.REDIS_PORT,
  password: process.env.REDIS_PASSWORD
});
```

---

## Conclusion

The Mist app type system provides a clean separation between different application patterns while maintaining a consistent user experience. Database and service apps are properly isolated on internal networks, while web apps can be exposed via Traefik or direct port mapping based on configuration.

All three app types coexist in the same infrastructure, can communicate internally, and have appropriate UI presentations based on their characteristics.
