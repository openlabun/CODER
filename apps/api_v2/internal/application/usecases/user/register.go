package user_usecases

import (
	"fmt"
	
	dtos "../../dtos/user"
	ports "../../ports/user"
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
	// Hash the provided password
	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}

	// Check if user already exists by email
	existingUser, _ := uc.userService.GetUserByEmail(email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Repository: validate user credentials
	user, err := uc.registerService.RegisterUserDirect(email, name, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}