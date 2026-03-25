package mapper

import (
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func MapCreateExamItemInputToExamItemEntity(input dtos.CreateExamItemInput) (*Entities.ExamItem, error) {
	examItem, err := factory.NewExamItem(
		input.ChallengeID,
		input.ExamID,
		input.Order,
		input.Points,
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

	return existingExamItem, nil
}