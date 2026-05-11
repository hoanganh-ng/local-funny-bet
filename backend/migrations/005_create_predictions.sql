CREATE TABLE predictions (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id),
  match_id   UUID NOT NULL REFERENCES matches(id),
  value      VARCHAR NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, match_id)
);
