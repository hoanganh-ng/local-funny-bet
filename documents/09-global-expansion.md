# Global Expansion — Backdoors Planted in v1

## Overview

The WC2026 office version is built with deliberate extension points
so that going global requires config changes and new files — not rewrites.

**Principle: Open/Closed — open for extension, closed for modification.**

---

## Backdoor 1 — tournaments Table

Present from migration #001. WC2026 is just the first row.

```sql
CREATE TABLE tournaments (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name        VARCHAR NOT NULL,      -- "FIFA World Cup 2026"
  season      VARCHAR NOT NULL,      -- "2026"
  logo_url    VARCHAR,
  status      VARCHAR NOT NULL,
  external_id VARCHAR
);
```

Going global = insert more rows. No schema change.

---

## Backdoor 2 — Config-Driven Email Domain

```go
// config/config.go
type Config struct {
  AllowedEmailDomains []string  // ["yourcompany.com"] → [] = any domain
  AuthProviders       []string  // ["google"] → ["google","email"]
}
```

```bash
# Office version
ALLOWED_EMAIL_DOMAINS=yourcompany.com

# Global version
ALLOWED_EMAIL_DOMAINS=
```

Going global = change one env var. Zero code change.

---

## Backdoor 3 — AuthProvider Interface

```go
// domain/auth/provider.go
type Provider interface {
  Name()        string
  VerifyToken() (*User, error)
}

// infrastructure/google/provider.go  ← exists now
// infrastructure/email/provider.go   ← add this file for global
```

The /auth/callback/{provider} route plugs straight into this.
Adding email auth = one new file. No handler changes.

---

## Backdoor 4 — AdSlot Vue Component

```vue
<!-- components/base/AdSlot.vue -->
<template>
  <div v-if="adsEnabled" class="ad-slot" :data-position="position">
    <!-- AdSense script drops here in global version -->
  </div>
</template>

<script setup>
defineProps<{ position: 'top' | 'sidebar' | 'between-matches' }>()
const adsEnabled = import.meta.env.VITE_ADS_ENABLED === 'true'
</script>
```

```bash
# Office version
VITE_ADS_ENABLED=false

# Global version
VITE_ADS_ENABLED=true
```

Component is placed in layouts now. Invisible. One env var to activate.

---

## Backdoor 5 — Redis in docker-compose

```yaml
# Uncomment when going global
# redis:
#   image: redis:7-alpine
#   restart: unless-stopped
```

And in leaderboard service:

```go
func (s *LeaderboardService) GetScores(id string) ([]Score, error) {
  // TODO(global): check Redis cache here before hitting DB
  return s.repo.CalculateScores(id)
}
```

---

## Global Expansion Phases

### Phase 1 — Ship WC2026 (current)

- Office PC hosting
- Google OAuth, company email domain only
- ~50 users, one tournament
- Ads disabled

### Phase 2 — Go Global (before/during WC2026)

- Migrate to Railway / Render (~$20/month)
- Remove email domain restriction (one env var)
- Add magic link email auth (one new file)
- Enable AdSense (one env var + GDPR consent banner)
- Buy a domain, set up proper SSL

### Phase 3 — Scale (if traction)

- Uncomment Redis, add caching to leaderboard query
- CDN for static assets
- Upgrade football-data.org to paid plan
- Add privacy-friendly analytics (Plausible / Fathom)

---

## Cost of Backdoors (planted during v1 build)

| Backdoor | Extra effort in v1 |
|----------|--------------------|
| tournaments table | 30 min |
| Config-driven email domain | 15 min |
| AuthProvider interface | 1 hour |
| AdSlot component | 20 min |
| Redis in docker-compose | 5 min |
| TODO(global): comments | 10 min |
| **Total** | **~2.5 hours** |

2.5 hours now saves weeks of refactoring later.
