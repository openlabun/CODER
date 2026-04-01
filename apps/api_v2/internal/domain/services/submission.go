package services

import (
	"context"
	
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
)

func RemoveSubmission (ctx context.Context,
	submissionID string,
	submissionRepository submissionRepository.SubmissionRepository,
	resultsRepository submissionRepository.SubmissionResultRepository,
	) error {
		// [STEP 1] Get all Submission results for the submission
		results, err := resultsRepository.GetResultsBySubmissionID(ctx, submissionID)
		if err != nil {
			return err
		}

		// [STEP 2] Delete all existing Submission results for the submission
		for _, result := range results {
			if result != nil {
				err = resultsRepository.DeleteResult(ctx, result.ID)
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