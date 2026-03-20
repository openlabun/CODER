package mapper

import (
	dtos "../" 
	Entities "../../../../domain/entities/user"
)

func MapToUserAccessDTO(user *Entities.User, accessToken, refreshToken string) *dtos.UserAccess {
	return &dtos.UserAccess{
		UserData: user,
		Token: &dtos.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}