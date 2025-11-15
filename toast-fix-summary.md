# Toast Fix - Only Show for Live Deployments

## Problem
Toasts were showing when viewing old/completed deployment logs via REST API, even though the user was just viewing historical data.

## Solution
Modified the hook so that `onComplete()` and `onError()` callbacks are ONLY triggered for live WebSocket deployments, NOT for completed REST API fetches.

## Changes Made

### File: `dash/src/features/applications/hooks/useDeploymentMonitor.ts`

#### ❌ Before (lines 68-71):
```typescript
if (deployment.status === 'failed' && deployment.error_message) {
  setError(deployment.error_message);
  onError?.(deployment.error_message);  // ❌ Called for historical deployments
} else if (deployment.status === 'success') {
  onComplete?.();  // ❌ Called for historical deployments
}
```

#### ✅ After (lines 68-71):
```typescript
if (deployment.status === 'failed' && deployment.error_message) {
  setError(deployment.error_message);
  // ✅ No callback - just viewing history
} else if (deployment.status === 'success') {
  // ✅ No callback - just viewing history
}
```

#### ✅ WebSocket Handler (lines 112-119) - UNCHANGED:
```typescript
// Handle completion
if (statusData.status === 'success') {
  onComplete?.();  // ✅ ONLY called for live deployments
}

// Handle errors
if (statusData.status === 'failed' && statusData.error_message) {
  setError(statusData.error_message);
  onError?.(statusData.error_message);  // ✅ ONLY called for live deployments
}
```

## Behavior Now

### Viewing Completed Deployment (REST API)
```
User clicks "View Logs" on old deployment
         ↓
    REST API fetch
         ↓
   Returns 200 OK
         ↓
  Display logs instantly
         ↓
  ❌ NO toast shown
  ❌ NO onComplete() callback
  ❌ NO onError() callback
```

### Watching Live Deployment (WebSocket)
```
User clicks "Deploy Now"
         ↓
   Opens monitor immediately
         ↓
   REST returns 400 (in progress)
         ↓
   Connect WebSocket
         ↓
   Stream logs in real-time
         ↓
   Deployment completes
         ↓
   WebSocket sends status event
         ↓
  ✅ onComplete() called
  ✅ Toast shown: "Deployment completed successfully!"
```

## Testing Checklist

- [x] Build succeeds
- [ ] Open old completed deployment → No toast
- [ ] Open old failed deployment → No toast
- [ ] Start new deployment → Shows toast when complete
- [ ] Start new deployment that fails → Shows error toast

## Files Modified

1. `dash/src/features/applications/hooks/useDeploymentMonitor.ts`
   - Removed `onComplete()` and `onError()` calls from REST API path (lines 68-71)
   - Kept callbacks in WebSocket handler only (lines 112-119)

## Parent Components (No Changes Needed)

Both `Info.tsx` and `Deployments.tsx` pass `onComplete` callback that shows toast:
```typescript
onComplete={() => {
  toast.success("Deployment completed successfully!")
}}
```

These will now ONLY fire for live deployments, not historical views. ✅
