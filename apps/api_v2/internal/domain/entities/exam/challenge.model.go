package challenge_entities

import (
	"time"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	sub_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
)

type Challenge struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`

	// State and Access Control
	Status     constants.ChallengeStatus     `json:"status"` // state machine: draft -> published|private -> archived
	Difficulty constants.ChallengeDifficulty `json:"difficulty"`

	// Worker Constraints
	WorkerTimeLimit   int `json:"worker_time_limit"`   // in ms
	WorkerMemoryLimit int `json:"worker_memory_limit"` // in MB

	// Code Templates
	CodeTemplates []CodeTemplate `json:"code_templates"`

	// I/O Considerations
	InputVariables []IOVariable `json:"input_variables"`
	OutputVariable IOVariable   `json:"output_variable"`
	Constraints    string       `json:"constraints"`

	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
}

func (c *Challenge) GetLanguageTemplate(language sub_constants.ProgrammingLanguage) *CodeTemplate {
	for _, template := range c.CodeTemplates {
		if template.Language == language {
			return &template
		}
	}
	return nil
}
