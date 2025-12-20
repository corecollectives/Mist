# Domains & SSL

Configure custom domains for your applications.

## Adding a Domain

1. Navigate to your application
2. Go to **Domains** tab
3. Click **"Add Domain"**
4. Enter your domain name (e.g., `app.example.com`)
5. Click **"Add"**

## DNS Configuration

Point your domain to your Mist server:

### A Record

```
Type: A
Name: app (or @ for root)
Value: YOUR_SERVER_IP
TTL: 3600
```

### CNAME Record (Subdomain)

```
Type: CNAME
Name: app
Value: yourdomain.com
TTL: 3600
```

## Verify DNS Propagation

Check if DNS has propagated:

```bash
dig app.example.com
nslookup app.example.com
```

Propagation can take 5 minutes to 48 hours.

## SSL/TLS Certificates

<div class="coming-soon-banner">
  <h4>🚧 Let's Encrypt Integration Coming Soon</h4>
  <p>Automatic SSL certificate provisioning and renewal is planned for the next major release.</p>
</div>

### Manual SSL Configuration

Currently, SSL certificates must be configured manually in Traefik. See [Traefik Configuration](/deployment/traefik) for details.

### Coming Features

- **Automatic Let's Encrypt** - Free SSL certificates
- **Auto-renewal** - Certificates renewed before expiry
- **Wildcard Support** - `*.example.com` certificates
- **Custom Certificates** - Upload your own SSL certs
- **Force HTTPS** - Automatic HTTP to HTTPS redirect

## Multiple Domains

Add multiple domains to the same application:

```
app.example.com
www.app.example.com
app.yourdomain.com
```

All domains will route to the same container.

## Domain Management

### Edit Domain

Update domain configuration in the domains list.

### Delete Domain

Remove a domain from your application:

1. Click the delete icon next to the domain
2. Confirm deletion

## Troubleshooting

### Domain Not Resolving

- Verify DNS records are correct
- Wait for DNS propagation
- Clear DNS cache: `sudo systemd-resolve --flush-caches`

### SSL Certificate Errors

- Check Traefik logs: `docker logs traefik`
- Verify domain points to correct IP
- Ensure ports 80 and 443 are open
