# System Requirements

Hardware and software requirements for running Mist.

## Minimum Requirements

### Small Deployments (1-5 applications)

- **CPU**: 1 vCPU
- **RAM**: 1GB
- **Disk**: 10GB SSD
- **Network**: 100 Mbps
- **Cost**: ~$5-10/month VPS

### Medium Deployments (5-20 applications)

- **CPU**: 2 vCPU
- **RAM**: 2-4GB
- **Disk**: 20-50GB SSD
- **Network**: 500 Mbps
- **Cost**: ~$12-24/month VPS

### Large Deployments (20+ applications)

- **CPU**: 4+ vCPU
- **RAM**: 8GB+
- **Disk**: 100GB+ SSD
- **Network**: 1 Gbps
- **Cost**: ~$40+/month VPS

## Software Requirements

### Operating System

- Ubuntu 20.04 LTS or later
- Debian 11 or later
- Any Linux distribution with:
  - systemd
  - Docker support
  - Modern kernel (5.x+)

### Docker

- Docker Engine 20.10 or later
- Docker Compose (optional, for Traefik)

### Network

- Ports 80, 443, 3000 must be available
- Outbound internet access for:
  - GitHub API
  - Docker Hub
  - Package downloads

## Recommended Cloud Providers

### DigitalOcean

- **Droplet**: Basic ($12/month, 2GB RAM)
- Pre-installed Docker images available
- Good for most use cases

### Hetzner Cloud

- **CX21**: (€5/month, 2GB RAM)
- Excellent price/performance
- European data centers

### Linode

- **Nanode**: ($5/month, 1GB RAM)
- Good for small deployments
- Global data centers

### Vultr

- **Regular Performance**: ($6/month, 1GB RAM)
- Multiple locations worldwide

## Storage Considerations

### Database

- SQLite database grows slowly (~10MB per 1000 deployments)
- Recommend at least 5GB for database

### Git Repositories

- Cloned repos consume disk space
- Average repo: 50-500MB
- Cleaned up after successful builds

### Docker Images

- Images can be large (100MB - 2GB)
- Keep last 3-5 deployments: 1-10GB per app
- Implement cleanup policy

### Logs

- Build logs stored in database
- Container logs rotate automatically
- Recommend 5-10GB for logs

## Performance Benchmarks

### API Response Time

- Average: 20-50ms
- p95: < 100ms
- p99: < 200ms

### Concurrent Deployments

- Small (1 CPU): 1-2 concurrent
- Medium (2 CPU): 3-5 concurrent
- Large (4+ CPU): 10+ concurrent

### WebSocket Connections

- Hundreds of concurrent connections supported
- Minimal resource overhead

## Scaling Considerations

<div class="coming-soon-banner">
  <h4>🚧 Multi-Node Support Coming Soon</h4>
  <p>Horizontal scaling planned for Phase 4.</p>
</div>

### Current Limitations

- Single server deployment only
- Vertical scaling (upgrade server) required
- No load balancing across nodes

### Future Plans

- Docker Swarm support
- Kubernetes option
- Multi-server deployments
- Shared storage for builds
