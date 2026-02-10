-- 001_create_challenges.sql
CREATE TABLE IF NOT EXISTS public.challenges (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  status VARCHAR(16) NOT NULL CHECK (status IN ('draft','published','archived')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Actualiza updated_at en cada UPDATE
CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_challenges_updated_at ON public.challenges;

CREATE TRIGGER trg_challenges_updated_at
BEFORE UPDATE ON public.challenges
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- Índice útil para listados recientes
CREATE INDEX IF NOT EXISTS idx_challenges_created_at
  ON public.challenges (created_at DESC);
