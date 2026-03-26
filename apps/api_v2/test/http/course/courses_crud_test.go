package course_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestCoursesCRUDHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}
	t.Log("[OK] App inicializada")

	t.Log("[STEP 2] Login/registro de profesor por HTTP")
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := httputils.AuthHeaders(teacherAccess)
	t.Logf("[OK] Profesor autenticado. teacherID=%s", teacherAccess.UserData.ID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-HTTP-CRUD-%d", now.UnixNano())
	courseCode := fmt.Sprintf("HTTP-CRUD-%d", now.Unix()%100000)

	t.Logf("[STEP 3] Crear curso via POST /courses code=%s", courseCode)
	createBody := map[string]any{
		"name":            "HTTP Course CRUD",
		"description":     "course created by HTTP test",
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

	created := httputils.DecodeMap(t, body, "create course")
	courseID := httputils.MapString(t, created, "ID", "create course")
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	t.Log("[STEP 4] Actualizar curso via POST /courses/:id")
	updatedName := "HTTP Course CRUD Updated"
	updatedDescription := "course updated by HTTP test"
	updateBody := map[string]any{
		"name":        updatedName,
		"description": updatedDescription,
	}

	status, body, err = httputils.PostCoursesById(teacherHeaders, map[string]any{"id": courseID}, updateBody)
	if err != nil {
		t.Fatalf("update course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected update status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	updated := httputils.DecodeMap(t, body, "update course")
	if httputils.MapString(t, updated, "name", "update course") != updatedName {
		t.Fatalf("expected updated name=%q, got body=%s", updatedName, string(body))
	}
	t.Logf("[OK] Curso actualizado. name=%q", updatedName)

	t.Log("[STEP 5] Consultar curso via GET /courses/:id")
	status, body, err = httputils.GetCoursesById(teacherHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("get course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected get-by-id status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	reloaded := httputils.DecodeMap(t, body, "get course by id")
	if httputils.MapString(t, reloaded, "ID", "get course by id") != courseID {
		t.Fatalf("expected course id=%s, got body=%s", courseID, string(body))
	}
	if httputils.MapString(t, reloaded, "name", "get course by id") != updatedName {
		t.Fatalf("expected updated name=%q, got body=%s", updatedName, string(body))
	}
	t.Log("[OK] Detalles del curso validados")

	t.Log("[STEP 6] Consultar cursos own via GET /courses?scope=owned")
	status, body, err = httputils.GetCourses(teacherHeaders, map[string]any{"scope": "owned", "teacherId": teacherAccess.UserData.ID})
	if err != nil {
		t.Fatalf("list owned courses request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	if !containsCourseID(t, body, courseID) {
		t.Fatalf("expected created courseID=%s in owned list, got body=%s", courseID, string(body))
	}
	t.Log("[OK] Curso aparece en listado owned")

	t.Log("[STEP 7] Eliminar curso via DELETE /courses/:id")
	status, body, err = httputils.DeleteCoursesById(teacherHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	deleted := httputils.DecodeMap(t, body, "delete course")
	removed, ok := deleted["removed"].(bool)
	if !ok || !removed {
		t.Fatalf("expected delete response with removed=true, got body=%s", string(body))
	}
	t.Log("[OK] Curso eliminado")
}

func containsCourseID(t *testing.T, raw []byte, courseID string) bool {
	t.Helper()

	var list []map[string]any
	if err := json.Unmarshal(raw, &list); err != nil {
		t.Fatalf("decode courses list failed: %v body=%s", err, string(raw))
	}

	for _, item := range list {
		if id, ok := item["ID"].(string); ok && id == courseID {
			return true
		}
	}

	return false
}
