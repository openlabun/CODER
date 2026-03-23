package exam_test

import (
	"fmt"
	"net/http"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamCRUDHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app := initExamHTTPApp(t)
	t.Log("[OK] App inicializada")

	t.Log("[STEP 2] Login/registro de profesor por HTTP")
	teacherAccess := ensureExamHTTPAuthUserAccess(t, app, "test@test.com", "Testing123!", "Teacher Test")
	teacherHeaders := authHeaders(teacherAccess)
	t.Logf("[OK] Profesor autenticado. teacherID=%s", teacherAccess.UserData.ID)

	t.Log("[STEP 3] Crear curso para examenes")
	courseID := createCourseHTTP(t, app, teacherAccess, "exam-crud")
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	t.Log("[STEP 4] Crear examen 1 via POST /exams")
	examID1 := createExamHTTP(t, app, teacherAccess, courseID, "HTTP Exam 1")
	t.Logf("[OK] Examen 1 creado. examID=%s", examID1)

	t.Log("[STEP 5] Actualizar examen 1 via PATCH /exams/:id")
	updatedTitle := "HTTP Exam 1 Updated"
	updateBody := map[string]any{
		"title":       updatedTitle,
		"description": "Updated by HTTP test",
		"try_limit":   3,
	}
	status, body, err := httputils.DoJSONRequest(app, http.MethodPatch, "/exams/"+examID1, updateBody, teacherHeaders)
	if err != nil {
		t.Fatalf("update exam request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected update status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	updated := decodeMap(t, body, "update exam")
	if mapString(t, updated, "Title", "update exam") != updatedTitle {
		t.Fatalf("expected updated title=%q, got body=%s", updatedTitle, string(body))
	}
	t.Logf("[OK] Examen 1 actualizado. title=%q", updatedTitle)

	t.Log("[STEP 6] Obtener examen 1 via GET /exams/:id")
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, "/exams/"+examID1, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("get exam request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected get-by-id status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	reloaded := decodeMap(t, body, "get exam by id")
	if mapString(t, reloaded, "ID", "get exam by id") != examID1 {
		t.Fatalf("expected exam id=%s, got body=%s", examID1, string(body))
	}
	if mapString(t, reloaded, "Title", "get exam by id") != updatedTitle {
		t.Fatalf("expected title=%q, got body=%s", updatedTitle, string(body))
	}
	t.Log("[OK] Detalle de examen 1 validado")

	t.Log("[STEP 7] Cambiar visibilidad del examen 1")
	visibilityBody := map[string]any{"visibility": "private"}
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/exams/"+examID1+"/visibility", visibilityBody, teacherHeaders)
	if err != nil {
		t.Fatalf("change visibility request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected visibility status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	t.Log("[OK] Visibilidad actualizada")

	t.Log("[STEP 8] Cerrar examen 1")
	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/exams/"+examID1+"/close", nil, teacherHeaders)
	if err != nil {
		t.Fatalf("close exam request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected close status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	t.Log("[OK] Examen cerrado")

	t.Log("[STEP 9] Crear examen 2")
	examID2 := createExamHTTP(t, app, teacherAccess, courseID, "HTTP Exam 2")
	t.Logf("[OK] Examen 2 creado. examID=%s", examID2)

	t.Log("[STEP 10] Listar examenes por curso via GET /exams/course/:courseId")
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, "/exams/course/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("get exams by course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	exams := decodeSliceMap(t, body, "get exams by course")
	if !containsID(exams, examID1) {
		t.Fatalf("expected examID1=%s in course exam list, got body=%s", examID1, string(body))
	}
	if !containsID(exams, examID2) {
		t.Fatalf("expected examID2=%s in course exam list, got body=%s", examID2, string(body))
	}
	t.Logf("[OK] Listado validado. totalExams=%d", len(exams))

	t.Log("[STEP 11] Cleanup via DELETE /courses/:id")
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/courses/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected course delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	deleted := decodeMap(t, body, "delete course")
	if !mapBool(t, deleted, "removed", "delete course") {
		t.Fatalf("expected removed=true for course delete, got body=%s", string(body))
	}
	t.Logf("[OK] Curso eliminado tras flujo de examenes. courseID=%s", courseID)

	_ = fmt.Sprintf("%s", examID2)
}
