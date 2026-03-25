package exam_validations

import (
	"fmt"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func ValidateExamItem(examItem *Entities.ExamItem) error {
	if examItem == nil {
		return fmt.Errorf("exam item is nil")
	}

	if strings.TrimSpace(examItem.ID) == "" {
		return fmt.Errorf("exam item id is required")
	}

	if strings.TrimSpace(examItem.ChallengeID) == "" {
		return fmt.Errorf("exam item challenge id is required")
	}

	if strings.TrimSpace(examItem.ExamID) == "" {
		return fmt.Errorf("exam item exam id is required")
	}

	if examItem.Order < 0 {
		return fmt.Errorf("exam item order must be non-negative")
	}

	if examItem.Points < 0 {
		return fmt.Errorf("exam item points must be non-negative")
	}
	

	return nil
}