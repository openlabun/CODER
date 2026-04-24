package submission_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSubmissionTryLimitHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Submission Try Limit HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var examID string
	var challengeID string
	var testCaseOneID string
	var testCaseTwoID string
	var examItemID string
	var sessionID string

	defer func() {
		if sessionID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando sesion %s", sessionID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, sessionID)
		}
		if examItemID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando punto de examen %s", examItemID)
			_ = httputils.DeleteExamItemByID(t, app, teacherAccess, examItemID)
		}
		if testCaseTwoID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando caso de prueba %s", testCaseTwoID)
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, testCaseTwoID)
		}
		if testCaseOneID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando caso de prueba %s", testCaseOneID)
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, testCaseOneID)
		}
		if challengeID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto %s", challengeID)
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
		}
		if examID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examID)
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	process.EndStep()

	// [STEP 2] Create a public exam with try_limit = 2
	process.StartStep("Crear examen publico")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Submission Try Limit Exam HTTP",
		"description":            "Examen para validar limite de intentos por revision",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             300,
		"try_limit":              3,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create a challenge for the exam item
	process.StartStep("Crear un reto")
	resp = httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Submission Try Limit Challenge HTTP",
		"description":         "Challenge para validar limite de revisiones",
		"tags":                []string{"submission", "try-limit"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve():\n    pass",
		},
		"input_variables": []map[string]any{{"name": "n", "type": "int", "value": "2"}},
		"output_variable": map[string]any{"name": "out", "type": "int", "value": "4"},
		"constraints":     "1 <= n <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 4] Create test cases for the challenge
	process.StartStep("Crear casos de prueba")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "submission_try_limit_case_1_http",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "2"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "4"},
		"is_sample":       false,
		"points":          5,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create test case 1")
	testCaseOneID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	// [STEP 5] Create second test case
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "submission_try_limit_case_2_http",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "5"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "10"},
		"is_sample":       false,
		"points":          5,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create test case 2")
	testCaseTwoID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 6] Create an exam item with the challenge
	process.StartStep("Crear un punto de examen (con TryLimit == 2)")
	resp = httputils.PostExamItemCreate(t, app, teacherAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        1,
		"points":       100,
		"try_limit":    2,
	})
	httputils.RequireStatus(t, resp, 201, "create exam item")
	examItemID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 7] Login as student
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	process.EndStep()

	// [STEP 8] Create a session for the student
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create session")
	sessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 9] Create first submission for the session
	process.StartStep("Crear una revision")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		"language":     "python",
		"challenge_id": challengeID,
		"session_id":   sessionID,
	})
	httputils.RequireStatus(t, resp, 201, "create first submission")
	process.EndStep()

	// [STEP 10] Create second submission for the session
	process.StartStep("Crear una revision")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		"language":     "python",
		"challenge_id": challengeID,
		"session_id":   sessionID,
	})
	httputils.RequireStatus(t, resp, 201, "create second submission")
	process.EndStep()

	// [STEP 11] Create third submission for the session - should fail with try limit exceeded error
	process.StartStep("Crear una revision (espera error)")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		"language":     "python",
		"challenge_id": challengeID,
		"session_id":   sessionID,
	})
	if resp.StatusCode == 201 {
		process.Fail("create third submission", fmt.Errorf("expected error when submission try limit is exceeded"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
