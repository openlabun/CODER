package exam_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	user_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

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

	status, body, err := httputils.PostCourses(httputils.AuthHeaders(teacherAccess), createBody)
	if err != nil {
		t.Fatalf("create course request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create course status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := httputils.DecodeMap(t, body, "create course")
	return httputils.MapString(t, created, "ID", "create course")
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

	status, body, err := httputils.PostExams(httputils.AuthHeaders(teacherAccess), bodyReq)
	if err != nil {
		t.Fatalf("create exam request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create exam status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := httputils.DecodeMap(t, body, "create exam")
	return httputils.MapString(t, created, "ID", "create exam")
}

func createChallengeHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, title string) string {
	t.Helper()

	bodyReq := map[string]any{
		"title":             title,
		"description":       "Challenge for HTTP tests",
		"tags":              []string{"http", "exam"},
		"status":            "draft",
		"difficulty":        "easy",
		"worker_time_limit":   1500,
		"worker_memory_limit": 256,
		"input_variables": []map[string]any{
			{"name": "a", "type": "int", "value": "2"},
			{"name": "b", "type": "int", "value": "3"},
		},
		"output_variable": map[string]any{"name": "sum", "type": "int", "value": "5"},
		"constraints":    "1 <= a,b <= 1000",
		"user_id":         teacherAccess.UserData.ID,
	}

	status, body, err := httputils.PostChallenges(httputils.AuthHeaders(teacherAccess), bodyReq)
	if err != nil {
		t.Fatalf("create challenge request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create challenge status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := httputils.DecodeMap(t, body, "create challenge")
	return httputils.MapString(t, created, "id", "create challenge")
}

func createExamItemHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, examID, challengeID string, order int) string {
	t.Helper()

	bodyReq := map[string]any{
		"exam_id":      examID,
		"challenge_id":  challengeID,
		"order":        order,
		"points":       100,
	}

	status, body, err := httputils.PostExamItems(httputils.AuthHeaders(teacherAccess), bodyReq)
	if err != nil {
		t.Fatalf("create exam item request failed: %v", err)
	}

	if status != http.StatusCreated {
		t.Fatalf("expected create exam item status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := httputils.DecodeMap(t, body, "create exam item")
	return httputils.MapString(t, created, "id", "create exam item")
}

func createTestCaseHTTP(t *testing.T, app *fiber.App, teacherAccess *user_dtos.UserAccess, challengeID, name string, isSample bool) string {
	t.Helper()

	bodyReq := map[string]any{
		"name": name,
		"input": []map[string]any{
			{"name": "a", "type": "int", "value": "2"},
			{"name": "b", "type": "int", "value": "3"},
		},
		"expected_output": map[string]any{"name": "sum", "type": "int", "value": "5"},
		"is_sample":       isSample,
		"points":         10,
		"challenge_id":    challengeID,
	}

	status, body, err := httputils.PostTestCases(httputils.AuthHeaders(teacherAccess), bodyReq)
	if err != nil {
		t.Fatalf("create test-case request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create test-case status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := httputils.DecodeMap(t, body, "create test-case")
	return httputils.MapString(t, created, "ID", "create test-case")
}
