# Deployment Monitor Logic Flow

## Decision Tree: REST API vs WebSocket

The hook uses a **smart detection** pattern - it tries REST first, then falls back to WebSocket if needed.

```
┌─────────────────────────────────┐
│   User Opens Deployment Monitor │
└────────────┬────────────────────┘
             │
             ▼
┌─────────────────────────────────┐
│  Call fetchCompletedDeployment() │
│  (REST API GET /api/deployments/logs?id=X) │
└────────────┬────────────────────┘
             │
             ▼
       ┌─────────────┐
       │  Response?   │
       └─────┬───┬───┘
             │   │
       200 OK│   │400 Bad Request
             │   │
             ▼   ▼
    ┌────────┐ ┌──────────┐
    │COMPLETED│ │  LIVE    │
    │         │ │          │
    │Show logs│ │Set isLive│
    │instantly│ │= true    │
    │         │ │          │
    │isLive = │ │          │
    │false    │ │          │
    │         │ └────┬─────┘
    │         │      │
    │         │      ▼
    │         │ ┌──────────────────┐
    │         │ │connectWebSocket()│
    │         │ │WS /api/deployments/logs/stream?id=X│
    │         │ └────┬─────────────┘
    │         │      │
    │         │      ▼
    │         │ ┌──────────────┐
    │         │ │Stream logs & │
    │         │ │status updates│
    │         │ │in real-time  │
    └─────────┴─┴──────────────┘
```

## Key Points

### ✅ Mutually Exclusive
- **Completed deployment**: Uses REST API ONLY (isLive = false)
- **Live deployment**: Uses WebSocket ONLY (isLive = true)
- The REST call is just a "probe" to determine deployment state

### ✅ No Double Fetching
- `hasFetchedRef` prevents REST API from being called multiple times
- REST call happens ONCE when component mounts
- If REST returns 400 → Switch to WebSocket mode permanently

### ✅ State Transitions
```
Initial State: isLive = false, hasFetchedRef = false
     │
     ▼
REST API Call
     │
     ├─ 200 OK ──→ isLive = false, hasFetchedRef = true (DONE - Show logs)
     │
     └─ 400 Error ──→ isLive = true, hasFetchedRef = true
                           │
                           ▼
                     Connect WebSocket (Stream until complete)
```

### ✅ No Mixed Methods
- Once `isLive = true`, we NEVER call REST API again
- Once `isLive = false` (completed), we NEVER open WebSocket
- `hasFetchedRef` ensures REST is only called once per mount

## Code Reference

### REST API Detection (lines 29-81)
```typescript
const fetchCompletedDeployment = useCallback(async () => {
  if (hasFetchedRef.current) return; // Prevent duplicate calls
  
  const response = await fetch(`/api/deployments/logs?id=${deploymentId}`);
  
  if (response.status === 400) {
    // Deployment is live - switch to WebSocket mode
    setIsLive(true);
    return;
  }
  
  // Deployment is completed - show logs from REST response
  // isLive remains false - no WebSocket connection
});
```

### WebSocket Connection (lines 83-166)
```typescript
const connectWebSocket = useCallback(() => {
  if (!enabled || !isLive) return; // Only connect if isLive = true
  
  const ws = new WebSocket(`wss://host/api/deployments/logs/stream?id=${deploymentId}`);
  // Stream logs and status updates
});
```

### Effect Hooks
```typescript
// Hook 1: Try REST API first (runs once on mount)
useEffect(() => {
  if (enabled) {
    fetchCompletedDeployment(); // Will set isLive if needed
  }
}, [enabled, fetchCompletedDeployment]);

// Hook 2: Connect WebSocket ONLY if isLive = true
useEffect(() => {
  if (isLive && enabled) {
    connectWebSocket(); // Only runs if REST API returned 400
  }
}, [isLive, enabled, connectWebSocket]);
```

## Summary

**The hook does NOT use both methods.** It uses:
1. REST API as a "probe" to check deployment state (one-time)
2. If completed → Display REST response (no WebSocket)
3. If live → Connect WebSocket (ignore REST response)

This is an **either-or** pattern, not a **both** pattern.
