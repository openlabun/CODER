package getlist

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type QueryDTO struct { 
	ExamID string `query:"examId"` 
}

func ToInput(q QueryDTO) examDtos.GetChallengesByExamInput { 
	return examDtos.GetChallengesByExamInput{
		ExamID: q.ExamID,
	} 
}
