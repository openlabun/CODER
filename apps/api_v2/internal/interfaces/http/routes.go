package http_interfaces

import (
	"github.com/gofiber/fiber/v2"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

func RegisterRoutes(app *fiber.App) {
	registerDocsRoutes(app)
	registerAuthRoutes(app)
	registerAIRoutes(app)
	registerChallengesRoutes(app)
	registerTestCasesRoutes(app)
	registerCoursesRoutes(app)
	registerExamsRoutes(app)
	registerSubmissionsRoutes(app)
	registerLeaderboardRoutes(app)
	registerMetricsRoutes(app)
	registerHealthRoutes(app)
	registerCacheRoutes(app)
	registerDBRoutes(app)
}

func registerDocsRoutes(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(
			&scalar.Options{
				SpecURL: "./docs/openapi.yaml",
				CustomOptions: scalar.CustomOptions{
					PageTitle: "Artemisa Source-Search Service",
				},
				DarkMode: true,
			},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error generando la referencia")
		}

		return c.Type("html").SendString(htmlContent)
	})
}

func registerAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", mockHandler("auth/post-register.json", fiber.StatusCreated))
	auth.Post("/login", mockHandler("auth/post-login.json", fiber.StatusOK))
	auth.Get("/me", mockHandler("auth/get-me.json", fiber.StatusOK))
}

func registerAIRoutes(app *fiber.App) {
	ai := app.Group("/ai")
	ai.Post("/generate-challenge-ideas", mockHandler("ai/post-generate-challenge-ideas.json", fiber.StatusOK))
	ai.Post("/generate-test-cases", mockHandler("ai/post-generate-test-cases.json", fiber.StatusOK))
}

func registerChallengesRoutes(app *fiber.App) {
	challenges := app.Group("/challenges")
	challenges.Post("/", mockHandler("challenges/post-create.json", fiber.StatusCreated))
	challenges.Get("/", mockHandler("challenges/get-list.json", fiber.StatusOK))
	challenges.Get("/:id", mockHandler("challenges/get-by-id.json", fiber.StatusOK))
	challenges.Patch("/:id", mockHandler("challenges/patch-update.json", fiber.StatusOK))
	challenges.Post("/:id/publish", mockHandler("challenges/post-publish.json", fiber.StatusOK))
	challenges.Post("/:id/archive", mockHandler("challenges/post-archive.json", fiber.StatusOK))
}

func registerTestCasesRoutes(app *fiber.App) {
	testCases := app.Group("/test-cases")
	testCases.Post("/", mockHandler("test-cases/post-create.json", fiber.StatusCreated))
	testCases.Get("/challenge/:challengeId", mockHandler("test-cases/get-by-challenge-id.json", fiber.StatusOK))
	testCases.Delete("/:id", mockHandler("test-cases/delete-by-id.json", fiber.StatusOK))
}

func registerCoursesRoutes(app *fiber.App) {
	courses := app.Group("/courses")
	courses.Post("/enroll", mockHandler("courses/post-enroll.json", fiber.StatusOK))
	courses.Post("/", mockHandler("courses/post-create.json", fiber.StatusCreated))
	courses.Get("/browse", mockHandler("courses/get-browse.json", fiber.StatusOK))
	courses.Get("/", mockHandler("courses/get-list.json", fiber.StatusOK))
	courses.Get("/:id", mockHandler("courses/get-by-id.json", fiber.StatusOK))
	courses.Post("/:id", mockHandler("courses/post-update.json", fiber.StatusOK))
	courses.Post("/:id/students", mockHandler("courses/post-add-student.json", fiber.StatusOK))
	courses.Delete("/:id/students/:studentId", mockHandler("courses/delete-student.json", fiber.StatusOK))
	courses.Post("/:id/challenges", mockHandler("courses/post-assign-challenge.json", fiber.StatusOK))
	courses.Get("/:id/students", mockHandler("courses/get-students.json", fiber.StatusOK))
	courses.Get("/:id/challenges", mockHandler("courses/get-challenges.json", fiber.StatusOK))
}

func registerExamsRoutes(app *fiber.App) {
	exams := app.Group("/exams")
	exams.Post("/", mockHandler("exams/post-create.json", fiber.StatusCreated))
	exams.Get("/course/:courseId", mockHandler("exams/get-by-course-id.json", fiber.StatusOK))
	exams.Get("/:id", mockHandler("exams/get-by-id.json", fiber.StatusOK))
}

func registerSubmissionsRoutes(app *fiber.App) {
	submissions := app.Group("/submissions")
	submissions.Post("/", mockHandler("submissions/post-create.json", fiber.StatusCreated))
	submissions.Get("/:id", mockHandler("submissions/get-by-id.json", fiber.StatusOK))
	submissions.Get("/", mockHandler("submissions/get-list.json", fiber.StatusOK))
}

func registerLeaderboardRoutes(app *fiber.App) {
	leaderboard := app.Group("/leaderboard")
	leaderboard.Get("/challenge/:id", mockHandler("leaderboard/get-challenge-id.json", fiber.StatusOK))
	leaderboard.Get("/course/:id", mockHandler("leaderboard/get-course-id.json", fiber.StatusOK))
}

func registerMetricsRoutes(app *fiber.App) {
	app.Get("/metrics", mockHandler("metrics/get.json", fiber.StatusOK))
}

func registerHealthRoutes(app *fiber.App) {
	app.Get("/health", mockHandler("health/get.json", fiber.StatusOK))
}

func registerCacheRoutes(app *fiber.App) {
	cache := app.Group("/cache")
	cache.Get("/health", mockHandler("cache/get-health.json", fiber.StatusOK))
}

func registerDBRoutes(app *fiber.App) {
	db := app.Group("/db")
	db.Get("/health", mockHandler("db/get-health.json", fiber.StatusOK))
}
