CREATE TABLE matches (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tournament_id UUID NOT NULL REFERENCES tournaments(id),
  home_team     VARCHAR NOT NULL,
  away_team     VARCHAR NOT NULL,
  home_score    INT,
  away_score    INT,
  kickoff_at    TIMESTAMP NOT NULL,
  status        VARCHAR NOT NULL,
  external_id   VARCHAR
);
