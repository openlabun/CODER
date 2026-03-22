package services

import (
	"context"
	"fmt"
	"strings"
)

type ContextKey string

const UserEmailContextKey ContextKey = "userEmail"

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
