package getbychallengeid

import (
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	"github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/shared"
)

func Handler(appContainer *container.Application) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := MapPath(c.Params("challengeId"))

		// exam_id es opcional
		var examIDPtr *string
		examID := c.Query("exam_id", "")
		if examID != "" {
			examIDPtr = &examID
		}

		input := ToInput(path)
		input.ExamID = examIDPtr

		result, err := appContainer.TestCaseModule.GetTestCasesByChallenge.Execute(shared.BuildRequestContext(c), input)
		if err != nil {
			return shared.HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(result)
	}
}
