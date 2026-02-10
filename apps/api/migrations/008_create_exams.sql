-- 008_create_exams.sql

-- Table: exams
CREATE TABLE IF NOT EXISTS public.exams (
  id UUID PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  course_id UUID NOT NULL REFERENCES public.courses(id) ON DELETE CASCADE,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  duration_minutes INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table: exam_challenges (Many-to-Many with points)
CREATE TABLE IF NOT EXISTS public.exam_challenges (
  exam_id UUID NOT NULL REFERENCES public.exams(id) ON DELETE CASCADE,
  challenge_id UUID NOT NULL REFERENCES public.challenges(id) ON DELETE CASCADE,
  points INTEGER NOT NULL DEFAULT 100,
  "order" INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (exam_id, challenge_id)
);

-- Add exam_id to submissions
ALTER TABLE public.submissions 
ADD COLUMN IF NOT EXISTS exam_id UUID REFERENCES public.exams(id) ON DELETE SET NULL;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_exams_course ON public.exams (course_id);
CREATE INDEX IF NOT EXISTS idx_submissions_exam ON public.submissions (exam_id);

-- Trigger for updated_at
DROP TRIGGER IF EXISTS trg_exams_updated_at ON public.exams;

CREATE TRIGGER trg_exams_updated_at
BEFORE UPDATE ON public.exams
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
