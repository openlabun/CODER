package services

import (
	"context"
	"fmt"

	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
)

func CreateTestCase(
	ctx context.Context,
	testCase *examEntities.TestCase,
	testCaseRepository examRepository.TestCaseRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) (*examEntities.TestCase, error) {
	if testCase == nil {
		return nil, fmt.Errorf("test case is nil")
	}

	hydrated, err := hydrateTestCaseIOVariables(ctx, testCase, ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return testCaseRepository.CreateTestCase(ctx, hydrated)
}

func UpdateTestCase(
	ctx context.Context,
	testCase *examEntities.TestCase,
	testCaseRepository examRepository.TestCaseRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) (*examEntities.TestCase, error) {
	if testCase == nil {
		return nil, fmt.Errorf("test case is nil")
	}

	hydrated, err := hydrateTestCaseIOVariables(ctx, testCase, ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return testCaseRepository.UpdateTestCase(ctx, hydrated)
}

func RemoveTestCase(
	ctx context.Context,
	testCaseID string,
	testCaseRepository examRepository.TestCaseRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) error {
	testCase, err := testCaseRepository.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		return err
	}
	if testCase == nil {
		return fmt.Errorf("test case with id %q does not exist", testCaseID)
	}

	for _, ioVariable := range testCase.Input {
		if err := ioVariableRepository.DeleteIOVariable(ctx, ioVariable.ID); err != nil {
			return err
		}
	}

	if err := ioVariableRepository.DeleteIOVariable(ctx, testCase.ExpectedOutput.ID); err != nil {
		return err
	}

	return testCaseRepository.DeleteTestCase(ctx, testCaseID)
}
