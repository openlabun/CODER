package services

import (
	"context"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
)

func RemoveSubmission(ctx context.Context,
	submissionID string,
	submissionRepository submissionRepository.SubmissionRepository,
	resultsRepository submissionRepository.SubmissionResultRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) error {
	// [STEP 1] Get all Submission results for the submission
	results, err := resultsRepository.GetResultsBySubmissionID(ctx, submissionID)
	if err != nil {
		return err
	}

	// [STEP 2] Delete all existing Submission results for the submission
	for _, result := range results {
		if result != nil {
			err = RemoveSubmissionResult(ctx, result.ID, resultsRepository, ioVariableRepository)
			if err != nil {
				return err
			}
		}
	}

	// [STEP 3] Delete the submission itself
	err = submissionRepository.DeleteSubmission(ctx, submissionID)
	if err != nil {
		return err
	}

	return nil
}
