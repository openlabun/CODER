package challenge_entities

import (
	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
)

type CodeTemplate struct {
	Language   constants.ProgrammingLanguage 	`json:"language"`
	Template   string 							`json:"template"`
}