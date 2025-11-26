# 📡 Real-Time Deployment Monitoring - Complete Data Flow

## Overview

The system uses **WebSockets** with **two concurrent goroutines** to stream both deployment status updates and build logs in real-time to the frontend.

---

## 🔄 Complete Architecture Flow

```
Frontend                    Backend                     Database/Filesystem
   │                           │                              │
   │  1. User clicks "Deploy"  │                              │
   ├──────────────────────────►│                              │
   │  POST /api/deployments    │                              │
   │  { appId: 123 }           │                              │
   │                           │                              │
   │                           │  2. Create Deployment        │
   │                           ├─────────────────────────────►│
   │                           │  INSERT INTO deployments     │
   │                           │  (status='pending')          │
   │                           │                              │
   │                           │  3. Add to Queue             │
   │  4. Return deployment ID  │  queue.AddJob(depID)         │
   │◄──────────────────────────┤                              │
   │  { id: 456, status: ... } │                              │
   │                           │                              │
   │  5. Open WebSocket        │                              │
   ├──────────────────────────►│                              │
   │  ws://host/api/           │                              │
   │  deployments/logs?id=456  │                              │
   │                           │                              │
   │                           │  6. Upgrade to WebSocket     │
   │◄──────────────────────────┤     (logsHandler.go:37)     │
   │  [WebSocket Connected]    │                              │
   │                           │                              │
   │                           │  ┌──────────────────────────┐│
   │                           │  │ 7. Start TWO Goroutines: ││
   │                           │  └──────────────────────────┘│
   │                           │                              │
   │                           │  ┌─ Goroutine 1 ────────────┐│
   │                           │  │ WatchDeploymentStatus()  ││
   │                           │  │ (statusWatcher.go)       ││
   │                           │  │                          ││
   │                           │  │ • Polls DB every 500ms   ││
   │                           │  │ • Checks for changes in: ││
   │                           │  │   - status               ││
   │                           │  │   - stage                ││
   │                           │  │   - progress             ││
   │                           │  │   - error_message        ││
   │                           │  │                          ││
   │                           │  │ • Sends to events channel││
   │                           │  └──────────────────────────┘│
   │                           │                              │
   │                           │  ┌─ Goroutine 2 ────────────┐│
   │                           │  │ WatcherLogs()            ││
   │                           │  │ (logWatcher.go)          ││
   │                           │  │                          ││
   │                           │  │ • Waits for log file     ││
   │                           │  │ • Tails log file         ││
   │                           │  │ • Reads new lines        ││
   │                           │  │   every 500ms            ││
   │                           │  │                          ││
   │                           │  │ • Sends to events channel││
   │                           │  └──────────────────────────┘│
   │                           │                              │
   │                           │  ┌──────────────────────────┐│
   │                           │  │ 8. Worker processes job  ││
   │                           │  └──────────────────────────┘│
   │                           │                              │
   │                           │  handleWork(456)             │
   │                           │    │                         │
   │                           │    ├─ Mark started          │
   │                           │    │  UPDATE deployments    │
   │                           │    │  SET started_at=NOW()  │
   │                           │    │                         │
   │                           │    ├─ Clone repo            │
   │                           │    │  stage='cloning'       │
   │                           │    │  progress=20           │
   │                           │    │  ────────────────────►DB│
   │  9. Status Event          │    │                         │
   │◄──────────────────────────┤────┘ (detected by          │
   │  {                        │      Goroutine 1)           │
   │    type: "status",        │                             │
   │    data: {                │                             │
   │      status: "cloning",   │                             │
   │      stage: "cloning",    │                             │
   │      progress: 20         │                             │
   │    }                      │                             │
   │  }                        │                             │
   │                           │                             │
   │                           │    ├─ Build image           │
   │                           │    │  stage='building'      │
   │                           │    │  progress=50           │
   │                           │    │  ────────────────────►DB│
   │  10. Status Event         │    │                         │
   │◄──────────────────────────┤────┘                         │
   │  {                        │                             │
   │    type: "status",        │                             │
   │    data: {                │                             │
   │      status: "building",  │      Docker writes to:     │
   │      stage: "building",   │      /logs/abc123456_build_│
   │      progress: 50         │                        logs │
   │    }                      │                             │
   │  }                        │                             │
   │                           │                             │
   │  11. Log Events (stream)  │      (detected by          │
   │◄──────────────────────────┤────  Goroutine 2)          │
   │  {                        │                             │
   │    type: "log",           │                             │
   │    data: {                │                             │
   │      line: "Step 1/5..."  │                             │
   │    }                      │                             │
   │  }                        │                             │
   │                           │                             │
   │  {                        │                             │
   │    type: "log",           │                             │
   │    data: {                │                             │
   │      line: "Step 2/5..."  │                             │
   │    }                      │                             │
   │  }                        │                             │
   │                           │                             │
   │                           │    ├─ Deploy container     │
   │                           │    │  stage='deploying'    │
   │                           │    │  progress=80          │
   │                           │    │  ────────────────────►DB│
   │  12. Status Event         │    │                        │
   │◄──────────────────────────┤────┘                        │
   │                           │                             │
   │                           │    ├─ Success!             │
   │                           │    │  stage='success'      │
   │                           │    │  progress=100         │
   │                           │    │  finished_at=NOW()    │
   │                           │    │  duration=125         │
   │                           │    │  ────────────────────►DB│
   │  13. Final Status Event   │    │                        │
   │◄──────────────────────────┤────┘ (Goroutine 1 detects │
   │  {                        │      and then EXITS)       │
   │    type: "status",        │                            │
   │    data: {                │                            │
   │      status: "success",   │                            │
   │      progress: 100,       │                            │
   │      duration: 125        │                            │
   │    }                      │                            │
   │  }                        │                            │
   │                           │                            │
   │  14. WebSocket closes     │  Both goroutines exit     │
   │◄──────────────────────────┤  events channel closes    │
   │                           │  Connection closed        │
```

