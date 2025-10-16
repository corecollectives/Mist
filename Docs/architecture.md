# Project Architecture Overview

This document describes the overall architecture for the Minimal Self-Hostable Platform-as-a-Service (PaaS).  
It is designed to be lightweight, resource-efficient, and operable on a single VPS with minimal dependencies.

---

## Core Vision

- A monolithic Go backend compiled into a single binary.
- Frontend built with React and Vite providing a modern UI.
- Embedded SQLite database for persistent application state.
- No external dependencies like Redis or message brokers.
- Users authenticate to create projects, deploy apps from Git repos, and manage deployments.
- Deployment queue and workers implemented as Go goroutines.
- Docker socket interaction for building and managing app containers.
- Dynamic reverse proxy configuration for domain-based routing.

---

## System Components

- **Backend API Server**  
  Handles RESTful HTTP API calls, user authentication, project/app management, deployments, live logs and webhook interactions (e.g., GitHub).

- **Frontend SPA**  
  React + Vite client served by the backend; interacts via the REST API.

- **SQLite Database**  
  File-based embedded database storing users, projects, apps, deployments, configs, and domains.

- **In-Memory Deployment Queue**  
  Simple non-persistent queue managed using Go channels.

- **Deployment Daemon Worker**  
  - goroutine that processes queued deployment jobs:  
  - Clones Git repos  
  - Executes build/start commands  
  - Builds Docker images and manages containers via Docker socket  
  - Updates reverse proxy settings for routing

- **Reverse Proxy (Traefik)**  
  Routes incoming HTTP(S) requests to the correct app container based on configured domains and SSL.

---

## Deployment Workflow

1. User creates a project and app via the frontend.
2. User configures Git repo, build/start commands, environment variables, and domains.
3. User triggers a deployment (or via github webhooks).
4. Deployment job is added to the in-memory queue.
5. Deployment worker processes the job:
    - Clones the Git repo.
    - Runs build commands.
    - Builds a Docker image.
    - Starts the container with environment variables (remove previous one if exists).
    - Updates Traefik configuration for routing.
6. User can view live logs and deployment status in the frontend.
7. User can manage apps, domains, and view deployment history.

---

## Deployment Strategy

- **Production**
  - Single VPS (e.g., DigitalOcean, Linode) with Docker installed.
  - Deploy the Go binary and SQLite database on the host.
  - Run Traefik as a Docker container for reverse proxying.
  - Use a process manager (e.g., systemd) to ensure the Go backend is always running.

- **Development**
  - Run the Go backend locally with a local SQLite file.
  - Use Docker Compose to run Traefik and any other services needed for testing.
  - Frontend can be run with Vite's development server for hot-reloading.


