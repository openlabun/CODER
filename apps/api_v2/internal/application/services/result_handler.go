package services

import (
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/submission"

	submission_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

func CheckSubmissionResult(submission *submissionEntities.SubmissionResult, testCase *examEntities.TestCase) (*submissionEntities.SubmissionResult, error) {
	if submission.ActualOutput.Value == testCase.ExpectedOutput.Value {
		if submission.ActualOutput.Type == testCase.ExpectedOutput.Type {
			state_machine.ApplyTransition(submission, submission_constants.SubmissionStatusAccepted)
			return submission, nil
		}
	} 
	
	state_machine.ApplyTransition(submission, submission_constants.SubmissionStatusWrongAnswer)
	return submission, nil
}