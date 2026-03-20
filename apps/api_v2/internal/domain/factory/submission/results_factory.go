package submission_factory

import (
	"strings"

	Entities "../../entities/submission"
	ExamEntities "../../entities/exam"
	Validations "../../validations/submission"
)

func NewSubmissionResult(id, submissionID, testCaseID string) (*Entities.SubmissionResult, error) {
	result := &Entities.SubmissionResult{
		ID:           strings.TrimSpace(id),
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
