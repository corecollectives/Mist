# Testing Guide: App Type System

## Quick Verification Commands

### Check Docker Containers and Port Mappings
```bash
docker ps -a --format "table {{.Names}}\t{{.Image}}\t{{.Ports}}\t{{.Status}}"
```

**Expected Results:**

| Container | Type | Ports | Network |
|-----------|------|-------|---------|
| app-xxx (web with domain) | myapp:v1 | *empty* or 80/tcp | traefik-net |
| app-xxx (web no domain) | myapp:v1 | 0.0.0.0:3000->3000/tcp | default |
| app-xxx (service) | mybot:v1 | *empty* | traefik-net |
| app-xxx (database) | postgres:16 | 5432/tcp (internal only) | traefik-net |

### Inspect Container Network
```bash
docker inspect app-xxx | grep -A 10 "Networks"
```

### Check if Port is Exposed to Host
```bash
# This should show NOTHING for database/service apps
netstat -tlnp | grep <port>
# OR
ss -tlnp | grep <port>
```

---

## Test Scenarios

### Scenario 1: Create PostgreSQL Database

**Steps:**
1. Navigate to project page
2. Click "Create Application"
3. Select "Database / Pre-made"
4. Select "PostgreSQL"
5. Enter name: "my-postgres"
6. Click "Generate" for password
7. Note the password (e.g., `Xy7!mK9@pL3#`)
8. Click "Deploy PostgreSQL"

**Verify Backend:**
```bash
# Check database
sqlite3 mist.db "SELECT id, name, app_type, template_name, port FROM apps WHERE name='my-postgres';"

# Check env vars
sqlite3 mist.db "SELECT key, value FROM envs WHERE app_id=(SELECT id FROM apps WHERE name='my-postgres');"
```

**Expected:**
- `app_type` = 'database'
- `template_name` = 'postgres'
- `port` = 5432
- Env var: `POSTGRES_PASSWORD` = (generated password)

**Deploy and Verify:**
```bash
# Deploy the app via UI, then check:
docker ps | grep my-postgres

# Should show:
# app-xxxxx   postgres:16-alpine   5432/tcp   Up X minutes

# Should NOT show:
# 0.0.0.0:5432->5432/tcp
```

**Verify UI:**
- [ ] Git tab is hidden (only 5 tabs visible)
- [ ] App info shows "Service Template: postgres"
- [ ] No git repository/branch/commit shown
- [ ] Deployment shows "Version: 16-alpine"

---

### Scenario 2: Create Web App with Domain

**Steps:**
1. Create web app: "my-frontend"
2. Set port: 3000
3. Connect git repository
4. Add domain: "myapp.local" (via Settings > Domains tab)
5. Deploy

**Verify:**
```bash
docker ps | grep my-frontend

# Should show:
# app-xxxxx   my-frontend:abc1234   3000/tcp   Up X minutes

# Should have Traefik labels
docker inspect app-xxxxx | grep traefik.enable
docker inspect app-xxxxx | grep "traefik.http.routers"

# Should NOT have port mapping
docker ps | grep my-frontend | grep "0.0.0.0"
# (should return nothing)
```

**Access Test:**
```bash
curl -H "Host: myapp.local" http://localhost
# Should return your app's response
```

---

### Scenario 3: Create Web App WITHOUT Domain

