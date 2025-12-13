# UI Restrictions for Service and Database Apps

## Overview

Service apps (background workers) and database apps don't have external HTTP access and therefore shouldn't have UI options related to external accessibility.

---

## Restrictions Applied

### 1. View Live App Link - HIDDEN ❌

**Files Modified:**
- `dash/src/components/applications/app-info.tsx`
- `dash/src/components/applications/app-stats.tsx`

**Condition:**
```tsx
{previewUrl && app.status === "running" && app.appType === 'web' && (
  <a href={previewUrl}>View Live App</a>
)}
```

**Result:**
- ✅ Shows for web apps
- ❌ Hidden for service apps
- ❌ Hidden for database apps

---

### 2. Domains Section - HIDDEN ❌

**File Modified:**
- `dash/src/features/applications/AppPage.tsx`

**Change:**
```tsx
<TabsContent value="settings" className="space-y-6">
  <ContainerControls appId={app.id} onStatusChange={fetchAppDetails} />
  <AppSettings app={app} onUpdate={fetchAppDetails} />
  {app.appType === 'web' && <Domains appId={app.id} />}
</TabsContent>
```

**Result:**
- ✅ Domains section shows for web apps
- ❌ Domains section hidden for service apps
- ❌ Domains section hidden for database apps

**Reason:** Service and database apps don't use Traefik routing and aren't accessible via domains.

---

## Complete UI Comparison

### Web Apps

**Info Tab:**
- ✅ Status (with "View Live App" link if running)
- ✅ Deployment Strategy
- ✅ Git Repository
- ✅ Git Branch
- ✅ Port
- ✅ Root Directory
- ✅ Build/Start Commands
- ✅ Latest Commit

**Settings Tab:**
- ✅ Container Controls (Start/Stop/Restart)
- ✅ App Settings (Git, Port, Commands, Resources)
- ✅ **Domains Section** (Add/Remove domains)

---

### Service Apps (Background Workers)

**Info Tab:**
- ✅ Status (NO "View Live App" link)
- ✅ Deployment Strategy
- ✅ Git Repository
- ✅ Git Branch
- ✅ Port (internal only)
- ✅ Root Directory
- ✅ Build/Start Commands
- ✅ Latest Commit

**Settings Tab:**
- ✅ Container Controls (Start/Stop/Restart)
- ✅ App Settings (Git, Commands, Resources)
- ❌ **Domains Section** (Hidden)

---

### Database Apps

**Info Tab:**
- ✅ Status (NO "View Live App" link)
- ✅ Deployment Strategy
- ✅ Service Template (instead of git info)
- ✅ Port (internal only)
- ❌ Git Repository
- ❌ Git Branch
- ❌ Build/Start Commands
- ❌ Latest Commit

**Settings Tab:**
- ✅ Container Controls (Start/Stop/Restart)
- ✅ App Settings (Resource limits, restart policy)
- ⚠️ Git fields disabled (with info alert)
- ❌ **Domains Section** (Hidden)

---

## Backend Behavior

Even if domains are added via API directly, they won't be applied:

### Web Apps
```go
case models.AppTypeWeb:
    if len(domains) > 0 {
        // Apply Traefik routing with domains
    } else {
        // Direct port mapping
    }
```

### Service/Database Apps
```go
case models.AppTypeService:
case models.AppTypeDatabase:
    // Always internal network only
    // Domains are ignored even if they exist in database
    --network traefik-net
```

---

## Testing Checklist

### Web App
- [ ] "View Live App" link appears when running
- [ ] Domains section visible in Settings tab
- [ ] Can add domains
- [ ] Domains work with Traefik routing

### Service App
- [ ] NO "View Live App" link
- [ ] NO Domains section in Settings tab
- [ ] Port field visible but not exposed externally
- [ ] Git tab visible

### Database App
- [ ] NO "View Live App" link
- [ ] NO Domains section in Settings tab
- [ ] NO Git tab
- [ ] Template name shown instead of git info
- [ ] Port shown but internal only

---

## UI Screenshots (Conceptual)

### Web App Settings Tab
```
┌────────────────────────────────────────┐
│ Container Controls                     │
│ [Start] [Stop] [Restart]              │
└────────────────────────────────────────┘

┌────────────────────────────────────────┐
│ Application Settings                   │
│ Git Repository: owner/repo             │
│ Port: 3000                             │
│ ...                                    │
└────────────────────────────────────────┘

┌────────────────────────────────────────┐
│ Domains                            ✅  │
│ myapp.com                    [Remove]  │
│ [+ Add Domain]                         │
└────────────────────────────────────────┘
```

### Service/Database App Settings Tab
```
┌────────────────────────────────────────┐
│ Container Controls                     │
│ [Start] [Stop] [Restart]              │
└────────────────────────────────────────┘

┌────────────────────────────────────────┐
│ Application Settings                   │
│ ...settings fields...                  │
└────────────────────────────────────────┘

(NO Domains section) ❌
```

---

## Summary of Changes

| Feature | Web | Service | Database |
|---------|-----|---------|----------|
| View Live App Link | ✅ | ❌ | ❌ |
| Domains Section | ✅ | ❌ | ❌ |
| Git Tab | ✅ | ✅ | ❌ |
| Git Info in App Info | ✅ | ✅ | ❌ |
| Port Configuration | ✅ External | ✅ Internal | ✅ Internal |
| Container Controls | ✅ | ✅ | ✅ |
| Resource Limits | ✅ | ✅ | ✅ |
| Environment Variables | ✅ | ✅ | ✅ Auto-created |

---

## Files Modified

1. ✅ `dash/src/features/applications/AppPage.tsx` - Hide domains for non-web apps
2. ✅ `dash/src/components/applications/app-info.tsx` - Hide "View Live App" for non-web apps
3. ✅ `dash/src/components/applications/app-stats.tsx` - Hide "View Live App" button for non-web apps

---

## Related Documentation

- See `Docs/app-type-system.md` for complete app type system overview
- See `Docs/view-live-app-fix.md` for "View Live App" link fix details
- See `Docs/testing-app-types.md` for testing scenarios

---

## Future Considerations

### Possible Additional Restrictions

**Service Apps:**
- Could hide "Port" field entirely (since it's not externally accessible)
- Could add info banner: "This is a background service with no external access"

**Database Apps:**
- Could show connection string examples in a dedicated section
- Could add "Copy Connection String" button
- Could show internal network info (container name, internal port)

### Possible Additional Features

**Web Apps:**
- SSL/TLS certificate management (for domains)
- Custom domain verification
- Domain health checks

**All Types:**
- Network topology visualization
- Inter-app connection mapping
- Service discovery UI

---

## Conclusion

The UI now properly reflects the capabilities of each app type:
- Web apps show all external access features (domains, view live app)
- Service apps hide external access features
- Database apps hide both external access and git-related features

This creates a cleaner, more intuitive UX that guides users based on what's actually possible with each app type.
