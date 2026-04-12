package mapper

import (
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/submission"
	exam_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/submission"
)

func MapCreateSubmissionInputToSubmissionEntity(userID string, input dtos.CreateSubmissionInput) (*Entities.Submission, error) {
	submission, err := factory.NewSubmission(
		input.Code,
		constants.ProgrammingLanguage(input.Language),
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

func MapResultInputToSubmissionResultEntity(input dtos.UpdateResultInput, submissionResult *Entities.SubmissionResult, testCase *examEntities.TestCase) (*Entities.SubmissionResult, error) {
	if submissionResult == nil {
		return nil, fmt.Errorf("submission result cannot be nil")
	}
	
	if input.Output != nil {
		output_name := testCase.ExpectedOutput.Name
		output_type := testCase.ExpectedOutput.Type

		if submissionResult.ActualOutput == nil {
			submissionResult.ActualOutput, _ = exam_factory.NewIOVariable(
				output_name,
				output_type,
				*input.Output,
			)
		}
		submissionResult.ActualOutput.Value = *input.Output
	}
	submissionResult.ErrorMessage = input.Error

	status := constants.SubmissionStatus(input.Status)

	valid := state_machine.IsValidState(status)
	if !valid {
		return nil, fmt.Errorf("invalid submission status: %s", status)
	}

	err := state_machine.ApplyTransition(submissionResult, status)
	if err != nil {
		return nil, fmt.Errorf("failed to apply transition: %w", err)
	}

	return submissionResult, nil
}

func MapSubmissionResultToPublishedDTO(
	submission Entities.Submission, 
	result Entities.SubmissionResult, 
	test_case examEntities.TestCase, 
	challenge examEntities.Challenge,
) *dtos.SubmissionResultPublishedDTO {
	input := services.ExtractInputFromTestCase(test_case)

	return &dtos.SubmissionResultPublishedDTO{
		SubmissionID: submission.ID,
		Code: submission.Code,
		Input: input,
		ResultID: result.ID,
		TimeLimitMs: challenge.WorkerTimeLimit,
		MemoryLimitMb: challenge.WorkerMemoryLimit,
		Status: string(result.Status),
		Type: string(test_case.ExpectedOutput.Type),
		Language: string(submission.Language),
	}
}