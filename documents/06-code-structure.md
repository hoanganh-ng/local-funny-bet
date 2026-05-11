# Stage 6 — Code Structure

## Architectural Pattern: Clean Architecture + DDD

Organise by domain, not by technical role.

```
❌ Technical role (avoid)     ✅ Domain-driven (used here)
/controllers                  /leaderboard
/models                       /match
/services                     /prediction
/repositories                 /auth
```

---

## Go Backend

```
wc2026/backend/
├── cmd/
│   └── api/
│       └── main.go              # Entry point — wires everything together
└── internal/
    ├── domain/                  # Pure entities, interfaces, domain errors
    │   ├── leaderboard/         # Leaderboard entity + repository interface
    │   ├── match/               # Match entity + repository interface
    │   ├── prediction/          # Prediction entity + CalculatePoint() logic
    │   ├── tournament/          # Tournament entity + repository interface
    │   └── user/                # User entity
    ├── application/             # Use cases — orchestrates domain + infrastructure
    │   ├── leaderboard/         # CreateLeaderboard, JoinLeaderboard, GenerateInvite
    │   ├── match/               # ListMatches, FetchFromAPI
    │   └── prediction/          # UpsertPrediction, ListPredictions, GetLeaderboardScores
    ├── infrastructure/          # Implements domain interfaces
    │   ├── postgres/            # SQL repository implementations
    │   ├── football/            # football-data.org API client
    │   └── google/              # Google OAuth client
    ├── adapters/
    │   ├── http/                # HTTP handlers (one file per domain)
    │   │   ├── auth.go
    │   │   ├── leaderboard.go
    │   │   ├── match.go
    │   │   └── prediction.go
    │   └── ws/                  # WebSocket hub — manages connections + broadcasts
    │       └── hub.go
    ├── config/                  # Env vars, app config struct
    │   └── config.go
    └── middleware/              # JWT validation, logging, CORS
        ├── auth.go
        ├── cors.go
        └── logger.go
```

## Dependency Rule (never break this)

```
domain/         → imports nothing internal
application/    → imports domain only
infrastructure/ → imports domain only
adapters/       → imports application only

If domain/ ever imports from infrastructure/, the architecture is broken.
Go's compiler won't catch this — your discipline must.
```

---

## Vue Frontend

```
wc2026/frontend/src/
├── App.vue                      # Root component
├── main.js                      # App entry point
├── router/                      # Vue Router — route definitions
│   └── index.js
├── views/                       # Page-level components (one per route)
│   ├── LoginView.vue
│   ├── HomeView.vue
│   ├── LeaderboardView.vue
│   └── MatchView.vue
├── components/
│   ├── base/                    # Primitives: Button, Input, Modal, Card
│   ├── common/                  # Shared across domains: MatchCard, UserAvatar
│   ├── leaderboard/             # LeaderboardTable, LeaderboardHeader, ScoreRow
│   ├── match/                   # MatchList, MatchRow, MatchStatusBadge
│   └── prediction/              # PredictionForm, PredictionHistory, OutcomePicker
├── composables/                 # Reusable logic (Vue 3 Composition API)
│   ├── useAuth.js
│   ├── useWebSocket.js
│   └── useLeaderboard.js
├── services/                    # All API calls live here, never in components
│   ├── leaderboard.service.js
│   ├── match.service.js
│   └── prediction.service.js
├── store/                       # Pinia stores — global state
│   ├── auth.store.js
│   └── leaderboard.store.js
└── assets/
    └── css/
        ├── variables.css        # Design tokens: colors, spacing, typography
        └── global.css           # Reset, base styles
```

## Frontend Rules

- Composition API with `<script setup>` only — no Options API
- No inline styles — all values from variables.css tokens
- No direct fetch() in components — always through services/
- One composable per domain concept
- Component naming: PascalCase, domain-prefixed (LeaderboardTable, MatchCard)

---

## Full Project Root

```
wc2026/
├── backend/
├── frontend/
├── docker-compose.yml
├── .env                         # Never committed — local only
├── .gitignore
├── CLAUDE.md                    # Claude Code instructions
└── .github/
    └── workflows/
        └── deploy.yml
```
