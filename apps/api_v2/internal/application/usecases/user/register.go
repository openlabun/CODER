package user_usecases

import (
	"strings"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"
)

type RegisterUseCase struct {
	registerService    ports.RegisterPort
	userService 	  ports.UserServicePort
	passwordHasher ports.PasswordHasherPort
}

func NewRegisterUseCase(
	registerService ports.RegisterPort,
	userService ports.UserServicePort,
	passwordHasher ports.PasswordHasherPort,
) *RegisterUseCase {
	return &RegisterUseCase{
		registerService:    registerService,
		userService: 	  userService,
		passwordHasher: passwordHasher,
	}
}

func (uc *RegisterUseCase) Execute(email, name, password string) (*dtos.UserAccess, error) {
	// [STEP 1] Trim password
	password = strings.TrimSpace(password)
	
	// [STEP 2] Hash the provided password
	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	// [STEP 3] Register the user
	user, err := uc.registerService.RegisterUserDirect(email, hashedPassword, name)
	if err != nil {
		return nil, err
	}

	return user, nil
}