package postcreate

import (
	examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	common "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/common"
)

func MapRequestToInput(req RequestDTO) examDtos.CreateChallengeInput {
	return examDtos.CreateChallengeInput{
		Title:             req.Title,
		Description:       req.Description,
		Tags:              req.Tags,
		Status:            req.Status,
		Difficulty:        req.Difficulty,
		WorkerTimeLimit:   req.WorkerTimeLimit,
		WorkerMemoryLimit: req.WorkerMemoryLimit,
		CodeTemplates:     common.MapCodeTemplateMapToDTOs(req.CodeTemplates),
		InputVariables:    req.InputVariables,
		OutputVariable:    req.OutputVariable,
		Constraints:       req.Constraints,
	}
}
