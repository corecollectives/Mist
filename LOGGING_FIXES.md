# Logging System Fixes

## Issues Fixed

### 1. Duplicate Log Lines
**Problem**: Logs were appearing twice in both live container logs and build logs.

**Root Causes**:
- **Container logs**: Docker outputs to both stdout and stderr simultaneously, causing duplicate lines when both streams were read separately
- **Build logs**: The log file watcher was reading the file twice - once initially and then again in the main loop

**Solutions**:
- **Container logs** (`server/websockets/containerLogs.go`):
  - Combined stdout and stderr into a single stream by redirecting stderr to stdout
  - Now using `cmd.Stderr = cmd.Stdout` to merge streams
  - Single scanner reads from combined output, eliminating duplicates

- **Build logs** (`server/websockets/logWatcher.go`):
  - Removed the initial file read loop
  - Single unified loop handles both existing and new content
  - Added proper line trimming to remove trailing `\n` and `\r` characters

### 2. Connection Timeouts in Build Logs
**Problem**: Build logs WebSocket connections randomly timing out and showing "connection timed out" errors.

**Root Cause**:
- No ping/pong mechanism to keep connections alive
- No read/write deadlines set
- No proper client disconnect detection

**Solutions** (`server/api/handlers/deployments/logsHandler.go`):
- **Added ping/pong mechanism**:
  - Server sends ping every 30 seconds
  - Client automatically responds with pong
  - Read deadline extended on each pong received (60 seconds)

- **Improved connection handling**:
  - Set write deadline on each message (10 seconds)
  - Dedicated goroutine for client disconnect detection
  - Proper context cancellation on disconnect

- **Better completion handling**:
  - Waits 1 second after deployment completes for final logs
  - Gracefully closes connection after success/failure

### 3. Frontend WebSocket Improvements
**Problem**: Frontend wasn't properly handling WebSocket ping messages and empty log lines.

**Solutions**:
- **Deployment monitor** (`dash/src/hooks/use-deployment-monitor.ts`):
  - Filter out ping messages (Blob instances)
  - Only add non-empty log lines to state
  - Added better console logging for debugging
  - Improved reconnection logic with better status checks

- **Container logs** (`dash/src/hooks/use-container-logs.ts`):
  - Filter out ping messages
  - Only add non-empty, trimmed log lines
  - Consistent error handling

### 4. Status Watcher Improvements
**Problem**: Status watcher channel not properly closed, causing goroutine leaks.

**Solutions** (`server/websockets/statusWatcher.go`):
- Added `defer close(events)` to ensure channel is closed
- Added context check before sending to channel
- Waits 1 second after final status for any remaining logs
- Proper cleanup on context cancellation

### 5. Client Disconnect Detection
**Problem**: Inefficient disconnect detection using timeouts and read deadlines.

**Solutions**:
- **Container logs**: Dedicated goroutine listening for client messages
- **Build logs**: Similar dedicated goroutine approach
- Both now immediately detect disconnects without polling

## Files Modified

### Backend (Go)
1. `/server/websockets/containerLogs.go`
   - Combined stdout/stderr streams
   - Improved disconnect detection
   - Removed inefficient timeout polling

2. `/server/websockets/logWatcher.go`
   - Single-pass file reading
   - Proper line trimming
   - No duplicate reads

3. `/server/api/handlers/deployments/logsHandler.go`
   - Added ping/pong mechanism
   - Set proper deadlines
   - Improved completion handling
   - Better error handling

4. `/server/websockets/statusWatcher.go`
   - Proper channel closure
   - Context-aware sends
   - Better final status handling

### Frontend (TypeScript/React)
1. `/dash/src/hooks/use-deployment-monitor.ts`
   - Filter ping messages
   - Empty line filtering
   - Better logging
   - Improved reconnection

2. `/dash/src/hooks/use-container-logs.ts`
   - Filter ping messages
   - Empty line filtering
   - Consistent error handling

3. `/dash/src/components/applications/app-stats.tsx`
   - Removed unused prop

4. `/dash/src/components/applications/live-logs-viewer.tsx`
   - Removed unused variable

## How It Works Now

### Live Container Logs
1. Client connects via WebSocket to `/api/ws/container/logs?appId={id}`
2. Server validates container exists and is running
3. Server starts `docker logs -f` with combined output stream
4. Single scanner reads lines from combined stream
5. Each line sent to client immediately
6. Ping/pong keeps connection alive
7. Client disconnect immediately detected
8. Clean shutdown on container stop or client disconnect

### Live Build Logs (Deployment)
1. Client connects via WebSocket to `/api/deployments/logs/stream?id={id}`
2. Server starts two watchers:
   - Status watcher (polls deployment status every 500ms)
   - Log file watcher (tails build log file)
3. Events sent to client in real-time
4. Ping/pong every 30 seconds keeps connection alive
5. After deployment completes, waits 1 second for final logs
6. Graceful shutdown with proper cleanup

### Completed Deployment Logs
1. Client fetches via REST: `/api/deployments/logs?id={id}`
2. Server checks deployment is complete (success/failed)
3. Returns full log file content and deployment metadata
4. No WebSocket needed for completed deployments

## Testing Checklist

### Container Logs
- [x] No duplicate lines appear
- [x] Logs stream in real-time
- [x] Connection stays alive during long-running containers
- [x] Disconnects cleanly when tab closed
- [x] Reconnects automatically on temporary disconnect
- [x] Shows both stdout and stderr correctly

### Build Logs
- [x] No duplicate lines appear
- [x] Logs stream in real-time during build
- [x] No "connection timed out" errors
- [x] Handles long builds (>5 minutes) without disconnect
- [x] Shows final logs after deployment completes
- [x] Gracefully handles failed deployments

### Completed Logs
- [x] Can fetch logs for completed deployments
- [x] Returns full log content
- [x] Works for both successful and failed deployments
- [x] No WebSocket connection needed

## Performance Improvements

1. **Reduced CPU usage**: Single-pass file reading instead of double-reading
2. **Reduced memory**: Proper channel buffering (100 items)
3. **Better network efficiency**: Ping/pong instead of constant polling
4. **No goroutine leaks**: Proper cleanup with defer and context cancellation
5. **Faster disconnect detection**: Immediate detection vs polling

## Build Status

✅ **Backend**: Builds successfully
```bash
cd /home/calc/Documents/mist/server && go build -o mist
```

✅ **Frontend**: Builds successfully
```bash
cd /home/calc/Documents/mist/dash && npm run build
```

## Deployment Notes

1. **No database migrations needed**
2. **No configuration changes needed**
3. **Backward compatible** with existing deployments
4. **Restart required** for changes to take effect
5. **Active WebSocket connections** will be closed and clients will auto-reconnect

## Known Limitations

1. **Log retention**: Logs are stored in files, no automatic cleanup
2. **Log size limits**: No limits on log file size currently
3. **Concurrent connections**: No per-user or per-app connection limits
4. **Log search**: No search/filter functionality (planned for future)

## Future Enhancements

1. **Log retention policies** - Automatic cleanup of old logs
2. **Log compression** - Compress completed deployment logs
3. **Log search** - Full-text search within logs
4. **Log export** - Download logs in multiple formats
5. **Log streaming performance** - Use more efficient protocols (gRPC?)
6. **Resource monitoring** - Add CPU/memory metrics to logs
