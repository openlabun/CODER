package mapper

import (
	"time"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
)

func MapCreateExamInputToExamEntity(input dtos.CreateExamInput) (*Entities.Exam, error) {
	// Convert start and end time from string to time.Time
	startTime, err := time.Parse(time.RFC3339, input.StartTime)
	if err != nil {
		return nil, err
	}

	var endTime *time.Time
	if input.EndTime != nil {
		parsedEndTime, err := time.Parse(time.RFC3339, *input.EndTime)
		if err != nil {
			return nil, err
		}
		endTime = &parsedEndTime
	}

	// Map CreateExamInput to Exam entity
	exam, err := factory.NewExam(
		input.Title,
		input.Description,
		Entities.Visibility(input.Visibility),
		startTime,
		endTime,
		input.AllowLateSubmissions,
		input.TimeLimit,
		input.TryLimit,
		input.ProfessorID,
		input.CourseID,
	)
	if err != nil {
		return nil, err
	}
	return exam, nil
}

func MapUpdateExamInputToExamEntity(existingExam *Entities.Exam, input dtos.UpdateExamInput) (*Entities.Exam, error) {
	// Update fields only if they are provided in the input
	if input.Title != nil {
		existingExam.Title = *input.Title
	}

	if input.Description != nil {
		existingExam.Description = *input.Description
	}

	if input.Visibility != nil {
		existingExam.Visibility = Entities.Visibility(*input.Visibility)
	}

	if input.StartTime != nil {
		startTime, err := time.Parse(time.RFC3339, *input.StartTime)
		if err != nil {
			return nil, err
		}
		existingExam.StartTime = startTime
	}

	if input.EndTime != nil {
		endTime, err := time.Parse(time.RFC3339, *input.EndTime)
		if err != nil {
			return nil, err
		}
		existingExam.EndTime = &endTime
	}

	if input.AllowLateSubmissions != nil {
		existingExam.AllowLateSubmissions = *input.AllowLateSubmissions
	}

	if input.TimeLimit != nil {
		existingExam.TimeLimit = *input.TimeLimit
	}

	if input.TryLimit != nil {
		existingExam.TryLimit = *input.TryLimit
	}

	return existingExam, nil
}

func MapExamVisibilityInputToExamEntity (existingExam *Entities.Exam, input dtos.ChangeExamVisibilityInput) (*Entities.Exam, error) {
	existingExam.Visibility = Entities.Visibility(input.Visibility)
	return existingExam, nil
}

func MapExamEndTimeInputToExamEntity (existingExam *Entities.Exam, now time.Time) (*Entities.Exam, error) {
	existingExam.EndTime = &now
	return existingExam, nil
}