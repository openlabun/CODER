package postchangevisibility

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type RequestDTO struct { Visibility string `json:"visibility"` }

type PathDTO struct { ID string }

func ToInput(path PathDTO, body RequestDTO) examDtos.ChangeExamVisibilityInput {
return examDtos.ChangeExamVisibilityInput{ExamID: path.ID, Visibility: body.Visibility}
}
