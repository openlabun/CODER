package exam_validations

import (
	"fmt"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func ValidateTestCase(testCase *Entities.TestCase) error {
	if testCase == nil {
		return fmt.Errorf("test case is nil")
	}

	if strings.TrimSpace(testCase.ID) == "" {
		return fmt.Errorf("test case id is required")
	}

	if strings.TrimSpace(testCase.Name) == "" {
		return fmt.Errorf("test case name is required")
	}

	if strings.TrimSpace(testCase.ChallengeID) == "" {
		return fmt.Errorf("test case challenge id is required")
	}

	if len(testCase.Input) == 0 {
		return fmt.Errorf("test case must define at least one input variable")
	}

	for _, input := range testCase.Input {
		if err := ValidateIOVariable(input); err != nil {
			return fmt.Errorf("invalid test case input variable: %w", err)
		}
	}

	if err := ValidateIOVariable(testCase.ExpectedOutput); err != nil {
		return fmt.Errorf("invalid test case expected output: %w", err)
	}

	if testCase.IsSample && testCase.Points < 0 {
		return fmt.Errorf("test case points cannot be negative")
	}

	return nil
}
