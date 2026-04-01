package getbychallengeid

func MapPath(challengeId string) PathDTO {
	return PathDTO{ChallengeID: challengeId}
}

func MapQuery(status, testId string) QueryDTO {
	return QueryDTO{Status: status, TestID: testId}
}
