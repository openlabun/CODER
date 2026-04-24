package exam_factory

import (
	"github.com/google/uuid"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"
)

func NewExamItem(
	challengeID, examID string,
	order, points int,
	try_limit *int,
) (*Entities.ExamItem, error) {
	if try_limit == nil {
		defaultTryLimit := -1
		try_limit = &defaultTryLimit
	}

	examItem := &Entities.ExamItem{
		ID:          uuid.New().String(),
		ChallengeID: challengeID,
		ExamID:      examID,
		Order:       order,
		Points:      points,
		TryLimit:    *try_limit,
	}

	if err := Validations.ValidateExamItem(examItem); err != nil {
		return nil, err
	}

	return examItem, nil
}

func ExistingExamItem(
	id, challengeID, examID string,
	order, points, try_limit int,
) (*Entities.ExamItem, error) {
	examItem := &Entities.ExamItem{
		ID:          id,
		ChallengeID: challengeID,
		ExamID:      examID,
		Order:       order,
		Points:      points,
		TryLimit:    try_limit,
	}

	if err := Validations.ValidateExamItem(examItem); err != nil {
		return nil, err
	}

	return examItem, nil
}
