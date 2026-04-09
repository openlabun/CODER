package mapper

import (
	challenge_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
)

func MapCreateTestCaseInputToTestCaseEntity(input dtos.CreateTestCaseInput) (*challenge_entities.TestCase, error) {
	outputVariable, err := MapIOVariableDTOToIOVariableEntity(input.ExpectedOutput)
	if err != nil {
		return nil, err
	}

	inputVariables, err := MapIOVariablesDTOToIOVariablesEntity(input.Input)
	if err != nil {
		return nil, err
	}

	testCase, err := factory.NewTestCase(
		input.Name,
		inputVariables,
		*outputVariable,
		input.IsSample,
		input.Points,
		input.ChallengeID,
	)

	if err != nil {
		return nil, err
	}
	
	return testCase, nil
}

func MapUpdateTestCaseInputToTestCaseEntity(existingTestCase *challenge_entities.TestCase, input dtos.UpdateTestCaseInput) (*challenge_entities.TestCase, error) {
	if input.Name != nil {
		existingTestCase.Name = *input.Name
	}

	if input.Input != nil {
		inputVariables, err := MapIOVariablesDTOToIOVariablesEntity(*input.Input)
		if err != nil {
			return nil, err
		}
		existingTestCase.Input = inputVariables
	}

	if input.ExpectedOutput != nil {
		outputVariable, err := MapIOVariableDTOToIOVariableEntity(*input.ExpectedOutput)
		if err != nil {
			return nil, err
		}
		existingTestCase.ExpectedOutput = *outputVariable
	}

	if input.IsSample != nil {
		existingTestCase.IsSample = *input.IsSample
	}

	if input.Points != nil {
		existingTestCase.Points = *input.Points
	}

	return existingTestCase, nil
}