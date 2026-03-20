package course_factory

import (
	"strings"
	"time"

	Entities "../../entities/course"
	Validations "../../validations/course"
)

func NewCourse(
	id, name, description string,
	visibility Entities.CourseVisibility,
	visualIdentity Entities.CourseColour,
	code string,
	period *Entities.Period,
	enrollmentCode, enrollmentURL, professorID string,
) (*Entities.Course, error) {
	now := time.Now()

	if visualIdentity == "" {
		visualIdentity = Entities.CourseColourBlue
	}

	course := &Entities.Course{
		ID:             strings.TrimSpace(id),
		Name:           strings.TrimSpace(name),
		Description:    strings.TrimSpace(description),
		Visibility:     visibility,
		VisualIdentity: visualIdentity,
		Code:           strings.TrimSpace(code),
		Period:         period,
		EnrollmentCode: strings.TrimSpace(enrollmentCode),
		EnrollmentURL:  strings.TrimSpace(enrollmentURL),
		CreatedAt:      now,
		UpdatedAt:      now,
		ProfessorID:    strings.TrimSpace(professorID),
	}

	if err := Validations.ValidateCourse(course); err != nil {
		return nil, err
	}

	return course, nil
}

func ExistingCourse(
	id, name, description string,
	visibility Entities.CourseVisibility,
	visualIdentity Entities.CourseColour,
	code string,
	period *Entities.Period,
	enrollmentCode, enrollmentURL, professorID string,
	createdAt, updatedAt time.Time,
) (*Entities.Course, error) {
	course := &Entities.Course{
		ID:             strings.TrimSpace(id),
		Name:           strings.TrimSpace(name),
		Description:    strings.TrimSpace(description),
		Visibility:     visibility,
		VisualIdentity: visualIdentity,
		Code:           strings.TrimSpace(code),
		Period:         period,
		EnrollmentCode: strings.TrimSpace(enrollmentCode),
		EnrollmentURL:  strings.TrimSpace(enrollmentURL),
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		ProfessorID:    strings.TrimSpace(professorID),
	}

	if err := Validations.ValidateCourse(course); err != nil {
		return nil, err
	}

	return course, nil
}