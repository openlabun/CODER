package postcreate

import (
"github.com/gofiber/fiber/v2"
container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
submissionUsecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/submission"
"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
uc := submissionUsecases.NewCreateSubmissionUseCase(
appContainer.Dependencies.UserRepository,
appContainer.Dependencies.SubmissionRepository,
appContainer.Dependencies.SessionRepository,
appContainer.Dependencies.ChallengeRepository,
appContainer.Dependencies.TestCaseRepository,
appContainer.Dependencies.SubmissionResultRepository,
)
return func(c *fiber.Ctx) error {
var req RequestDTO
if err := c.BodyParser(&req); err != nil { return shared.HandleError(c, err) }
result, err := uc.Execute(shared.BuildRequestContext(c), MapRequestToInput(req))
if err != nil { return shared.HandleError(c, err) }
return c.Status(fiber.StatusCreated).JSON(result)
}
}
