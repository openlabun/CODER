package user_usecases

import (
	"strings"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"
)

type LoginUseCase struct {
	userService    ports.LoginPort
	passwordHasher ports.PasswordHasherPort
}

func NewLoginUseCase(
	userService ports.LoginPort,
	passwordHasher ports.PasswordHasherPort,
) *LoginUseCase {
	return &LoginUseCase{
		userService:    userService,
		passwordHasher: passwordHasher,
	}
}

func (uc *LoginUseCase) Execute(email, password string) (*dtos.UserAccess, error) {
	// [STEP 1] Trim password
	password = strings.TrimSpace(password)
	
	// [STEP 2] Hash the provided password
	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	// [STEP 3] Validate user credentials
	user, err := uc.userService.LoginUser(email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}