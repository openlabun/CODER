package getbyid

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type RequestDTO struct { ID string }

func (r RequestDTO) ToInput() examDtos.GetExamDetailsInput { return examDtos.GetExamDetailsInput{ExamID: r.ID} }
