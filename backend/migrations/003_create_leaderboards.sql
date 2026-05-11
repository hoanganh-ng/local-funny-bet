CREATE TABLE leaderboards (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name       VARCHAR NOT NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE leaderboard_members (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  leaderboard_id UUID NOT NULL REFERENCES leaderboards(id),
  user_id        UUID NOT NULL REFERENCES users(id),
  role           VARCHAR NOT NULL,
  joined_at      TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(leaderboard_id, user_id)
);
