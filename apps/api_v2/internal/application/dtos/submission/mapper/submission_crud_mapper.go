package mapper

import (
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/submission"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

func MapCreateSubmissionInputToSubmissionEntity(userID string, input dtos.CreateSubmissionInput) (*Entities.Submission, error) {
	submission, err := factory.NewSubmission(
		input.Code,
		Entities.ProgrammingLanguage(input.Language),
		input.ChallengeID,
		input.SessionID,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return submission, nil
}

func MapSubmissionResultEntity (submissionID string, testCaseID string) (*Entities.SubmissionResult, error) {
	submissionResult, err := factory.NewSubmissionResult(
		submissionID,
		testCaseID,
	)

	if err != nil {
		return nil, err
	}
	return submissionResult, nil
}

func MapSubmissionOutputDTO(submission *Entities.Submission, results []Entities.SubmissionResult) *dtos.SubmissionOutputDTO {
	return &dtos.SubmissionOutputDTO{
		Submission: *submission,
		Results: results,
	}
}