package container

import (
	"fmt"

	user_usecases "../usecases/user"
)



// UserUseCases holds all user-related use cases available in the application.
type UserUseCases struct {
	Register     *user_usecases.RegisterUseCase
	Login        *user_usecases.LoginUseCase
	GetData      *user_usecases.GetDataUseCase
	RefreshToken *user_usecases.RefreshTokenUseCase
}

type Application struct {
	Dependencies ApplicationDependencies
	UserModule   UserUseCases
}

func NewApplication(deps ApplicationDependencies) (*Application, error) {

	if err := deps.CheckDependencies(); err != nil {
		return nil, fmt.Errorf("application dependencies check failed: %w", err)
	}
	
	app := &Application{Dependencies: deps}
	app.UserModule = UserUseCases{
		Register: user_usecases.NewRegisterUseCase(
			deps.RegisterService,
			deps.UserService,
			deps.PasswordHasher,
		),
		Login: user_usecases.NewLoginUseCase(
			deps.LoginService,
			deps.PasswordHasher,
		),
		GetData: user_usecases.NewGetDataUseCase(deps.LoginService),
		RefreshToken: user_usecases.NewRefreshTokenUseCase(
			deps.TokenService,
		),
	}

	return app, nil
}
