package exam_test

import (
	"net/http"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestTestCasesCRUDHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app := initExamHTTPApp(t)

	t.Log("[STEP 2] Autenticando profesor")
	teacherAccess := ensureExamHTTPAuthUserAccess(t, app, "test@test.com", "Testing123!", "Teacher Test")
	teacherHeaders := authHeaders(teacherAccess)

	t.Log("[STEP 3] Crear curso + examen + challenge")
	courseID := createCourseHTTP(t, app, teacherAccess, "tc-crud")
	examID := createExamHTTP(t, app, teacherAccess, courseID, "HTTP TestCase Exam")
	challengeID := createChallengeHTTP(t, app, teacherAccess, examID, "HTTP TestCase Challenge")

	t.Log("[STEP 4] Crear test case")
	testCaseID := createTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_add_1", true)
	t.Logf("[OK] Test case creado. testCaseID=%s", testCaseID)

	t.Log("[STEP 5] Obtener test cases por challenge")
	status, body, err := httputils.DoJSONRequest(app, http.MethodGet, "/test-cases/challenge/"+challengeID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("get test-cases by challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected get list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	testCases := decodeSliceMap(t, body, "get test-cases by challenge")
	if !containsID(testCases, testCaseID) {
		t.Fatalf("expected testCaseID=%s in challenge list, got body=%s", testCaseID, string(body))
	}
	t.Logf("[OK] Test case encontrado en listado. total=%d", len(testCases))

	t.Log("[STEP 6] Eliminar test case")
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/test-cases/"+testCaseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("delete test-case request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	deleted := decodeMap(t, body, "delete test-case")
	if !mapBool(t, deleted, "deleted", "delete test-case") {
		t.Fatalf("expected deleted=true, got body=%s", string(body))
	}
	t.Log("[OK] Test case eliminado")

	t.Log("[STEP 7] Validar ausencia del test case en el listado")
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, "/test-cases/challenge/"+challengeID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("get test-cases after delete request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected get list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	afterDelete := decodeSliceMap(t, body, "get test-cases after delete")
	if containsID(afterDelete, testCaseID) {
		t.Fatalf("expected deleted testCaseID=%s absent from list, got body=%s", testCaseID, string(body))
	}
	t.Log("[OK] Ausencia del test case validada")

	t.Log("[STEP 8] Cleanup via DELETE /courses/:id")
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/courses/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected course delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	if !mapBool(t, decodeMap(t, body, "delete course"), "removed", "delete course") {
		t.Fatalf("expected removed=true for course delete, got body=%s", string(body))
	}
	t.Log("[OK] Cleanup completado")
}

func TestTestCasesFromStudentViewHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app := initExamHTTPApp(t)

	t.Log("[STEP 2] Autenticando profesor")
	teacherAccess := ensureExamHTTPAuthUserAccess(t, app, "test@test.com", "Testing123!", "Teacher Test")
	teacherHeaders := authHeaders(teacherAccess)

	t.Log("[STEP 3] Crear curso + examen + challenge")
	courseID := createCourseHTTP(t, app, teacherAccess, "tc-student")
	examID := createExamHTTP(t, app, teacherAccess, courseID, "HTTP TestCase Student View Exam")
	challengeID := createChallengeHTTP(t, app, teacherAccess, examID, "HTTP TestCase Student View Challenge")

	t.Log("[STEP 4] Publicar challenge")
	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/challenges/"+challengeID+"/publish", nil, teacherHeaders)
	if err != nil {
		t.Fatalf("publish challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected publish status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 5] Crear test-cases (2 sample + 1 private)")
	publicTC1 := createTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_public_1", true)
	publicTC2 := createTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_public_2", true)
	privateTC := createTestCaseHTTP(t, app, teacherAccess, challengeID, "tc_private_1", false)

	t.Log("[STEP 6] Autenticando estudiante y matriculando")
	studentAccess := ensureExamHTTPAuthUserAccess(t, app, "stud@test.com", "Testing123!", "Student Test")
	studentHeaders := authHeaders(studentAccess)
	enrollBody := map[string]string{"course_id": courseID, "student_id": studentAccess.UserData.ID}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/courses/enroll", enrollBody, studentHeaders)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 7] Obtener test-cases como estudiante")
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, "/test-cases/challenge/"+challengeID, nil, studentHeaders)
	if err != nil {
		t.Fatalf("get test-cases as student failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected get list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	studentView := decodeSliceMap(t, body, "student test-cases list")
	foundPublic1 := containsID(studentView, publicTC1)
	foundPublic2 := containsID(studentView, publicTC2)
	foundPrivate := containsID(studentView, privateTC)

	if !foundPublic1 || !foundPublic2 {
		t.Fatalf("expected both public test-cases visible to student, got body=%s", string(body))
	}
	if foundPrivate {
		t.Fatalf("did not expect private test-case in student view, got body=%s", string(body))
	}
	t.Logf("[OK] Restricciones de vista estudiante validadas. visibleTestCases=%d", len(studentView))

	t.Log("[STEP 8] Cleanup via DELETE /courses/:id")
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/courses/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected course delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	if !mapBool(t, decodeMap(t, body, "delete course"), "removed", "delete course") {
		t.Fatalf("expected removed=true for course delete, got body=%s", string(body))
	}
	t.Log("[OK] Cleanup completado")
}
