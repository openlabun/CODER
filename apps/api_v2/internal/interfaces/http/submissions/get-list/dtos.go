package getlist

import submissionDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

type QueryDTO struct {
ChallengeID string `query:"challengeId"`
Status      string `query:"status"`
TestID      string `query:"testId"`
}

func ToInput(q QueryDTO) submissionDtos.GetChallengeSubmissionsInput {
var statusPtr *string
if q.Status != "" { statusPtr = &q.Status }
var testIDPtr *string
if q.TestID != "" { testIDPtr = &q.TestID }
return submissionDtos.GetChallengeSubmissionsInput{ChallengeID: q.ChallengeID, Status: statusPtr, TestID: testIDPtr}
}
