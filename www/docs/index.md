---
layout: home

hero:
  name: Mist
  text: Self-hostable Platform-as-a-Service
  tagline: Deploy applications with ease. Lightweight, fast, and built for developers.
  image:
    src: /mist.png
    alt: Mist Logo
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/corecollectives/mist

features:
  - icon: 🚀
    title: Easy Deployment
    details: Deploy Docker-based applications directly from Git repositories with automatic builds and real-time monitoring.
  
  - icon: 🔗
    title: Git Integration
    details: Seamlessly connect with GitHub repositories. Auto-deploy on push with webhooks.
  
  - icon: 🌐
    title: Domain & SSL Management
    details: Configure custom domains with Traefik reverse proxy integration. SSL automation coming soon.
  
  - icon: 📊
    title: Real-time Monitoring
    details: Watch deployment logs in real-time, monitor system metrics, and track container status with WebSocket-powered updates.
  
  - icon: 🔧
    title: Environment Variables
    details: Manage build-time and runtime variables with bulk import support for easy configuration.
  
  - icon: 👥
    title: Multi-User & Projects
    details: Organize applications into projects with team collaboration and role-based access control.
  
  - icon: 🪶
    title: Ultra Lightweight
    details: Single Go binary with embedded SQLite database. No external dependencies like Redis or PostgreSQL required.
  
  - icon: 🔒
    title: Secure by Default
    details: JWT authentication, HTTP-only cookies, audit logging, and comprehensive security best practices.
  
  - icon: ⚡
    title: Real-time Everything
    details: WebSocket-first architecture provides instant feedback for logs, metrics, and deployment status.
---

## Why Mist?

Mist is a **lightweight, self-hostable Platform-as-a-Service** designed for developers and small teams who want the simplicity of Heroku with the control of self-hosting.

### Key Advantages

- **🪶 Minimal Resource Usage** - Runs efficiently on a single VPS with just Docker installed
- **⚡ Built for Speed** - Go backend with WebSocket-powered real-time updates
- **🔧 Developer Friendly** - Intuitive UI, comprehensive API, and CLI tool (coming soon)
- **🔐 Production Ready** - JWT auth, audit logs, role-based access control
- **💰 Cost Effective** - No external services required, runs on any cloud provider

## Quick Start

```bash
# Install Mist on your server
curl -sSL https://raw.githubusercontent.com/corecollectives/mist/main/install.sh | bash

# Access the dashboard
# Navigate to http://your-server-ip:3000
```

## What's Included

### ✅ Available Now

- Docker-based application deployments
- GitHub integration with webhooks
- Real-time deployment logs and monitoring
- Environment variable management (with bulk import)
- Custom domain configuration
- Project organization and team collaboration
- User authentication and role-based access
- System metrics monitoring (CPU, RAM, Disk)
- Audit logging

### 🚧 Coming Soon

- Let's Encrypt SSL automation
- Database provisioning (PostgreSQL, MySQL, Redis, MongoDB)
- Deployment rollback
- Blue-green deployments
- Email notifications and webhooks
- CLI tool for command-line deployment
- Preview environments for pull requests
- One-click app templates

## Community

- [GitHub Repository](https://github.com/corecollectives/mist)
- [Report Issues](https://github.com/corecollectives/mist/issues)
- [Discussions](https://github.com/corecollectives/mist/discussions)

## License

Mist is open source software licensed under the [MIT License](https://github.com/corecollectives/mist/blob/main/LICENSE).
