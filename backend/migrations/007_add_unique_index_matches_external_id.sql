-- Deduplication anchor for football-data.org matches.
-- Without this index, ON CONFLICT has no conflict target and duplicates
-- will silently accumulate on every scheduler tick.
CREATE UNIQUE INDEX idx_matches_external_id
    ON matches (external_id)
    WHERE external_id IS NOT NULL;
