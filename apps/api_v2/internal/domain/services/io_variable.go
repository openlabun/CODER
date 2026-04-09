package services

import (
	"context"
	"fmt"

	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
)


func hydrateChallengeIOVariables(
	ctx context.Context,
	challenge *examEntities.Challenge,
	ioVariableRepository examRepository.IOVariableRepository,
) (*examEntities.Challenge, error) {
	if ioVariableRepository == nil {
		return nil, fmt.Errorf("io variable repository is nil")
	}

	inputs := make([]examEntities.IOVariable, 0, len(challenge.InputVariables))
	for i := range challenge.InputVariables {
		persistedInput, err := upsertIOVariable(ctx, challenge.InputVariables[i], ioVariableRepository)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, *persistedInput)
	}

	persistedOutput, err := upsertIOVariable(ctx, challenge.OutputVariable, ioVariableRepository)
	if err != nil {
		return nil, err
	}

	hydrated := *challenge
	hydrated.InputVariables = inputs
	hydrated.OutputVariable = *persistedOutput

	return &hydrated, nil
}

func hydrateTestCaseIOVariables(
	ctx context.Context,
	testCase *examEntities.TestCase,
	ioVariableRepository examRepository.IOVariableRepository,
) (*examEntities.TestCase, error) {
	if ioVariableRepository == nil {
		return nil, fmt.Errorf("io variable repository is nil")
	}

	inputs := make([]examEntities.IOVariable, 0, len(testCase.Input))
	for i := range testCase.Input {
		persistedInput, err := upsertIOVariable(ctx, testCase.Input[i], ioVariableRepository)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, *persistedInput)
	}

	persistedOutput, err := upsertIOVariable(ctx, testCase.ExpectedOutput, ioVariableRepository)
	if err != nil {
		return nil, err
	}

	hydrated := *testCase
	hydrated.Input = inputs
	hydrated.ExpectedOutput = *persistedOutput

	return &hydrated, nil
}

func hydrateSubmissionResultIOVariable(
	ctx context.Context,
	result *submissionEntities.SubmissionResult,
	ioVariableRepository examRepository.IOVariableRepository,
) (*submissionEntities.SubmissionResult, error) {
	if ioVariableRepository == nil {
		return nil, fmt.Errorf("io variable repository is nil")
	}

	hydrated := *result
	if result.ActualOutput == nil {
		return &hydrated, nil
	}

	persistedOutput, err := upsertIOVariable(ctx, *result.ActualOutput, ioVariableRepository)
	if err != nil {
		return nil, err
	}
	hydrated.ActualOutput = persistedOutput

	return &hydrated, nil
}

func upsertIOVariable(
	ctx context.Context,
	ioVariable examEntities.IOVariable,
	ioVariableRepository examRepository.IOVariableRepository,
) (*examEntities.IOVariable, error) {
	existing, err := ioVariableRepository.GetIOVariableByID(ctx, ioVariable.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get IO variable by ID: %w", err)
	}

	if existing == nil {
		return ioVariableRepository.CreateIOVariable(ctx, &ioVariable)
	}

	return ioVariableRepository.UpdateIOVariable(ctx, &ioVariable)
}
