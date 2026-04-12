package exam_validations

import (
	"fmt"
	"strings"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	StateMachine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/challenge"
)

func validateChallengeDifficulty(difficulty Entities.ChallengeDifficulty) error {
	switch difficulty {
	case constants.ChallengeDifficultyEasy, constants.ChallengeDifficultyMedium, constants.ChallengeDifficultyHard:
		return nil
	default:
		return fmt.Errorf("invalid challenge difficulty: %q", difficulty)
	}
}

func validateIOFormat(format Entities.VariableFormat) error {
	switch format {
	case constants.VariableFormatString, constants.VariableFormatInt, constants.VariableFormatFloat:
		return nil
	default:
		return fmt.Errorf("invalid io variable format: %q", format)
	}
}

func ValidateIOVariable(v Entities.IOVariable) error {
	if strings.TrimSpace(v.Name) == "" {
		return fmt.Errorf("io variable name is required")
	}

	if err := validateIOFormat(v.Type); err != nil {
		return fmt.Errorf("invalid io variable format: %w", err)
	}

	return nil
}

func ValidateChallenge(challenge *Entities.Challenge) error {
	if challenge == nil {
		return fmt.Errorf("challenge is nil")
	}

	if strings.TrimSpace(challenge.ID) == "" {
		return fmt.Errorf("challenge id is required")
	}
	if strings.TrimSpace(challenge.Title) == "" {
		return fmt.Errorf("challenge title is required")
	}
	if strings.TrimSpace(challenge.Description) == "" {
		return fmt.Errorf("challenge description is required")
	}

	if !StateMachine.IsValidState(challenge.Status) {
		return fmt.Errorf("invalid challenge status: %q", challenge.Status)
	}

	if err := validateChallengeDifficulty(challenge.Difficulty); err != nil {
		return fmt.Errorf("invalid challenge difficulty: %w", err)
	}

	if challenge.WorkerTimeLimit <= 0 {
		return fmt.Errorf("challenge worker time limit must be greater than 0")
	}
	if challenge.WorkerMemoryLimit <= 0 {
		return fmt.Errorf("challenge worker memory limit must be greater than 0")
	}

	if len(challenge.InputVariables) == 0 {
		return fmt.Errorf("challenge must define at least one input variable")
	}
	for _, input := range challenge.InputVariables {
		if err := ValidateIOVariable(input); err != nil {
			return fmt.Errorf("invalid challenge input variable: %w", err)
		}
	}

	if err := ValidateIOVariable(challenge.OutputVariable); err != nil {
		return fmt.Errorf("invalid challenge output variable: %w", err)
	}

	return nil
}
