# Smart Deployment Monitoring - REST vs WebSocket

## Overview

The deployment monitoring system intelligently chooses between **REST API** (for completed deployments) and **WebSocket** (for live deployments) to optimize performance and resource usage.

---

## 🎯 Why This Matters

### Before (WebSocket Only)
```
User views completed deployment
    ↓
Opens WebSocket connection
    ↓
Server starts 2 goroutines:
  - Status watcher (polls DB every 500ms) ❌ Unnecessary
  - Log watcher (tails file) ❌ File is complete
    ↓
Sends all logs line by line over WebSocket ❌ Slow
    ↓
Connection stays open briefly ❌ Resource waste
```

**Problems:**
- ❌ Slow initial load (logs sent one by one)
- ❌ Wasted server resources (unnecessary goroutines)
- ❌ Unnecessary database polling
- ❌ Overhead of WebSocket handshake

### After (Smart Detection)
```
User views completed deployment
    ↓
Frontend fetches deployment status ✅ One HTTP request
    ↓
Backend checks: status = "success" or "failed"
    ↓
Returns ALL logs + status in ONE response ✅ Fast
    ↓
Frontend displays instantly ✅ No waiting
    ↓
No WebSocket connection needed ✅ Clean
```

**Benefits:**
- ✅ **Instant load** - All data in one request
- ✅ **No server overhead** - No goroutines spawned
- ✅ **Better UX** - Logs appear immediately
- ✅ **Resource efficient** - Only use WebSocket when needed

---

## 🔄 Flow Diagram

```
┌─────────────────────────────────────────────────────────┐
│ User clicks "View Logs" on deployment                   │
└────────────────────────┬────────────────────────────────┘
                         │
                         ↓
        ┌────────────────────────────────┐
        │ Frontend: DeploymentMonitor    │
        │ opens with deploymentId        │
        └────────────────┬───────────────┘
                         │
                         ↓
        ┌────────────────────────────────────────┐
        │ useDeploymentMonitor hook executes     │
        │ Step 1: Try REST API first            │
        └────────────────┬───────────────────────┘
                         │
                         ↓
    ╔════════════════════════════════════════════════╗
    ║ GET /api/deployments/logs?id=456               ║
    ║ (with authentication)                          ║
    ╚════════════════════════════════════════════════╝
                         │
         ┌───────────────┴───────────────┐
         │                               │
         ↓                               ↓
    ┌─────────┐                   ┌──────────┐
    │ SUCCESS │                   │ 400 BAD  │
    │ 200 OK  │                   │ REQUEST  │
    └────┬────┘                   └────┬─────┘
         │                             │
         │                             ↓
         │              ┌──────────────────────────────┐
         │              │ "deployment is still in      │
         │              │  progress, use WebSocket"    │
         │              └──────────────┬───────────────┘
         │                             │
         │                             ↓
         │              ┌──────────────────────────────┐
         │              │ Frontend: setIsLive(true)    │
         │              └──────────────┬───────────────┘
         │                             │
         │                             ↓
         │              ╔══════════════════════════════╗
         │              ║ WebSocket Connection         ║
         │              ║ ws://host/api/deployments/  ║
         │              ║ logs/stream?id=456          ║
         │              ╚══════════════════════════════╝
         │                             │
         │                             ↓
         │              ┌──────────────────────────────┐
         │              │ Live streaming:              │
         │              │ - Status updates (500ms poll)│
         │              │ - Log lines as they appear   │
         │              └──────────────────────────────┘
         │
         ↓
┌─────────────────────────────────────────┐
│ Response (Completed Deployment):        │
│ {                                       │
│   "success": true,                      │
│   "data": {                             │
│     "deployment": {                     │
│       "id": 456,                        │
│       "status": "success",              │
│       "stage": "success",               │
│       "progress": 100,                  │
│       "duration": 125,                  │
│       "error_message": null,            │
│       ...                               │
│     },                                  │
│     "logs": "Step 1/5...\nStep 2/5..." │
│   }                                     │
│ }                                       │
└────────────────┬────────────────────────┘
                 │
                 ↓
┌─────────────────────────────────────────┐
│ Frontend:                               │
│ - Parse logs (split by \n)             │
│ - Display all logs instantly           │
│ - Show final status badge               │
│ - Show "Completed" indicator            │
│ - No WebSocket connection               │
└─────────────────────────────────────────┘
```

---

## 📋 Implementation Details

### Backend: REST Endpoint

**File:** `server/api/handlers/deployments/getCompletedLogs.go`

