package shared

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	roble "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

func BuildRequestContext(c *fiber.Ctx) context.Context {
	ctx := context.Background()

	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		token := strings.TrimSpace(authHeader[7:])
		if token != "" {
			ctx = roble.WithAccessToken(ctx, token)
		}
	}

	email := strings.TrimSpace(c.Get("X-User-Email"))
	if email == "" {
		email = strings.TrimSpace(c.Query("userEmail"))
	}
	if email != "" {
		ctx = services.WithUserEmail(ctx, email)
	}

	return ctx
}
