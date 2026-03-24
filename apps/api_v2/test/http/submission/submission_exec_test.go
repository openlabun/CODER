package submission_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSubmissionExecHTTP(t *testing.T) {
	t.Log("[STEP 1] Initialize app and dependencies")
	app := initSubmissionHTTPApp(t)
	workerKey := "worker-key-http-test"
	t.Setenv("WORKER_KEY", workerKey)

	studentEmail := "stud@test.com"

	t.Log("[STEP 2] Ensure teacher and student access")
	teacherAccess := ensureSubmissionHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := submissionAuthHeaders(teacherAccess)
	studentAccess := ensureSubmissionHTTPAuthUserAccess(t, app, studentEmail, "Password123!", "Student Exec")
	studentHeaders := submissionAuthHeaders(studentAccess)

	t.Log("[STEP 3] Create course and exam")
	courseID := createSubmissionCourseHTTP(t, app, teacherAccess, "exec")
	examID := createSubmissionExamHTTP(t, app, teacherAccess, courseID, "HTTP Submission Exec Exam")

	enrollBody := map[string]string{"course_id": courseID, "student_id": studentAccess.UserData.ID}
	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/courses/enroll", enrollBody, studentHeaders)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 4] Create challenge with test cases")
	challengeID := createSubmissionChallengeHTTP(t, app, teacherAccess, examID, "HTTP Submission Exec Challenge")
	_ = createSubmissionTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_exec_1", "5")
	_ = createSubmissionTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_exec_2", "5")

	t.Log("[STEP 5] Publish challenge")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/challenges/"+challengeID+"/publish", nil, teacherHeaders)
	if err != nil {
		t.Fatalf("publish challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected publish challenge status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 6] Create submission for the challenge")
	createSessionBody := map[string]string{"userID": studentAccess.UserData.ID, "examID": examID}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/", createSessionBody, studentHeaders)
	if err != nil {
		t.Fatalf("create session request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create session status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	sessionID := submissionMapString(t, decodeSubmissionMap(t, body, "create session"), "id", "create session")

	createSubmissionBody := map[string]any{
		"code":        "def solve(a, b):\n    return a + b",
		"language":    "python",
		"challengeID": challengeID,
		"sessionID":   sessionID,
	}

	// Current usecase permissions require professor headers for create submission.
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/", createSubmissionBody, teacherHeaders)
	if err != nil {
		t.Fatalf("create submission request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create submission status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	createdSubmission := decodeSubmissionMap(t, body, "create submission")
	submissionID := submissionMapString(t, createdSubmission, "id", "create submission")

	// Simulate worker updates over internal endpoint so execution flow can be validated in HTTP tests.
	states := execPollSubmissionStates(t, app, submissionID, teacherHeaders, 10*time.Second)
	if len(states) == 0 {
		t.Fatalf("expected at least one submission result, got none for submission=%s", submissionID)
	}

	for _, state := range states {
		execUpdateResultStatus(t, app, state.ID, "running", "", workerKey)
		execUpdateResultStatus(t, app, state.ID, "accepted", "5", workerKey)
	}

	t.Log("[STEP 7] Poll submission result until completion")
	finalStates := execPollSubmissionStates(t, app, submissionID, teacherHeaders, 10*time.Second)
	if len(finalStates) == 0 {
		t.Fatalf("expected submission results after execution updates for submission=%s", submissionID)
	}

	for _, state := range finalStates {
		if state.Status != "accepted" {
			t.Fatalf("expected result=%s status=accepted, got=%s", state.ID, state.Status)
		}
	}

	t.Log("[STEP 8] Validate submission result matches expected output")
	for _, state := range finalStates {
		if state.Output != "5" {
			t.Fatalf("expected result=%s output=5, got=%q", state.ID, state.Output)
		}
		if state.ErrorMessage != "" {
			t.Fatalf("expected result=%s without error message, got=%q", state.ID, state.ErrorMessage)
		}
	}

	t.Log("[STEP 9] Close student session")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/"+sessionID+"/close", nil, studentHeaders)
	if err != nil {
		t.Fatalf("close session request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected close session status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 10] Cleanup created data")
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

}

type execResultState struct {
	ID           string
	Status       string
	Output       string
	ErrorMessage string
}

func execPollSubmissionStates(t *testing.T, app *fiber.App, submissionID string, headers map[string]string, timeout time.Duration) []execResultState {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		status, body, err := httputils.DoJSONRequest(app, http.MethodGet, "/submissions/"+submissionID, nil, headers)
		if err != nil {
			t.Fatalf("get submission status request failed: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("expected submission status code=%d, got=%d body=%s", http.StatusOK, status, string(body))
		}

		payload := decodeSubmissionMap(t, body, "get submission status")
		states := execExtractResultStates(payload)
		if len(states) > 0 {
			return states
		}

		time.Sleep(200 * time.Millisecond)
	}

	t.Fatalf("timeout waiting for submission=%s results", submissionID)
	return nil
}

func execUpdateResultStatus(t *testing.T, app *fiber.App, resultID, statusValue, output, workerKey string) {
	t.Helper()

	headers := map[string]string{"WorkerKey": workerKey}
	bodyReq := map[string]any{
		"status":        statusValue,
		"timeExecution": 50,
		"error":         nil,
	}

	if statusValue == "accepted" {
		bodyReq["output"] = output
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPatch, "/submissions/results/"+resultID, bodyReq, headers)
	if err != nil {
		t.Fatalf("update submission result request failed for result=%s status=%s: %v", resultID, statusValue, err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected update result status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
}

func execExtractResultStates(payload map[string]any) []execResultState {
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

		state := execResultState{
			ID:           execMapString(resultMap, "id", "id"),
			Status:       execMapString(resultMap, "status", "status"),
			ErrorMessage: execMapString(resultMap, "errorMessage", "errorMessage"),
		}

		if rawActualOutput, ok := resultMap["ActualOutput"]; ok {
			if actualOutput, ok := rawActualOutput.(map[string]any); ok {
				state.Output = execMapString(actualOutput, "Value", "value")
			}
		}
		if rawActualOutput, ok := resultMap["actualOutput"]; ok {
			if actualOutput, ok := rawActualOutput.(map[string]any); ok {
				state.Output = execMapString(actualOutput, "Value", "value")
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
