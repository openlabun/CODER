package getbyuserid

func MapPath(userID string) PathDTO { return PathDTO{UserID: userID} }

func MapQuery(status, testID, challengeID string) QueryDTO {
	return QueryDTO{Status: status, TestID: testID, ChallengeID: challengeID}
}
