package roble_auth_test

import (
	"net/http"
	"testing"
	"time"
	"fmt"

	test "github.com/openlabun/CODER/apps/api_v2/test"

	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
)

func TestUserRegistration(t *testing.T) {
	process := test.StartTest(t, "[Authentication Test] Roble User Registration")
	random_value := time.Now().UnixNano()
	email := fmt.Sprintf("test%d@test.com", random_value)
	password := "Password123!"
	name := fmt.Sprintf("Test User %d", random_value)
	
	// [STEP 1] Initialize Roble client and auth adapter
	process.StartStep("Initialize Roble client and auth adapter")
	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		t.Fatalf("initialize roble client: %v", err)
	}

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)

	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	process.EndStep()

	// [STEP 2] Attempting direct registration with email, password, and name
	process.StartStep("Attempting direct registration with email, password, and name")
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

	process.EndStep()

	process.End()
}
