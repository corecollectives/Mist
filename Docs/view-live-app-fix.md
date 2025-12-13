# View Live App Link - App Type Filtering

## Issue
The "View Live App" link was appearing for all app types, including background services and databases, which don't have external web access.

## Fix Applied

### Files Modified

#### 1. `dash/src/components/applications/app-info.tsx`

**Before:**
```tsx
{previewUrl && app.status === "running" && (
  <a href={previewUrl} target="_blank" rel="noopener noreferrer">
    View Live App
  </a>
)}
```

**After:**
```tsx
{previewUrl && app.status === "running" && app.appType === 'web' && (
  <a href={previewUrl} target="_blank" rel="noopener noreferrer">
    View Live App
  </a>
)}
```

#### 2. `dash/src/components/applications/app-stats.tsx`

**Changes:**
- Added `App` type import
- Updated interface to accept `app?: App`
- Added app type check to Preview URL button

**Before:**
```tsx
{previewUrl && containerStatus.state === "running" && (
  <a href={previewUrl}>View Live App</a>
)}
```

**After:**
```tsx
{previewUrl && containerStatus.state === "running" && app?.appType === 'web' && (
  <a href={previewUrl}>View Live App</a>
)}
```

#### 3. `dash/src/features/applications/AppPage.tsx`

**Change:**
- Pass `app` object to AppStats component

**Updated:**
```tsx
<AppStats appId={app.id} appStatus={app.status} app={app} previewUrl={previewUrl} />
```

## Behavior

### Web Apps (`app.appType === 'web'`)
- ✅ "View Live App" link appears when app is running
- ✅ Link opens the application in a new tab
- ✅ Visible in both app-info and app-stats components

### Service Apps (`app.appType === 'service'`)
- ❌ "View Live App" link is hidden
- ℹ️ Background services don't have HTTP endpoints
- ℹ️ Not accessible from external network

### Database Apps (`app.appType === 'database'`)
- ❌ "View Live App" link is hidden
- ℹ️ Databases only accessible via internal network
- ℹ️ Use connection strings from other apps

## Testing

### Test Web App
1. Create a web app with port 3000
2. Deploy the app
3. Navigate to app page
4. Verify "View Live App" link appears in:
   - Status badge area (app-info)
   - Container Status card (app-stats)
5. Click link and verify it opens the application

### Test Service App
1. Create a background service app
2. Deploy the app
3. Navigate to app page
4. Verify "View Live App" link does NOT appear

### Test Database App
1. Create PostgreSQL or Redis database
2. Deploy the app
3. Navigate to app page
4. Verify "View Live App" link does NOT appear

## UI Examples

### Web App - Shows Link ✅
```
┌─────────────────────────────────────────┐
│ Status                                  │
│ ● Running | 🔗 View Live App           │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│ Container Status                        │
│ Status: ● Running                       │
│ Uptime: 2h 30m                         │
│ ┌─────────────────────────────────┐   │
│ │  🔗 View Live App               │   │
│ └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

### Service/Database App - No Link ✅
```
┌─────────────────────────────────────────┐
│ Status                                  │
│ ● Running                               │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│ Container Status                        │
│ Status: ● Running                       │
│ Uptime: 2h 30m                         │
│ (No View Live App button)              │
└─────────────────────────────────────────┘
```

## Summary

The "View Live App" link now correctly appears only for web applications. Background services and databases, which don't have external HTTP access, no longer show this misleading link.

This improves UX by:
- Preventing confusion about service/database accessibility
- Clearly indicating which apps are web-accessible
- Maintaining consistency with the app type system

## Related Files

- `dash/src/components/applications/app-info.tsx` ✅ Fixed
- `dash/src/components/applications/app-stats.tsx` ✅ Fixed
- `dash/src/features/applications/AppPage.tsx` ✅ Updated
- `dash/src/types/app.ts` (contains AppType definition)
