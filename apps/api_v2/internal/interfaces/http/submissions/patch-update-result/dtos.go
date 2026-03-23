package patchupdateresult

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type RequestDTO struct {
	Status        string  `json:"status"`
	TimeExecution int     `json:"timeExecution"`
	Output        *string `json:"output"`
	Error         *string `json:"error"`
}

type PathDTO struct {
	ResultID string
}

func ToInput(path PathDTO, body RequestDTO) submissionDtos.UpdateResultInput {
	return submissionDtos.UpdateResultInput{
		ResultID:      path.ResultID,
		Status:        body.Status,
		TimeExecution: body.TimeExecution,
		Output:        body.Output,
		Error:         body.Error,
	}
}