**Steps:**
1. Create web app: "my-api"
2. Set port: 8080
3. Connect git repository
4. Deploy (don't add any domains)

**Verify:**
```bash
docker ps | grep my-api

# Should show:
# app-xxxxx   my-api:abc1234   0.0.0.0:8080->8080/tcp   Up X minutes

# Should NOT have Traefik labels
docker inspect app-xxxxx | grep traefik.enable
# (should return nothing)
```

**Access Test:**
```bash
curl http://localhost:8080
# Should return your app's response
```

---

### Scenario 4: Create Service App (Background Worker)

**Steps:**
1. Create application
2. Select "Background Service"
3. Enter name: "my-worker"
4. Description: "Queue processor"
5. Connect git repository
6. Deploy

**Verify:**
```bash
docker ps | grep my-worker

# Should show:
# app-xxxxx   my-worker:abc1234   (no ports listed)   Up X minutes

# Check network
docker inspect app-xxxxx | grep -A 5 "Networks"
# Should show: "traefik-net"

# Verify NO external port
netstat -tlnp | grep 3000
# Should return nothing (or not show this container)
```

**Verify UI:**
- [ ] Git tab is visible
- [ ] Shows git repository and commit info
- [ ] No external access possible
- [ ] Logs still visible

---

### Scenario 5: Cross-App Communication

**Setup:**
1. Create PostgreSQL database: "app-db-1"
2. Create web app: "app-web-1"
3. Add env var to web app:
   ```
   DATABASE_URL=postgresql://postgres:PASSWORD@app-db-1:5432/mydb
   ```

**Test:**
```bash
# From web app container, try to connect to database
docker exec app-web-1 ping app-db-1
# Should succeed

# Try to connect to PostgreSQL port
docker exec app-web-1 nc -zv app-db-1 5432
# Should show: "app-db-1 (172.x.x.x:5432) open"
```

**From Host Machine:**
```bash
# Try to access PostgreSQL directly
nc -zv localhost 5432
# Should FAIL (connection refused)

# Try to access via container name
psql -h app-db-1 -U postgres
# Should FAIL (cannot resolve host)
```

**Expected:** Internal communication works, external access blocked.

---

### Scenario 6: Create Redis with Password

**Steps:**
1. Select "Database / Pre-made"
2. Choose "Redis"
3. Name: "my-redis"
4. Generate password: `Kp9@xL2#mN8!`
5. Deploy

**Verify Env Vars:**
```bash
sqlite3 mist.db "SELECT key, value FROM envs WHERE app_id=(SELECT id FROM apps WHERE name='my-redis');"
```

**Expected:**
- `REDIS_PASSWORD` = (generated password)

**Test Connection from Another App:**
```bash
# Add to web app env vars:
REDIS_URL=redis://:Kp9@xL2#mN8!@app-xxxxx:6379

# From web app:
docker exec app-web-1 redis-cli -h app-xxxxx -a 'Kp9@xL2#mN8!' PING
# Should return: PONG
```

---

## Verification Checklist

### Database Apps
- [ ] No external port mapping (check `docker ps`)
- [ ] Connected to `traefik-net`
- [ ] Git tab hidden in UI
- [ ] Shows template name in app info
- [ ] Environment variables created (password)
- [ ] Deployment shows "Version: X" not commit hash
- [ ] Cannot access from host machine
- [ ] Can access from other containers on traefik-net

### Service Apps
- [ ] No external port mapping
- [ ] Connected to `traefik-net`
- [ ] Git tab visible
- [ ] Shows git info in UI
- [ ] No Traefik labels
- [ ] Deployment shows commit hash
- [ ] Cannot access from host machine
- [ ] Can communicate with other containers

### Web Apps (with domains)
- [ ] No direct port mapping (no `-p` flag)
- [ ] Connected to `traefik-net`
- [ ] Has Traefik labels
- [ ] Accessible via domain name
- [ ] Git tab visible
- [ ] Shows git info in UI

### Web Apps (without domains)
- [ ] Has external port mapping (`-p HOST:CONTAINER`)
- [ ] May or may not be on traefik-net
- [ ] No Traefik labels
- [ ] Accessible via `http://localhost:PORT`
- [ ] Git tab visible
- [ ] Shows git info in UI

---

## Common Test Commands

### Check All Running Containers
```bash
docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Ports}}"
```

### Check Container Network
```bash
docker network inspect traefik-net
```

### Check Container Environment Variables
```bash
docker exec app-xxxxx env | grep -E "POSTGRES|REDIS|MONGO|MYSQL"
```

### Check Container Logs
```bash
docker logs app-xxxxx
```

### Test Internal Network Connectivity
```bash
# From one container to another
docker exec app-web-1 ping app-db-1
docker exec app-web-1 nc -zv app-db-1 5432
```

### Test External Access (Should Fail for DB/Service)
```bash
curl http://localhost:5432  # PostgreSQL - should fail
curl http://localhost:6379  # Redis - should fail
nc -zv localhost 5432       # Should fail for databases
```

---

## Troubleshooting

### Issue: Database accessible from host
```bash
# Check if port is exposed
docker ps | grep app-xxxxx

# If you see 0.0.0.0:6379->6379/tcp, it's WRONG
# Redeploy the app after code fix
```

### Issue: Web app not accessible
```bash
# Check if domain is configured
sqlite3 mist.db "SELECT * FROM domains WHERE app_id=X;"

# Check Traefik labels
docker inspect app-xxxxx | grep traefik

# Check if port is mapped (for apps without domains)
docker ps | grep app-xxxxx
```

### Issue: Apps can't communicate
```bash
# Check if both are on same network
docker network inspect traefik-net

# Verify DNS resolution
docker exec app-web-1 nslookup app-db-1
docker exec app-web-1 ping app-db-1
```

### Issue: Wrong app type after creation
```bash
# Check database
sqlite3 mist.db "SELECT name, app_type FROM apps;"

# Can't change app type after creation - delete and recreate
```

---

## Expected Behavior Summary

| Feature | Web | Service | Database |
|---------|-----|---------|----------|
| External Port | Yes (if no domains) | No | No |
| Traefik Routing | Yes (if domains) | No | No |
| Git Operations | Yes | Yes | No |
| Git Tab in UI | Yes | Yes | No |
| Template | No | No | Yes |
| Env Vars Auto-created | No | No | Yes (password) |
| Internal Network | Yes | Yes | Yes |
| Build from Git | Yes | Yes | No |
| Pull Docker Image | No | No | Yes |

---

## Quick Test Script

```bash
#!/bin/bash

echo "=== Mist App Type System Test ==="

echo -e "\n1. Checking running containers..."
docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Ports}}"

echo -e "\n2. Checking traefik-net members..."
docker network inspect traefik-net | grep -E "Name|IPv4"

echo -e "\n3. Checking database apps (should have NO host port mapping)..."
docker ps | grep -E "postgres|redis|mysql|mongo" | grep -E "0.0.0.0|:::"
if [ $? -eq 0 ]; then
    echo "❌ ERROR: Database apps have external port mappings!"
else
    echo "✅ PASS: Database apps properly isolated"
fi

echo -e "\n4. Checking apps table..."
sqlite3 mist.db "SELECT id, name, app_type, template_name, port FROM apps;"

echo -e "\n5. Checking environment variables for database apps..."
sqlite3 mist.db "
SELECT a.name, e.key, e.value 
FROM apps a 
JOIN envs e ON a.id = e.app_id 
WHERE a.app_type = 'database';
"

echo -e "\n=== Test Complete ==="
```

Save as `test-app-types.sh` and run:
```bash
chmod +x test-app-types.sh
./test-app-types.sh
```
