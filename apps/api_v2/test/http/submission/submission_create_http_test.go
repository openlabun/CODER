package submission_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSubmissionCreateAndReadHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Submission Create and Read HTTP")
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
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
		}
		if examID != "" && teacherAccess != nil {
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
		"title":                  "Submission Create Exam",
		"description":            "Examen para creacion y consulta de revisiones",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3600,
		"try_limit":              2,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create a challenge
	process.StartStep("Crear un reto")
	resp = httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Submission Create Challenge",
		"description":         "Challenge para pruebas de submissions",
		"tags":                []string{"submission", "create"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
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
		"name":            "sample_submission_case",
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
	process.StartStep("Crear sesion con examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create student session")
	sessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 8] Create a submission for the challenge in the exam item
	process.StartStep("Crear una revision")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "def solve(n):\n    return n",
		"language":     "python",
		"score":        0,
		"challenge_id": challengeID,
		"session_id":   sessionID,
	})
	httputils.RequireStatus(t, resp, 201, "create submission")
	submissionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 9] Get the submission status
	process.StartStep("Obtener revisiones a partir del ID del reto")
	resp = httputils.GetSubmissionsByChallenge(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "get challenge submissions")
	var challengeSubmissions []map[string]any
	if err := json.Unmarshal(resp.Body, &challengeSubmissions); err != nil {
		process.Fail("get challenge submissions", err)
	}
	if len(challengeSubmissions) == 0 {
		process.Fail("get challenge submissions", fmt.Errorf("expected at least one submission for challenge"))
	}
	process.EndStep()

	// [STEP 10] Get the submission status by user ID
	process.StartStep("Obtener revisiones a partir del ID del usuario")
	resp = httputils.GetSubmissionsByUser(t, app, teacherAccess, teacherAccess.UserID)
	httputils.RequireStatus(t, resp, 200, "get user submissions")
	var userSubmissions []map[string]any
	if err := json.Unmarshal(resp.Body, &userSubmissions); err != nil {
		process.Fail("get user submissions", err)
	}
	if userSubmissions == nil {
		process.Fail("get user submissions", fmt.Errorf("expected non-nil user submissions slice"))
	}
	process.EndStep()

	 // [STEP 11] Get the submission status by submission ID
	process.StartStep("Obtener el status de la revision")
	resp = httputils.GetSubmissionByID(t, app, studentAccess, submissionID)
	httputils.RequireStatus(t, resp, 200, "get submission status")
	status := httputils.MustJSONMap(t, resp)
	submission, ok := status["submission"].(map[string]any)
	if !ok || httputils.StringField(submission, "id") != submissionID {
		process.Fail("get submission status", fmt.Errorf("expected status for submission %s", submissionID))
	}
	process.EndStep()

	process.End()
}
