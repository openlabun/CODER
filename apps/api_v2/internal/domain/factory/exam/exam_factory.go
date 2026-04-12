package exam_factory

import (
	"strings"
	"time"

	"github.com/google/uuid"

	exam_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"
)

func NewExam(
	title, description string,
	visibility Entities.Visibility,
	startTime time.Time,
	endTime *time.Time,
	allowLateSubmissions bool,
	timeLimit, tryLimit int,
	professorID string,
	courseID *string,
) (*Entities.Exam, error) {
	now := time.Now()
	if visibility == "" {
		visibility = exam_constants.VisibilityPrivate
	}
	if startTime.IsZero() {
		startTime = now
	}
	if courseID != nil && strings.TrimSpace(*courseID) != "" {
		trimmedCourseID := strings.TrimSpace(*courseID)
		courseID = &trimmedCourseID
	}

	exam := &Entities.Exam{
		ID:                   uuid.New().String(),
		Title:                strings.TrimSpace(title),
		Description:          strings.TrimSpace(description),
		Visibility:           visibility,
		StartTime:            startTime,
		EndTime:              endTime,
		AllowLateSubmissions: allowLateSubmissions,
		TimeLimit:            timeLimit,
		TryLimit:             tryLimit,
		CreatedAt:            now,
		UpdatedAt:            now,
		ProfessorID:          strings.TrimSpace(professorID),
		CourseID:             courseID,
	}

	if err := Validations.ValidateExamEndTime(exam, now); err != nil {
		return nil, err
	}

	if err := Validations.ValidateExam(exam); err != nil {
		return nil, err
	}

	return exam, nil
}

func ExistingExam(
	id, title, description string,
	visibility Entities.Visibility,
	startTime time.Time,
	endTime *time.Time,
	allowLateSubmissions bool,
	timeLimit, tryLimit int,
	professorID string,
	courseID *string,
	createdAt, updatedAt time.Time,
) (*Entities.Exam, error) {
	exam := &Entities.Exam{
		ID:                   strings.TrimSpace(id),
		Title:                strings.TrimSpace(title),
		Description:          strings.TrimSpace(description),
		Visibility:           visibility,
		StartTime:            startTime,
		EndTime:              endTime,
		AllowLateSubmissions: allowLateSubmissions,
		TimeLimit:            timeLimit,
		TryLimit:             tryLimit,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
		ProfessorID:          strings.TrimSpace(professorID),
		CourseID:             courseID,
	}

	if err := Validations.ValidateExam(exam); err != nil {
		return nil, err
	}

	return exam, nil
}

