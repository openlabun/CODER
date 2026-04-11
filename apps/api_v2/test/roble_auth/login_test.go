package roble_auth_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	test "github.com/openlabun/CODER/apps/api_v2/test"

	hasher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
)

func TestUserLogin(t *testing.T) {
	process := test.StartTest(t, "[Authentication Test] Roble Login Infrastructure")
	email := "test@test.com"
	password := "Password123!"

	// [STEP 1] Initialize Roble client and auth adapter
	process.StartStep("Initialize Roble client and auth adapter")
	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		process.Fail("initialize roble client", err)
	}

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)

	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	process.EndStep()

	// [STEP 2] Initialize hasher and hash password
	process.StartStep("Initialize hasher and hash password")
	adapter := hasher.NewSecurityAdapter()
	hashedPassword, err := adapter.Hash(password)
	if err != nil {
		process.Fail("hash password", err)
	}
	process.EndStep()

	// [STEP 3] Attempting login with email and password
	process.StartStep("Attempting login with email and password")

	access, err := authAdapter.LoginUser(email, hashedPassword)
	if err != nil {
		process.Fail("login user", err)
	}

	if access == nil {
		process.Fail("login user", fmt.Errorf("expected user access, got nil"))
	}

	if access.Token == nil || access.Token.AccessToken == "" {
		process.Fail("login user", fmt.Errorf("expected non-empty access token"))
	}

	if access.UserData == nil {
		process.Fail("login user", fmt.Errorf("expected user data in login response"))
	}

	process.EndStep()

	process.End()
}
