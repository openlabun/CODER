package getbyid

import (
"github.com/gofiber/fiber/v2"
container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
return func(c *fiber.Ctx) error {
req := MapPathToRequest(c.Params("id"))
result, err := appContainer.ExamModule.GetExamDetails.Execute(shared.BuildRequestContext(c), req.ToInput())
if err != nil { return shared.HandleError(c, err) }
return c.Status(fiber.StatusOK).JSON(result)
}
}
