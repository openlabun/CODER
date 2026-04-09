package mapper

import (
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user" 
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
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