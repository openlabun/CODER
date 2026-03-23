package patchupdate

import (
"github.com/gofiber/fiber/v2"
container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
return func(c *fiber.Ctx) error {
var body RequestDTO
if err := c.BodyParser(&body); err != nil { return shared.HandleError(c, err) }
input := ToInput(MapPath(c.Params("id")), body)
result, err := appContainer.ExamModule.UpdateExam.Execute(shared.BuildRequestContext(c), input)
if err != nil { return shared.HandleError(c, err) }
return c.Status(fiber.StatusOK).JSON(result)
}
}
