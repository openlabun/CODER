package exam_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestTestCaseFromStudentViewHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "TestCase Student View HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var challengeID string
	var sampleTestCaseID string
	var privateTestCaseID string
	var examID string
	var examItemID string

	defer func() {
		if examItemID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando punto de examen %s", examItemID)
			_ = httputils.DeleteExamItemByID(t, app, teacherAccess, examItemID)
		}
		if examID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examID)
		}
		if privateTestCaseID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando test case privado %s", privateTestCaseID)
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, privateTestCaseID)
		}
		if sampleTestCaseID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando test case sample %s", sampleTestCaseID)
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, sampleTestCaseID)
		}
		if challengeID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto %s", challengeID)
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	process.EndStep()

	// [STEP 2] Create a challenge to associate with test cases
	process.StartStep("Crear un reto")
	resp := httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Challenge for TestCase Student View",
		"description":         "Challenge auxiliar para vista estudiante",
		"tags":                []string{"testcase", "student-view"},
		"status":              "draft",
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

	// [STEP 3] Create a sample test case (isSample == true)
	process.StartStep("Crear un caso de prueba (isSample == true)")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "sample_case",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "10"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "10"},
		"is_sample":       true,
		"points":          0,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create sample test case")
	sampleTestCaseID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 4] Create a private test case (isSample == false)
	process.StartStep("Crear un caso de prueba (isSample == false)")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "hidden_case",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "11"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "11"},
		"is_sample":       false,
		"points":          10,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create hidden test case")
	privateTestCaseID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Publish the challenge and create an exam with the challenge to be able to access test cases from student view
	process.StartStep("Obtener casos de prueba con vista de Docente")
	resp = httputils.PostChallengePublish(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "publish challenge")

	now := time.Now().UTC()
	resp = httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "TestCase Student View Exam",
		"description":            "Exam publico para acceder a test cases",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3600,
		"try_limit":              1,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	resp = httputils.PostExamItemCreate(t, app, teacherAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        1,
		"points":       100,
	})
	httputils.RequireStatus(t, resp, 201, "create exam item")
	examItemID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	resp = httputils.GetTestCasesByChallenge(t, app, teacherAccess, challengeID, "")
	httputils.RequireStatus(t, resp, 200, "teacher get test cases")
	var teacherCases []map[string]any
	if err := json.Unmarshal(resp.Body, &teacherCases); err != nil {
		process.Fail("teacher get test cases", fmt.Errorf("decode teacher test cases: %w", err))
	}
	if len(teacherCases) != 2 {
		process.Fail("teacher get test cases", fmt.Errorf("expected 2 test cases, got %d", len(teacherCases)))
	}
	process.EndStep()

	// [STEP 6] Login as student
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	process.EndStep()

	// [STEP 7] Get test cases by challenge as student and validate only the sample test case is returned
	process.StartStep("Obtener casos de prueba con vista de Estudiante (espera solo 1)")
	resp = httputils.GetTestCasesByChallenge(t, app, studentAccess, challengeID, examID)
	httputils.RequireStatus(t, resp, 200, "student get test cases")
	var studentCases []map[string]any
	if err := json.Unmarshal(resp.Body, &studentCases); err != nil {
		process.Fail("student get test cases", fmt.Errorf("decode student test cases: %w", err))
	}
	if len(studentCases) != 1 {
		process.Fail("student get test cases", fmt.Errorf("expected 1 public sample test case, got %d", len(studentCases)))
	}
	if httputils.StringField(studentCases[0], "id") != sampleTestCaseID {
		process.Fail("student get test cases", fmt.Errorf("expected only sample test case %s", sampleTestCaseID))
	}
	process.EndStep()

	process.End()
}