---

## 🔍 Detailed Code Walkthrough

### 1. **Frontend Initiates WebSocket Connection**

```typescript
// dash/src/features/applications/hooks/useDeploymentMonitor.ts

const wsUrl = `ws://localhost:8080/api/deployments/logs?id=${deploymentId}`;
const ws = new WebSocket(wsUrl);

ws.onmessage = (event) => {
  const deploymentEvent: DeploymentEvent = JSON.parse(event.data);
  
  if (deploymentEvent.type === 'status') {
    setStatus(deploymentEvent.data); // Update status bar
  } else if (deploymentEvent.type === 'log') {
    setLogs(prev => [...prev, deploymentEvent.data.line]); // Append log
  }
};
```

---

### 2. **Backend Upgrades HTTP to WebSocket**

```go
// server/api/handlers/deployments/logsHandler.go:24

func LogsHandler(w http.ResponseWriter, r *http.Request) {
    depId := getDeploymentID(r)
    
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()
    
    // Create a buffered channel for events
    events := make(chan websockets.DeploymentEvent, 100)
    
    // Start TWO concurrent goroutines...
}
```

---

### 3. **Goroutine 1: Status Watcher (Polls Database)**

```go
// server/websockets/statusWatcher.go:31

func WatchDeploymentStatus(ctx context.Context, depID int64, events chan<- DeploymentEvent) {
    ticker := time.NewTicker(500 * time.Millisecond) // Poll every 500ms
    defer ticker.Stop()
    
    var lastStatus, lastStage string
    var lastProgress int
    
    for {
        select {
        case <-ticker.C:
            // Query database for current deployment state
            dep, err := models.GetDeploymentByID(depID)
            if err != nil {
                continue
            }
            
            // Check if anything changed
            if dep.Status != lastStatus || 
               dep.Stage != lastStage || 
               dep.Progress != lastProgress {
                
                // Update tracking variables
                lastStatus = dep.Status
                lastStage = dep.Stage
                lastProgress = dep.Progress
                
                // Send status event to channel
                events <- DeploymentEvent{
                    Type: "status",
                    Data: StatusUpdate{
                        Status:   dep.Status,
                        Stage:    dep.Stage,
                        Progress: dep.Progress,
                        Message:  utils.GetStageMessage(dep.Stage),
                    },
                }
            }
            
            // Exit if deployment finished
            if dep.Status == "success" || dep.Status == "failed" {
                return
            }
        }
    }
}
```

**What it does:**
- Polls database every **500ms**
- Compares current state with last known state
- Only sends event if something changed (status, stage, or progress)
- Automatically exits when deployment finishes

---

### 4. **Goroutine 2: Log Watcher (Tails Log File)**

```go
// server/websockets/logWatcher.go:11

