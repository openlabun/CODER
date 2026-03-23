package getlist

import (
"github.com/gofiber/fiber/v2"
container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
submissionUsecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/submission"
"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
uc := submissionUsecases.NewGetChallengeSubmissionsUseCase(
appContainer.Dependencies.UserRepository,
appContainer.Dependencies.ChallengeRepository,
appContainer.Dependencies.ExamRepository,
appContainer.Dependencies.SubmissionRepository,
appContainer.Dependencies.SubmissionResultRepository,
)
return func(c *fiber.Ctx) error {
query := MapQuery(c.Query("challengeId"), c.Query("status"), c.Query("testId"))
result, err := uc.Execute(shared.BuildRequestContext(c), ToInput(query))
if err != nil { return shared.HandleError(c, err) }
return c.Status(fiber.StatusOK).JSON(result)
}
}
