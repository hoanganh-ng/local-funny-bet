# Stage 3 — Technology Stack

## Locked Decisions

| Layer | Choice | Rationale |
|-------|--------|-----------|
| Frontend framework | Vue 3 (Composition API) | Team familiarity, productive |
| Styling | Pure CSS + variables.css | No framework overhead, full control |
| HTTP router | Go net/http stdlib (1.22+) | Path params supported natively, zero dependencies |
| WebSocket | nhooyr.io/websocket | Modern, context-aware, successor to Gorilla |
| Database access | Raw SQL | Full visibility into queries and indexes |
| Auth | Google OAuth 2.0 → JWT | Company email enforcement, mobile-ready, future-proof |
| Football data | football-data.org (free tier) | Clean REST, WC2026 coverage, 10 req/min free |
| Deployment | Docker + GitHub Actions | Reproducible builds, automated deploy on tag push |
| State management | Pinia | Official Vue 3 standard, replaces Vuex |

## Key Principles Applied

- **Familiarity over novelty** — Vue chosen because team knows it
- **Stdlib over framework** — Go 1.22+ net/http handles path params natively
- **Raw SQL over ORM** — full visibility, no hidden queries, no magic
- **JWT over sessions** — stateless, mobile-ready, no server-side session store

## Football API Comparison

| API | Free tier | WC2026 | Notes |
|-----|-----------|--------|-------|
| football-data.org ✅ | 10 req/min | Yes | Clean REST, clear docs |
| API-Football | 100 req/day | Yes | More data, RapidAPI friction |
| TheSportsDB | Unlimited | Basic only | Less reliable for live scores |

## What NOT to Use

- ~~GORM~~ — ORM hides complexity that bites at scale
- ~~Gin / Echo~~ — stdlib is sufficient, no added dependencies
- ~~Vuex~~ — replaced by Pinia in Vue 3
- ~~Gorilla WebSocket~~ — maintenance mode, use nhooyr.io/websocket
