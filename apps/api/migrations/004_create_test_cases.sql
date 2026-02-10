-- 004_create_test_cases.sql
CREATE TABLE IF NOT EXISTS public.test_cases (
  id UUID PRIMARY KEY,
  challenge_id UUID NOT NULL REFERENCES public.challenges(id) ON DELETE CASCADE,
  name VARCHAR(100) NOT NULL,
  input TEXT NOT NULL,
  expected_output TEXT NOT NULL,
  is_sample BOOLEAN DEFAULT false,
  points INTEGER DEFAULT 10,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for faster lookups by challenge
CREATE INDEX IF NOT EXISTS idx_test_cases_challenge 
  ON public.test_cases (challenge_id);

-- Index for filtering sample cases
CREATE INDEX IF NOT EXISTS idx_test_cases_sample 
  ON public.test_cases (challenge_id, is_sample);
