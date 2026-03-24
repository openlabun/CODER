package getchallenges

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		courseID := c.Params("id")
		
		challenges, err := appContainer.CourseModule.GetCourseChallenges.Execute(
			shared.BuildRequestContext(c), 
			courseID,
		)
		if err != nil { 
			return shared.HandleError(c, err) 
		}
		
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"challenges": challenges,
		})
	}
}
