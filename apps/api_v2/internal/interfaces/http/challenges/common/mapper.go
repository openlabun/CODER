package common

import (
	examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	challengeEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

type ChallengeResponseDTO struct {
	ID                string                         `json:"id"`
	Title             string                         `json:"title"`
	Description       string                         `json:"description"`
	Tags              []string                       `json:"tags"`
	Status            string                         `json:"status"`
	Difficulty        string                         `json:"difficulty"`
	WorkerTimeLimit   int                            `json:"worker_time_limit"`
	WorkerMemoryLimit int                            `json:"worker_memory_limit"`
	CodeTemplates     map[string]string              `json:"code_templates"`
	InputVariables    []challengeEntities.IOVariable `json:"input_variables"`
	OutputVariable    challengeEntities.IOVariable   `json:"output_variable"`
	Constraints       string                         `json:"constraints"`
	CreatedAt         string                         `json:"created_at"`
	UpdatedAt         string                         `json:"updated_at"`
	UserID            string                         `json:"user_id"`
}

func MapCodeTemplateMapToDTOs(input map[string]string) []examDtos.CodeTemplateDTO {
	if len(input) == 0 {
		return nil
	}

	out := make([]examDtos.CodeTemplateDTO, 0, len(input))
	for language, template := range input {
		out = append(out, examDtos.CodeTemplateDTO{
			Language: language,
			Template: template,
		})
	}

	return out
}

func MapCodeTemplateDTOsToMap(input []challengeEntities.CodeTemplate) map[string]string {
	if len(input) == 0 {
		return map[string]string{}
	}

	out := make(map[string]string, len(input))
	for _, tpl := range input {
		out[string(tpl.Language)] = tpl.Template
	}

	return out
}

func MapChallengeToResponse(challenge *challengeEntities.Challenge) ChallengeResponseDTO {
	if challenge == nil {
		return ChallengeResponseDTO{}
	}

	return ChallengeResponseDTO{
		ID:                challenge.ID,
		Title:             challenge.Title,
		Description:       challenge.Description,
		Tags:              challenge.Tags,
		Status:            string(challenge.Status),
		Difficulty:        string(challenge.Difficulty),
		WorkerTimeLimit:   challenge.WorkerTimeLimit,
		WorkerMemoryLimit: challenge.WorkerMemoryLimit,
		CodeTemplates:     MapCodeTemplateDTOsToMap(challenge.CodeTemplates),
		InputVariables:    challenge.InputVariables,
		OutputVariable:    challenge.OutputVariable,
		Constraints:       challenge.Constraints,
		CreatedAt:         challenge.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         challenge.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UserID:            challenge.UserID,
	}
}

func MapChallengesToResponse(challenges []*challengeEntities.Challenge) []ChallengeResponseDTO {
	out := make([]ChallengeResponseDTO, 0, len(challenges))
	for _, challenge := range challenges {
		out = append(out, MapChallengeToResponse(challenge))
	}

	return out
}
