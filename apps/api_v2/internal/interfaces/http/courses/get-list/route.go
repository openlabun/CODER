package getlist

import (
"github.com/gofiber/fiber/v2"
container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := MapQuery(c.Query("scope"), c.Query("studentId"), c.Query("teacherId"))
		ctx := shared.BuildRequestContext(c)
		if query.Scope == "owned" {
			result, err := appContainer.CourseModule.GetOwnedCourses.Execute(ctx, ToOwnedInput(query))
			if err != nil { 
				return shared.HandleError(c, err) 
			}
			return c.Status(fiber.StatusOK).JSON(result)
		}
		result, err := appContainer.CourseModule.GetEnrolledCourses.Execute(ctx, ToEnrolledInput(query))
		if err != nil { 
			return shared.HandleError(c, err) 
		}
		return c.Status(fiber.StatusOK).JSON(result)
	}
}