```go
func GetCompletedDeploymentLogsHandler(w http.ResponseWriter, r *http.Request) {
    depId := getDeploymentIdFromQuery(r)
    
    // Get deployment from database
    dep, err := models.GetDeploymentByID(depId)
    
    // Check if deployment is completed
    if dep.Status != "success" && dep.Status != "failed" {
        // Return 400 - tells frontend to use WebSocket instead
        handlers.SendResponse(w, http.StatusBadRequest, false, nil, 
            "deployment is still in progress, use WebSocket endpoint", "")
        return
    }
    
    // Read entire log file
    logPath := docker.GetLogsPath(dep.CommitHash, depId)
    logContent := readEntireFile(logPath)
    
    // Return everything in one response
    response := GetDeploymentLogsResponse{
        Deployment: dep,    // Full deployment object
        Logs:       logContent, // Entire log file as string
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
```

**Key Points:**
- Returns **400 Bad Request** if deployment is in progress
- Reads **entire log file** at once (no streaming needed)
- Returns **complete deployment object** with all metadata
- Single HTTP request/response cycle

---

### Frontend: Smart Detection

**File:** `dash/src/features/applications/hooks/useDeploymentMonitor.ts`

```typescript
export const useDeploymentMonitor = ({ deploymentId, enabled }) => {
  const [isLive, setIsLive] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Step 1: Try REST API first
  const fetchCompletedDeployment = async () => {
    const response = await fetch(`/api/deployments/logs?id=${deploymentId}`);
    
    if (response.status === 400) {
      // Deployment is in progress, use WebSocket
      setIsLive(true);
      return;
    }
    
    if (response.ok) {
      const result = await response.json();
      const deployment = result.data.deployment;
      const logsContent = result.data.logs;
      
      // Set all logs at once (instant display)
      setLogs(logsContent.split('\n').filter(line => line.length > 0));
      
      // Set final status
      setStatus({
        status: deployment.status,
        stage: deployment.stage,
        progress: deployment.progress,
        duration: deployment.duration,
        ...
      });
      
      setIsLoading(false);
    }
  };

  // Step 2: Connect WebSocket if deployment is live
  const connectWebSocket = () => {
    if (!isLive) return;
    
    const ws = new WebSocket(`/api/deployments/logs/stream?id=${deploymentId}`);
    
    ws.onmessage = (event) => {
      const deploymentEvent = JSON.parse(event.data);
      
      if (deploymentEvent.type === 'log') {
        setLogs(prev => [...prev, deploymentEvent.data.line]);
      } else if (deploymentEvent.type === 'status') {
        setStatus(deploymentEvent.data);
      }
    };
  };

  useEffect(() => {
    if (enabled) {
      fetchCompletedDeployment(); // Try REST first
    }
  }, [enabled]);

  useEffect(() => {
    if (isLive && enabled) {
      connectWebSocket(); // Use WebSocket if needed
    }
  }, [isLive, enabled]);
};
```

**Logic Flow:**
1. **Always try REST API first** when monitor opens
2. If **400 response** → Deployment is live → Use WebSocket
3. If **200 response** → Deployment is complete → Display immediately
4. WebSocket only connects if `isLive === true`

---

## 🎨 UI Indicators

The frontend shows different indicators based on the mode:

### Completed Deployment (REST)
```
┌─────────────────────────────────────┐
│ Deployment Monitor        🔵 Completed │
│ #456                                 │
├─────────────────────────────────────┤
│ ✅ Success                          │
│ █████████████████████ 100%          │
├─────────────────────────────────────┤
│ ✅ Deployment Successful!           │
│    Completed in 125s                │
├─────────────────────────────────────┤
│ [ALL LOGS DISPLAYED INSTANTLY]      │
│ Step 1/5 : FROM node:18             │
│ Step 2/5 : COPY . /app              │
│ ...                                 │
└─────────────────────────────────────┘
```

### Live Deployment (WebSocket)
```
┌─────────────────────────────────────┐
│ Deployment Monitor        🟢 Live   │
│ #456                                │
├─────────────────────────────────────┤
│ 🔵 Building                         │
│ █████████████░░░░░░░░ 50%          │
├─────────────────────────────────────┤
│ [LOGS STREAMING IN REAL-TIME]      │
│ Step 1/5 : FROM node:18             │
│ Step 2/5 : COPY . /app              │
│ ... [streaming] ...                 │
└─────────────────────────────────────┘
```

---

## 📊 Performance Comparison

