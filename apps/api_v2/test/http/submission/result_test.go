package submission_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestUpdateSubmissionResultWithWorker(t *testing.T) {
	t.Log("[STEP 1] Initialize the app and authenticate a user")
	app := initSubmissionHTTPApp(t)
	workerKey := "worker-key-http-result-test"
	t.Setenv("WORKER_KEY", workerKey)
	t.Setenv("INTERNAL_USER_EMAIL", "test@test.com")
	t.Setenv("INTERNAL_USER_PASSWORD", "Testing123!")

	runSuffix := fmt.Sprintf("%d", time.Now().UnixNano())
	studentEmail := "student.result." + runSuffix + "@test.com"

	t.Log("[STEP 2] Login as teacher")
	teacherAccess := ensureSubmissionHTTPAuthUserAccess(t, app, "test@test.com", "Testing123!", "Teacher Test")
	teacherHeaders := submissionAuthHeaders(teacherAccess)

	t.Log("[STEP 3] Create course, exam, challenge and a test case")
	courseID := createSubmissionCourseHTTP(t, app, teacherAccess, "result")
	examID := createSubmissionExamHTTP(t, app, teacherAccess, courseID, "HTTP Submission Result Exam")
	challengeID := createSubmissionChallengeHTTP(t, app, teacherAccess, examID, "HTTP Submission Result Challenge")
	_ = createSubmissionTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_result_1", "5")

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/challenges/"+challengeID+"/publish", nil, teacherHeaders)
	if err != nil {
		t.Fatalf("publish challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected publish challenge status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 4] Login as student and create a submission for the challenge")
	studentAccess := ensureSubmissionHTTPAuthUserAccess(t, app, studentEmail, "Testing123!", "Student Result")
	studentHeaders := submissionAuthHeaders(studentAccess)

	enrollBody := map[string]string{"course_id": courseID, "student_id": studentAccess.UserData.ID}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/courses/enroll", enrollBody, studentHeaders)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	createSessionBody := map[string]string{"userID": studentAccess.UserData.ID, "examID": examID}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/sessions/", createSessionBody, studentHeaders)
	if err != nil {
		t.Fatalf("create session request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create session status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	sessionID := submissionMapString(t, decodeSubmissionMap(t, body, "create session"), "ID", "create session")

	createSubmissionBody := map[string]any{
		"code":        "def solve(a, b):\n    return a + b",
		"function":	    "solve(a, b)",
		"language":    "python",
		"challengeID": challengeID,
		"sessionID":   sessionID,
	}

	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/", createSubmissionBody, studentHeaders)
	if err != nil {
		t.Fatalf("create submission request failed: %v", err)
	}
	if status != http.StatusCreated {
		// Current usecase permissions may require professor headers.
		status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/submissions/", createSubmissionBody, teacherHeaders)
		if err != nil {
			t.Fatalf("create submission fallback request failed: %v", err)
		}
		if status != http.StatusCreated {
			t.Fatalf("expected create submission status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
		}
	}

	createdSubmission := decodeSubmissionMap(t, body, "create submission")
	submissionID := submissionMapString(t, createdSubmission, "ID", "create submission")

	states := execPollSubmissionStates(t, app, submissionID, teacherHeaders, 10*time.Second)
	if len(states) == 0 {
		t.Fatalf("expected at least one submission result for submission=%s", submissionID)
	}

	t.Log("[STEP 5] Request submission status until its \"executed\" (max. 3 tries with 2s wait)")
	executed := false
	for attempt := 1; attempt <= 3; attempt++ {
		status, body, err = httputils.DoJSONRequest(app, http.MethodGet, "/submissions/"+submissionID, nil, teacherHeaders)
		if err != nil {
			t.Fatalf("get submission status request failed on attempt=%d: %v", attempt, err)
		}
		if status != http.StatusOK {
			t.Fatalf("expected submission status code=%d, got=%d body=%s", http.StatusOK, status, string(body))
		}

		payload := decodeSubmissionMap(t, body, "get submission status")
		currentStates := execExtractResultStates(payload)
		if len(currentStates) == 0 {
			if attempt < 3 {
				time.Sleep(2 * time.Second)
				continue
			}
			t.Fatalf("expected submission results in status payload, got body=%s", string(body))
		}

		allExecuted := true
		for _, state := range currentStates {
			if state.Status != "executed" {
				allExecuted = false
				break
			}
		}

		if allExecuted {
			executed = true
			break
		}

		if attempt < 3 {
			time.Sleep(2 * time.Second)
		}
	}

	if !executed {
		t.Fatalf("expected submission=%s results to reach executed within 3 attempts", submissionID)
	}

	t.Log("[STEP 6] Delete the course")
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/courses/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete course status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
}
