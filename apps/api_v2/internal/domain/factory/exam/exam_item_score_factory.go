package exam_factory

import (
	"strings"
	"time"

	"github.com/google/uuid"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"
)

func NewExamItemScore(examItemID, examScoreID string) (*Entities.ExamItemScore, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	examItemScore := &Entities.ExamItemScore{
		ID: 			 uuid.New().String(),
		ExamItemID:      strings.TrimSpace(examItemID),
		ExamScoreID:     strings.TrimSpace(examScoreID),
		Score:           0,
		Tries:           0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := Validations.ValidateExamItemScore(examItemScore); err != nil {
		return nil, err
	}

	return examItemScore, nil

}
func ExistingExamItemScore(
	ID, examItemID, examScoreID string,
	score, tries int,
	createdAt, updatedAt string,
) (*Entities.ExamItemScore, error) {
	examItemScore := &Entities.ExamItemScore{
		ID: 			 strings.TrimSpace(ID),
		ExamItemID:      strings.TrimSpace(examItemID),
		ExamScoreID:     strings.TrimSpace(examScoreID),
		Score:           score,
		Tries:           tries,
		CreatedAt:       strings.TrimSpace(createdAt),
		UpdatedAt:       strings.TrimSpace(updatedAt),
	}

	if err := Validations.ValidateExamItemScore(examItemScore); err != nil {
		return nil, err
	}

	return examItemScore, nil
}