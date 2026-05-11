# World Cup 2026 Prediction App

A real-time prediction and leaderboard platform for the FIFA World Cup 2026.

## Tech Stack

- **Backend**: Go 1.22 (stdlib net/http), PostgreSQL, WebSockets
- **Frontend**: Vue 3 (Composition API), Pinia, Vue Router
- **Infrastructure**: Docker, nginx

## Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)
- Node.js 18+ (for local development)

## Quick Start

1. Clone the repository
2. Copy `.env.example` to `.env` and fill in your secrets
3. Start all services:
   ```bash
   docker-compose up -d
   ```
4. Access the app at `http://localhost`

## Development

### Backend

```bash
cd backend
go mod tidy
go run cmd/api/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Project Structure

```
backend/
  cmd/api/              - Application entry point
  internal/
    domain/             - Entities and interfaces (no external imports)
    application/        - Use cases (imports domain only)
    infrastructure/     - External services (imports domain only)
    adapters/           - HTTP/WS handlers (imports application only)
    config/             - Configuration
    middleware/         - HTTP middleware

frontend/
  src/
    components/         - Vue components
    composables/        - Reusable logic
    services/           - API client
    store/              - Pinia stores
    views/              - Page components
    router/             - Vue Router config
    assets/             - CSS and static assets
```

## Clean Architecture Layers

- **domain**: Core business logic and entities
- **application**: Use cases and business workflows
- **infrastructure**: External services (database, APIs)
- **adapters**: HTTP/WebSocket handlers

**Dependency rule**: Arrows point inward only.

## License

MIT
