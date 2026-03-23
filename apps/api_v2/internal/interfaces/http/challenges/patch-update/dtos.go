package patchupdate

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type RequestDTO struct {
	Title             *string                `json:"title"`
	Description       *string                `json:"description"`
	Tags              *[]string              `json:"tags"`
	Status            *string                `json:"status"`
	Difficulty        *string                `json:"difficulty"`
	WorkerTimeLimit   *int                   `json:"workerTimeLimit"`
	WorkerMemoryLimit *int                   `json:"workerMemoryLimit"`
	InputVariables    *[]examDtos.IOVariableDTO `json:"inputVariables"`
	OutputVariable    *examDtos.IOVariableDTO   `json:"outputVariable"`
	Constraints       *string                `json:"constraints"`
}

type PathDTO struct { 
	ID string 
}

func ToInput(path PathDTO, body RequestDTO) examDtos.UpdateChallengeInput {
	return examDtos.UpdateChallengeInput{
		ChallengeID: path.ID,
		Title: body.Title,
		Description: body.Description,
		Tags: body.Tags,
		Status: body.Status,
		Difficulty: body.Difficulty,
		WorkerTimeLimit: body.WorkerTimeLimit,
		WorkerMemoryLimit: body.WorkerMemoryLimit,
		InputVariables: body.InputVariables,
		OutputVariable: body.OutputVariable,
		Constraints: body.Constraints,
	}
}
