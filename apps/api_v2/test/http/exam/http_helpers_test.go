package exam_test

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

func initExamHTTPApp(t *testing.T) *fiber.App {
	t.Helper()
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}
	return app
}

func ensureExamHTTPAuthUserAccess(t *testing.T, app *fiber.App, email, password, name string) *user_dtos.UserAccess {
	t.Helper()

	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/auth/login", loginBody, nil)
	if err != nil {
		t.Fatalf("login request failed for %s: %v", email, err)
	}

	if status == http.StatusOK {
		access := decodeUserAccess(t, body, "login")
		validateUserAccess(t, access, email)
		return access
	}

	t.Logf("login failed for %s with status=%d, trying register", email, status)
	registerBody := map[string]string{
		"email":    email,
		"name":     name,
		"password": password,
	}

	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/auth/register", registerBody, nil)
	if err != nil {
		t.Fatalf("register request failed for %s: %v", email, err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected register status=%d for %s, got=%d body=%s", http.StatusCreated, email, status, string(body))
	}

	access := decodeUserAccess(t, body, "register")
	validateUserAccess(t, access, email)
	return access
}

func decodeUserAccess(t *testing.T, raw []byte, source string) *user_dtos.UserAccess {
	t.Helper()

	var access user_dtos.UserAccess
	if err := json.Unmarshal(raw, &access); err != nil {
		t.Fatalf("decode %s response failed: %v body=%s", source, err, string(raw))
	}
	return &access
}

func validateUserAccess(t *testing.T, access *user_dtos.UserAccess, expectedEmail string) {
	t.Helper()

	if access == nil || access.UserData == nil || access.Token == nil {
		t.Fatalf("expected valid access payload for %s", expectedEmail)
	}
	if access.UserData.ID == "" {
		t.Fatalf("expected user ID in access payload for %s", expectedEmail)
	}
	if access.Token.AccessToken == "" {
		t.Fatalf("expected access token in payload for %s", expectedEmail)
	}
	if access.UserData.Email != expectedEmail {
		t.Fatalf("expected email=%s in payload, got=%s", expectedEmail, access.UserData.Email)
	}
}

func authHeaders(access *user_dtos.UserAccess) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + access.Token.AccessToken,
		"X-User-Email":  access.UserData.Email,
	}
}

func decodeMap(t *testing.T, raw []byte, source string) map[string]any {
	t.Helper()

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode %s map response failed: %v body=%s", source, err, string(raw))
	}
	return out
}

func decodeSliceMap(t *testing.T, raw []byte, source string) []map[string]any {
	t.Helper()

	var out []map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode %s list response failed: %v body=%s", source, err, string(raw))
	}
	return out
}

func mapString(t *testing.T, m map[string]any, key, source string) string {
	t.Helper()

	v, ok := m[key]
	if !ok {
		t.Fatalf("missing key=%s in %s response", key, source)
	}
	s, ok := v.(string)
	if !ok {
		t.Fatalf("key=%s in %s response is not string (type=%T)", key, source, v)
	}
	return s
}

func mapBool(t *testing.T, m map[string]any, key, source string) bool {
	t.Helper()

	v, ok := m[key]
	if !ok {
		t.Fatalf("missing key=%s in %s response", key, source)
	}
	b, ok := v.(bool)
	if !ok {
		t.Fatalf("key=%s in %s response is not bool (type=%T)", key, source, v)
	}
	return b
}

func containsID(list []map[string]any, id string) bool {
	for _, item := range list {
		if v, ok := item["ID"].(string); ok && v == id {
			return true
		}
		if v, ok := item["id"].(string); ok && v == id {
			return true
		}
	}
	return false
}

func createCourseHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, suffix string) string {
	t.Helper()

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-HTTP-EX-%s-%d", suffix, now.UnixNano())
	courseCode := fmt.Sprintf("HTTP-EX-%s-%d", suffix, now.Unix()%100000)

	createBody := map[string]any{
		"name":            "HTTP Exam Course " + suffix,
		"description":     "Course for exam HTTP tests",
		"visibility":      "public",
		"visual_identity": "#ff0055",
		"code":            courseCode,
		"year":            now.Year(),
		"semester":        "01",
		"enrollment_code": enrollmentCode,
		"enrollment_url":  "https://example.test/enroll/" + enrollmentCode,
		"teacher_id":      teacherAccess.UserData.ID,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/courses/", createBody, authHeaders(teacherAccess))
	if err != nil {
		t.Fatalf("create course request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create course status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeMap(t, body, "create course")
	return mapString(t, created, "ID", "create course")
}

func createExamHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, courseID, title string) string {
	t.Helper()

	// Keep start_time in the past so close operation (which sets end_time=now) stays valid.
	startTime := time.Now().UTC().Add(-2 * time.Hour)
	endTime := startTime.Add(3 * time.Hour)

	bodyReq := map[string]any{
		"course_id":              courseID,
		"title":                  title,
		"description":            "Exam for HTTP tests",
		"visibility":             "course",
		"start_time":             startTime.Format(time.RFC3339),
		"end_time":               endTime.Format(time.RFC3339),
		"allow_late_submissions": false,
		"time_limit":             5400,
		"try_limit":              2,
		"professor_id":           teacherAccess.UserData.ID,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/exams/", bodyReq, authHeaders(teacherAccess))
	if err != nil {
		t.Fatalf("create exam request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create exam status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeMap(t, body, "create exam")
	return mapString(t, created, "ID", "create exam")
}

func createChallengeHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, examID, title string) string {
	t.Helper()

	bodyReq := map[string]any{
		"title":             title,
		"description":       "Challenge for HTTP tests",
		"tags":              []string{"http", "exam"},
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

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/challenges/", bodyReq, authHeaders(teacherAccess))
	if err != nil {
		t.Fatalf("create challenge request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create challenge status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeMap(t, body, "create challenge")
	return mapString(t, created, "ID", "create challenge")
}

func createTestCaseHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, challengeID, name string, isSample bool) string {
	t.Helper()

	bodyReq := map[string]any{
		"name": name,
		"input": []map[string]any{
			{"name": "a", "type": "int", "value": "2"},
			{"name": "b", "type": "int", "value": "3"},
		},
		"expectedOutput": map[string]any{"name": "sum", "type": "int", "value": "5"},
		"isSample":       isSample,
		"points":         10,
		"challengeID":    challengeID,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/test-cases/", bodyReq, authHeaders(teacherAccess))
	if err != nil {
		t.Fatalf("create test-case request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create test-case status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeMap(t, body, "create test-case")
	return mapString(t, created, "ID", "create test-case")
}
