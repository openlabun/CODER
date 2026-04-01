package getbysessionid

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// `session_id` path param removed. Accept optional `user_id` via query.
		path := MapPath(c.Query("user_id"))
		result, err := appContainer.SessionModule.GetActiveSession.Execute(shared.BuildRequestContext(c), ToInput(path))
		if err != nil {
			return shared.HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(result)
	}
}
