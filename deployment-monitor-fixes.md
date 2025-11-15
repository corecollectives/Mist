# Deployment Monitor Fixes

## Issues Fixed

### 1. WebSocket Connection Issue
**Problem**: Frontend was trying to connect to `/api/ws/logs` but the backend endpoint is `/api/deployments/logs/stream`

**Fix**: Updated `useDeploymentMonitor.ts` line 84 to use correct WebSocket URL with proper protocol handling:
```typescript
const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
const host = window.location.host;
const ws = new WebSocket(`${protocol}//${host}/api/deployments/logs/stream?id=${deploymentId}`)
```

### 2. Infinite Fetching Loop for Completed Deployments
**Problem**: When viewing a completed deployment, the `fetchCompletedDeployment()` function would run on every re-render, causing constant API calls.

**Fix**: 
- Added `hasFetchedRef` to track if deployment has already been fetched
- Check `hasFetchedRef.current` at the start of `fetchCompletedDeployment()`
- Set `hasFetchedRef.current = true` after successful or failed fetch
- Reset `hasFetchedRef.current = false` in the `reset()` function

### 3. Toast Notifications for Old Deployments
**Problem**: When viewing old/completed deployment logs, the success/error callbacks (`onComplete`, `onError`) would trigger, showing toasts for historical events.

**Fix**: 
- Removed `onComplete?.()` and `onError?.()` calls from `fetchCompletedDeployment()` 
- Only call these callbacks for live WebSocket events (real-time deployments)
- Added comments explaining why callbacks are not called for completed deployments

### 4. Success/Error Banners for Historical Deployments
**Problem**: Success and error banners would show when viewing old deployment logs, even though these are historical views.

**Fix**: Modified `DeploymentMonitor.tsx` to only show banners for live deployments:
```typescript
{error && isLive && ( ... )}  // Only show error banner for live deployments
{status?.status === 'success' && isLive && ( ... )}  // Only show success banner for live deployments
```

## Testing Checklist

- [ ] Open a completed successful deployment → Should load instantly (REST API), no toast shown
- [ ] Open a completed failed deployment → Should load instantly, no error toast shown
- [ ] Create a new deployment → Should connect via WebSocket and stream logs in real-time
- [ ] Live deployment completes successfully → Should show success banner and toast
- [ ] Live deployment fails → Should show error banner and toast
- [ ] Close and reopen deployment monitor → Should not cause multiple API calls
- [ ] Check browser console → No WebSocket connection errors
- [ ] Check Network tab → Completed deployments use REST API, live use WebSocket

## Backend Endpoints (Reference)

- **REST API** (completed deployments): `GET /api/deployments/logs?id={id}`
  - Returns 200 with logs if deployment is completed
  - Returns 400 if deployment is still in progress

- **WebSocket** (live deployments): `WS /api/deployments/logs/stream?id={id}`
  - Streams real-time logs and status updates
  - Sends events: `log`, `status`, `error`

## Files Modified

1. `dash/src/features/applications/hooks/useDeploymentMonitor.ts`
   - Fixed WebSocket URL
   - Added fetch tracking to prevent infinite loops
   - Removed toast callbacks for completed deployments

2. `dash/src/features/applications/components/DeploymentMonitor.tsx`
   - Added `isLive` condition to success/error banners

