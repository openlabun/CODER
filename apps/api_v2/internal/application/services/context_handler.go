package services

import (
	"context"
	"fmt"
	"strings"
)

type ContextKey string

const (
	UserEmailContextKey ContextKey = "userEmail"
	AccessTokenContextKey ContextKey = "accessToken"
	InternalServiceContextKey ContextKey = "workerKey"
)

func WithUserEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, UserEmailContextKey, strings.TrimSpace(email))
}

func UserEmailFromContext(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context is required")
	}

	email, ok := ctx.Value(UserEmailContextKey).(string)
	if !ok || strings.TrimSpace(email) == "" {
		return "", fmt.Errorf("user email is required in context")
	}

	return strings.TrimSpace(email), nil
}


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

func WithInternalServiceKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, InternalServiceContextKey, strings.TrimSpace(key))
}

func AccessInternalServiceKeyFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	if key, ok := ctx.Value(InternalServiceContextKey).(string); ok {
		key = strings.TrimSpace(key)
		if key != "" {
			return key, true
		}
	}

	return "", false
}