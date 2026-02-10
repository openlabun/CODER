-- 002_create_submissions.sql
CREATE TABLE IF NOT EXISTS public.submissions (
  id UUID PRIMARY KEY,
  challenge_id UUID NOT NULL,
  user_id TEXT NOT NULL,
  code TEXT NOT NULL,
  language VARCHAR(32) NOT NULL,
  status VARCHAR(16) NOT NULL CHECK (status IN ('queued','running','accepted','wrong_answer','error')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- trigger updated_at
CREATE OR REPLACE FUNCTION public.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_submissions_updated_at ON public.submissions;

CREATE TRIGGER trg_submissions_updated_at
BEFORE UPDATE ON public.submissions
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE INDEX IF NOT EXISTS idx_submissions_challenge ON public.submissions (challenge_id);
CREATE INDEX IF NOT EXISTS idx_submissions_created ON public.submissions (created_at DESC);