func WatcherLogs(ctx context.Context, filePath string, send chan<- string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()
    
    reader := bufio.NewReader(file)
    
    // PHASE 1: Read existing logs (catch-up)
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break // Reached end of file
        }
        send <- line // Send existing line
    }
    
    // PHASE 2: Tail mode (follow new logs)
    for {
        select {
        case <-ctx.Done():
            return nil
        default:
            line, err := reader.ReadString('\n')
            if err == io.EOF {
                time.Sleep(500 * time.Millisecond) // Wait for new data
                continue
            }
            send <- line // Send new line
        }
    }
}
```

**What it does:**
- Opens the log file (e.g., `/logs/abc123456_build_logs`)
- **Phase 1**: Reads all existing logs (if user opens monitor mid-deployment)
- **Phase 2**: Tails the file for new logs (like `tail -f`)
- Checks every **500ms** for new lines
- Exits when context is cancelled

---

### 5. **Main Handler Merges Both Channels**

```go
// server/api/handlers/deployments/logsHandler.go:48-96

func LogsHandler(w http.ResponseWriter, r *http.Request) {
    // ... setup code ...
    
    events := make(chan websockets.DeploymentEvent, 100)
    
    // Start status watcher goroutine
    go websockets.WatchDeploymentStatus(ctx, depId, events)
    
    // Start log watcher goroutine
    go func() {
        // Wait for log file to exist (max 10 seconds)
        for i := 0; i < 20; i++ {
            if _, err := os.Stat(logPath); err == nil {
                break
            }
            time.Sleep(500 * time.Millisecond)
        }
        
        // Start tailing logs
        send := make(chan string)
        go websockets.WatcherLogs(ctx, logPath, send)
        
        // Convert log strings to DeploymentEvents
        for line := range send {
            events <- DeploymentEvent{
                Type: "log",
                Data: LogUpdate{
                    Line:      line,
                    Timestamp: time.Now(),
                },
            }
        }
    }()
    
    // Main loop: send all events to WebSocket
    for event := range events {
        msg, _ := json.Marshal(event)
        conn.WriteMessage(websocket.TextMessage, msg)
    }
}
```

**What it does:**
- Creates a **unified events channel**
- Both goroutines send to this channel
- Main loop reads from channel and sends to WebSocket
- Frontend receives both status and logs interleaved

---

## 📊 Event Flow Timeline

```
Time    Worker Thread           Status Watcher          Log Watcher         Frontend
────────────────────────────────────────────────────────────────────────────────────
0ms     Create deployment      [Polling DB...]         [Waiting...]        Deploy clicked
        status='pending'
        
500ms   Add to queue           Detects: pending        [Still waiting]     Connected!
        ───────────────►        ──────────────►                             Status: pending
                               Send "pending"
                               
2000ms  Worker starts          Detects: cloning        Log file created    Status: cloning
        Clone repo...          ──────────────►         Start tailing       
        stage='cloning'        Send "cloning"          ──────────►         Logs start...
        progress=20                                    Send log lines
        
5000ms  Cloning done           Detects: building       ──────────►         Status: building
        Build image...         ──────────────►         "Step 1/5..."       Logs: Step 1...
        stage='building'       Send "building"         "Step 2/5..."       Logs: Step 2...
        progress=50                                    ──────────►
                                                       
