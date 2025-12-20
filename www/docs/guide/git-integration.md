# Git Integration

Connect Mist with GitHub to enable automatic deployments.

## GitHub App Setup

### 1. Create GitHub App

Follow the [GitHub App Setup Guide](/deployment/github-app) to create and install the Mist GitHub App.

### 2. Install on Repositories

Install the GitHub App on repositories you want to deploy:

1. Visit GitHub App settings
2. Click **"Install App"**
3. Select repositories
4. Grant required permissions

### 3. Connect in Mist

Once installed, repositories appear in the application creation form.

## Auto-Deploy on Push

Enable automatic deployments:

1. Create application with GitHub repository
2. Select branch to deploy
3. Push code to that branch
4. Mist receives webhook and deploys automatically

## Webhook Events

Mist listens for these GitHub events:

- **Push** - Code pushed to branch
- **Pull Request** (coming soon) - PR opened/updated

## Commit Tracking

Each deployment tracks:

- Commit hash (SHA)
- Commit message
- Author information
- Timestamp

View commit details in the deployments list.

## Branch Management

### Selecting Branch

Choose which branch to deploy:

- `main` - Production deployments
- `develop` - Staging environment
- `feature/*` - Feature testing

### Changing Branch

Update the deployment branch:

1. Go to application settings
2. Change **Branch** field
3. Save changes
4. Next deployment uses new branch

## Coming Soon

<div class="coming-soon-banner">
  <h4>🚧 Git Features in Development</h4>
</div>

- **GitLab Integration** - Connect GitLab repositories
- **Bitbucket Support** - Bitbucket repositories
- **Pull Request Previews** - Deploy PRs to preview URLs
- **Commit Status Updates** - Update GitHub with deployment status
- **Multiple Git Providers** - Use multiple providers per app
- **Self-hosted Git** - Gitea, Gogs support

## Troubleshooting

### Webhook Not Triggering

Check:
- GitHub App is installed on repository
- Webhook URL is correct
- Repository permissions are granted
- Firewall allows GitHub webhook IPs

### Wrong Branch Deploying

Verify:
- Branch name matches exactly (case-sensitive)
- Application settings have correct branch
- Webhook is configured for correct events
