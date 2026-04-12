package course_factory

import (
	"strings"
	"time"
	"github.com/google/uuid"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/course"
	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
)

func NewCourse(
	name, description string,
	visibility consts.CourseVisibility,
	visualIdentity consts.CourseColour,
	code string,
	period *Entities.Period,
	enrollmentCode, enrollmentURL, professorID string,
) (*Entities.Course, error) {
	now := time.Now()

	if visualIdentity == "" {
		visualIdentity = consts.CourseColourBlue
	}

	course := &Entities.Course{
		ID:             uuid.New().String(),
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
	visibility consts.CourseVisibility,
	visualIdentity consts.CourseColour,
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