# GitHub App Setup

Configure GitHub App integration for automatic deployments.

## Create GitHub App

### 1. Navigate to GitHub Settings

- Go to GitHub → Settings → Developer settings → GitHub Apps
- Click **"New GitHub App"**

### 2. Configure Basic Information

- **GitHub App name**: `Mist (Your Server)`
- **Homepage URL**: `https://your-mist-instance.com`
- **Webhook URL**: `https://your-mist-instance.com/api/webhooks/github`
- **Webhook secret**: Generate a random secret

### 3. Set Permissions

**Repository permissions:**
- Contents: Read
- Metadata: Read
- Webhooks: Read & Write

**Subscribe to events:**
- Push

### 4. Generate Private Key

- Scroll down and click **"Generate a private key"**
- Download the `.pem` file
- Save it to `/etc/mist/github-app.pem`

### 5. Note App ID

- Copy the App ID from the GitHub App page
- You'll need this for configuration

## Configure Mist

Add to `/etc/mist/config.yml` or environment variables:

```yaml
github:
  app_id: "123456"
  private_key_path: "/etc/mist/github-app.pem"
  webhook_secret: "your-webhook-secret"
```

Or via environment:

```bash
export GITHUB_APP_ID=123456
export GITHUB_APP_PRIVATE_KEY=/etc/mist/github-app.pem
export GITHUB_WEBHOOK_SECRET=your-webhook-secret
```

## Install GitHub App

### On Your Repositories

1. Go to your GitHub App settings
2. Click **"Install App"**
3. Select account/organization
4. Choose repositories (all or select)
5. Click **"Install"**

### Verify Installation

In Mist dashboard:
1. Create new application
2. Repository dropdown should show installed repos
3. Select repository and branch

## Test Webhook

1. Push code to your repository
2. Check Mist dashboard for automatic deployment
3. If no deployment, check:
   - Webhook delivery in GitHub App settings
   - Mist logs: `sudo journalctl -u mist -f`

## Troubleshooting

### Webhooks Not Working

- Verify webhook URL is correct and accessible
- Check webhook secret matches configuration
- Ensure firewall allows GitHub webhook IPs
- Check webhook delivery logs in GitHub

### Private Key Errors

```bash
# Verify file permissions
chmod 600 /etc/mist/github-app.pem

# Verify file format (should start with -----BEGIN RSA PRIVATE KEY-----)
head -n 1 /etc/mist/github-app.pem
```

For more information, see [GitHub Apps Documentation](https://docs.github.com/en/apps).
