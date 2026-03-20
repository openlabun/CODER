package course_validations

import (
	"fmt"
	"strings"

	Entities "../../entities/course"
)

func ValidateCourseStudent(courseStudent *Entities.CourseStudent) error {
	if courseStudent == nil {
		return fmt.Errorf("course student is nil")
	}

	if strings.TrimSpace(courseStudent.CourseID) == "" {
		return fmt.Errorf("course id is required")
	}

	if strings.TrimSpace(courseStudent.StudentID) == "" {
		return fmt.Errorf("student id is required")
	}

	return nil
}
