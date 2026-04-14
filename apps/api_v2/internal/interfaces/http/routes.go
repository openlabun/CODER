package http_interfaces

import (
	"os"
	"path/filepath"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	post_generate_exam "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/ai/post-generate-exam"
	post_generate_full_challenge "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/ai/post-generate-full-challenge"
	auth_get_me "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/auth/get-me"
	auth_post_login "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/auth/post-login"
	auth_post_refresh_token "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/auth/post-refresh-token"
	auth_post_register "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/auth/post-register"
	challenge_delete_by_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/delete-by-id"
	challenge_get_by_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/get-by-id"
	challenge_get_list "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/get-list"
	challenge_get_public "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/get-public"
	challenge_patch_update "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/patch-update"
	challenge_post_archive "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/post-archive"
	challenge_post_create "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/post-create"
	challenge_post_fork "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/post-fork"
	challenge_post_publish "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/post-publish"
	course_delete_course "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/delete-course"
	course_delete_student "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/delete-student"
	course_get_by_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/get-by-id"
	course_get_list "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/get-list"
	course_get_students "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/get-students"
	course_post_add_student "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/post-add-student"
	course_post_create "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/post-create"
	course_post_enroll "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/post-enroll"
	course_post_update "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/courses/post-update"
	exam_item_delete_by_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exam-items/delete-by-id"
	exam_item_patch_update "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exam-items/patch-update"
	exam_item_post_create "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exam-items/post-create"
	exam_delete_by_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/delete-by-id"
	exam_get_by_course_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/get-by-course-id"
	exam_get_by_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/get-by-id"
	exam_get_exam_items "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/get-exam-items"
	exam_get_public "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/get-public"
	exam_patch_update "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/patch-update"
	exam_post_change_visibility "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/post-change-visibility"
	exam_post_close "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/post-close"
	exam_post_create "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/exams/post-create"
	challenge_post_default_code_templates "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/challenges/post-default-code-templates"
	sub_get_by_challenge_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/get-by-challenge-id"
	sub_get_by_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/get-by-id"
	sub_get_by_user_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/get-by-user-id"
	sub_patch_update_result "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/patch-update-result"
	sub_post_create_custom "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/post-create-custom"
	sub_post_create "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/post-create"
	sub_post_create_without_score "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/post-create-without-score"
	sub_get_active_session "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/session/get-active-session"
	sub_post_block "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/session/post-block"
	sub_post_close "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/session/post-close"
	sub_post_heartbeat "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/session/post-heartbeat"
	sub_post_session "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/submissions/session/post-session"
	tc_delete_by_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/test-cases/delete-by-id"
	tc_get_by_challenge_id "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/test-cases/get-by-challenge-id"
	tc_patch_update "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/test-cases/patch-update"
	tc_post_create "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http/test-cases/post-create"
)

func RegisterRoutes(app *fiber.App, appContainer *container.Application) {
	registerDocsRoutes(app)
	registerAuthRoutes(app, appContainer)
	registerAIRoutes(app, appContainer)
	registerChallengesRoutes(app, appContainer)
	registerTestCasesRoutes(app, appContainer)
	registerCoursesRoutes(app, appContainer)
	registerExamsRoutes(app, appContainer)
	registerExamItemsRoutes(app, appContainer)
	registerSubmissionsRoutes(app, appContainer)
	registerLeaderboardRoutes(app)
	registerMetricsRoutes(app)
	registerHealthRoutes(app)
	registerCacheRoutes(app)
	registerDBRoutes(app)
}

func registerDocsRoutes(app *fiber.App) {
	app.Get("/docs/openapi.yaml", func(c *fiber.Ctx) error {
		specPath, err := resolveOpenAPISpecPath()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendFile(specPath)
	})

	app.Get("/docs", func(c *fiber.Ctx) error {
		specURL := c.BaseURL() + "/docs/openapi.yaml"
		htmlContent, err := scalar.ApiReferenceHTML(
			&scalar.Options{
				SpecURL: specURL,
				CustomOptions: scalar.CustomOptions{
					PageTitle: "Artemisa Source-Search Service",
				},
				DarkMode: true,
			},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error generando la referencia: " + err.Error())
		}

		return c.Type("html").SendString(htmlContent)
	})
}

func resolveOpenAPISpecPath() (string, error) {
	candidates := []string{
		filepath.Join("docs", "openapi.yaml"),
		filepath.Join("apps", "api_v2", "docs", "openapi.yaml"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			abs, absErr := filepath.Abs(candidate)
			if absErr != nil {
				return "", absErr
			}
			return abs, nil
		}
	}

	return "", os.ErrNotExist
}

func registerAuthRoutes(app *fiber.App, appContainer *container.Application) {
	auth := app.Group("/auth")
	auth.Post("/register", auth_post_register.Handler(appContainer))
	auth.Post("/login", auth_post_login.Handler(appContainer))
	auth.Post("/refresh-token", auth_post_refresh_token.Handler(appContainer))
	auth.Get("/me", auth_get_me.Handler(appContainer))
}

func registerAIRoutes(app *fiber.App, appContainer *container.Application) {
	ai := app.Group("/ai")
	ai.Post("/generate-full-challenge", post_generate_full_challenge.Handler(appContainer))
	ai.Post("/generate-exam", post_generate_exam.Handler(appContainer))
}

