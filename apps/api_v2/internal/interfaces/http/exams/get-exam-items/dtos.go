package getexamitems

import dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type QueryDTO struct {
	ExamID string `params:"examId"`
}

func ToInput(q QueryDTO) dtos.GetExamItemsInput {
	return dtos.GetExamItemsInput{
		ExamID: q.ExamID,
	}
}
