# Stage 7 — Infrastructure & Deployment

## Overview

| Concern | Solution |
|---------|----------|
| Containers | 3 — frontend (nginx), backend (Go), db (PostgreSQL) |
| Deployment | GitHub Actions self-hosted runner on office PC |
| Crash recovery | Docker restart: unless-stopped |
| Data persistence | Named Docker volume for PostgreSQL |
| Secrets | GitHub Actions secrets → .env file on deploy |

---

## docker-compose.yml

```yaml
services:
  frontend:
    build: ./frontend
    # nginx: serves built Vue SPA + reverse proxies /api/* to backend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

  backend:
    build: ./backend
    # Multi-stage Dockerfile: build Go binary → minimal runtime image
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - JWT_SECRET=${JWT_SECRET}
      - JWT_REFRESH_SECRET=${JWT_REFRESH_SECRET}
      - GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
      - GOOGLE_CLIENT_SECRET=${GOOGLE_CLIENT_SECRET}
      - FOOTBALL_API_KEY=${FOOTBALL_API_KEY}
      - ALLOWED_EMAIL_DOMAINS=${ALLOWED_EMAIL_DOMAINS}  # global backdoor
    restart: unless-stopped
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_DB=${DB_NAME}
      - POSTGRES_USER=${DB_USER}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  # redis:                        # uncomment when going global
  #   image: redis:7-alpine
  #   restart: unless-stopped

volumes:
  postgres_data:                  # survives docker-compose down and rebuilds

networks:
  default:
    name: wc2026-internal
```

## .github/workflows/deploy.yml

```yaml
name: Build and Deploy
on:
  push:
    tags:
      - 'v*.*.*'

jobs:
  deploy:
    runs-on: self-hosted             # your office PC is the runner

    steps:
      - uses: actions/checkout@v4

      - name: Write env file
        run: |
          echo "DATABASE_URL=${{ secrets.DATABASE_URL }}" > .env
          echo "JWT_SECRET=${{ secrets.JWT_SECRET }}" >> .env
          echo "JWT_REFRESH_SECRET=${{ secrets.JWT_REFRESH_SECRET }}" >> .env
          echo "GOOGLE_CLIENT_ID=${{ secrets.GOOGLE_CLIENT_ID }}" >> .env
          echo "GOOGLE_CLIENT_SECRET=${{ secrets.GOOGLE_CLIENT_SECRET }}" >> .env
          echo "FOOTBALL_API_KEY=${{ secrets.FOOTBALL_API_KEY }}" >> .env
          echo "ALLOWED_EMAIL_DOMAINS=${{ secrets.ALLOWED_EMAIL_DOMAINS }}" >> .env
          echo "DB_NAME=${{ secrets.DB_NAME }}" >> .env
          echo "DB_USER=${{ secrets.DB_USER }}" >> .env
          echo "DB_PASSWORD=${{ secrets.DB_PASSWORD }}" >> .env

      - name: Deploy
        run: |
          docker-compose pull
          docker-compose up --build -d
          docker-compose exec -T backend ./migrate up

      - name: Cleanup old images
        run: docker image prune -f
```

## How Secrets Flow

```
GitHub Actions Secrets (source of truth)
          ↓  (workflow writes on every deploy)
.env file on office PC
          ↓  (Docker Compose reads at startup)
Running containers (environment variables)
```

.env is in .gitignore — secrets never touch git.

## Key Decisions

- **restart: unless-stopped** — auto-recovery on crash/reboot, respects manual docker stop
- **Self-hosted runner** — runner polls GitHub outbound only, no inbound firewall rules needed
- **Named volume** — anything you can't regenerate must live outside the container
- **Multi-stage Go Dockerfile** — build stage compiles binary, runtime stage is minimal (scratch or alpine)
- **nginx double duty** — serves static Vue files AND reverse proxies /api/* to Go backend

## Go Dockerfile (multi-stage)

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api ./cmd/api

# Runtime stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/migrate .
EXPOSE 8080
CMD ["./api"]
```
