package services

import (
	"fmt"
	"context"

	submissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
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

// TODO: Cascade deletion for SubmissionResult with IOVariables