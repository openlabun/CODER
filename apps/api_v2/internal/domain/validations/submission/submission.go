package submission_validations

import (
	"fmt"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

func validateSubmissionLanguage(language Entities.ProgrammingLanguage) bool {
	switch language {
	case Entities.LanguageCPP, Entities.LanguagePython, Entities.LanguageJava:
		return true
	default:
		return false
	}
}

func ValidateSubmission(submission *Entities.Submission) error {
	if submission == nil {
		return fmt.Errorf("submission is nil")
	}

	if strings.TrimSpace(submission.ID) == "" {
		return fmt.Errorf("submission id is required")
	}

	if strings.TrimSpace(submission.Code) == "" {
		return fmt.Errorf("submission code is required")
	}

	if strings.TrimSpace(submission.ChallengeID) == "" {
		return fmt.Errorf("submission challenge id is required")
	}

	if strings.TrimSpace(submission.SessionID) == "" {
		return fmt.Errorf("submission session id is required")
	}

	if strings.TrimSpace(submission.UserID) == "" {
		return fmt.Errorf("submission user id is required")
	}

	if !validateSubmissionLanguage(submission.Language) {
		return fmt.Errorf("invalid submission language: %q", submission.Language)
	}

	if submission.Score < 0 {
		return fmt.Errorf("submission score cannot be negative")
	}

	if submission.TimeMsTotal < 0 {
		return fmt.Errorf("submission timeMsTotal cannot be negative")
	}

	return nil
}