15000ms Image built            Detects: deploying      ──────────►         Status: deploying
        Deploy container...    ──────────────►         "Starting..."       Logs: Starting
        stage='deploying'      Send "deploying"
        progress=80
        
20000ms Success!               Detects: success        Log file EOF        Status: success ✅
        stage='success'        ──────────────►         Exit goroutine      Connection closed
        progress=100           Send "success"
        finished_at=NOW()      Exit goroutine
        duration=125s
```

---

## 🎯 Key Design Decisions

### Why Two Goroutines?

1. **Status Watcher** - Tracks deployment lifecycle
   - Database polling (500ms interval)
   - Catches state changes even if logs fail
   - Provides progress percentage

2. **Log Watcher** - Provides detailed output
   - File tailing (500ms check interval)
   - Real-time Docker build output
   - Survives status updates

### Why Polling Instead of Push?

**Status Watcher polls database because:**
- Simple implementation
- No need for pub/sub infrastructure
- 500ms latency is acceptable
- Worker updates DB, watcher reads it (loose coupling)

### Why File Tailing for Logs?

- Docker commands write to file
- File acts as permanent record
- Can be re-read if WebSocket disconnects
- Standard Unix approach (`tail -f`)

---

## 💡 Event Examples

### Status Event
```json
{
  "type": "status",
  "timestamp": "2025-11-14T21:30:45Z",
  "data": {
    "deployment_id": 456,
    "status": "building",
    "stage": "building",
    "progress": 50,
    "message": "Building Docker image",
    "error_message": ""
  }
}
```

### Log Event
```json
{
  "type": "log",
  "timestamp": "2025-11-14T21:30:45Z",
  "data": {
    "line": "Step 1/5 : FROM node:18\n",
    "timestamp": "2025-11-14T21:30:45Z"
  }
}
```

### Error Event (on failure)
```json
{
  "type": "status",
  "timestamp": "2025-11-14T21:30:50Z",
  "data": {
    "deployment_id": 456,
    "status": "failed",
    "stage": "failed",
    "progress": 0,
    "message": "Deployment failed",
    "error_message": "Docker build failed with exit code 1: syntax error in Dockerfile"
  }
}
```

---

## 🔄 Update Mechanism

### How Worker Updates Are Detected

```go
// Worker thread (queue/handleWork.go)
models.UpdateDeploymentStatus(id, "building", "building", 50, nil)
                                    ↓
// Updates database
UPDATE deployments 
SET status='building', stage='building', progress=50 
WHERE id=456
                                    ↓
// Status watcher polls (500ms later)
dep, _ := models.GetDeploymentByID(456)
if dep.Stage != lastStage {  // "building" != "cloning"
    events <- StatusUpdate{...}  // Send to WebSocket
}
```

### How Logs Are Written and Read

```go
// Worker thread (docker/build.go)
cmd := exec.Command("docker", "build", ...)
cmd.Stdout = logfile  // Redirect to /logs/abc123_build_logs
cmd.Stderr = logfile
cmd.Run()
                                    ↓
// Docker writes:
// "Step 1/5 : FROM node:18\n"
// "Step 2/5 : COPY . /app\n"
                                    ↓
// Log watcher tails file
line, err := reader.ReadString('\n')
send <- line  // Send to channel
                                    ↓
// Converted to event
events <- DeploymentEvent{Type: "log", Data: {...}}
                                    ↓
// Sent to WebSocket
conn.WriteMessage(msg)
```

---

## ✅ Summary

### Data Sources
1. **Database** → Status, stage, progress, errors
2. **Log File** → Build output, container logs

### Delivery Mechanism
1. **WebSocket** → Single persistent connection
2. **Two Goroutines** → Status polling + log tailing
3. **Unified Channel** → Merges both streams
4. **JSON Events** → Typed messages to frontend

### Frontend Handling
1. **useDeploymentMonitor hook** → Manages WebSocket
2. **Event routing** → Updates UI based on type
3. **Auto-reconnection** → Handles disconnects
4. **Real-time UI** → Instant feedback

This architecture provides **real-time updates** with **minimal overhead** and **resilient connections**! 🚀
