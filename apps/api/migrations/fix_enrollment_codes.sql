-- Fix enrollment codes for existing courses
-- This migration updates enrollment codes to use the correct format
-- Format: COURSE-PERIODG# (e.g., CS101-20251G1)

UPDATE courses
SET enrollment_code = CONCAT(
    UPPER(code),
    '-',
    REPLACE(period, '-', ''),
    'G',
    group_number
)
WHERE enrollment_code IS NOT NULL;
