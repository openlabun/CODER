package test_utils

import (
	"context"
	"testing"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	user_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
)

func EnsureAuthUserAccess(t *testing.T, app *container.Application, email, password, name string) *user_dtos.UserAccess {
	t.Helper()
	access, err := app.UserModule.Login.Execute(email, password)
	if err != nil {
		t.Logf("login failed for %s: %v", email, err)
	}

	if err == nil && access != nil && access.UserData != nil && access.UserData.ID != "" && access.Token != nil && access.Token.AccessToken != "" {
		return access
	}

	registered, registerErr := app.UserModule.Register.Execute(email, name, password)
	if registerErr != nil {
		t.Fatalf("register user failed for %s: %v", email, registerErr)
	}
	if registered == nil || registered.UserData == nil || registered.UserData.ID == "" || registered.Token == nil || registered.Token.AccessToken == "" {
		t.Fatalf("expected registered user with valid access for %s", email)
	}

	return registered
}

func buildContext(token string, email string) context.Context {
	ctx := context.Background()
	ctx = services.WithAccessToken(ctx, token)
	ctx = services.WithUserEmail(ctx, email)
	return ctx
}

func BuildUserCtx(access *user_dtos.UserAccess) context.Context {
	return buildContext(access.Token.AccessToken, access.UserData.Email)
}

