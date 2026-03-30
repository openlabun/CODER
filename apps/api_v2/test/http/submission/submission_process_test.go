package submission_test

import (
	"testing"
	"time"
	"net/http"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSubmissionMockProcess(t *testing.T) {
	t.Log("[STEP 1] Initialize app and dependencies")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}

	workerKey := "worker-key-http-test"
	t.Setenv("WORKER_KEY", workerKey)

	t.Log("[STEP 2] Ensure teacher and student access")
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := httputils.AuthHeaders(teacherAccess)
	studentAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Exec")
	studentHeaders := httputils.AuthHeaders(studentAccess)

	t.Log("[STEP 3] Create course and exam")
	courseID := createSubmissionCourseHTTP(t, teacherAccess, "exec")
	examID := createSubmissionExamHTTP(t, teacherAccess, courseID, "HTTP Submission Exec Exam")

	enrollBody := map[string]string{"course_id": courseID, "student_id": studentAccess.UserData.ID}
	status, body, err := httputils.PostCoursesEnroll(studentHeaders, enrollBody)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 4] Create challenge with test cases")
	challengeID := createSubmissionChallengeHTTP(t, teacherAccess, examID, "HTTP Submission Exec Challenge")
	_ = createSubmissionTestCaseHTTP(t, teacherAccess, challengeID, "tc_exec_1", "5")
	_ = createSubmissionTestCaseHTTP(t, teacherAccess, challengeID, "tc_exec_2", "5")

	t.Log("[STEP 5] Publish challenge")
	status, body, err = httputils.PostChallengesPublish(teacherHeaders, map[string]any{"id": challengeID})
	if err != nil {
		t.Fatalf("publish challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected publish challenge status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 6] Create submission for the challenge")
	createSessionBody := map[string]string{"user_id": studentAccess.UserData.ID, "exam_id": examID}
	status, body, err = httputils.PostSubmissionsSessions(studentHeaders, createSessionBody)
	if err != nil {
		t.Fatalf("create session request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create session status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	sessionID := httputils.MapString(t, httputils.DecodeMap(t, body, "create session"), "id", "create session")

	createSubmissionBody := map[string]any{
		"code":        "def solve(a, b):\n    return a + b",
		"language":    "python",
		"challenge_id": challengeID,
		"session_id":   sessionID,
	}

	// Current usecase permissions require professor headers for create submission.
	status, body, err = httputils.PostSubmissions(studentHeaders, createSubmissionBody)
	if err != nil {
		t.Fatalf("create submission request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create submission status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	createdSubmission := httputils.DecodeMap(t, body, "create submission")
	submissionID := httputils.MapString(t, createdSubmission, "id", "create submission")

	t.Log("[STEP 7] Simulate worker processing updates and poll submission results")
	// Simulate worker updates over internal endpoint so execution flow can be validated in HTTP tests.
	states := execPollSubmissionStates(t, submissionID, teacherHeaders, 10*time.Second)
	if len(states) == 0 {
		t.Fatalf("expected at least one submission result, got none for submission=%s", submissionID)
	}

	for _, state := range states {
		execUpdateResultStatus(t, state.ID, "running", "", workerKey)
		execUpdateResultStatus(t, state.ID, "executed", "5", workerKey)
		execUpdateResultStatus(t, state.ID, "accepted", "5", workerKey)
	}

	t.Log("[STEP 8] Poll submission result until completion")
	finalStates := execPollSubmissionStates(t, submissionID, studentHeaders, 10*time.Second)
	if len(finalStates) == 0 {
		t.Fatalf("expected submission results after execution updates for submission=%s", submissionID)
	}

	for _, state := range finalStates {
		if state.Status != "accepted" {
			t.Fatalf("expected result=%s status=accepted, got=%s", state.ID, state.Status)
		}
	}

	t.Log("[STEP 9] Validate submission result matches expected output")
	for _, state := range finalStates {
		if state.Output != "5" {
			t.Fatalf("expected result=%s output=5, got=%q", state.ID, state.Output)
		}
		if state.ErrorMessage != "" {
			t.Fatalf("expected result=%s without error message, got=%q", state.ID, state.ErrorMessage)
		}
	}

	t.Log("[STEP 10] Cleanup created data")
	status, body, err = httputils.DeleteCoursesById(teacherHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete course status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	removed := httputils.DecodeMap(t, body, "delete course")
	if !httputils.MapBool(t, removed, "removed", "delete course") {
		t.Fatalf("expected removed=true, got body=%s", string(body))
	}
}

func TestSubmissionProcessWithWorkers (t *testing.T) {
	t.Log("[STEP 1] Initialize the app and authenticate a user")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}

	workerKey := "worker-key-http-result-test"
	t.Setenv("WORKER_KEY", workerKey)
	t.Setenv("INTERNAL_USER_EMAIL", "test@test.com")
	t.Setenv("INTERNAL_USER_PASSWORD", "Password123!")

	t.Log("[STEP 2] Login as teacher")
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := httputils.AuthHeaders(teacherAccess)

	t.Log("[STEP 3] Create course, exam, challenge and a test case")
	courseID := createSubmissionCourseHTTP(t, teacherAccess, "result")
	examID := createSubmissionExamHTTP(t, teacherAccess, courseID, "HTTP Submission Result Exam")
	challengeID := createSubmissionChallengeHTTP(t, teacherAccess, examID, "HTTP Submission Result Challenge")
	_ = createSubmissionTestCaseHTTP(t, teacherAccess, challengeID, "tc_result_1", "5")

	t.Log("[STEP 3.1] Publish challenge and create exam Item")
	status, body, err := httputils.PostChallengesPublish(teacherHeaders, map[string]any{"id": challengeID})
	if err != nil {
		t.Fatalf("publish challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected publish challenge status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	status, body, err = httputils.PostExamItems(teacherHeaders, map[string]any{
		"exam_id":     examID,
		"challenge_id": challengeID,
		"order":       1,
		"points":      10,
	})
	if err != nil {
		t.Fatalf("create exam item request failed: %v", err)
	}

	t.Log("[STEP 4] Login as student and create a submission for the challenge")
	studentAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Result")
	studentHeaders := httputils.AuthHeaders(studentAccess)

	enrollBody := map[string]string{"course_id": courseID, "student_id": studentAccess.UserData.ID}
	status, body, err = httputils.PostCoursesEnroll(teacherHeaders, enrollBody)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	createSessionBody := map[string]string{"user_id": studentAccess.UserData.ID, "exam_id": examID}
	status, body, err = httputils.PostSubmissionsSessions(studentHeaders, createSessionBody)
	if err != nil {
		t.Fatalf("create session request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create session status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	sessionID := httputils.MapString(t, httputils.DecodeMap(t, body, "create session"), "id", "create session")

	createSubmissionBody := map[string]any{
		"code":        "def solve(a, b):\n    return a + b",
		"function":	    "solve(a, b)",
		"language":    "python",
		"challenge_id": challengeID,
		"session_id":   sessionID,
	}

	status, body, err = httputils.PostSubmissions(studentHeaders, createSubmissionBody)
	if err != nil {
		t.Fatalf("create submission request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create submission status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	createdSubmission := httputils.DecodeMap(t, body, "create submission")
	submissionID := httputils.MapString(t, createdSubmission, "id", "create submission")

	states := execPollSubmissionStates(t, submissionID, teacherHeaders, 10*time.Second)
	if len(states) == 0 {
		t.Fatalf("expected at least one submission result for submission=%s", submissionID)
	}

	t.Log("[STEP 5] Request submission status until its \"accepted\" (max. 3 tries with 2s wait)")
	accepted := false
	for attempt := 1; attempt <= 3; attempt++ {
		status, body, err = httputils.GetSubmissionsById(studentHeaders, map[string]any{"id": submissionID})
		if err != nil {
			t.Fatalf("get submission status request failed on attempt=%d: %v", attempt, err)
		}
		if status != http.StatusOK {
			t.Fatalf("expected submission status code=%d, got=%d body=%s", http.StatusOK, status, string(body))
		}

		payload := httputils.DecodeMap(t, body, "get submission status")
		currentStates := execExtractResultStates(t, payload)
		if len(currentStates) == 0 {
			if attempt < 3 {
				time.Sleep(2 * time.Second)
				continue
			}
			t.Fatalf("expected submission results in status payload, got body=%s", string(body))
		}

		allAccepted := true
		for _, state := range currentStates {
			if state.Status != "accepted" {
				allAccepted = false
				break
			}
		}

		if allAccepted {
			accepted = true
			t.Log("[OK] Submission results reached accepted status")
			break
		}

		if attempt < 3 {
			time.Sleep(2 * time.Second)
		}
	}

	if !accepted {
		t.Fatalf("expected submission=%s results to reach accepted within 3 attempts", submissionID)
	}

	t.Log("[STEP 6] Delete the course")
	status, body, err = httputils.DeleteCoursesById(teacherHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete course status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 7] Cleanup Challenge via DELETE /challenges/:id")
	status, body, err = httputils.DeleteChallengesById(teacherHeaders, map[string]any{"id": challengeID})
	if err != nil {
		t.Fatalf("delete challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete challenge status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
}
