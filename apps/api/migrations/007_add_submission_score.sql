-- 007_add_submission_score.sql
ALTER TABLE public.submissions 
ADD COLUMN IF NOT EXISTS score INTEGER DEFAULT 0,
ADD COLUMN IF NOT EXISTS time_ms_total INTEGER DEFAULT 0;

-- Index for leaderboard queries
CREATE INDEX IF NOT EXISTS idx_submissions_challenge_score 
  ON public.submissions (challenge_id, score DESC, created_at ASC) 
  WHERE status = 'accepted';

CREATE INDEX IF NOT EXISTS idx_submissions_user_challenge 
  ON public.submissions (user_id, challenge_id, score DESC);
