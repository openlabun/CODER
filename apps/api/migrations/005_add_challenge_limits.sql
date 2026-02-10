-- 005_add_challenge_limits.sql
ALTER TABLE public.challenges 
ADD COLUMN IF NOT EXISTS time_limit INTEGER DEFAULT 1500,
ADD COLUMN IF NOT EXISTS memory_limit INTEGER DEFAULT 256,
ADD COLUMN IF NOT EXISTS difficulty VARCHAR(16) DEFAULT 'medium' CHECK (difficulty IN ('easy', 'medium', 'hard')),
ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}';

-- Index for filtering by difficulty
CREATE INDEX IF NOT EXISTS idx_challenges_difficulty 
  ON public.challenges (difficulty);
