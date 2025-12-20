# SSL Automation <span class="coming-soon-badge">Coming Soon</span>

Automatic SSL/TLS certificate provisioning with Let's Encrypt.

<div class="coming-soon-banner">
  <h4>🚧 Feature in Development</h4>
  <p>Let's Encrypt integration is a high-priority feature for Phase 1.</p>
</div>

## Planned Features

### Automatic Certificate Issuance

- Let's Encrypt ACME integration
- Automatic certificate generation on domain add
- Wildcard certificate support
- Multi-domain certificates (SAN)

### Auto-Renewal

- Certificates renewed 30 days before expiry
- Automatic renewal process
- Zero-downtime renewal
- Renewal notifications

### Security Options

- **Force HTTPS** - Automatic HTTP to HTTPS redirect
- **HSTS** - HTTP Strict Transport Security headers
- **TLS 1.3** - Modern encryption standards
- **Certificate Pinning** - Enhanced security

### Custom Certificates

- Upload your own SSL certificates
- Support for wildcard certificates
- Certificate chain validation
- Private key encryption

## How It Will Work

```bash
# 1. Add domain
mist domains add app.example.com

# 2. Mist automatically:
#    - Verifies domain ownership (HTTP-01 or DNS-01)
#    - Requests certificate from Let's Encrypt
#    - Installs certificate in Traefik
#    - Configures HTTPS routing

# 3. Auto-renewal runs 30 days before expiry
```

## Expected Release

Target: Q1 2025 - Quick Win Priority

[View implementation plan →](https://github.com/corecollectives/mist/blob/main/roadmap.md#ssl-automation)
