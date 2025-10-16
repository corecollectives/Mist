# Data Handling Overview

This document outlines how data is structured, stored, and managed in the Minimal Self-Hostable PaaS.  
The system prioritizes simplicity and performance by using **SQLite** as a local embedded database, accessed directly from the Go backend.

---

## Core Principles

- **Embedded & Self-contained**  
  The platform uses SQLite, stored as a single file (`mist.db`) on the host’s filesystem (in /var/lib/mist/mist.db).  
  This ensures that mist can run on a single VPS with no external dependencies.

- **No External Services**  
  No Redis, message brokers, or additional databases are required.

- **Lightweight Schema Management**  
  Database migrations are handled with `golang-migrate`, ensuring versioned, repeatable schema updates.

---

## Entities & Relationships

- **Users**
  - individual user accounts

- **Projects**
  - one project can have multiple users (one owner and multiple collaborators)
  - one project can have multiple apps/services running

- **Apps/Services**
  - each app is associated with a single project
  - each app has its own environment variables, domains, and deployment history

- **Domains**
  - each domain is linked to a specific app
  - domains can be custom (user-provided) or auto-generated
  - SSL certificates are managed per domain

- **Environment Variables**
  - key-value pairs associated with apps

- **Deployments**
  - each deployment is linked to a specific app
  - stores metadata about the deployment (timestamp, status, logs)


