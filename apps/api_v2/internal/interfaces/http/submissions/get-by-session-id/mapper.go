package getbysessionid

func MapPath(sessionId string) PathDTO {
	return PathDTO{SessionID: sessionId}
}

func MapQuery(status, testId, challengeId string) QueryDTO {
	return QueryDTO{Status: status, TestID: testId, ChallengeID: challengeId}
}
