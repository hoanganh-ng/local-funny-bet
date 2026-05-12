-- Add unique constraint on external_id to prevent duplicate tournaments
ALTER TABLE tournaments
  ADD CONSTRAINT tournaments_external_id_unique UNIQUE (external_id);

-- Insert required reference data: FIFA World Cup 2026
-- This is NOT dev-only seed data — the app requires this row to function.
-- ON CONFLICT ensures idempotency: safe to run multiple times.
INSERT INTO tournaments (id, name, season, logo_url, status, external_id)
VALUES (
  gen_random_uuid(),
  'FIFA World Cup 2026',
  '2026',
  NULL,
  'scheduled',
  'WC'
)
ON CONFLICT (external_id) DO NOTHING;
