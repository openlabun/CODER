package getbysessionid

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type PathDTO struct{ SessionID string }

type QueryDTO struct {
	Status      string `query:"status"`
	TestID      string `query:"testId"`
	ChallengeID string `query:"challengeId"`
}

func ToInput(path PathDTO, query QueryDTO) submissionDtos.GetSessionSubmissionsInput {
	var statusPtr *string
	if query.Status != "" {
		statusPtr = &query.Status
	}
	var testIDPtr *string
	if query.TestID != "" {
		testIDPtr = &query.TestID
	}
	var challengeIDPtr *string
	if query.ChallengeID != "" {
		challengeIDPtr = &query.ChallengeID
	}
	return submissionDtos.GetSessionSubmissionsInput{
		SessionID:   path.SessionID,
		Status:      statusPtr,
		TestID:      testIDPtr,
		ChallengeID: challengeIDPtr,
	}
}
