package submission_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	user_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestCreateSessionHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app := initSubmissionHTTPApp(t)
	runSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	studentEmail := "stud@test.com"
	observerEmail := "observer.session." + runSuffix + "@test.com"

	t.Log("[STEP 2] Login de profesor")
	teacherAccess := ensureSubmissionHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := submissionAuthHeaders(teacherAccess)
	teacherID := teacherAccess.UserData.ID

	t.Log("[STEP 3] Crear curso de prueba")
	courseID := createSubmissionCourseHTTP(t, app, teacherAccess, "sess")

	t.Log("[STEP 4] Crear examen de prueba")
	examID := createSubmissionExamHTTP(t, app, teacherAccess, courseID, "HTTP Submission Session Exam")

	t.Log("[STEP 5] Login de estudiante y matrícula")
	studentAccess := ensureSubmissionHTTPAuthUserAccess(t, app, studentEmail, "Password123!", "Student Test")
	studentHeaders := submissionAuthHeaders(studentAccess)
	studentID := studentAccess.UserData.ID
	enrollBody := map[string]string{"course_id": courseID, "student_id": studentID}
	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/courses/enroll", enrollBody, studentHeaders)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 6] Crear sesión para estudiante")
	createSessionBody := map[string]string{"userID": studentID, "examID": examID}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/", createSessionBody, studentHeaders)
	if err != nil {
		t.Fatalf("create student session request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create session status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	studentSession := decodeSubmissionMap(t, body, "create student session")
	studentSessionID := submissionMapString(t, studentSession, "id", "create student session")
	if submissionMapString(t, studentSession, "StudentID", "create student session") != studentID {
		t.Fatalf("expected student session StudentID=%s, got body=%s", studentID, string(body))
	}
	if submissionMapString(t, studentSession, "ExamID", "create student session") != examID {
		t.Fatalf("expected student session ExamID=%s, got body=%s", examID, string(body))
	}

	t.Log("[STEP 7] Validar rechazo de sesión duplicada")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/", createSessionBody, studentHeaders)
	if err != nil {
		t.Fatalf("duplicate session request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Fatalf("expected duplicate session status=%d, got=%d body=%s", http.StatusBadRequest, status, string(body))
	}

	t.Log("[STEP 8] Validar error con examen inexistente")
	observerAccess := ensureSubmissionHTTPAuthUserAccess(t, app, observerEmail, "Password123!", "Observer Test")
	observerHeaders := submissionAuthHeaders(observerAccess)
	nonExistingExamBody := map[string]string{"userID": studentID, "examID": "non-existing-exam-id"}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/", nonExistingExamBody, observerHeaders)
	if err != nil {
		t.Fatalf("create session non-existing exam request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Fatalf("expected non-existing exam status=%d, got=%d body=%s", http.StatusBadRequest, status, string(body))
	}

	t.Log("[STEP 9] Crear sesión de profesor y ejecutar heartbeat")
	teacherSessionBody := map[string]string{"userID": teacherID, "examID": examID}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/", teacherSessionBody, teacherHeaders)
	if err != nil {
		t.Fatalf("create teacher session request failed: %v", err)
	}
	var teacherSessionID string
	if status == http.StatusCreated {
		teacherSession := decodeSubmissionMap(t, body, "create teacher session")
		teacherSessionID = submissionMapString(t, teacherSession, "id", "create teacher session")
	} else if status == http.StatusBadRequest {
		status, body, err = httputils.DoJSONRequest(app, http.MethodGet, "/submissions/sessions/active", nil, teacherHeaders)
		if err != nil {
			t.Fatalf("get active teacher session request failed: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("expected active teacher session status=%d, got=%d body=%s", http.StatusOK, status, string(body))
		}
		teacherSession := decodeSubmissionMap(t, body, "get active teacher session")
		teacherSessionID = submissionMapString(t, teacherSession, "id", "get active teacher session")
	} else {
		t.Fatalf("unexpected teacher session create status=%d body=%s", status, string(body))
	}

	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/"+teacherSessionID+"/heartbeat", nil, teacherHeaders)
	if err != nil {
		t.Fatalf("heartbeat request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected heartbeat status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	heartbeatSession := decodeSubmissionMap(t, body, "heartbeat session")
	if submissionMapString(t, heartbeatSession, "id", "heartbeat session") != teacherSessionID {
		t.Fatalf("expected heartbeat session ID=%s, got body=%s", teacherSessionID, string(body))
	}

	t.Log("[STEP 10] Cerrar sesión de estudiante")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/"+studentSessionID+"/close", nil, studentHeaders)
	if err != nil {
		t.Fatalf("close student session request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected close student session status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 11] Cerrar sesión de profesor")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/"+teacherSessionID+"/close", nil, teacherHeaders)
	if err != nil {
		t.Fatalf("close teacher session request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected close teacher session status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 12] Cleanup por curso")
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/courses/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete course status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	removed := decodeSubmissionMap(t, body, "delete course")
	if !submissionMapBool(t, removed, "removed", "delete course") {
		t.Fatalf("expected removed=true, got body=%s", string(body))
	}

	_ = studentSessionID
}

func TestSubmissionsHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app := initSubmissionHTTPApp(t)
	runSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	studentEmail := "student.submission." + runSuffix + "@test.com"

	t.Log("[STEP 2] Login de profesor")
	teacherAccess := ensureSubmissionHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := submissionAuthHeaders(teacherAccess)
	teacherID := teacherAccess.UserData.ID

	t.Log("[STEP 3] Crear curso y examen")
	courseID := createSubmissionCourseHTTP(t, app, teacherAccess, "sub")
	examID := createSubmissionExamHTTP(t, app, teacherAccess, courseID, "HTTP Submission Exam")

	t.Log("[STEP 4] Login estudiante y matrícula")
	studentAccess := ensureSubmissionHTTPAuthUserAccess(t, app, studentEmail, "Password123!", "Student Test")
	studentHeaders := submissionAuthHeaders(studentAccess)
	studentID := studentAccess.UserData.ID
	enrollBody := map[string]string{"course_id": courseID, "student_id": studentID}
	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/courses/enroll", enrollBody, studentHeaders)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 5] Crear challenge y test-cases")
	challengeID := createSubmissionChallengeHTTP(t, app, teacherAccess, examID, "HTTP Submission Challenge")
	_ = createSubmissionTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_submission_1", "5")
	_ = createSubmissionTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_submission_2", "5")

	t.Log("[STEP 6] Crear sesión para el estudiante")
	createSessionBody := map[string]string{"userID": studentID, "examID": examID}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/", createSessionBody, studentHeaders)
	if err != nil {
		t.Fatalf("create session request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create session status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	sessionID := submissionMapString(t, decodeSubmissionMap(t, body, "create session"), "id", "create session")

	t.Log("[STEP 7] Crear submission")
	createSubmissionBody := map[string]any{
		"code":        "def solve(a, b):\n    return a + b",
		"language":    "python",
		"challengeID": challengeID,
		"sessionID":   sessionID,
	}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/", createSubmissionBody, teacherHeaders)
	if err != nil {
		t.Fatalf("create submission request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create submission status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	createdSubmission := decodeSubmissionMap(t, body, "create submission")
	submissionID := submissionMapString(t, createdSubmission, "id", "create submission")

	t.Log("[STEP 8] Listar submissions por challenge y validar inclusión")
	listPath := "/submissions/?challengeId=" + challengeID
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, listPath, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("list challenge submissions request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected list submissions status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	outputs := decodeSubmissionSliceMap(t, body, "list challenge submissions")
	found := false
	for _, item := range outputs {
		rawSubmission, ok := item["Submission"]
		if !ok {
			rawSubmission, ok = item["submission"]
		}
		if !ok {
			continue
		}
		subMap, ok := rawSubmission.(map[string]any)
		if !ok {
			continue
		}
		id, ok := subMap["id"].(string)
		if !ok {
			id, _ = subMap["id"].(string)
		}
		sid, ok := subMap["SessionID"].(string)
		if !ok {
			sid, _ = subMap["session_id"].(string)
		}
		if id == submissionID && sid == sessionID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected created submission=%s to be present in challenge list, got body=%s", submissionID, string(body))
	}

	t.Log("[STEP 9] Cerrar sesión de estudiante")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/"+sessionID+"/close", nil, studentHeaders)
	if err != nil {
		t.Fatalf("close student session request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected close student session status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 10] Cleanup por curso")
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/courses/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete course status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	removed := decodeSubmissionMap(t, body, "delete course")
	if !submissionMapBool(t, removed, "removed", "delete course") {
		t.Fatalf("expected removed=true, got body=%s", string(body))
	}

	_ = teacherID
}

func initSubmissionHTTPApp(t *testing.T) *fiber.App {
	t.Helper()
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}
	return app
}

func ensureSubmissionHTTPAuthUserAccess(t *testing.T, app *fiber.App, email, password, name string) *user_dtos.UserAccess {
	t.Helper()

	loginBody := map[string]string{"email": email, "password": password}
	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/auth/login", loginBody, nil)
	if err != nil {
		t.Fatalf("login request failed for %s: %v", email, err)
	}
	if status == http.StatusOK {
		access := decodeSubmissionUserAccess(t, body, "login")
		validateSubmissionUserAccess(t, access, email)
		return access
	}

	registerBody := map[string]string{"email": email, "name": name, "password": password}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/auth/register", registerBody, nil)
	if err != nil {
		t.Fatalf("register request failed for %s: %v", email, err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected register status=%d for %s, got=%d body=%s", http.StatusCreated, email, status, string(body))
	}

	access := decodeSubmissionUserAccess(t, body, "register")
	validateSubmissionUserAccess(t, access, email)
	return access
}

func decodeSubmissionUserAccess(t *testing.T, raw []byte, source string) *user_dtos.UserAccess {
	t.Helper()
	var access user_dtos.UserAccess
	if err := json.Unmarshal(raw, &access); err != nil {
		t.Fatalf("decode %s response failed: %v body=%s", source, err, string(raw))
	}
	return &access
}

func validateSubmissionUserAccess(t *testing.T, access *user_dtos.UserAccess, expectedEmail string) {
	t.Helper()
	if access == nil || access.UserData == nil || access.Token == nil {
		t.Fatalf("expected valid access payload for %s", expectedEmail)
	}
	if access.UserData.ID == "" {
		t.Fatalf("expected user ID in access payload for %s", expectedEmail)
	}
	if access.Token.AccessToken == "" {
		t.Fatalf("expected access token in access payload for %s", expectedEmail)
	}
}

func submissionAuthHeaders(access *user_dtos.UserAccess) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + access.Token.AccessToken,
		"X-User-Email":  access.UserData.Email,
	}
}

func createSubmissionCourseHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, suffix string) string {
	t.Helper()
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-HTTP-SUB-%s-%d", suffix, now.UnixNano())
	courseCode := fmt.Sprintf("HTTP-SUB-%s-%d", suffix, now.Unix()%100000)

	bodyReq := map[string]any{
		"name":            "HTTP Submission Course " + suffix,
		"description":     "Course for HTTP submission tests",
		"visibility":      "public",
		"visual_identity": "#0055ff",
		"code":            courseCode,
		"year":            now.Year(),
		"semester":        "01",
		"enrollment_code": enrollmentCode,
		"enrollment_url":  "https://example.test/enroll/" + enrollmentCode,
		"teacher_id":      teacherAccess.UserData.ID,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/courses/", bodyReq, submissionAuthHeaders(teacherAccess))
	if err != nil {
		t.Fatalf("create course request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create course status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeSubmissionMap(t, body, "create course")
	return submissionMapString(t, created, "id", "create course")
}

func createSubmissionExamHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, courseID, title string) string {
	t.Helper()
	startTime := time.Now().UTC().Add(-2 * time.Hour)
	endTime := startTime.Add(3 * time.Hour)

	bodyReq := map[string]any{
		"course_id":              courseID,
		"title":                  title,
		"description":            "Exam for HTTP submission tests",
		"visibility":             "course",
		"start_time":             startTime.Format(time.RFC3339),
		"end_time":               endTime.Format(time.RFC3339),
		"allow_late_submissions": false,
		"time_limit":             3600,
		"try_limit":              2,
		"professor_id":           teacherAccess.UserData.ID,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/exams/", bodyReq, submissionAuthHeaders(teacherAccess))
	if err != nil {
		t.Fatalf("create exam request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create exam status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeSubmissionMap(t, body, "create exam")
	return submissionMapString(t, created, "id", "create exam")
}

func createSubmissionChallengeHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, examID, title string) string {
	t.Helper()
	bodyReq := map[string]any{
		"title":             title,
		"description":       "Challenge for HTTP submission tests",
		"tags":              []string{"submission", "http"},
		"status":            "draft",
		"difficulty":        "easy",
		"workerTimeLimit":   1500,
		"workerMemoryLimit": 256,
		"inputVariables": []map[string]any{
			{"name": "a", "type": "int", "value": "2"},
			{"name": "b", "type": "int", "value": "3"},
		},
		"outputVariable": map[string]any{"name": "sum", "type": "int", "value": "5"},
		"constraints":    "1 <= a,b <= 1000",
		"examID":         examID,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/challenges/", bodyReq, submissionAuthHeaders(teacherAccess))
	if err != nil {
		t.Fatalf("create challenge request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create challenge status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeSubmissionMap(t, body, "create challenge")
	return submissionMapString(t, created, "id", "create challenge")
}

func createSubmissionTestCaseHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, challengeID, name, expectedOutput string) string {
	t.Helper()
	bodyReq := map[string]any{
		"name": name,
		"input": []map[string]any{
			{"name": "a", "type": "int", "value": "2"},
			{"name": "b", "type": "int", "value": "3"},
		},
		"expectedOutput": map[string]any{"name": "sum", "type": "int", "value": expectedOutput},
		"isSample":       true,
		"points":         10,
		"challengeID":    challengeID,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/test-cases/", bodyReq, submissionAuthHeaders(teacherAccess))
	if err != nil {
		t.Fatalf("create test-case request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create test-case status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeSubmissionMap(t, body, "create test-case")
	return submissionMapString(t, created, "id", "create test-case")
}

func decodeSubmissionMap(t *testing.T, raw []byte, source string) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode %s map failed: %v body=%s", source, err, string(raw))
	}
	return out
}

func decodeSubmissionSliceMap(t *testing.T, raw []byte, source string) []map[string]any {
	t.Helper()
	var out []map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode %s list failed: %v body=%s", source, err, string(raw))
	}
	return out
}

func submissionMapString(t *testing.T, m map[string]any, key, source string) string {
	t.Helper()
	v, ok := m[key]
	if !ok {
		t.Fatalf("missing key=%s in %s response", key, source)
	}
	s, ok := v.(string)
	if !ok {
		t.Fatalf("key=%s in %s is not string (type=%T)", key, source, v)
	}
	return s
}

func submissionMapBool(t *testing.T, m map[string]any, key, source string) bool {
	t.Helper()
	v, ok := m[key]
	if !ok {
		t.Fatalf("missing key=%s in %s response", key, source)
	}
	b, ok := v.(bool)
	if !ok {
		t.Fatalf("key=%s in %s is not bool (type=%T)", key, source, v)
	}
	return b
}
