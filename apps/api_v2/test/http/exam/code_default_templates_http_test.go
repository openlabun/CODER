package exam_test

import (
	"fmt"
	"strings"
	"testing"

	sub_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestChallengeDefaultCodeTemplatesHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Challenge Default Code Templates HTTP")
	email := "test@test.com"
	password := "Password123!"

	var teacherID string
	var teacherAccess *httputils.HTTPAccess
	var challengeID string

	defer func() {
		if challengeID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto %s", challengeID)
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
		}
	}()

	// [STEP 1] Login Teacher user
	process.StartStep("Iniciar sesion con usuario de docente (creador)")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, email, password, "Teacher Test")
	teacherID = teacherAccess.UserID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))
	process.EndStep()

	// [STEP 2] Create Challenge
	process.StartStep("Crear un reto")
	resp := httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Challenge Templates HTTP Test",
		"description":         "Challenge para validar plantillas por defecto por HTTP",
		"tags":                []string{"templates", "http"},
		"status":              "draft",
		"difficulty":          "easy",
		"worker_time_limit":   1500,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "# custom template",
		},
		"input_variables": []map[string]any{
			{"name": "n", "type": "int", "value": "5"},
		},
		"output_variable": map[string]any{"name": "result", "type": "int", "value": "25"},
		"constraints":     "1 <= n <= 10^6",
	})
	httputils.RequireStatus(t, resp, 201, "create challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	if challengeID == "" {
		process.Fail("create challenge", fmt.Errorf("expected created challenge with ID"))
	}
	process.Log(fmt.Sprintf("challengeID=%s", challengeID))
	process.EndStep()

	// [STEP 3] Create test case
	process.StartStep("Crear casos de prueba")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name": "default_template_http_case",
		"input": []map[string]any{
			{"name": "n", "type": "int", "value": "5"},
		},
		"expected_output": map[string]any{"name": "result", "type": "int", "value": "25"},
		"is_sample":       true,
		"points":          10,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create test case")
	testCaseID := httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	if testCaseID == "" {
		process.Fail("create test case", fmt.Errorf("expected created test case with ID"))
	}
	process.Log(fmt.Sprintf("testCaseID=%s", testCaseID))
	process.EndStep()

	// [STEP 4] Get default templates
	process.StartStep("Obtener plantillas por defecto para el reto")
	resp = httputils.PostExamDefaultCodeTemplates(t, app, teacherAccess, map[string]any{
		"input_variables": []map[string]any{
			{"name": "n", "type": "int", "value": "5"},
			{"name": "nums", "type": "array", "value": "1 2 3"},
			{"name": "enabled", "type": "boolean", "value": "true"},
		},
		"output_variable": map[string]any{"name": "result", "type": "int", "value": "25"},
	})
	httputils.RequireStatus(t, resp, 200, "get default code templates")
	templates := httputils.MustJSONMap(t, resp)
	if len(templates) == 0 {
		process.Fail("get default code templates", fmt.Errorf("expected at least one template"))
	}
	process.EndStep()

	// [STEP 5] Validate expected variables and output print
	process.StartStep("Validar que se reciban todas las variables esperadas y el print con el output")
	for _, language := range sub_consts.SupportedProgrammingLanguages {
		raw := templates[string(language)]
		template, ok := raw.(string)
		if !ok {
			process.Fail("validate default templates", fmt.Errorf("expected string template for language %s", language))
		}
		if strings.TrimSpace(template) == "" {
			process.Fail("validate default templates", fmt.Errorf("template for language %s is empty", language))
		}
	}

	pythonRaw := templates[string(sub_consts.LanguagePython)]
	pythonTemplate, ok := pythonRaw.(string)
	if !ok {
		process.Fail("validate default templates", fmt.Errorf("expected python template"))
	}
	if !strings.Contains(pythonTemplate, "n = int(input().strip())") {
		process.Fail("validate default templates", fmt.Errorf("expected input variable assignment in python template"))
	}
	if !strings.Contains(pythonTemplate, "if _raw_nums.startswith('['):") || !strings.Contains(pythonTemplate, "nums = ast.literal_eval(_raw_nums)") {
		process.Fail("validate default templates", fmt.Errorf("expected bracket-style array parsing in python template"))
	}
	if !strings.Contains(pythonTemplate, "nums = list(map(int, _raw_nums.split()))") {
		process.Fail("validate default templates", fmt.Errorf("expected space-separated array parsing fallback in python template"))
	}
	if !strings.Contains(pythonTemplate, "enabled = input().strip().lower() in ('true', '1', 'yes')") {
		process.Fail("validate default templates", fmt.Errorf("expected boolean input parsing in python template"))
	}
	if !strings.Contains(pythonTemplate, "result = 0") {
		process.Fail("validate default templates", fmt.Errorf("expected output declaration in python template"))
	}
	if !strings.Contains(pythonTemplate, "print(result)") {
		process.Fail("validate default templates", fmt.Errorf("expected output print in python template"))
	}
	process.EndStep()

	process.End()
}