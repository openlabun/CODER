package patchupdate

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type RequestDTO struct {
	Order  *int `json:"order"`
	Points *int `json:"points"`
}

type PathDTO struct{ ID string }

func ToInput(path PathDTO, body RequestDTO) examDtos.UpdateExamItemInput {
	return examDtos.UpdateExamItemInput{
		ID:     path.ID,
		Order:  body.Order,
		Points: body.Points,
	}
}
