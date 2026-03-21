package submission_factory

import (
	"strings"

	"github.com/google/uuid"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	ExamEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Validations "github.com/openlabun/CODER/apps/api_v2/internal/domain/validations/submission"
)

func NewSubmissionResult(submissionID, testCaseID string) (*Entities.SubmissionResult, error) {
	result := &Entities.SubmissionResult{
		ID:           uuid.New().String(),
		SubmissionID: strings.TrimSpace(submissionID),
		TestCaseID:   strings.TrimSpace(testCaseID),
		Status:       Entities.SubmissionStatusQueued,
		ActualOutput: nil,
		ErrorMessage: nil,
	}

	if err := Validations.ValidateSubmissionResult(result); err != nil {
		return nil, err
	}

	return result, nil
}

func ExistingSubmissionResult(
	id, submissionID, testCaseID string,
	status Entities.SubmissionStatus,
	actualOutput *ExamEntities.IOVariable,
	errorMessage *string,
) (*Entities.SubmissionResult, error) {
	result := &Entities.SubmissionResult{
		ID:           strings.TrimSpace(id),
		SubmissionID: strings.TrimSpace(submissionID),
		TestCaseID:   strings.TrimSpace(testCaseID),
		Status:       status,
		ActualOutput: actualOutput,
		ErrorMessage: errorMessage,
	}

	if err := Validations.ValidateSubmissionResult(result); err != nil {
		return nil, err
	}

	return result, nil
}
