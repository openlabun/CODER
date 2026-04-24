package mapper

import (
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
)

func MapCreateExamItemInputToExamItemEntity(input dtos.CreateExamItemInput) (*Entities.ExamItem, error) {
	examItem, err := factory.NewExamItem(
		input.ChallengeID,
		input.ExamID,
		input.Order,
		input.Points,
		input.TryLimit,
	)
	if err != nil {
		return nil, err
	}

	return examItem, nil
}

func MapUpdateExamItemInputToExamItemEntity(existingExamItem *Entities.ExamItem, input dtos.UpdateExamItemInput) (*Entities.ExamItem, error) {
	if input.Order != nil {
		existingExamItem.Order = *input.Order
	}

	if input.Points != nil {
		existingExamItem.Points = *input.Points
	}

	if input.TryLimit != nil {
		existingExamItem.TryLimit = *input.TryLimit
	}

	return existingExamItem, nil
}

func MapExamItemScore(examItem *Entities.ExamItem, examScore *Entities.ExamScore) (*Entities.ExamItemScore, error) {
	if examItem == nil {
		return nil, fmt.Errorf("exam item cannot be nil")
	}

	if examScore == nil {
		return nil, fmt.Errorf("exam score cannot be nil")
	}

	return factory.NewExamItemScore(
		examItem.ID,
		examScore.ID,
	)
}
