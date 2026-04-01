package services

import (
	"context"
	"fmt"

	submissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
)

func CreateSubmissionResult(
	ctx context.Context,
	result *submissionEntities.SubmissionResult,
	resultRepository submissionRepository.SubmissionResultRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) (*submissionEntities.SubmissionResult, error) {
	if result == nil {
		return nil, fmt.Errorf("submission result is nil")
	}

	hydrated, err := hydrateSubmissionResultIOVariable(ctx, result, ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return resultRepository.CreateResult(ctx, hydrated)
}

func UpdateSubmissionResult(
	ctx context.Context,
	result *submissionEntities.SubmissionResult,
	resultRepository submissionRepository.SubmissionResultRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) (*submissionEntities.SubmissionResult, error) {
	if result == nil {
		return nil, fmt.Errorf("submission result is nil")
	}

	hydrated, err := hydrateSubmissionResultIOVariable(ctx, result, ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return resultRepository.UpdateResult(ctx, hydrated)
}

func RemoveSubmissionResult(
	ctx context.Context,
	resultID string,
	resultRepository submissionRepository.SubmissionResultRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) error {
	result, err := resultRepository.GetResultByID(ctx, resultID)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("submission result with id %q does not exist", resultID)
	}

	if result.ActualOutput != nil {
		if err := ioVariableRepository.DeleteIOVariable(ctx, result.ActualOutput.ID); err != nil {
			return err
		}
	}

	return resultRepository.DeleteResult(ctx, resultID)
}
