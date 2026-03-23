package roble_auth_test

import (
	"testing"
	"net/http"
	"time"

	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
)

func TestUserRegistration(t *testing.T) {
	t.Skip("use TestRobleRegistrationInfrastructure for integration validation")
}

func TestUserDirectRegistration(t *testing.T) {
	t.Skip("use TestRobleRegistrationInfrastructure for integration validation")
}

func TestRobleRegistrationInfrastructure(t *testing.T) {
	email := "test@test.com"
	password := "Testing123!"
	name := "Test User"

	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		t.Fatalf("initialize roble client: %v", err)
	}

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)

	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)

	access, err := authAdapter.RegisterUserDirect(email, password, name)
	if err != nil {
		t.Fatalf("direct registration failed: %v", err)
	}

	if access == nil {
		t.Fatal("expected user access, got nil")
	}

	if access.Token == nil || access.Token.AccessToken == "" {
		t.Fatal("expected non-empty access token after registration")
	}
}