| Metric | REST (Completed) | WebSocket (Live) |
|--------|------------------|------------------|
| **Initial Load Time** | ~100ms | ~2-5s (wait for logs) |
| **Server Resources** | Minimal | 2 goroutines per connection |
| **Database Queries** | 1 query | Continuous (every 500ms) |
| **Memory Usage** | Low | Higher (connection state) |
| **Network Overhead** | 1 HTTP request | WebSocket handshake + messages |
| **User Experience** | Instant | Progressive |

---

## 🔧 API Endpoints

### REST Endpoint (Completed Deployments)
```
GET /api/deployments/logs?id={deploymentId}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "data": {
    "deployment": {
      "id": 456,
      "app_id": 123,
      "status": "success",
      "stage": "success",
      "progress": 100,
      "duration": 125,
      "error_message": null,
      "commit_hash": "abc123",
      "created_at": "2025-11-14T10:00:00Z",
      "finished_at": "2025-11-14T10:02:05Z"
    },
    "logs": "Step 1/5 : FROM node:18\nStep 2/5 : COPY . /app\n..."
  },
  "message": "Deployment logs retrieved successfully"
}
```

**Response (In Progress - 400 Bad Request):**
```json
{
  "success": false,
  "message": "deployment is still in progress, use WebSocket endpoint",
  "error": ""
}
```

---

### WebSocket Endpoint (Live Deployments)
```
ws://host/api/deployments/logs/stream?id={deploymentId}
```

**Messages (JSON):**
```json
// Status Update
{
  "type": "status",
  "timestamp": "2025-11-14T10:00:30Z",
  "data": {
    "deployment_id": 456,
    "status": "building",
    "stage": "building",
    "progress": 50,
    "message": "Building Docker image"
  }
}

// Log Line
{
  "type": "log",
  "timestamp": "2025-11-14T10:00:31Z",
  "data": {
    "line": "Step 3/5 : RUN npm install\n",
    "timestamp": "2025-11-14T10:00:31Z"
  }
}
```

---

## ✅ Testing Scenarios

### Test 1: View Completed Success Deployment
```bash
# 1. Create and complete a deployment
curl -X POST /api/deployments -d '{"appId": 123}'
# Wait for it to complete...

# 2. Open deployment monitor
# Expected: Loads instantly via REST
# Expected: Shows "Completed" indicator
# Expected: All logs visible immediately
# Expected: No WebSocket connection
```

### Test 2: View Completed Failed Deployment
```bash
# 1. Create a deployment that will fail
# 2. Open deployment monitor
# Expected: Loads instantly via REST
# Expected: Shows error banner
# Expected: Shows "Completed" indicator
# Expected: Error message displayed
```

### Test 3: View In-Progress Deployment
```bash
# 1. Create a deployment
# 2. Immediately open monitor
# Expected: REST returns 400
# Expected: Switches to WebSocket
# Expected: Shows "Live" indicator
# Expected: Logs stream in real-time
# Expected: Status updates every 500ms
```

### Test 4: Open Monitor Mid-Deployment
```bash
# 1. Start deployment
# 2. Wait 30 seconds
# 3. Open monitor
# Expected: REST returns 400
# Expected: WebSocket connects
# Expected: Receives existing logs first
# Expected: Then streams new logs
```

---

## 🎯 Key Benefits Summary

### For Users
- ✅ **Faster loading** - Completed deployments load instantly
- ✅ **Better feedback** - Clear "Live" vs "Completed" indicators
- ✅ **Reliable** - Both methods are robust and tested

### For Developers
- ✅ **Clean architecture** - Smart detection at hook level
- ✅ **Easy to maintain** - Clear separation of concerns
- ✅ **Well-typed** - TypeScript types for both flows

### For Infrastructure
- ✅ **Resource efficient** - No unnecessary WebSocket connections
- ✅ **Scalable** - Reduced goroutine count
- ✅ **Cost effective** - Less database polling

---

## 🚀 Future Enhancements

1. **Caching** - Cache completed deployment logs in memory
2. **Compression** - Gzip large log responses
3. **Pagination** - Paginate very large log files
4. **Search** - Add log search capability for completed deployments
5. **Download** - Add "Download Logs" button for completed deployments

---

## Conclusion

This smart detection system provides the **best of both worlds**:
- **Speed** for completed deployments (REST)
- **Real-time** updates for live deployments (WebSocket)

Users get instant feedback when viewing completed deployments, while live deployments still provide the real-time streaming experience they expect! 🎉
