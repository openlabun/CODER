package performance_tests

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSubmissionPerformance(t *testing.T) {
	process, app := httputils.StartConcurrentHTTPTestWithApp(t, "Submission Performance")
	const n = 10
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var examID string
	var challengeID string
	var testCaseID string
	var examItemID string
	accesses := make([]*httputils.HTTPAccess, n)
	sessionIDs := make([]string, n)
	submissionIDs := make([]string, n)

	defer func() {
		for _, sid := range sessionIDs {
			if sid != "" && teacherAccess != nil {
				_ = httputils.PostSessionClose(t, app, teacherAccess, sid)
			}
		}
		if examItemID != "" && teacherAccess != nil {
			_ = httputils.DeleteExamItemByID(t, app, teacherAccess, examItemID)
		}
		if testCaseID != "" && teacherAccess != nil {
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, testCaseID)
		}
		if challengeID != "" && teacherAccess != nil {
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
		}
		if examID != "" && teacherAccess != nil {
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examID)
		}
	}()

	// [STEP 1] Login as teacher user
	t.Log("[STEP 1] Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")

	// [STEP 2] Create public exam
	t.Log("[STEP 2] Crear examen publico")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Submission Performance Exam",
		"description":            "Examen para prueba de rendimiento de revisiones",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3600,
		"try_limit":              10,
		"professor_id":           teacherAccess.UserID,
	})
	if resp.StatusCode != 201 {
		t.Fatalf("[STEP 2] create exam: status=%d body=%s", resp.StatusCode, string(resp.Body))
	}
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	// [STEP 3] Create a challenge
	t.Log("[STEP 3] Crear un reto")
	resp = httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Submission Performance Challenge",
		"description":         "Reto para prueba de rendimiento de revisiones",
		"tags":                []string{"performance", "test"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve(n):\n    return n",
		},
		"input_variables": []map[string]any{{"name": "n", "type": "int", "value": "10"}},
		"output_variable": map[string]any{"name": "out", "type": "int", "value": "10"},
		"constraints":     "1 <= n <= 1000",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("[STEP 3] create challenge: status=%d body=%s", resp.StatusCode, string(resp.Body))
	}
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	// [STEP 4] Create test cases
	t.Log("[STEP 4] Crear casos de prueba")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "perf_test_case",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "5"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "5"},
		"is_sample":       true,
		"points":          10,
		"challenge_id":    challengeID,
	})
	if resp.StatusCode != 201 {
		t.Fatalf("[STEP 4] create test case: status=%d body=%s", resp.StatusCode, string(resp.Body))
	}
	testCaseID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	// [STEP 5] Create exam item
	t.Log("[STEP 5] Crear punto de examen")
	resp = httputils.PostExamItemCreate(t, app, teacherAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        1,
		"points":       100,
	})
	if resp.StatusCode != 201 {
		t.Fatalf("[STEP 5] create exam item: status=%d body=%s", resp.StatusCode, string(resp.Body))
	}
	examItemID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	// [STEP 6] Login with student users (critical)
	var loginCounter int64 = -1
	process.StartConcurrentStep(6, "Iniciar sesion con usuarios estudiantes", func() (int, []byte, error) {
		i := atomic.AddInt64(&loginCounter, 1)
		email := fmt.Sprintf("stud-%d@test.com", i)
		httputils.PostAuthRegister(t, app, email, fmt.Sprintf("Student %d", i), password)
		resp := httputils.PostAuthLogin(t, app, email, password)
		if resp.StatusCode == 200 && resp.JSON != nil {
			access := &httputils.HTTPAccess{Email: email}
			if tok, ok := resp.JSON["token"].(map[string]any); ok {
				access.AccessToken, _ = tok["access_token"].(string)
			}
			if user, ok := resp.JSON["user_data"].(map[string]any); ok {
				access.UserID, _ = user["id"].(string)
			}
			accesses[i] = access
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	// [STEP 7] Each student creates a session (critical)
	var sessionCounter int64 = -1
	process.StartConcurrentStep(7, "Cada estudiante crea una sesion", func() (int, []byte, error) {
		i := atomic.AddInt64(&sessionCounter, 1)
		if accesses[i] == nil {
			return 0, nil, fmt.Errorf("no access for student %d", i)
		}
		resp := httputils.PostSessionCreate(t, app, accesses[i], map[string]any{
			"user_id": accesses[i].UserID,
			"exam_id": examID,
		})
		if resp.StatusCode == 201 && resp.JSON != nil {
			sessionIDs[i], _ = resp.JSON["id"].(string)
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	// [STEP 8] Each student submits their answer (critical)
	var submitCounter int64 = -1
	process.StartConcurrentStep(8, "Cada estudiante sube una revision", func() (int, []byte, error) {
		i := atomic.AddInt64(&submitCounter, 1)
		if accesses[i] == nil || sessionIDs[i] == "" {
			return 0, nil, fmt.Errorf("no session for student %d", i)
		}
		resp := httputils.PostSubmissionCreate(t, app, accesses[i], map[string]any{
			"code":         "n=input() \nprint(n)",
			"language":     "python",
			"challenge_id": challengeID,
			"session_id":   sessionIDs[i],
		})
		if resp.StatusCode == 201 && resp.JSON != nil {
			submissionIDs[i], _ = resp.JSON["id"].(string)
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	// [STEP 9] Each student polls submission until "accepted" (critical)
	var pollCounter int64 = -1
	process.StartConcurrentStep(9, "Cada estudiante consulta estado hasta obtener 'accepted'", func() (int, []byte, error) {
		i := atomic.AddInt64(&pollCounter, 1)
		if accesses[i] == nil || submissionIDs[i] == "" {
			return 0, nil, fmt.Errorf("no submission for student %d", i)
		}

		var accepted bool
		attempts := 0
		for !accepted && attempts < 15 {
			resp := httputils.GetSubmissionByID(t, app, accesses[i], submissionIDs[i])
			if resp.JSON != nil {
				status := httputils.MustJSONMap(t, resp)
				resultsRaw, ok := status["results"].([]any)
				if !ok || len(resultsRaw) == 0 {
					time.Sleep(2 * time.Second)
					attempts++
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
					return 200, resp.Body, nil
				}
			}

			time.Sleep(2 * time.Second)
			attempts++
		}

		return 400, nil, fmt.Errorf("submission %s did not reach accepted status", submissionIDs[i])
	})
	process.RunStep(n)

	// [STEP 10] Each student closes their session (critical)
	var closeCounter int64 = -1
	process.StartConcurrentStep(10, "Cada estudiante cierra su sesion", func() (int, []byte, error) {
		i := atomic.AddInt64(&closeCounter, 1)
		if accesses[i] == nil || sessionIDs[i] == "" {
			return 0, nil, fmt.Errorf("no session for student %d", i)
		}
		resp := httputils.PostSessionClose(t, app, accesses[i], sessionIDs[i])
		if resp.StatusCode == 200 {
			sessionIDs[i] = ""
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	process.End()
}
