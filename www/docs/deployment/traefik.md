# Traefik Setup

Configure Traefik reverse proxy for Mist.

## Basic Configuration

Traefik is installed automatically with the Mist installation script.

### Docker Compose

`/etc/mist/traefik/docker-compose.yml`:

```yaml
version: '3'

services:
  traefik:
    image: traefik:v2.10
    container_name: traefik
    restart: unless-stopped
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"  # Traefik dashboard
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik.yml:/traefik.yml:ro
      - ./acme.json:/acme.json
    networks:
      - mist
    labels:
      - "traefik.enable=true"

networks:
  mist:
    external: true
```

### Static Configuration

`/etc/mist/traefik/traefik.yml`:

```yaml
global:
  checkNewVersion: true
  sendAnonymousUsage: false

api:
  dashboard: true
  insecure: true

entryPoints:
  web:
    address: ":80"
  websecure:
    address: ":443"

providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
    network: mist

log:
  level: INFO
```

## SSL/TLS Configuration (Manual)

<div class="coming-soon-banner">
  <h4>🚧 Automatic Let's Encrypt Coming Soon</h4>
  <p>For now, SSL must be configured manually.</p>
</div>

### Let's Encrypt (Manual Setup)

Add to `traefik.yml`:

```yaml
certificatesResolvers:
  letsencrypt:
    acme:
      email: your-email@example.com
      storage: /acme.json
      httpChallenge:
        entryPoint: web
```

Create `acme.json`:

```bash
touch /etc/mist/traefik/acme.json
chmod 600 /etc/mist/traefik/acme.json
```

## Dashboard Access

Traefik dashboard is available at `http://your-server:8080`

::: warning Security
Disable the insecure API in production or restrict access via firewall.
:::

## Troubleshooting

### Check Traefik Logs

```bash
docker logs traefik
```

### Verify Configuration

```bash
docker exec traefik traefik version
```

For more details, see [Traefik Documentation](https://doc.traefik.io/traefik/).
