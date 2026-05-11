# Stage 1 — Requirements Document

## WC2026 Internal Prediction App

---

## Functional Requirements

| # | Requirement |
|---|---|
| F1 | Users authenticate via Google OAuth (company email domain only) |
| F2 | Users predict each match outcome: Home Win / Draw / Away Win |
| F3 | Predictions can be edited until kickoff, then locked |
| F4 | Correct prediction = 1 point, added to all leaderboards the user belongs to |
| F5 | Users can belong to multiple leaderboards |
| F6 | Each leaderboard has its own ranking |
| F7 | Match results are fetched automatically via external football API |
| F8 | App is usable on mobile browsers (responsive UI) |
| F9 | Users join a leaderboard via invite link only (signed token, expires 10 min) |
| F10 | Leaderboard owner generates the invite link |

## Non-Functional Requirements

| # | Requirement |
|---|---|
| N1 | Max ~100 users, ~50 concurrent at peak (match day) |
| N2 | Low infra cost — single-developer maintainable |
| N3 | Google OAuth only — no custom password auth |
| N4 | Graceful handling of football API fetch failures |
| N5 | Real-time leaderboard updates via WebSocket |

## Out of Scope (v1)

- Admin dashboard for manual result entry
- Push notifications
- Complex scoring (exact score, goal difference)
- Payment / prize tracking

## Global Expansion Backdoors (planted in v1)

- `tournaments` table exists from day one
- `AllowedEmailDomains` config-driven (not hardcoded)
- `AuthProvider` interface for adding new auth methods
- `AdSlot` Vue component (inactive, env-flag controlled)
- Redis service commented in docker-compose
- `// TODO(global):` markers at extension points
