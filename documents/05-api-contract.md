# Stage 5 — API Contract

## Design Rules Applied

- Path parameters for identity, query parameters for filtering
- GET requests never have side effects
- Sensitive tokens always in request body, never in query strings
- PUT /predictions is an upsert — one prediction per user per match
- Invite link = Vue frontend URL (/join?token=xxx), API call = POST with token in body

---

## Auth

| Method | Path | Description |
|--------|------|-------------|
| GET | /auth/login | Serve login page (Vue SPA) |
| GET | /auth/callback/google | Handle OAuth callback, issue JWT |
| POST | /auth/logout | Invalidate refresh token |
| POST | /auth/refresh | Exchange refresh token → new access token |

### Invite Link Flow

```
Vue frontend URL (shareable):  https://app.com/join?token=xxx
                                        ↓
                               Vue reads token from URL
                                        ↓
API call (secure):             POST /leaderboards/join
                               body: { "invite_token": "xxx" }
```

---

## Leaderboards

| Method | Path | Description |
|--------|------|-------------|
| POST | /leaderboards | Create a new leaderboard |
| GET | /leaderboards | List leaderboards I belong to |
| GET | /leaderboards/{uuid} | Get one leaderboard |
| POST | /leaderboards/{uuid}/invite | Generate a signed invite token |
| POST | /leaderboards/join | Join via invite (token in request body) |

---

## Matches

| Method | Path | Description |
|--------|------|-------------|
| GET | /matches | List all matches (?status=scheduled\|live\|finished) |
| GET | /matches/{uuid} | Get one match with details |

---

## Predictions

| Method | Path | Description |
|--------|------|-------------|
| PUT | /predictions | Upsert my prediction (body: match_id, value) |
| GET | /leaderboards/{uuid}/matches/{uuid}/predictions | List predictions for a match in leaderboard context |

### Why GET /leaderboards/{uuid}/matches/{uuid}/predictions?

Path parameters for identity. The URL tells a complete story:
"give me the predictions, for this match, within this leaderboard."
A new developer reads the URL and immediately understands the context.

---

## JWT Structure

```json
{
  "sub": "user-uuid",
  "email": "user@company.com",
  "name": "User Name",
  "provider": "google",
  "exp": 1234567890
}
```

Access token: short-lived (15 min)
Refresh token: long-lived (7 days), stored in httpOnly cookie
