-- schema.sql

CREATE TABLE films (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  year INTEGER,
  source_key TEXT NOT NULL,          -- s3://raw bucket key
  transcoded_prefix TEXT,            -- prefix in hls bucket
  status TEXT NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
  -- poster
  -- thumbnails
  -- duration
);
