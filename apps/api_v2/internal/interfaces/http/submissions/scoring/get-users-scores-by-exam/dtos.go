package getusersscoresbyexam

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type PathDTO struct{ ExamID string }

type QueryDTO struct{ CourseID *string }

func ToInput(path PathDTO, query QueryDTO) submissionDtos.GetUserScoresByExamInput {
	return submissionDtos.GetUserScoresByExamInput{ExamID: path.ExamID, CourseID: query.CourseID}
}
