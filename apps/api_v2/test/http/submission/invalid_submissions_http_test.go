package submission_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestInvalidSubmissionsHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Invalid Submissions HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var examID string
	var challengeID string
	var testCaseID string
	var examItemID string
	var firstSessionID string
	var secondSessionID string
	var thirdSessionID string
	var blockedSessionID string

	defer func() {
		if thirdSessionID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando sesión %s", thirdSessionID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, thirdSessionID)
		}
		if secondSessionID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando sesión %s", secondSessionID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, secondSessionID)
		}
		if firstSessionID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando sesión %s", firstSessionID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, firstSessionID)
		}
		if examItemID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando punto de examen %s", examItemID)
			_ = httputils.DeleteExamItemByID(t, app, teacherAccess, examItemID)
		}
		if testCaseID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando caso de prueba %s", testCaseID)
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, testCaseID)
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

	// [STEP 2] Create a public exam with short time limit and a challenge with a test case to be able to create sessions and submissions
	process.StartStep("Crear examen publico (visibilidad public, sin curso y 60 segundos de tiempo para resolver)")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Invalid Submissions Exam",
		"description":            "Examen para validar revisiones invalidas",
		"visibility":             "public",
		"start_time":             now.Add(-2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": false,
		"time_limit":             60,
		"try_limit":              5,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create a challenge with a test case
	process.StartStep("Crear un reto")
	resp = httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Invalid Submissions Challenge",
		"description":         "Challenge para escenarios de revision invalida",
		"tags":                []string{"submission", "invalid"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables":     []map[string]any{{"name": "n", "type": "int", "value": "10"}},
		"output_variable":     map[string]any{"name": "out", "type": "int", "value": "10"},
		"constraints":         "1 <= n <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 4] Create a test case for the challenge
	process.StartStep("Crear casos de prueba")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "invalid_submission_case",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "10"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "10"},
		"is_sample":       true,
		"points":          10,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create test case")
	testCaseID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Create an exam item with the challenge
	process.StartStep("Crear un punto de examen")
	resp = httputils.PostExamItemCreate(t, app, teacherAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        1,
		"points":       100,
	})
	httputils.RequireStatus(t, resp, 201, "create exam item")
	examItemID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 6] Publish the exam
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	process.EndStep()

	// [STEP 7] Create a revision without session (expect error)	
	process.StartStep("Crear una revision sin sesion (espera error)")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(sys.stdin.read().strip())",
		"language":     "python",
		"challenge_id": challengeID,
		"session_id":   "session-not-found",
	})
	if resp.StatusCode == 201 {
		process.Fail("create submission without session", fmt.Errorf("expected error when creating submission without valid session"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 8] Create a revision with invalid session (expect error)
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create first session")
	firstSessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 9] Create a revision with valid session but invalid code (expect error)
	process.StartStep("Cerrar el examen desde la vista de docente")
	resp = httputils.PostExamClose(t, app, teacherAccess, examID)
	httputils.RequireStatus(t, resp, 200, "close exam")
	process.EndStep()

	// [STEP 10] Create a revision after exam close (expect error)
	process.StartStep("Crear una revision (espera error)")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(sys.stdin.read().strip())",
		"language":     "python",
		"challenge_id": challengeID,
		"session_id":   firstSessionID,
	})
	if resp.StatusCode == 201 {
		process.Fail("create submission after exam close", fmt.Errorf("expected error after exam close"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 11] Re-open the exam, wait for session timeout and try to create a revision (expect error)
	process.StartStep("Esperar 61 segundos y crear una revision (espera error)")
	time.Sleep(61 * time.Second)
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(sys.stdin.read().strip())",
		"language":     "python",
		"challenge_id": challengeID,
		"session_id":   firstSessionID,
	})
	if resp.StatusCode == 201 {
		process.Fail("create submission after timeout", fmt.Errorf("expected error after timeout"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 12] Publish the challenge and create an exam with the challenge to be able to access test cases from student view
	process.StartStep("Obtener la sesion y confirmar que ya no esta activa")
	resp = httputils.GetActiveSession(t, app, studentAccess, "")
	if resp.StatusCode == 200 {
		process.Fail("verify session expiration", fmt.Errorf("expected no active session after timeout"))
	}
	firstSessionID = ""
	process.EndStep()

	// [STEP 13] Create a new session, block it from teacher view and try to create a revision with the blocked session (expect error)
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create second session")
	secondSessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 14] Create a revision with valid session but invalid code (expect error)
	process.StartStep("Bloquear la sesion desde la vista de docente")
	resp = httputils.PostSessionBlock(t, app, teacherAccess, secondSessionID)
	httputils.RequireStatus(t, resp, 200, "block second session")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "status") != "blocked" {
		process.Fail("block second session", fmt.Errorf("expected blocked status"))
	}
	blockedSessionID = secondSessionID
	secondSessionID = ""
	process.EndStep()

	// [STEP 15] Create a revision with blocked session (expect error)
	process.StartStep("Crear una revision (espera error)")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(sys.stdin.read().strip())",
		"language":     "python",
		"challenge_id": challengeID,
		"session_id":   blockedSessionID,
	})
	if resp.StatusCode == 201 {
		process.Fail("create submission with blocked session", fmt.Errorf("expected error for blocked session"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 16] Create a teachers-only exam without course relation
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create third session")
	thirdSessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 17] Create a revision with valid session but invalid code (expect error)
	process.StartStep("Cerrar el examen")
	resp = httputils.PostExamClose(t, app, teacherAccess, examID)
	httputils.RequireStatus(t, resp, 200, "close exam again")
	process.EndStep()

	// [STEP 18] Create a revision after exam close (expect error)
	process.StartStep("Crear una revision (espera error)")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(sys.stdin.read().strip())",
		"language":     "python",
		"challenge_id": challengeID,
		"session_id":   thirdSessionID,
	})
	if resp.StatusCode == 201 {
		process.Fail("create submission after second close", fmt.Errorf("expected error after closing exam again"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	thirdSessionID = ""
	process.EndStep()

	process.End()
}
