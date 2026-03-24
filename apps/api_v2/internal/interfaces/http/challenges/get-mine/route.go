package getmine

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := appContainer.ChallengeModule.GetMyChallenges.Execute(shared.BuildRequestContext(c))
		if err != nil {
			return shared.HandleError(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
