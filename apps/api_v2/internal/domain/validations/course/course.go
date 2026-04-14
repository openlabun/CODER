package course_validations

import (
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
)

func validateCourseVisibility(visibility consts.CourseVisibility) bool {
	switch visibility {
	case consts.CourseVisibilityPublic, consts.CourseVisibilityPrivate, consts.CourseVisibilityBlocked:
		return true
	default:
		return false
	}
}

func validateActualYear(year int) bool {
	currentYear := time.Now().Year()
	return year >= currentYear-1
}

func validateCoursePeriod(period consts.AcademicPeriod) bool {
	switch period {
	case consts.AcademicFirstPeriod, consts.AcademicIntersemestral, consts.AcademicSecondPeriod:
	default:
		return false
	}
	return true
}

func ValidateCourseEnrollmentCode(course *Entities.Course, enrollmentCode string) (bool, error) {
	if course == nil {
		return false, fmt.Errorf("course is nil")
	}

	if course.Visibility == consts.CourseVisibilityBlocked {
		return false, fmt.Errorf("blocked courses cannot be enrolled in")
	}

	if course.Visibility == consts.CourseVisibilityPrivate {
		if strings.TrimSpace(enrollmentCode) == "" {
			return false, fmt.Errorf("enrollment code is required for private courses")
		}
	}

	if course.EnrollmentCode == enrollmentCode {
		return true, nil
	}

	return false, fmt.Errorf("invalid enrollment code")
}

func ValidateCourse(course *Entities.Course) error {
	if course == nil {
		return fmt.Errorf("course is nil")
	}

	if strings.TrimSpace(course.ID) == "" {
		return fmt.Errorf("course id is required")
	}
	if strings.TrimSpace(course.Name) == "" {
		return fmt.Errorf("course name is required")
	}
	if strings.TrimSpace(course.Code) == "" {
		return fmt.Errorf("course code is required")
	}
	if strings.TrimSpace(course.ProfessorID) == "" {
		return fmt.Errorf("professor id is required")
	}

	if !validateCourseVisibility(course.Visibility) {
		return fmt.Errorf("invalid course visibility: %q", course.Visibility)
	}

	if course.Period != nil {
		if !validateActualYear(course.Period.Year) {
			return fmt.Errorf("course period year is invalid, it must be an actual course")
		}
		if !validateCoursePeriod(course.Period.Semester) {
			return fmt.Errorf("invalid academic semester: %q", course.Period.Semester)
		}
	}

	if course.Visibility == consts.CourseVisibilityPrivate {
		hasCode := strings.TrimSpace(course.EnrollmentCode) != ""
		hasURL := strings.TrimSpace(course.EnrollmentURL) != ""
		if !hasCode && !hasURL {
			return fmt.Errorf("private courses require enrollment code or enrollment URL")
		}
	}

	return nil
}
