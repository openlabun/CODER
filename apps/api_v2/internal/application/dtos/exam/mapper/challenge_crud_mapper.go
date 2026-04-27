package mapper

import (
	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/challenge"
)

func MapIOVariableDTOToIOVariableEntity(input dtos.IOVariableDTO) (*Entities.IOVariable, error) {
	result, err := factory.NewIOVariable(
		input.Name,
		constants.VariableFormat(input.Type),
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

func MapCodeTemplateDTOToCodeTemplateEntity(input dtos.CodeTemplateDTO) (*Entities.CodeTemplate, error) {
	result, err := factory.NewCodeTemplate(
		input.Language,
		input.Template,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func MapCodeTemplatesListDTOToCodeTemplatesEntity(input []dtos.CodeTemplateDTO) ([]Entities.CodeTemplate, error) {
	var codeTemplates []Entities.CodeTemplate
	for _, templateDTO := range input {
		template, err := MapCodeTemplateDTOToCodeTemplateEntity(templateDTO)
		if err != nil {
			return nil, err
		}

		if template == nil {
			return nil, err
		}

		codeTemplates = append(codeTemplates, *template)
	}
	return codeTemplates, nil
}

func MapDefaultCodeTemplatesInputToEntities(input dtos.DefaultCodeTemplatesInput) ([]Entities.IOVariable, *Entities.IOVariable, error) {
	inputVariables, err := MapIOVariablesDTOToIOVariablesEntity(input.Inputs)
	if err != nil {
		return nil, nil, err
	}

	outputVariable, err := MapIOVariableDTOToIOVariableEntity(input.Output)
	if err != nil {
		return nil, nil, err
	}

	return inputVariables, outputVariable, nil
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

	codeTemplate, err := MapCodeTemplatesListDTOToCodeTemplatesEntity(input.CodeTemplates)
	if err != nil {
		return nil, err
	}

	challenge, err := factory.NewChallenge(
		input.Title,
		input.Description,
		input.Tags,
		constants.ChallengeStatus(input.Status),
		constants.ChallengeDifficulty(input.Difficulty),
		input.WorkerTimeLimit,
		input.WorkerMemoryLimit,
		codeTemplate,
	 	inputVariables,
		*outputVariable,
		input.Constraints,
		input.UserID,
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
		new_status := constants.ChallengeStatus(*input.Status)
		if new_status != existingChallenge.Status {
			err := state_machine.ApplyTransition(existingChallenge, new_status)
			if err != nil {
				return nil, err
			}
		}
    }

	if input.Difficulty != nil {
		existingChallenge.Difficulty = constants.ChallengeDifficulty(*input.Difficulty)
	}

	if input.WorkerTimeLimit != nil {
		existingChallenge.WorkerTimeLimit = *input.WorkerTimeLimit
	}

	if input.WorkerMemoryLimit != nil {
		existingChallenge.WorkerMemoryLimit = *input.WorkerMemoryLimit
	}

	if input.CodeTemplates != nil {
		codeTemplates, err := MapCodeTemplatesListDTOToCodeTemplatesEntity(*input.CodeTemplates)
		if err != nil {
			return nil, err
		}
		existingChallenge.CodeTemplates = codeTemplates
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

	return existingChallenge, nil
}

func MapPublishChallengeInputToChallengeEntity(existingChallenge *Entities.Challenge) (*Entities.Challenge, error) {
	err := state_machine.ApplyTransition(existingChallenge, constants.ChallengeStatusPublished)
	if err != nil {
		return nil, err
	}
	return existingChallenge, nil
}

func MapArchiveChallengeInputToChallengeEntity(existingChallenge *Entities.Challenge) (*Entities.Challenge, error) {
	err := state_machine.ApplyTransition(existingChallenge, constants.ChallengeStatusArchived)
	if err != nil {
		return nil, err
	}
	return existingChallenge, nil
}
