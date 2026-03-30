package submission_test

import (
	"fmt"
	"time"
	"testing"
	"net/http"

	user_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func createSubmissionCourseHTTP(t *testing.T, teacherAccess *user_dtos.UserAccess, suffix string) string {
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

	status, body, err := httputils.PostCourses(httputils.AuthHeaders(teacherAccess), bodyReq)
	if err != nil {
		t.Fatalf("create course request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create course status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := httputils.DecodeMap(t, body, "create course")
	return httputils.MapString(t, created, "id", "create course")
}

func createSubmissionExamHTTP(t *testing.T, teacherAccess *user_dtos.UserAccess, courseID, title string) string {
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

	status, body, err := httputils.PostExams(httputils.AuthHeaders(teacherAccess), bodyReq)
	if err != nil {
		t.Fatalf("create exam request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create exam status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := httputils.DecodeMap(t, body, "create exam")
	return httputils.MapString(t, created, "id", "create exam")
}

func createSubmissionChallengeHTTP(t *testing.T, teacherAccess *user_dtos.UserAccess, examID, title string) string {
	t.Helper()
	bodyReq := map[string]any{
		"title":             title,
		"description":       "Challenge for HTTP submission tests",
		"tags":              []string{"submission", "http"},
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
		"exam_id":         examID,
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

func createSubmissionTestCaseHTTP(t *testing.T, teacherAccess *user_dtos.UserAccess, challengeID, name, expectedOutput string) string {
	t.Helper()
	bodyReq := map[string]any{
		"name": name,
		"input": []map[string]any{
			{"name": "a", "type": "int", "value": "2"},
			{"name": "b", "type": "int", "value": "3"},
		},
		"expected_output": map[string]any{"name": "sum", "type": "int", "value": expectedOutput},
		"is_sample":       true,
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
	return httputils.MapString(t, created, "id", "create test-case")
}


func execPollSubmissionStates(t *testing.T, submissionID string, headers map[string]string, timeout time.Duration) []execResultState {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		status, body, err := httputils.GetSubmissionsById(headers, map[string]any{"id": submissionID})
		if err != nil {
			t.Fatalf("get submission status request failed: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("expected submission status code=%d, got=%d body=%s", http.StatusOK, status, string(body))
		}

		payload := httputils.DecodeMap(t, body, "get submission status")
		states := execExtractResultStates(t, payload)
		if len(states) > 0 {
			return states
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatalf("timeout waiting for submission=%s results", submissionID)
	return nil
}

func execUpdateResultStatus(t *testing.T, resultID, statusValue, output, workerKey string) {
	t.Helper()

	headers := map[string]string{"WorkerKey": workerKey}
	bodyReq := map[string]any{
		"status":        statusValue,
		"time_execution": 50,
		"error":         nil,
	}

	if statusValue == "executed" {
		bodyReq["output"] = output
	}

	status, body, err := httputils.PatchSubmissionsResult(headers, map[string]any{"result_id": resultID}, bodyReq)
	if err != nil {
		t.Fatalf("update submission result request failed for result=%s status=%s: %v", resultID, statusValue, err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected update result status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
}

func execExtractResultStates(t *testing.T, payload map[string]any) []execResultState {
	body := payload
	if wrapped, ok := payload["body"].(map[string]any); ok {
		body = wrapped
	}

	rawResults, ok := body["Results"]
	if !ok {
		rawResults, ok = body["results"]
		if !ok {
			return nil
		}
	}

	resultsSlice, ok := rawResults.([]any)
	if !ok {
		return nil
	}

	states := make([]execResultState, 0, len(resultsSlice))
	for _, rawResult := range resultsSlice {
		resultMap, ok := rawResult.(map[string]any)
		if !ok {
			continue
		}

		if resultMap["error_message"] == nil {
			resultMap["error_message"] = ""
		}

		state := execResultState{
			ID:           httputils.MapString(t, resultMap, "id", "id"),
			Status:       httputils.MapString(t, resultMap, "status", "status"),
			ErrorMessage: httputils.MapString(t, resultMap, "error_message", "error_message"),
		}

		if rawActualOutput, ok := resultMap["actual_output"]; ok {
			if actualOutput, ok := rawActualOutput.(map[string]any); ok {
				state.Output = httputils.MapString(t, actualOutput, "value", "value")
			}
		}

		states = append(states, state)
	}

	return states
}

func execMapString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := m[key]
		if !ok || value == nil {
			continue
		}

		s, ok := value.(string)
		if ok {
			return s
		}
	}

	return ""
}


type execResultState struct {
	ID           string
	Status       string
	Output       string
	ErrorMessage string
}
