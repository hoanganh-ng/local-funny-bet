CREATE TABLE tournaments (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name        VARCHAR NOT NULL,
  season      VARCHAR NOT NULL,
  logo_url    VARCHAR,
  status      VARCHAR NOT NULL,
  external_id VARCHAR
);
