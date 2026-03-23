package getbyuserid

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := MapPath(c.Params("userId"))
		query := MapQuery(c.Query("status"), c.Query("testId"), c.Query("challengeId"))
		result, err := appContainer.SubmissionUseCases.GetUserSubmissions.Execute(shared.BuildRequestContext(c), ToInput(path, query))
		if err != nil {
			return shared.HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(result)
	}
}
