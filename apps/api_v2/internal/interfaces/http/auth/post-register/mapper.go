package postregister

import userDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"

func MapRequestToInput(req RequestDTO) userDtos.UserRegisterInput {
	return userDtos.UserRegisterInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}
}
