package postrefreshtoken

func MapRequestToInput(req RequestDTO) string {
	return req.RefreshToken
}
