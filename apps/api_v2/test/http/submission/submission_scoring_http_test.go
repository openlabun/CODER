package submission_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSubmissionScoringHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Submission Scoring HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var examID string
	var challengeID string
	var testCaseOneID string
	var testCaseTwoID string
	var testCaseThreeID string
	var examItemID string
	var sessionID string
	var submissionID string

	defer func() {
		if sessionID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando sesión %s", sessionID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, sessionID)
		}
		if examItemID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando punto de examen %s", examItemID)
			_ = httputils.DeleteExamItemByID(t, app, teacherAccess, examItemID)
		}
		if testCaseThreeID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando caso de prueba %s", testCaseThreeID)
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, testCaseThreeID)
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

	// [STEP 2] Create an exam
	process.StartStep("Crear examen publico (visibilidad public y sin curso)")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Submission Scoring Exam",
		"description":            "Examen para validar puntaje de submissions",
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
		"title":               "Submission Scoring Challenge",
		"description":         "Challenge para validar score",
		"tags":                []string{"submission", "score"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"input_variables":     []map[string]any{{"name": "n", "type": "int", "value": "2"}},
		"output_variable":     map[string]any{"name": "out", "type": "int", "value": "4"},
		"constraints":         "1 <= n <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 4] Create test cases for the challenge
	process.StartStep("Crear 2 casos de prueba con valor de 3 puntos")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "score_case_1",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "2"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "4"},
		"is_sample":       false,
		"points":          3,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create score test case 1")
	testCaseOneID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "score_case_2",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "5"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "10"},
		"is_sample":       false,
		"points":          3,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create score test case 2")
	testCaseTwoID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Create an exam item with the challenge
	process.StartStep("Crear un caso de prueba con valor de 6 puntos (debe ser imposible de cumplir)")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "score_case_impossible",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "7"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "999"},
		"is_sample":       false,
		"points":          6,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create impossible test case")
	testCaseThreeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 6] Create an exam item with the challenge
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

	// [STEP 7] Login as student
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	process.EndStep()

	// [STEP 8] Create a session with the exam
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create session")
	sessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 9] Create a submission for the challenge in the exam item
	process.StartStep("Crear una revision")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		"language":     "python",
		"score":        0,
		"challenge_id": challengeID,
		"session_id":   sessionID,
	})
	httputils.RequireStatus(t, resp, 201, "create submission")
	submissionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 10] Get the submission status until it is accepted
	process.StartStep("Obtener el status de la revision hasta que su estado sea accepted o wrong_answer")
	deadline := time.Now().Add(2 * time.Minute)
	var lastStatus map[string]any
	timedOut := true
	for time.Now().Before(deadline) {
		resp = httputils.GetSubmissionByID(t, app, studentAccess, submissionID)
		if resp.StatusCode != 200 {
			process.Log(fmt.Sprintf("Polling status code: %d", resp.StatusCode))
			time.Sleep(2 * time.Second)
			continue
		}

		status := httputils.MustJSONMap(t, resp)
		resultsRaw, ok := status["results"].([]any)
		if !ok || len(resultsRaw) != 3 {
			time.Sleep(2 * time.Second)
			continue
		}

		allTerminal := true
		for _, item := range resultsRaw {
			result, ok := item.(map[string]any)
			if !ok {
				allTerminal = false
				break
			}
			s := httputils.StringField(result, "status")
			if s != "accepted" && s != "wrong_answer" {
				allTerminal = false
				break
			}
		}

		lastStatus = status
		if allTerminal {
			timedOut = false
			break
		}

		time.Sleep(2 * time.Second)
	}
	if timedOut {
		process.Fail("wait scoring terminal status", fmt.Errorf("submission %s did not reach terminal accepted/wrong_answer status before timeout", submissionID))
	}
	process.EndStep()

	// [STEP 11] Confirm the submission score corresponds to the test cases passed (6 points)
	process.StartStep("Confirmar valor del atributo Score de la revision corresponde a 6")
	if lastStatus == nil {
		process.Fail("verify submission score", fmt.Errorf("expected submission status output"))
	}
	submission, ok := lastStatus["submission"].(map[string]any)
	if !ok {
		process.Fail("verify submission score", fmt.Errorf("expected submission field in status output"))
	}
	score, ok := submission["score"].(float64)
	if !ok || int(score) != 6 {
		process.Fail("verify submission score", fmt.Errorf("expected submission score 6, got %v", submission["score"]))
	}
	process.EndStep()

	process.End()
}
