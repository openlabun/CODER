package getbyid

import (
"github.com/gofiber/fiber/v2"
container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
submissionUsecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/submission"
"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
uc := submissionUsecases.NewGetSubmissionStatusUseCase(
appContainer.Dependencies.UserRepository,
appContainer.Dependencies.SubmissionResultRepository,
appContainer.Dependencies.SubmissionRepository,
)
return func(c *fiber.Ctx) error {
result, err := uc.Execute(shared.BuildRequestContext(c), ToInput(MapPath(c.Params("id"))))
if err != nil { return shared.HandleError(c, err) }
return c.Status(fiber.StatusOK).JSON(result)
}
}
