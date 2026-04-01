package submission_test

import (
	"net/http"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSubmissionsCRUDHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}

	t.Log("[STEP 2] Login de profesor")
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := httputils.AuthHeaders(teacherAccess)
	teacherID := teacherAccess.UserData.ID

	t.Log("[STEP 3] Crear curso y examen")
	courseID := createSubmissionCourseHTTP(t, teacherAccess, "sub")
	examID := createSubmissionExamHTTP(t, teacherAccess, courseID, "HTTP Submission Exam")

	t.Log("[STEP 4] Login estudiante y matrícula")
	studentAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Test")
	studentHeaders := httputils.AuthHeaders(studentAccess)
	studentID := studentAccess.UserData.ID
	enrollBody := map[string]string{"course_id": courseID, "student_id": studentID}
	status, body, err := httputils.PostCoursesEnroll(studentHeaders, enrollBody)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 5] Crear challenge y test-cases")
	challengeID := createSubmissionChallengeHTTP(t, teacherAccess, examID, "HTTP Submission Challenge")
	_ = createSubmissionTestCaseHTTP(t, teacherAccess, challengeID, "tc_submission_1", "5")
	_ = createSubmissionTestCaseHTTP(t, teacherAccess, challengeID, "tc_submission_2", "5")

	t.Log("[STEP 6] Publicar challenge y crear ExamItem")
	status, body, err = httputils.PostChallengesPublish(teacherHeaders, map[string]any{"id": challengeID})
	if err != nil {
		t.Fatalf("publish challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected publish challenge status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	status, body, err = httputils.PostExamItems(teacherHeaders, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        1,
		"points":       10,
	})
	if err != nil {
		t.Fatalf("create exam item request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create exam item status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	t.Log("[STEP 7] Crear sesión para el estudiante")
	createSessionBody := map[string]string{"user_id": studentID, "exam_id": examID}
	status, body, err = httputils.PostSubmissionsSessions(studentHeaders, createSessionBody)
	if err != nil {
		t.Fatalf("create session request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create session status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	sessionID := httputils.MapString(t, httputils.DecodeMap(t, body, "create session"), "id", "create session")

	defer func() {
		if sessionID == "" {
			return
		}
		status, body, err := httputils.PostSubmissionsSessionsClose(teacherHeaders, map[string]any{"id": sessionID})
		if err != nil {
			t.Logf("close session request failed: %v", err)
			return
		}
		if status != http.StatusOK {
			t.Logf("unexpected close session status=%d body=%s", status, string(body))
		}
	}()

	t.Log("[STEP 8] Crear submission")
	createSubmissionBody := map[string]any{
		"code":         "def solve(a, b):\n    return a + b",
		"language":     "python",
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

	t.Log("[STEP 9] Listar submissions por challenge y validar inclusión")
	status, body, err = httputils.GetSubmissionsChallenge(studentHeaders, map[string]any{"challenge_id": challengeID}, nil)
	if err != nil {
		t.Fatalf("list challenge submissions request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected list submissions status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	outputs := httputils.DecodeSliceMap(t, body, "list challenge submissions")
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
		id, ok := subMap["ID"].(string)
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

	t.Log("[STEP 10] Cleanup por curso")
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/courses/"+courseID, nil, teacherHeaders)
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

	t.Log("[STEP 11] Cleanup por Challenge")
	status, body, err = httputils.DeleteChallengesById(teacherHeaders, map[string]any{"id": challengeID})
	if err != nil {
		t.Fatalf("delete challenge request failed: %v", err)
	}

	t.Log("[STEP 12] Cerrar Session")
	status, body, err = httputils.PostSubmissionsSessionsClose(teacherHeaders, map[string]any{"id": sessionID})
	if err != nil {
		t.Logf("close session request failed: %v", err)
		return
	}
	if status != http.StatusOK {
		t.Logf("unexpected close session status=%d body=%s", status, string(body))
	}

	t.Log("[OK] Cleanup completado")

	_ = teacherID
}
