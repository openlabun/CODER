package mapper

import (
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/challenge"
)

func MapIOVariableDTOToIOVariableEntity(input dtos.IOVariableDTO) (*Entities.IOVariable, error) {
	result, err := factory.NewIOVariable(
		input.Name,
		Entities.VariableFormat(input.Type),
		input.Value,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func MapIOVariablesDTOToIOVariablesEntity(input []dtos.IOVariableDTO) ([]Entities.IOVariable, error) {
	var inputVariables []Entities.IOVariable

	for _, inputVar := range input {
		variable, err := MapIOVariableDTOToIOVariableEntity(inputVar)
		if err != nil {
			return nil, err
		}

		if variable == nil {
			return nil, err
		}

		inputVariables = append(inputVariables, *variable)
	}

	return inputVariables, nil
}

func MapCreateChallengeInputToChallengeEntity(input dtos.CreateChallengeInput) (*Entities.Challenge, error) {
	outputVariable, err := MapIOVariableDTOToIOVariableEntity(input.OutputVariable)
	if err != nil {
		return nil, err
	}

	inputVariables, err := MapIOVariablesDTOToIOVariablesEntity(input.InputVariables)
	if err != nil {
		return nil, err
	}

	challenge, err := factory.NewChallenge(
		input.Title,
		input.Description,
		input.Tags,
		Entities.ChallengeStatus(input.Status),
		Entities.ChallengeDifficulty(input.Difficulty),
		input.WorkerTimeLimit,
		input.WorkerMemoryLimit,
	 	inputVariables,
		*outputVariable,
		input.Constraints,
		input.ExamID,
		input.CourseID,
	)
	if err != nil {
		return nil, err
	}

	return challenge, nil
}

func MapUpdateChallengeInputToChallengeEntity(existingChallenge *Entities.Challenge, input dtos.UpdateChallengeInput) (*Entities.Challenge, error) {
	if input.Title != nil {
		existingChallenge.Title = *input.Title
	}

	if input.Description != nil {
		existingChallenge.Description = *input.Description
	}

	if input.Tags != nil {
		existingChallenge.Tags = *input.Tags
	}

	if input.Status != nil {
        err := state_machine.ApplyTransition(existingChallenge, Entities.ChallengeStatus(*input.Status))
        if err != nil {
            return nil, err
        }
    }

	if input.Difficulty != nil {
		existingChallenge.Difficulty = Entities.ChallengeDifficulty(*input.Difficulty)
	}

	if input.WorkerTimeLimit != nil {
		existingChallenge.WorkerTimeLimit = *input.WorkerTimeLimit
	}

	if input.WorkerMemoryLimit != nil {
		existingChallenge.WorkerMemoryLimit = *input.WorkerMemoryLimit
	}

	if input.InputVariables != nil {
		inputVariables, err := MapIOVariablesDTOToIOVariablesEntity(*input.InputVariables)
		if err != nil {
			return nil, err
		}
		existingChallenge.InputVariables = inputVariables
	}

	if input.OutputVariable != nil {
		outputVariable, err := MapIOVariableDTOToIOVariableEntity(*input.OutputVariable)
		if err != nil {
			return nil, err
		}
		existingChallenge.OutputVariable = *outputVariable
	}


	if input.Constraints != nil {
		existingChallenge.Constraints = *input.Constraints
	}

	if input.ExamID != nil {
		existingChallenge.ExamID = *input.ExamID
	}

	if input.CourseID != nil {
		existingChallenge.CourseID = *input.CourseID
	}

	return existingChallenge, nil
}

func MapPublishChallengeInputToChallengeEntity(existingChallenge *Entities.Challenge) (*Entities.Challenge, error) {
	err := state_machine.ApplyTransition(existingChallenge, Entities.ChallengeStatusPublished)
	if err != nil {
		return nil, err
	}
	return existingChallenge, nil
}

func MapArchiveChallengeInputToChallengeEntity(existingChallenge *Entities.Challenge) (*Entities.Challenge, error) {
	err := state_machine.ApplyTransition(existingChallenge, Entities.ChallengeStatusArchived)
	if err != nil {
		return nil, err
	}
	return existingChallenge, nil
}
