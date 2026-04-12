package submission_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSubmissionExecutionHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Submission Execution HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var examID string
	var challengeID string
	var testCaseID string
	var examItemID string
	var sessionID string
	var submissionID string

	defer func() {
		if sessionID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando sesión %s", sessionID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, sessionID)
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

	// [STEP 2] Create an exam
	process.StartStep("Crear examen publico (visibilidad public y sin curso)")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Submission Execution Exam",
		"description":            "Examen para validacion de ejecucion",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             120,
		"try_limit":              3,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create a challenge
	process.StartStep("Crear un reto")
	resp = httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Submission Execution Challenge",
		"description":         "Challenge para validar ejecucion asincrona",
		"tags":                []string{"submission", "execution"},
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
		"name":            "execution_case",
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

	// [STEP 6] Login as student
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	process.EndStep()

	// [STEP 7] Create a session with the exam
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create session")
	sessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 8] Create a submission for the challenge in the exam item
	process.StartStep("Crear una revision")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(int(sys.stdin.read().strip()))",
		"language":     "python",
		"score":        0,
		"challenge_id": challengeID,
		"session_id":   sessionID,
	})
	httputils.RequireStatus(t, resp, 201, "create submission")
	submissionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 9] Get the submission status until it is accepted
	process.StartStep("Obtener el status de la revision hasta que su estado sea accepted")
	deadline := time.Now().Add(2 * time.Minute)
	accepted := false
	for time.Now().Before(deadline) {
		resp = httputils.GetSubmissionByID(t, app, studentAccess, submissionID)
		if resp.StatusCode != 200 {
			process.Log(fmt.Sprintf("Polling status code: %d", resp.StatusCode))
			time.Sleep(2 * time.Second)
			continue
		}
		status := httputils.MustJSONMap(t, resp)
		resultsRaw, ok := status["results"].([]any)
		if !ok || len(resultsRaw) == 0 {
			time.Sleep(2 * time.Second)
			continue
		}

		allAccepted := true
		for _, item := range resultsRaw {
			result, ok := item.(map[string]any)
			if !ok || httputils.StringField(result, "status") != "accepted" {
				allAccepted = false
				break
			}
		}
		if allAccepted {
			accepted = true
			break
		}

		time.Sleep(2 * time.Second)
	}
	if !accepted {
		process.Fail("wait accepted submission status", fmt.Errorf("submission %s did not reach accepted status before timeout", submissionID))
	}
	process.EndStep()

	process.End()
}
