package getme

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := shared.BuildRequestContext(c)
		userID := ResolveUserIDFromRequest(c.Query("userId"), c.Get("X-User-Email"))
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "userId or X-User-Email header is required"})
		}

		user, err := appContainer.UserModule.GetData.Execute(ctx, userID)
		if err != nil {
			return shared.HandleError(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(user)
	}
}
