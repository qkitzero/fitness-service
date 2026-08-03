CREATE TABLE rank_standards (
  rank VARCHAR(1) PRIMARY KEY CHECK (rank IN ('A', 'B', 'C', 'D', 'E')),
  z_score_min NUMERIC(4, 2),
  z_score_max NUMERIC(4, 2),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  CHECK (z_score_min IS NOT NULL OR z_score_max IS NOT NULL),
  CHECK (z_score_min IS NULL OR z_score_max IS NULL OR z_score_min < z_score_max)
);
