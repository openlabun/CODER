package postcreate

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	common "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/common"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RequestDTO
		if err := c.BodyParser(&req); err != nil {
			return shared.HandleError(c, err)
		}
		result, err := appContainer.ChallengeModule.CreateChallenge.Execute(
			shared.BuildRequestContext(c),
			MapRequestToInput(req),
		)
		if err != nil {
			return shared.HandleError(c, err)
		}
		return c.Status(fiber.StatusCreated).JSON(common.MapChallengeToResponse(result))
	}
}
