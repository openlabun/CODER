package container

import (
	"fmt"

	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"
)


// ApplicationDependencies groups all contract-based dependencies required by
// application use cases.
type ApplicationDependencies struct {
	RegisterService ports.RegisterPort
	LoginService    ports.LoginPort
	UserService     ports.UserServicePort
	TokenService    ports.TokenServicePort
	PasswordHasher  ports.PasswordHasherPort
}

func NewApplicationDependencies(
	registerService ports.RegisterPort,
	loginService ports.LoginPort,
	userService ports.UserServicePort,
	tokenService ports.TokenServicePort,
	passwordHasher ports.PasswordHasherPort,
) ApplicationDependencies {
	return ApplicationDependencies{
		RegisterService: registerService,
		LoginService:    loginService,
		UserService:     userService,
		TokenService:    tokenService,
		PasswordHasher:  passwordHasher,
	}
}

func (deps ApplicationDependencies) CheckDependencies() error {
	if deps.RegisterService == nil {
		return fmt.Errorf("RegisterService dependency is not provided")
	}

	if deps.LoginService == nil {
		return fmt.Errorf("LoginService dependency is not provided")
	}

	if deps.UserService == nil {
		return fmt.Errorf("UserService dependency is not provided")
	}

	if deps.TokenService == nil {
		return fmt.Errorf("TokenService dependency is not provided")
	}

	if deps.PasswordHasher == nil {
		return fmt.Errorf("PasswordHasher dependency is not provided")
	}
	
	return  nil
}