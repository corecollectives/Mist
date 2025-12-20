# Security Best Practices

Secure your Mist installation in production.

## Server Security

### Firewall Configuration

```bash
# UFW example
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 3000/tcp
sudo ufw enable
```

### SSH Hardening

```bash
# Disable root login
sudo sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config

# Disable password auth (use SSH keys only)
sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config

# Restart SSH
sudo systemctl restart sshd
```

### Keep System Updated

```bash
# Ubuntu/Debian
sudo apt update && sudo apt upgrade -y

# Enable automatic security updates
sudo apt install unattended-upgrades
sudo dpkg-reconfigure -plow unattended-upgrades
```

## Mist Security

### Strong Passwords

- Minimum 12 characters
- Mix of uppercase, lowercase, numbers, symbols
- Use password manager

### JWT Secret

Generate strong JWT secret:

```bash
# Generate random 64-character secret
openssl rand -hex 32
```

Add to configuration:

```bash
export JWT_SECRET=your-generated-secret
```

### GitHub Webhook Secret

Generate webhook secret:

```bash
openssl rand -hex 32
```

Use in GitHub App configuration and Mist config.

## Docker Security

### Socket Permission

Limit Docker socket access:

```bash
# Create docker group if not exists
sudo groupadd docker

# Add mist user to docker group
sudo usermod -aG docker mist-user

# Restrict socket permissions
sudo chmod 660 /var/run/docker.sock
```

### Image Security

- Use official base images
- Scan images for vulnerabilities (coming soon)
- Keep images updated
- Minimize image layers

## Network Security

### Reverse Proxy

- Always use Traefik as reverse proxy
- Never expose application ports directly
- Configure proper CORS headers

### SSL/TLS

<div class="coming-soon-banner">
  <h4>🚧 Automatic SSL Coming Soon</h4>
</div>

Currently:
- Configure SSL manually in Traefik
- Use Let's Encrypt certificates
- Force HTTPS redirects

Future:
- Automatic SSL provisioning
- Auto-renewal
- HSTS headers

## Database Security

### SQLite Security

```bash
# Restrict database file permissions
chmod 600 /var/lib/mist/mist.db
chown mist:mist /var/lib/mist/mist.db
```

### Backups

```bash
# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
cp /var/lib/mist/mist.db /var/lib/mist/backups/mist_$DATE.db

# Keep last 7 days
find /var/lib/mist/backups -name "mist_*.db" -mtime +7 -delete
```

Schedule with cron:

```bash
# Daily backup at 2 AM
0 2 * * * /usr/local/bin/mist-backup.sh
```

## Application Security

### Environment Variables

- Never commit secrets to Git
- Use environment variables for all sensitive data
- Rotate API keys regularly

### Audit Logging

- Review audit logs regularly
- Monitor for suspicious activity
- Track failed login attempts

## Monitoring

### Log Monitoring

```bash
# Watch Mist logs
sudo journalctl -u mist -f

# Watch Docker logs
docker logs -f traefik
```

### System Monitoring

- Monitor disk space
- Track memory usage
- Watch CPU utilization
- Set up alerts (coming soon)

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Security Enhancements</h4>
</div>

- **Two-Factor Authentication** - TOTP support
- **Rate Limiting** - API abuse prevention
- **IP Whitelisting** - Restrict admin access
- **Secrets Encryption** - Encrypt env vars at rest
- **Vulnerability Scanning** - Scan Docker images
- **Security Headers** - CSP, HSTS, X-Frame-Options
- **Audit Log UI** - Browse security events

## Security Checklist

- [ ] Firewall configured and enabled
- [ ] SSH hardened (key-only auth, no root)
- [ ] System updates automated
- [ ] Strong passwords enforced
- [ ] JWT secret generated and configured
- [ ] GitHub webhook secret set
- [ ] Docker socket permissions restricted
- [ ] Database file permissions locked down
- [ ] Automated backups configured
- [ ] SSL/TLS certificates installed
- [ ] Monitoring and alerting set up

## Reporting Security Issues

Found a security vulnerability? Please email security@mistpaas.com (coming soon) or report privately on GitHub.
