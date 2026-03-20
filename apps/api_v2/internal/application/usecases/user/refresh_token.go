package user_usecases

import (
	dtos "../../dtos/user"
	ports "../../ports/user"
)

type RefreshTokenUseCase struct {
	tokenService ports.TokenServicePort
}

func NewRefreshTokenUseCase(tokenService ports.TokenServicePort) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{tokenService: tokenService}
}

func (uc *RefreshTokenUseCase) Execute(refresh_token string) (*dtos.Token, error) {

	// Repository: validate refresh token and issue new access token
	token, err := uc.tokenService.RefreshUserToken(refresh_token)
	if err != nil {
		return nil, err
	}

	return token, nil
}