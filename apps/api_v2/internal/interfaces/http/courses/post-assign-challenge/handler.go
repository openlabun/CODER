package postassignchallenge

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

type RequestDTO struct {
	ChallengeID string `json:"challengeId"`
}

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courseID := c.Params("id")
		var req RequestDTO
		if err := c.BodyParser(&req); err != nil {
			return shared.HandleError(c, err)
		}
		
		err := appContainer.CourseModule.AssignChallengeToCourse.Execute(
			shared.BuildRequestContext(c), 
			courseID,
			req.ChallengeID,
		)
		if err != nil { 
			return shared.HandleError(c, err) 
		}
		
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Challenge assigned to course successfully",
		})
	}
}
