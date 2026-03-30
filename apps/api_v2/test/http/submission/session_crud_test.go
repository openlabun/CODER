package submission_test

import (
	"net/http"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestCreateSessionHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}

	t.Log("[STEP 2] Login de profesor")
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := httputils.AuthHeaders(teacherAccess)

	t.Log("[STEP 3] Crear curso de prueba")
	courseID := createSubmissionCourseHTTP(t, teacherAccess, "sess")

	t.Log("[STEP 4] Crear examen de prueba")
	examID := createSubmissionExamHTTP(t, teacherAccess, courseID, "HTTP Submission Session Exam")

	t.Log("[STEP 5] Login de estudiante y matrícula")
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

	t.Log("[STEP 6] Crear sesión para estudiante")
	createSessionBody := map[string]string{"user_id": studentID, "exam_id": examID}
	status, body, err = httputils.PostSubmissionsSessions(studentHeaders, createSessionBody)
	if err != nil {
		t.Fatalf("create student session request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create session status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}
	t.Log("Session:", string(body))
	studentSession := httputils.DecodeMap(t, body, "create student session")
	studentSessionID := httputils.MapString(t, studentSession, "id", "create student session")
	if httputils.MapString(t, studentSession, "user_id", "create student session") != studentID {
		t.Fatalf("expected student session UserID=%s, got body=%s", studentID, string(body))
	}
	if httputils.MapString(t, studentSession, "exam_id", "create student session") != examID {
		t.Fatalf("expected student session ExamID=%s, got body=%s", examID, string(body))
	}

	t.Log("[STEP 7] Validar rechazo de sesión duplicada")
	status, body, err = httputils.PostSubmissionsSessions(studentHeaders, createSessionBody)
	if err != nil {
		t.Fatalf("duplicate session request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Fatalf("expected duplicate session status=%d, got=%d body=%s", http.StatusBadRequest, status, string(body))
	}

	t.Log("[STEP 8] Validar error con examen inexistente")
	observerAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Observer Test")
	observerHeaders := httputils.AuthHeaders(observerAccess)
	nonExistingExamBody := map[string]string{"user_id": studentID, "exam_id": "non-existing-exam-id"}
	status, body, err = httputils.PostSubmissionsSessions(observerHeaders, nonExistingExamBody)
	if err != nil {
		t.Fatalf("create session non-existing exam request failed: %v", err)
	}
	if status != http.StatusBadRequest {
		t.Fatalf("expected non-existing exam status=%d, got=%d body=%s", http.StatusBadRequest, status, string(body))
	}

	t.Log("[STEP 9] Ejecutar heartbeat en sesión de estudiante")
	status, body, err = httputils.PostSubmissionsSessionsHeartbeat(studentHeaders, map[string]any{"id": studentSessionID})
	if err != nil {
		t.Fatalf("heartbeat request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected heartbeat status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	heartbeatSession := httputils.DecodeMap(t, body, "heartbeat session")
	if httputils.MapString(t, heartbeatSession, "id", "heartbeat session") != studentSessionID {
		t.Fatalf("expected heartbeat session ID=%s, got body=%s", studentSessionID, string(body))
	}

	t.Log("[STEP 10] Cleanup por curso")
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

	_ = studentSessionID
}