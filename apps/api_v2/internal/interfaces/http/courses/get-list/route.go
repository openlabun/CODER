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
		switch query.Scope {
			case "owned":
				result, err := appContainer.CourseModule.GetOwnedCourses.Execute(ctx)
				if err != nil { 
					return shared.HandleError(c, err) 
				}
				return c.Status(fiber.StatusOK).JSON(result)
			case "enrolled":
				result, err := appContainer.CourseModule.GetEnrolledCourses.Execute(ctx)
				if err != nil { 
					return shared.HandleError(c, err) 
				}
				return c.Status(fiber.StatusOK).JSON(result)
		}
		// TODO: Not implemented get all public courses
		result := []string{}
		
		return c.Status(fiber.StatusNotImplemented).JSON(result)
	}
}
