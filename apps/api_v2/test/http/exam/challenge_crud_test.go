package exam_test

import (
	"fmt"
	"net/http"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestChallengeCRUDHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app := initExamHTTPApp(t)

	t.Log("[STEP 2] Autenticando profesor")
	teacherAccess := ensureExamHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := authHeaders(teacherAccess)

	t.Log("[STEP 3] Crear curso + examen para challenges")
	courseID := createCourseHTTP(t, app, teacherAccess, "challenge-crud")
	examID := createExamHTTP(t, app, teacherAccess, courseID, "HTTP Challenge Exam")

	t.Log("[STEP 4] Crear challenge")
	challengeID := createChallengeHTTP(t, app, teacherAccess, examID, "HTTP Challenge")
	t.Logf("[OK] Challenge creado. challengeID=%s", challengeID)

	t.Log("[STEP 5] Actualizar challenge")
	updatedTitle := "HTTP Challenge Updated"
	updateBody := map[string]any{
		"title":           updatedTitle,
		"difficulty":      "medium",
		"workerTimeLimit": 2000,
	}
	status, body, err := httputils.DoJSONRequest(app, http.MethodPatch, "/challenges/"+challengeID, updateBody, teacherHeaders)
	if err != nil {
		t.Fatalf("update challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected update status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	updated := decodeMap(t, body, "update challenge")
	if mapString(t, updated, "title", "update challenge") != updatedTitle {
		t.Fatalf("expected title=%q after update, got body=%s", updatedTitle, string(body))
	}
	t.Log("[OK] Challenge actualizado")

	t.Log("[STEP 6] Publicar challenge")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/challenges/"+challengeID+"/publish", nil, teacherHeaders)
	if err != nil {
		t.Fatalf("publish challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected publish status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	published := decodeMap(t, body, "publish challenge")
	if mapString(t, published, "status", "publish challenge") != "published" {
		t.Fatalf("expected Status=published, got body=%s", string(body))
	}
	t.Log("[OK] Challenge publicado")

	t.Log("[STEP 7] Listar challenges por examen")
	path := fmt.Sprintf("/challenges?examId=%s", examID)
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, path, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("list challenges request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	challenges := decodeSliceMap(t, body, "list challenges")
	if !containsID(challenges, challengeID) {
		t.Fatalf("expected challengeID=%s in list, got body=%s", challengeID, string(body))
	}
	t.Logf("[OK] Challenge encontrado en listado. total=%d", len(challenges))

	t.Log("[STEP 8] Archivar challenge")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/challenges/"+challengeID+"/archive", nil, teacherHeaders)
	if err != nil {
		t.Fatalf("archive challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected archive status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	archived := decodeMap(t, body, "archive challenge")
	if mapString(t, archived, "status", "archive challenge") != "archived" {
		t.Fatalf("expected Status=archived, got body=%s", string(body))
	}
	t.Log("[OK] Challenge archivado")

	t.Log("[STEP 9] Cleanup via DELETE /courses/:id")
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

func TestChallengeFromStudentViewHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app := initExamHTTPApp(t)

	t.Log("[STEP 2] Autenticando profesor")
	teacherAccess := ensureExamHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := authHeaders(teacherAccess)

	t.Log("[STEP 3] Crear curso + examen")
	courseID := createCourseHTTP(t, app, teacherAccess, "challenge-student")
	examID := createExamHTTP(t, app, teacherAccess, courseID, "HTTP Challenge Student View Exam")

	t.Log("[STEP 4] Crear 3 challenges en draft")
	publishedID := createChallengeHTTP(t, app, teacherAccess, examID, "Challenge Published")
	archivedID := createChallengeHTTP(t, app, teacherAccess, examID, "Challenge Archived")
	draftID := createChallengeHTTP(t, app, teacherAccess, examID, "Challenge Draft")

	t.Log("[STEP 5] Publicar uno y archivar otro")
	_, _, _ = httputils.DoJSONRequest(app, http.MethodPost, "/challenges/"+publishedID+"/publish", nil, teacherHeaders)
	_, _, _ = httputils.DoJSONRequest(app, http.MethodPost, "/challenges/"+archivedID+"/publish", nil, teacherHeaders)
	_, _, _ = httputils.DoJSONRequest(app, http.MethodPost, "/challenges/"+archivedID+"/archive", nil, teacherHeaders)

	t.Log("[STEP 6] Autenticando estudiante y matriculando")
	studentAccess := ensureExamHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Test")
	studentHeaders := authHeaders(studentAccess)
	enrollBody := map[string]string{"course_id": courseID, "student_id": studentAccess.UserData.ID}
	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/courses/enroll", enrollBody, studentHeaders)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 7] Listar challenges como estudiante")
	path := fmt.Sprintf("/challenges?examId=%s", examID)
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, path, nil, studentHeaders)
	if err != nil {
		t.Fatalf("list challenges as student failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	studentChallenges := decodeSliceMap(t, body, "student challenge list")
	foundPublished := containsID(studentChallenges, publishedID)
	foundArchived := containsID(studentChallenges, archivedID)
	foundDraft := containsID(studentChallenges, draftID)

	if !foundPublished {
		t.Fatal("expected published challenge in student view")
	}
	if foundArchived {
		t.Fatal("did not expect archived challenge in student view")
	}
	if foundDraft {
		t.Fatal("did not expect draft challenge in student view")
	}
	t.Logf("[OK] Vista estudiante validada. visibleChallenges=%d", len(studentChallenges))

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
