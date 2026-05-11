CREATE TABLE users (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email      VARCHAR NOT NULL UNIQUE,
  name       VARCHAR NOT NULL,
  avatar_url VARCHAR,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
