-- 006_create_courses.sql
CREATE TABLE IF NOT EXISTS public.courses (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(50) NOT NULL,
  period VARCHAR(20) NOT NULL,
  group_number INTEGER NOT NULL,
  professor_id UUID NOT NULL REFERENCES public.users(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Relation: course <-> students (many-to-many)
CREATE TABLE IF NOT EXISTS public.course_students (
  course_id UUID NOT NULL REFERENCES public.courses(id) ON DELETE CASCADE,
  student_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
  enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (course_id, student_id)
);

-- Relation: course <-> challenges (many-to-many)
CREATE TABLE IF NOT EXISTS public.course_challenges (
  course_id UUID NOT NULL REFERENCES public.courses(id) ON DELETE CASCADE,
  challenge_id UUID NOT NULL REFERENCES public.challenges(id) ON DELETE CASCADE,
  assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (course_id, challenge_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_courses_professor 
  ON public.courses (professor_id);

CREATE INDEX IF NOT EXISTS idx_course_students_student 
  ON public.course_students (student_id);

CREATE INDEX IF NOT EXISTS idx_course_challenges_challenge 
  ON public.course_challenges (challenge_id);

-- Trigger for updated_at
DROP TRIGGER IF EXISTS trg_courses_updated_at ON public.courses;

CREATE TRIGGER trg_courses_updated_at
BEFORE UPDATE ON public.courses
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();
