-- Add enrollment_code column to courses table
ALTER TABLE courses ADD COLUMN IF NOT EXISTS enrollment_code VARCHAR(50) UNIQUE;

-- Create index for faster lookups by enrollment code
CREATE INDEX IF NOT EXISTS idx_courses_enrollment_code ON courses(enrollment_code);
