package postclose

import (
"github.com/gofiber/fiber/v2"
container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
examUsecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam/exam_crud"
"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
uc := examUsecases.NewCloseExamUseCase(appContainer.Dependencies.UserRepository, appContainer.Dependencies.ExamRepository)
return func(c *fiber.Ctx) error {
result, err := uc.Execute(shared.BuildRequestContext(c), ToInput(MapPath(c.Params("id"))))
if err != nil { return shared.HandleError(c, err) }
return c.Status(fiber.StatusOK).JSON(result)
}
}
