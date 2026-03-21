package roble_infrastructure

import (
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user/mapper"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	UserFactory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/user"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

type RobleAuthAdapter struct {
	adapter    *infrastructure.RobleDatabaseAdapter
	repository *UserRepository
}

func NewRobleAuthAdapter(adapter *infrastructure.RobleDatabaseAdapter, repository *UserRepository) *RobleAuthAdapter {
	return &RobleAuthAdapter{adapter: adapter, repository: repository}
}

func (a *RobleAuthAdapter) GetUserData(email string) (*Entities.User, error) {
	return a.repository.GetUserByEmail(email)
}

func (a *RobleAuthAdapter) LoginUser(email, password string) (*dtos.UserAccess, error) {
	client := a.adapter.GetClient()

	// Authenticate user with Roble and get tokens
	tokens, err := client.Login(email, password)
	if err != nil {
		return nil, err
	}

	// Set access token for subsequent requests
	a.adapter.SetAccessToken(tokens.AccessToken)

	user, err := a.GetUserData(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user is not registered")
	}

	user_access := mapper.MapToUserAccessDTO(
		user,
		tokens.AccessToken,
		tokens.RefreshToken,
	)

	return user_access, nil
}

func (a *RobleAuthAdapter) RegisterUser(email, name, password string) (bool, error) {
	client := a.adapter.GetClient()

	// Register user in Roble and expect success message
	_, err := client.Signup(email, password, name)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *RobleAuthAdapter) RegisterUserDirect(email, password, name string) (*dtos.UserAccess, error) {
	client := a.adapter.GetClient()

	// Register user in Roble and expect success message
	_, err := client.SignupDirect(email, password, name)
	if err != nil {
		return nil, err
	}

	tokens, err := client.Login(email, password)
	if err != nil {
		return nil, err
	}

	// Set access token for subsequent requests
	a.adapter.SetAccessToken(tokens.AccessToken)

	// Create user entity (validations handled in factory)
	user, err := UserFactory.NewUser(
		tokens.User.ID,
		name,
		email,
		password,
	)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("failed to create user")
	}

	// Save user data
	_, err = a.repository.SaveUser(user)
	if err != nil {
		return nil, err
	}

	user_access := mapper.MapToUserAccessDTO(
		user,
		tokens.AccessToken,
		tokens.RefreshToken,
	)

	return user_access, nil
}

func (a *RobleAuthAdapter) RefreshUserToken(refreshToken string) (*dtos.Token, error) {
	client := a.adapter.GetClient()
	tokens, err := client.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return &dtos.Token{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (a *RobleAuthAdapter) VerifyEmail(email, code string) (bool, error) {
	client := a.adapter.GetClient()
	if err := client.VerifyEmail(email, code); err != nil {
		return false, err
	}

	return true, nil
}

func (a *RobleAuthAdapter) ForgotPassword(email string) error {
	client := a.adapter.GetClient()
	return client.ForgotPassword(email)
}

func (a *RobleAuthAdapter) ResetPassword(token, newPassword string) error {
	client := a.adapter.GetClient()
	return client.ResetPassword(token, newPassword)
}

func (a *RobleAuthAdapter) Logout(accessToken string) error {
	client := a.adapter.GetClient()
	return client.Logout(accessToken)
}

func (a *RobleAuthAdapter) VerifyToken(accessToken string) (*infrastructure.RobleVerifyTokenResponse, error) {
	client := a.adapter.GetClient()
	return client.VerifyToken(accessToken)
}
