package exam_factory

import (
	"strings"
	"time"

	"github.com/google/uuid"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"
)

func NewExamScore(examID, sessionID, studentID string) (*Entities.ExamScore, error) {
	now := time.Now().UTC().Format(time.RFC3339)

	examScore := &Entities.ExamScore{
		ID: uuid.New().String(),
		ExamID:      strings.TrimSpace(examID),
		SessionID:   strings.TrimSpace(sessionID),
		Score:       0,
		CreatedAt:   now,
		UpdatedAt:   now,
		StudentID:   strings.TrimSpace(studentID),
	}

	if err := Validations.ValidateExamScore(examScore); err != nil {
		return nil, err
	}

	return examScore, nil

}
func ExistingExamScore(
	ID, examID, sessionID string,
	score int,
	createdAt, updatedAt, studentID string,
) (*Entities.ExamScore, error) {
	examScore := &Entities.ExamScore{
		ID: strings.TrimSpace(ID),
		ExamID:      strings.TrimSpace(examID),
		SessionID:   strings.TrimSpace(sessionID),
		Score:       score,
		CreatedAt:   strings.TrimSpace(createdAt),
		UpdatedAt:   strings.TrimSpace(updatedAt),
		StudentID:   strings.TrimSpace(studentID),
	}

	if err := Validations.ValidateExamScore(examScore); err != nil {
		return nil, err
	}

	return examScore, nil
}