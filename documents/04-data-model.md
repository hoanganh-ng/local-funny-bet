# Stage 4 — Data Model

## Tables

### users

```sql
CREATE TABLE users (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email      VARCHAR NOT NULL UNIQUE,
  name       VARCHAR NOT NULL,
  avatar_url VARCHAR,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### tournaments ← global expansion backdoor, present from day one

```sql
CREATE TABLE tournaments (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name        VARCHAR NOT NULL,      -- "FIFA World Cup 2026"
  season      VARCHAR NOT NULL,      -- "2026"
  logo_url    VARCHAR,
  status      VARCHAR NOT NULL,      -- scheduled|active|finished
  external_id VARCHAR                -- from football-data.org
);
```

### leaderboards

```sql
CREATE TABLE leaderboards (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name       VARCHAR NOT NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### leaderboard_members (pivot)

```sql
CREATE TABLE leaderboard_members (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  leaderboard_id UUID NOT NULL REFERENCES leaderboards(id),
  user_id        UUID NOT NULL REFERENCES users(id),
  role           VARCHAR NOT NULL,  -- owner | member
  joined_at      TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(leaderboard_id, user_id)
);
```

### matches

```sql
CREATE TABLE matches (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tournament_id UUID NOT NULL REFERENCES tournaments(id),
  home_team     VARCHAR NOT NULL,
  away_team     VARCHAR NOT NULL,
  home_score    INT,               -- NULL until finished
  away_score    INT,               -- NULL until finished
  kickoff_at    TIMESTAMP NOT NULL,
  status        VARCHAR NOT NULL,  -- scheduled|live|finished
  external_id   VARCHAR            -- from football-data.org
);
```

### predictions

```sql
CREATE TABLE predictions (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id),
  match_id   UUID NOT NULL REFERENCES matches(id),
  value      VARCHAR NOT NULL,    -- home_win|draw|away_win
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, match_id)       -- one prediction per user per match, enforced at DB level
);
```

## Key Design Decisions

- **No invite_links table** — tokens are stateless HMAC-signed payloads (leaderboard_id + expires_at)
- **home/away_score nullable** — NULL means unknown, not zero
- **UNIQUE(user_id, match_id)** — enforced at DB level, not just application logic
- **tournaments from day one** — global expansion backdoor, WC2026 is just the first row

## The Key Query — Leaderboard Score

```sql
SELECT u.id, u.name, COUNT(*) AS points
FROM predictions p
JOIN matches m ON p.match_id = m.id
JOIN leaderboard_members lm
  ON lm.user_id = p.user_id AND lm.leaderboard_id = $1
JOIN users u ON u.id = p.user_id
WHERE m.status = 'finished'
  AND (
    (m.home_score > m.away_score AND p.value = 'home_win') OR
    (m.home_score = m.away_score AND p.value = 'draw')     OR
    (m.home_score < m.away_score AND p.value = 'away_win')
  )
GROUP BY u.id, u.name
ORDER BY points DESC;
```

Scores are always calculated on the fly.
Never cache a score in a column — it will go stale.

## SA Rule
>
> Design your schema around your queries, not your intuition.
> The most important query is: rank all users in a leaderboard by correct predictions.
> Everything in the schema serves that query.
