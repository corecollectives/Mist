# Getting Started

This guide will walk you through installing Mist and deploying your first application.

## Prerequisites

Before installing Mist, ensure you have:

- A Linux server (Ubuntu 20.04+ or Debian 11+ recommended)
- Docker installed and running
- At least 1GB RAM and 10GB disk space
- Root or sudo access
- A domain name (optional, but recommended for production)

## Installation

### Quick Install

The easiest way to install Mist is using the installation script:

```bash
curl -sSL https://raw.githubusercontent.com/corecollectives/mist/main/install.sh | bash
```

This script will:
1. Download the latest Mist binary
2. Set up Traefik reverse proxy
3. Create systemd service for auto-start
4. Configure firewall rules (ports 80, 443, 3000)

### Manual Installation

If you prefer to install manually:

```bash
# Download the latest release
wget https://github.com/corecollectives/mist/releases/latest/download/mist-linux-amd64

# Make it executable
chmod +x mist-linux-amd64

# Move to /usr/local/bin
sudo mv mist-linux-amd64 /usr/local/bin/mist

# Create data directory
sudo mkdir -p /var/lib/mist

# Run Mist
mist
```

[Learn more about installation →](/deployment/installation)

## First-Time Setup

After installation, Mist will be available at `http://your-server-ip:3000`.

### 1. Create Admin Account

On first visit, you'll see the setup page:

1. Enter admin email and password
2. Click "Create Admin Account"
3. You'll be automatically logged in

### 2. Configure GitHub Integration (Optional)

To enable Git deployments:

1. Go to **Settings** → **GitHub Integration**
2. Follow the instructions to create a GitHub App
3. Install the app on your repositories

[Learn more about GitHub setup →](/deployment/github-app)

## Deploy Your First Application

### Step 1: Create a Project

Projects organize your applications:

1. Click **"New Project"** in the dashboard
2. Enter a name (e.g., "My Portfolio")
3. Add tags (optional)
4. Click **"Create Project"**

### Step 2: Create an Application

1. Open your project
2. Click **"New Application"**
3. Fill in the details:
   - **Name**: Your app name
   - **Git Repository**: Select from connected repos
   - **Branch**: Choose branch to deploy
   - **Build Command**: `npm run build` (or your build command)
   - **Start Command**: `npm start` (or your start command)
   - **Port**: Application port (e.g., 3000)

### Step 3: Add Environment Variables

1. Go to the **Environment** tab
2. Click **"Add Variable"** or use **"Bulk Paste"**
3. Add your environment variables:
   ```
   NODE_ENV=production
   DATABASE_URL=your-database-url
   API_KEY=your-api-key
   ```

### Step 4: Deploy

1. Click **"Deploy"** button
2. Watch real-time build logs
3. Wait for deployment to complete
4. Access your app via the generated URL

## Configure Custom Domain

### Step 1: Add Domain

1. Go to **Domains** tab in your application
2. Click **"Add Domain"**
3. Enter your domain (e.g., `app.example.com`)

### Step 2: Configure DNS

Point your domain to your server:

```
Type: A Record
Name: app (or @ for root domain)
Value: YOUR_SERVER_IP
TTL: 3600 (or auto)
```

### Step 3: Wait for DNS Propagation

DNS changes can take 5 minutes to 48 hours to propagate. Check with:

```bash
dig app.example.com
```

::: tip SSL/TLS Certificate
Currently, SSL certificates need to be configured manually in Traefik. Automatic Let's Encrypt integration is coming soon!

[Learn more about SSL →](/guide/ssl-automation)
:::

## Next Steps

Now that you have Mist running, explore these topics:

- [**Projects**](./projects) - Organize applications
- [**Deployments**](./deployments) - Deployment strategies
- [**Environment Variables**](./environment-variables) - Managing configuration
- [**Monitoring**](./logs) - Logs and metrics
- [**Git Integration**](./git-integration) - Auto-deploy on push

## Getting Help

- [GitHub Issues](https://github.com/corecollectives/mist/issues) - Report bugs
- [GitHub Discussions](https://github.com/corecollectives/mist/discussions) - Ask questions
- [Documentation](/) - Read more guides

## What's Next?

### Essential Setup

- ✅ Configure GitHub integration for auto-deployments
- ✅ Set up custom domains for your applications  
- ✅ Configure backup strategy for SQLite database

### Coming Soon Features

- 🚧 Let's Encrypt SSL automation
- 🚧 Database provisioning (PostgreSQL, MySQL, Redis)
- 🚧 Email notifications for deployments
- 🚧 CLI tool for terminal deployments
