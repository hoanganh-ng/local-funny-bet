# Stage 2 — High-Level Architecture

## Decision: Monolith

For 50 users on an office server, a monolith is the correct call.
Microservices would be showing off, not thinking.

## Architecture Diagram

```
┌─────────────────────────────────────┐
│         Vue SPA (Browser)           │
│       REST ↕        WebSocket ↕     │
└──────────┬──────────────┬───────────┘
           │              │
┌──────────▼──────────────▼───────────────────────────┐
│                Go Backend (Monolith)                 │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────┐ │
│  │  REST API   │  │  WebSocket   │  │Job Scheduler│ │
│  │ Predictions │  │     Hub      │  │Fetch results│ │
│  │   & rooms   │  │ Live updates │  │             │ │
│  └──────┬──────┘  └──────────────┘  └──────┬──────┘ │
└─────────┼──────────────────────────────────┼─────────┘
          │                                  │
  ┌───────▼──────┐  ┌─────────────┐  ┌──────▼──────┐
  │  PostgreSQL  │  │Google OAuth │  │Football API │
  │  (primary    │  │   (auth)    │  │football-data│
  │  database)   │  │             │  │    .org     │
  └──────────────┘  └─────────────┘  └─────────────┘
```

## Key Architectural Decisions

| Decision | Choice | Reason |
|----------|--------|--------|
| Pattern | Monolith | 50 users, single developer, no distributed overhead |
| Real-time | WebSocket | Leaderboard pushed to clients on match result |
| Invite | Stateless HMAC token | No invite_links table, no storage needed |
| Results | Job Scheduler polls API | Automated, no admin needed |

## Data Flow — Real-time Leaderboard

```
football-data.org
      ↓ (job scheduler polls after kickoff)
PostgreSQL (match result written)
      ↓ (trigger broadcast)
WebSocket Hub
      ↓ (push to all connected clients)
Vue SPA (leaderboard re-renders)
```
