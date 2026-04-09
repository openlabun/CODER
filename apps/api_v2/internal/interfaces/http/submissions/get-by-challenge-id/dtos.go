package getbychallengeid

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type PathDTO struct{ ChallengeID string }

type QueryDTO struct {
	Status string `query:"status"`
	TestID string `query:"testId"`
}

func ToInput(path PathDTO, query QueryDTO) submissionDtos.GetChallengeSubmissionsInput {
	var statusPtr *string
	if query.Status != "" {
		statusPtr = &query.Status
	}
	var testIDPtr *string
	if query.TestID != "" {
		testIDPtr = &query.TestID
	}
	return submissionDtos.GetChallengeSubmissionsInput{
		ChallengeID: path.ChallengeID,
		Status:      statusPtr,
		TestID:      testIDPtr,
	}
}
