package getitemscoring

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type PathDTO struct{ ExamScoreID string }

func ToInput(path PathDTO) submissionDtos.GetItemScoringInput {
	return submissionDtos.GetItemScoringInput{ExamScoreID: path.ExamScoreID}
}
