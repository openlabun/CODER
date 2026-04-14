package exam_validations

import (
	"fmt"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func ValidateExamScore(examScore *Entities.ExamScore) error {
	if examScore == nil {
		return fmt.Errorf("exam score is nil")
	}

	if examScore.ID == "" {
		return fmt.Errorf("exam score id is required")
	}

	if examScore.Score < 0 {
		return fmt.Errorf("exam score cannot be negative")
	}

	if examScore.ExamID == "" {
		return fmt.Errorf("exam score exam id is required")
	}

	if examScore.StudentID == "" {
		return fmt.Errorf("exam score student id is required")
	}

	if examScore.SessionID == "" {
		return fmt.Errorf("exam score session id is required")
	}

	return nil
}

func ValidateExamItemScore(examItemScore *Entities.ExamItemScore) error {
	if examItemScore == nil {
		return fmt.Errorf("exam item score is nil")
	}

	if examItemScore.ID == "" {
		return fmt.Errorf("exam item score id is required")
	}

	if examItemScore.Score < 0 {
		return fmt.Errorf("exam item score cannot be negative")
	}

	if examItemScore.ExamItemID == "" {
		return fmt.Errorf("exam item score exam item id is required")
	}

	if examItemScore.ExamScoreID == "" {
		return fmt.Errorf("exam item score exam score id is required")
	}

	if examItemScore.Tries < 0 {
		return fmt.Errorf("exam item score tries cannot be negative")
	}

	return nil
}