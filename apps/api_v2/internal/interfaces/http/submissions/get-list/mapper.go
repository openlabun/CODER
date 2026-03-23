package getlist

func MapQuery(challengeID, status, testID string) QueryDTO {
return QueryDTO{ChallengeID: challengeID, Status: status, TestID: testID}
}
