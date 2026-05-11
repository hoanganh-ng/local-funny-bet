-- Development seed data
-- DO NOT run in production

-- Insert sample tournament
INSERT INTO tournaments (id, name, season, logo_url, status, external_id)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  'FIFA World Cup 2026',
  '2026',
  NULL,
  'scheduled',
  NULL
);

-- Insert sample matches
INSERT INTO matches (tournament_id, home_team, away_team, home_score, away_score, kickoff_at, status, external_id)
VALUES
  (
    '11111111-1111-1111-1111-111111111111',
    'Brazil',
    'Argentina',
    NULL,
    NULL,
    '2026-06-11 15:00:00',
    'scheduled',
    NULL
  ),
  (
    '11111111-1111-1111-1111-111111111111',
    'Germany',
    'France',
    NULL,
    NULL,
    '2026-06-12 18:00:00',
    'scheduled',
    NULL
  ),
  (
    '11111111-1111-1111-1111-111111111111',
    'Spain',
    'Italy',
    NULL,
    NULL,
    '2026-06-13 21:00:00',
    'scheduled',
    NULL
  );