func registerChallengesRoutes(app *fiber.App, appContainer *container.Application) {
	challenges := app.Group("/challenges")
	challenges.Post("/", challenge_post_create.Handler(appContainer))
	challenges.Get("/", challenge_get_list.Handler(appContainer))
	challenges.Get("/public", challenge_get_public.Handler(appContainer))
	challenges.Get("/:id", challenge_get_by_id.Handler(appContainer))
	challenges.Post("/default-code-templates", challenge_post_default_code_templates.Handler(appContainer))
	challenges.Patch("/:id", challenge_patch_update.Handler(appContainer))
	challenges.Delete("/:id", challenge_delete_by_id.Handler(appContainer))
	challenges.Post("/:id/publish", challenge_post_publish.Handler(appContainer))
	challenges.Post("/:id/fork", challenge_post_fork.Handler(appContainer))
	challenges.Post("/:id/archive", challenge_post_archive.Handler(appContainer))
}

func registerTestCasesRoutes(app *fiber.App, appContainer *container.Application) {
	testCases := app.Group("/test-cases")
	testCases.Post("/", tc_post_create.Handler(appContainer))
	testCases.Get("/challenge/:challengeId", tc_get_by_challenge_id.Handler(appContainer))
	testCases.Patch("/:id", tc_patch_update.Handler(appContainer))
	testCases.Delete("/:id", tc_delete_by_id.Handler(appContainer))
}

func registerCoursesRoutes(app *fiber.App, appContainer *container.Application) {
	courses := app.Group("/courses")
	courses.Post("/enroll", course_post_enroll.Handler(appContainer))
	courses.Post("/", course_post_create.Handler(appContainer))
	courses.Get("/", course_get_list.Handler(appContainer))
	courses.Get("/:id", course_get_by_id.Handler(appContainer))
	courses.Post("/:id", course_post_update.Handler(appContainer))
	courses.Delete("/:id", course_delete_course.Handler(appContainer))
	courses.Post("/:id/students", course_post_add_student.Handler(appContainer))
	courses.Delete("/:id/students/:studentId", course_delete_student.Handler(appContainer))
	courses.Get("/:id/students", course_get_students.Handler(appContainer))
}

func registerExamsRoutes(app *fiber.App, appContainer *container.Application) {
	exams := app.Group("/exams")
	exams.Post("/", exam_post_create.Handler(appContainer))
	exams.Get("/public", exam_get_public.Handler(appContainer))
	exams.Get("/course/:courseId", exam_get_by_course_id.Handler(appContainer))
	exams.Get("/public", exam_get_public.Handler(appContainer))
	exams.Get("/:id", exam_get_by_id.Handler(appContainer))
	exams.Patch("/:id", exam_patch_update.Handler(appContainer))
	exams.Delete("/:id", exam_delete_by_id.Handler(appContainer))
	exams.Post("/:id/visibility", exam_post_change_visibility.Handler(appContainer))
	exams.Post("/:id/close", exam_post_close.Handler(appContainer))
	exams.Get("/:examId/items", exam_get_exam_items.Handler(appContainer))
}

func registerExamItemsRoutes(app *fiber.App, appContainer *container.Application) {
	examItems := app.Group("/exam-items")
	examItems.Post("/", exam_item_post_create.Handler(appContainer))
	examItems.Patch("/:id", exam_item_patch_update.Handler(appContainer))
	examItems.Delete("/:id", exam_item_delete_by_id.Handler(appContainer))
}

func registerSubmissionsRoutes(app *fiber.App, appContainer *container.Application) {
	submissions := app.Group("/submissions")
	submissions.Post("/", sub_post_create.Handler(appContainer))
	submissions.Post("/execute", sub_post_create_without_score.Handler(appContainer))
	submissions.Post("/execute-custom", sub_post_create_custom.Handler(appContainer))
	submissions.Patch("/results/:resultId", sub_patch_update_result.Handler(appContainer))
	submissions.Get("/user/:userId", sub_get_by_user_id.Handler(appContainer))
	submissions.Get("/challenge/:challengeId", sub_get_by_challenge_id.Handler(appContainer))
	submissions.Get("/:id", sub_get_by_id.Handler(appContainer))
	sessions := submissions.Group("/sessions")
	sessions.Post("/", sub_post_session.Handler(appContainer))
	sessions.Get("/active", sub_get_active_session.Handler(appContainer))
	sessions.Post("/:id/heartbeat", sub_post_heartbeat.Handler(appContainer))
	sessions.Post("/:id/block", sub_post_block.Handler(appContainer))
	sessions.Post("/:id/close", sub_post_close.Handler(appContainer))
}

func registerLeaderboardRoutes(app *fiber.App) {
	leaderboard := app.Group("/leaderboard")
	leaderboard.Get("/challenge/:id", mockHandler("leaderboard/get-challenge-id/mockup/output.json", fiber.StatusOK))
	leaderboard.Get("/course/:id", mockHandler("leaderboard/get-course-id/mockup/output.json", fiber.StatusOK))
}

func registerMetricsRoutes(app *fiber.App) {
	app.Get("/metrics", mockHandler("metrics/get/mockup/output.json", fiber.StatusOK))
}

func registerHealthRoutes(app *fiber.App) {
	app.Get("/health", mockHandler("health/get/mockup/output.json", fiber.StatusOK))
}

func registerCacheRoutes(app *fiber.App) {
	cache := app.Group("/cache")
	cache.Get("/health", mockHandler("cache/get-health/mockup/output.json", fiber.StatusOK))
}

func registerDBRoutes(app *fiber.App) {
	db := app.Group("/db")
	db.Get("/health", mockHandler("db/get-health/mockup/output.json", fiber.StatusOK))
}
