package postlogin

import userDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"

func MapRequestToInput(req RequestDTO) userDtos.UserLoginInput {
	return userDtos.UserLoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}
