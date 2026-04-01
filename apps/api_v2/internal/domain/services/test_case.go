package services

import (
	"fmt"
	"context"

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

// TODO: Cascade deletion for TestCase with IOVariables