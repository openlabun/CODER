package patchupdate

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type RequestDTO struct {
	Name           *string                   `json:"name"`
	Input          *[]examDtos.IOVariableDTO `json:"input"`
	ExpectedOutput *examDtos.IOVariableDTO   `json:"expectedOutput"`
	IsSample       *bool                     `json:"isSample"`
	Points         *int                      `json:"points"`
}

type PathDTO struct{ ID string }

func ToInput(path PathDTO, body RequestDTO) examDtos.UpdateTestCaseInput {
	return examDtos.UpdateTestCaseInput{
		ID:             path.ID,
		Name:           body.Name,
		Input:          body.Input,
		ExpectedOutput: body.ExpectedOutput,
		IsSample:       body.IsSample,
		Points:         body.Points,
	}
}
