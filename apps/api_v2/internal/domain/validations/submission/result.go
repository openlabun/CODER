package submission_validations

import (
	"fmt"
	"strings"

	StateMachine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/submission"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	ExamValidations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/exam"
)

func validateSubmissionResultsDetails(result *Entities.SubmissionResult) error {
	hasActualOutput := result.ActualOutput != nil
	hasErrorMessage := result.ErrorMessage != nil && strings.TrimSpace(*result.ErrorMessage) != ""

	switch result.Status {
	case Entities.SubmissionStatusQueued, Entities.SubmissionStatusRunning:
		if hasActualOutput || hasErrorMessage {
			return fmt.Errorf("queued/running result cannot include output or error message")
		}
	case Entities.SubmissionStatusAccepted, Entities.SubmissionStatusWrongAnswer:
		if !hasActualOutput {
			return fmt.Errorf("accepted/wrong_answer result requires actual output")
		}
		if hasErrorMessage {
			return fmt.Errorf("accepted/wrong_answer result cannot include error message")
		}
		if err := ExamValidations.ValidateIOVariable(*result.ActualOutput); err != nil {
			return fmt.Errorf("invalid actual output: %w", err)
		}
	case Entities.SubmissionStatusError:
		if hasActualOutput {
			return fmt.Errorf("error result cannot include actual output")
		}
		if !hasErrorMessage {
			return fmt.Errorf("error result requires error message")
		}
	}

	return nil
}

func ValidateSubmissionResult(result *Entities.SubmissionResult) error {
	if result == nil {
		return fmt.Errorf("submission result is nil")
	}

	if strings.TrimSpace(result.ID) == "" {
		return fmt.Errorf("submission result id is required")
	}

	if strings.TrimSpace(result.SubmissionID) == "" {
		return fmt.Errorf("submission id is required")
	}

	if strings.TrimSpace(result.TestCaseID) == "" {
		return fmt.Errorf("test case id is required")
	}

	if !StateMachine.IsValidState(result.Status) {
		return fmt.Errorf("invalid submission result status: %q", result.Status)
	}

	if err := validateSubmissionResultsDetails(result); err != nil {
		return fmt.Errorf("invalid submission result details: %w", err)
	}

	return nil
}
