package postgeneratefullchallenge

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input dtos.GenerateFullChallengeInput
		if err := c.BodyParser(&input); err != nil {
			return shared.HandleError(c, err)
		}

		result, err := appContainer.AIModule.GenerateFullChallenge.Execute(shared.BuildRequestContext(c), input)
		if err != nil {
			return shared.HandleError(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
