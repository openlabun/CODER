package patchupdate

import (
	examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	common "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/common"
)

type RequestDTO struct {
	Title             *string                   `json:"title"`
	Description       *string                   `json:"description"`
	Tags              *[]string                 `json:"tags"`
	Status            *string                   `json:"status"`
	Difficulty        *string                   `json:"difficulty"`
	WorkerTimeLimit   *int                      `json:"worker_time_limit"`
	WorkerMemoryLimit *int                      `json:"worker_memory_limit"`
	CodeTemplates     *map[string]string        `json:"code_templates"`
	InputVariables    *[]examDtos.IOVariableDTO `json:"input_variables"`
	OutputVariable    *examDtos.IOVariableDTO   `json:"output_variable"`
	Constraints       *string                   `json:"constraints"`
}

type PathDTO struct {
	ID string
}

func ToInput(path PathDTO, body RequestDTO) examDtos.UpdateChallengeInput {
	var codeTemplates *[]examDtos.CodeTemplateDTO
	if body.CodeTemplates != nil {
		mapped := common.MapCodeTemplateMapToDTOs(*body.CodeTemplates)
		codeTemplates = &mapped
	}

	return examDtos.UpdateChallengeInput{
		ChallengeID:       path.ID,
		Title:             body.Title,
		Description:       body.Description,
		Tags:              body.Tags,
		Status:            body.Status,
		Difficulty:        body.Difficulty,
		WorkerTimeLimit:   body.WorkerTimeLimit,
		WorkerMemoryLimit: body.WorkerMemoryLimit,
		CodeTemplates:     codeTemplates,
		InputVariables:    body.InputVariables,
		OutputVariable:    body.OutputVariable,
		Constraints:       body.Constraints,
	}
}
