# Database Services <span class="coming-soon-badge">Coming Soon</span>

One-click database provisioning for your applications.

<div class="coming-soon-banner">
  <h4>🚧 Feature in Development</h4>
  <p>Managed database services are planned for Phase 2 of the roadmap.</p>
</div>

## Planned Databases

### PostgreSQL
- Multiple versions (12, 13, 14, 15, 16)
- Automatic backups
- Connection string auto-injection
- pgAdmin integration
- User management

### MySQL / MariaDB
- Version selection
- Automated backups with mysqldump
- phpMyAdmin integration
- Privilege management

### Redis
- Persistence options (RDB, AOF)
- Password protection
- Redis Commander UI
- Pub/sub support

### MongoDB
- Replica set support
- Mongo Express UI
- Backup and restore
- User authentication

## Features

- **One-Click Provisioning** - Create databases instantly
- **Auto-Configuration** - Environment variables injected automatically
- **Backup & Restore** - Automated backup scheduling
- **Management UIs** - Built-in admin interfaces
- **Connection Pooling** - Optimize database connections
- **Replication** - High availability setups

## Use Cases

```bash
# Automatic environment variable injection
DATABASE_URL=postgresql://user:pass@postgres-db:5432/myapp
REDIS_URL=redis://redis-cache:6379
MONGODB_URI=mongodb://mongo-db:27017/myapp
```

## Expected Release

Target: Q2 2025

[View full roadmap →](https://github.com/corecollectives/mist/blob/main/roadmap.md#phase-2-database--services)
