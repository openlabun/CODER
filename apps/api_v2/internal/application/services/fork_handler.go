package services

import (
	"context"
	"fmt"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"

	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
)

func ForkChallenge(ctx context.Context,
	challenge Entities.Challenge,
	userID string,
	challengeRepository examRepository.ChallengeRepository,
	testCaseRepository examRepository.TestCaseRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) (*Entities.Challenge, error) {

	// [STEP 1] Get Challenge TestCases
	testCases, err := testCaseRepository.GetTestCasesByChallengeID(ctx, challenge.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases for challenge: %w", err)
	}

	// [STEP 2] Fork IO Variables of Challenge
	output, err := forkIOVariable(challenge.OutputVariable)
	if err != nil {
		return nil, fmt.Errorf("failed to fork challenge output variable: %w", err)
	}

	input, err := forkIOVariables(challenge.InputVariables)
	if err != nil {
		return nil, fmt.Errorf("failed to fork challenge input variables: %w", err)
	}

	// [STEP 3] Create new Challenge with forked IO Variables
	forkedChallenge, err := factory.NewChallenge(
		challenge.Title,
		challenge.Description,
		challenge.Tags,
		challenge.Status,
		challenge.Difficulty,
		challenge.WorkerTimeLimit,
		challenge.WorkerMemoryLimit,
		input,
		*output,
		challenge.Constraints,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create forked challenge: %w", err)
	}

	// [STEP 4] Fork TestCases and associate them with the new Challenge
	forkedTestCases, err := forkTestCases(testCases, forkedChallenge.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fork test cases: %w", err)
	}

	// [STEP 5] Save forked Challenge and TestCases in repository
	new_challenge, err := domain_services.CreateChallenge(ctx, forkedChallenge, challengeRepository, ioVariableRepository)
	if err != nil {
		return nil, fmt.Errorf("failed to save forked challenge: %w", err)
	}

	// [STEP 6] Save forked TestCases in repository
	for _, tc := range forkedTestCases {
		_, err = domain_services.CreateTestCase(ctx, &tc, testCaseRepository, ioVariableRepository)
		if err != nil {
			return nil, fmt.Errorf("failed to save forked test case: %w", err)
		}
	}

	return new_challenge, nil
}

func forkTestCase(testCase *Entities.TestCase,
	newChallengeID string,
) (*Entities.TestCase, error) {
	if testCase == nil {
		return nil, fmt.Errorf("testCase is nil")
	}

	input, err := forkIOVariables(testCase.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to fork test case input: %w", err)
	}

	output, err := forkIOVariable(testCase.ExpectedOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to fork test case output: %w", err)
	}

	return factory.NewTestCase(
		testCase.Name,
		input,
		*output,
		testCase.IsSample,
		testCase.Points,
		newChallengeID,
	)
}

func forkTestCases(testCases []*Entities.TestCase,
	newChallengeID string,
) ([]Entities.TestCase, error) {
	forkedTestCases := make([]Entities.TestCase, len(testCases))
	for i, tc := range testCases {
		forked, err := forkTestCase(tc, newChallengeID)
		if err != nil {
			return nil, fmt.Errorf("failed to fork test case %s: %w", tc.ID, err)
		}
		forkedTestCases[i] = *forked
	}

	return forkedTestCases, nil
}

func forkIOVariable(ioVar Entities.IOVariable) (*Entities.IOVariable, error) {
	return factory.NewIOVariable(
		ioVar.Name,
		ioVar.Type,
		ioVar.Value,
	)
}

func forkIOVariables(ioVars []Entities.IOVariable) ([]Entities.IOVariable, error) {
	forkedVars := make([]Entities.IOVariable, len(ioVars))
	for i, ioVar := range ioVars {
		forkedVar, err := forkIOVariable(ioVar)
		if err != nil {
			return nil, fmt.Errorf("failed to fork IO variable %s: %w", ioVar.Name, err)
		}
		forkedVars[i] = *forkedVar
	}
	return forkedVars, nil
}
