package postgenerateexam

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/ai"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dtos.GenerateExamInput
		if err := c.BodyParser(&req); err != nil {
			return shared.HandleError(c, err)
		}

		result, err := appContainer.AIModule.GenerateExam.Execute(
			shared.BuildRequestContext(c),
			req,
		)
		if err != nil {
			return shared.HandleError(c, err)
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
