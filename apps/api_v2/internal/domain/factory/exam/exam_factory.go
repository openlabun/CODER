package exam_factory

import (
	"strings"
	"time"

	Entities "../../entities/exam"
	Validations "../../validations/exam"
)

func NewExam(
	id, title, description string,
	visibility Entities.Visibility,
	startTime time.Time,
	endTime *time.Time,
	allowLateSubmissions bool,
	timeLimit, tryLimit int,
	professorID, courseID string,
) (*Entities.Exam, error) {
	now := time.Now()
	if visibility == "" {
		visibility = Entities.VisibilityPrivate
	}
	if startTime.IsZero() {
		startTime = now
	}

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
		CreatedAt:            now,
		UpdatedAt:            now,
		ProfessorID:          strings.TrimSpace(professorID),
		CourseID:             strings.TrimSpace(courseID),
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
	professorID, courseID string,
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
		CourseID:             strings.TrimSpace(courseID),
	}

	if err := Validations.ValidateExam(exam); err != nil {
		return nil, err
	}

	return exam, nil
}

