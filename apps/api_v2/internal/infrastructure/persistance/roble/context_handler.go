package roble_infrastructure

import (
	"context"
	"fmt"
	service "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

func SetAdapterTokenFromContext(ctx context.Context, adapter *RobleDatabaseAdapter) error {
	if adapter == nil {
		return fmt.Errorf("roble adapter is nil")
	}

	token, ok := service.AccessTokenFromContext(ctx)
	if !ok {
		return fmt.Errorf("access token is required in context")
	}

	adapter.SetAccessToken(token)
	return nil
}
