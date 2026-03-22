package roble_infrastructure

import (
	"context"
	"fmt"
	"strings"
)

type ContextKey string

const (
	AccessTokenContextKey ContextKey = "accessToken"
)

func WithAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, AccessTokenContextKey, strings.TrimSpace(accessToken))
}

func AccessTokenFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	if token, ok := ctx.Value(AccessTokenContextKey).(string); ok {
		token = strings.TrimSpace(token)
		if token != "" {
			return token, true
		}
	}

	return "", false
}

func SetAdapterTokenFromContext(ctx context.Context, adapter *RobleDatabaseAdapter) error {
	if adapter == nil {
		return fmt.Errorf("roble adapter is nil")
	}

	token, ok := AccessTokenFromContext(ctx)
	if !ok {
		return fmt.Errorf("access token is required in context")
	}

	adapter.SetAccessToken(token)
	return nil
}
