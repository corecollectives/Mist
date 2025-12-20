# Traefik Configuration

This document explains how Mist manages Traefik configuration dynamically.

## Overview

Mist uses Traefik as a reverse proxy for routing traffic to applications and the Mist dashboard. The configuration is split into two parts:

1. **Static Configuration** (`traefik-static.yml`) - Provided from the repository and mounted as `/etc/traefik/traefik.yml`
2. **Dynamic Configuration** (`/var/lib/mist/traefik/dynamic.yml`) - Generated at runtime

## Static Configuration

The static configuration (`traefik-static.yml`) defines:
- Entry points (HTTP on port 80, HTTPS on port 443)
- Docker provider settings (connects to Docker socket)
- File provider for dynamic configuration (watches `/etc/traefik/dynamic/`)
- Let's Encrypt certificate resolver with HTTP challenge
- Traefik dashboard settings (accessible on port 8081)
- Logging configuration

This file is:
- Provided from the repository
- Mounted as read-only to `/etc/traefik/traefik.yml` in the container
- Email address is configured during installation via `install.sh`

### Static Configuration Content

```yaml
api:
  dashboard: true
  insecure: true

providers:
  docker:
    exposedByDefault: false
    network: traefik-net
    endpoint: "unix:///var/run/docker.sock"
  file:
    directory: /etc/traefik/dynamic
    watch: true

entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

certificatesResolvers:
  le:
    acme:
      email: admin@example.com  # Configured during installation
      storage: /letsencrypt/acme.json
      httpChallenge:
        entryPoint: web

log:
  level: INFO
```

## Dynamic Configuration

The dynamic configuration is automatically generated and updated when:
- The server starts up
- The wildcard domain is changed in settings
- The Mist app name is changed in settings

### Generated Content

When a wildcard domain is configured (e.g., `*.example.com` or `example.com`), the system generates:

1. **Mist Dashboard Router (HTTPS)**
   - Rule: `Host(mist.example.com)` (or whatever `mistAppName` is set to)
   - Entry point: websecure (port 443)
   - TLS: Let's Encrypt automatic certificate
   - Service: Proxies to the Mist dashboard container

2. **Mist Dashboard Router (HTTP)**
   - Rule: `Host(mist.example.com)`
   - Entry point: web (port 80)
   - Middleware: Redirects to HTTPS

3. **HTTPS Redirect Middleware**
   - Redirects all HTTP traffic to HTTPS

### Example Generated Configuration

```yaml
http:
  routers:
    mist-dashboard:
      rule: "Host(`mist.example.com`)"
      entryPoints:
        - websecure
      service: mist-dashboard
      tls:
        certResolver: le
    mist-dashboard-http:
      rule: "Host(`mist.example.com`)"
      entryPoints:
        - web
      middlewares:
        - https-redirect
      service: mist-dashboard

  services:
    mist-dashboard:
      loadBalancer:
        servers:
          - url: "http://mist:5173"

  middlewares:
    https-redirect:
      redirectScheme:
        scheme: https
        permanent: true
```

## File Locations

- **Static Config (in repo)**: `traefik-static.yml` → mounted to `/etc/traefik/traefik.yml` in container
- **Dynamic Config (runtime)**: `/var/lib/mist/traefik/dynamic.yml` → mounted to `/etc/traefik/dynamic/dynamic.yml` in container
- **Let's Encrypt certificates**: `./letsencrypt/acme.json`

## Installation

During installation, the `install.sh` script:

1. Prompts for your email address for Let's Encrypt certificates
2. Updates the `traefik-static.yml` file with your email
3. Creates the Traefik configuration directory (`/var/lib/mist/traefik`)
4. Creates the `traefik-net` Docker network if it doesn't exist
5. Starts Traefik using Docker Compose

The static configuration is only modified once during installation to set the Let's Encrypt email address.

## Docker Compose Integration

The `traefik-compose.yml` file mounts:
- Docker socket (for discovering containers)
- Let's Encrypt certificate directory
- Dynamic configuration directory (read-only)
- Static configuration file (read-only)

```yaml
services:
  traefik:
    image: traefik:v3.1
    container_name: traefik
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
      - "8081:8080"  
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "./letsencrypt:/letsencrypt"
      - "/var/lib/mist/traefik:/etc/traefik/dynamic:ro"
      - "./traefik-static.yml:/etc/traefik/traefik.yml:ro"
    networks:
      - traefik-net
```

**Note:** Traefik automatically reads `/etc/traefik/traefik.yml` as its static configuration file.

## Application Routing

Individual applications use Docker labels to configure their routing (managed in `server/docker/build.go`):

```go
"-l", "traefik.enable=true",
"-l", fmt.Sprintf("traefik.http.routers.%s.rule=%s", containerName, hostRule),
"-l", fmt.Sprintf("traefik.http.routers.%s.entrypoints=websecure", containerName),
"-l", fmt.Sprintf("traefik.http.routers.%s.tls=true", containerName),
"-l", fmt.Sprintf("traefik.http.routers.%s.tls.certresolver=le", containerName),
```

## Automatic Updates

The dynamic configuration is automatically regenerated when system settings are updated via the API:

```
POST /api/settings/system
{
  "wildcardDomain": "example.com",
  "mistAppName": "mist"
}
```

Traefik watches the dynamic configuration directory and automatically reloads when files change (no restart needed).

## Troubleshooting

### Check Dynamic Configuration
```bash
cat /var/lib/mist/traefik/dynamic.yml
```

### View Traefik Logs
```bash
docker logs traefik
```

### Access Traefik Dashboard
```
http://your-server-ip:8081/dashboard/
```

### Verify Configuration in Logs
When settings are updated, check the server logs for:
```
Generated Traefik dynamic config path=/var/lib/mist/traefik/dynamic.yml
```
