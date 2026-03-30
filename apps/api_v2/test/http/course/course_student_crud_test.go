package course_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestCoursesWithStudentsCRUDHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}
	t.Log("[OK] App inicializada")

	t.Log("[STEP 2] Autenticando profesor")
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := httputils.AuthHeaders(teacherAccess)
	t.Logf("[OK] Profesor autenticado. teacherID=%s", teacherAccess.UserData.ID)

	t.Log("[STEP 3] Autenticando estudiante")
	studentAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Test")
	studentHeaders := httputils.AuthHeaders(studentAccess)
	t.Logf("[OK] Estudiante autenticado. studentID=%s", studentAccess.UserData.ID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-HTTP-STU-%d", now.UnixNano())
	courseCode := fmt.Sprintf("HTTP-STU-%d", now.Unix()%100000)

	t.Log("[STEP 4] Crear curso para flujo de estudiantes")
	createBody := map[string]any{
		"name":            "HTTP Course Student CRUD",
		"description":     "course for student HTTP flow",
		"visibility":      "public",
		"visual_identity": "#ff0055",
		"code":            courseCode,
		"year":            2026,
		"semester":        "01",
		"enrollment_code": enrollmentCode,
		"enrollment_url":  "https://example.test/enroll/" + enrollmentCode,
		"teacher_id":      teacherAccess.UserData.ID,
	}

	status, body, err := httputils.PostCourses(teacherHeaders, createBody)
	if err != nil {
		t.Fatalf("create course request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	courseID := httputils.MapString(t, httputils.DecodeMap(t, body, "create course"), "id", "create course")
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	t.Log("[STEP 5] Matricular estudiante vía POST /courses/enroll")
	enrollBody := map[string]string{
		"course_id":  courseID,
		"student_id": studentAccess.UserData.ID,
	}
	status, body, err = httputils.PostCoursesEnroll(studentHeaders, enrollBody)
	if err != nil {
		t.Fatalf("enroll request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	t.Log("[OK] Estudiante matriculado")

	t.Log("[STEP 6] Verificar estudiante en GET /courses/:id/students")
	status, body, err = httputils.GetCoursesStudents(studentHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("get students request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected get students status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	if !containsStudentID(t, body, studentAccess.UserData.ID) {
		t.Fatalf("expected studentID=%s in students list, got body=%s", studentAccess.UserData.ID, string(body))
	}
	t.Log("[OK] Estudiante presente en lista")

	t.Log("[STEP 7] Remover estudiante via DELETE /courses/:id/students/:studentId")
	status, body, err = httputils.DeleteCoursesStudent(studentHeaders, map[string]any{"id": courseID, "student_id": studentAccess.UserData.ID})
	if err != nil {
		t.Fatalf("remove student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected remove status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	removed := httputils.DecodeMap(t, body, "remove student")
	if flag, ok := removed["removed"].(bool); !ok || !flag {
		t.Fatalf("expected removed=true, got body=%s", string(body))
	}
	t.Log("[OK] Estudiante removido")

	t.Log("[STEP 8] Validar que ya no aparece en lista de estudiantes")
	status, body, err = httputils.GetCoursesStudents(studentHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("get students after removal request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected get students status=%d after removal, got=%d body=%s", http.StatusOK, status, string(body))
	}
	if containsStudentID(t, body, studentAccess.UserData.ID) {
		t.Fatalf("expected studentID=%s absent after removal, got body=%s", studentAccess.UserData.ID, string(body))
	}
	t.Log("[OK] Validación post-remoción completada")

	t.Log("[STEP 9] Eliminar curso via DELETE /courses/:id")
	status, body, err = httputils.DeleteCoursesById(teacherHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	deleted := httputils.DecodeMap(t, body, "delete course")
	if flag, ok := deleted["removed"].(bool); !ok || !flag {
		t.Fatalf("expected removed=true, got body=%s", string(body))
	}
	t.Log("[OK] Curso eliminado")
}

func containsStudentID(t *testing.T, raw []byte, studentID string) bool {
	t.Helper()

	var list []map[string]any
	if err := json.Unmarshal(raw, &list); err != nil {
		t.Fatalf("decode students list failed: %v body=%s", err, string(raw))
	}

	for _, item := range list {
		if id, ok := item["id"].(string); ok && id == studentID {
			return true
		}
	}

	return false
}
