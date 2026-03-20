package user_usecases

import (
	dtos "../../dtos/user"
	ports "../../ports/user"
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
	// Hash the provided password
	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	// Repository: validate user credentials
	user, err := uc.userService.LoginUser(email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}