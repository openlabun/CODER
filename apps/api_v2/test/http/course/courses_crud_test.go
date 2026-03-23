package course_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	user_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
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
	teacherAccess := ensureCourseHTTPAuthUserAccess(t, app, "test@test.com", "Testing123!", "Teacher Test")
	teacherHeaders := authHeaders(teacherAccess)
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

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/courses/", createBody, teacherHeaders)
	if err != nil {
		t.Fatalf("create course request failed: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected create status=%d, got=%d body=%s", http.StatusCreated, status, string(body))
	}

	created := decodeMap(t, body, "create course")
	courseID := mapString(t, created, "ID", "create course")
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	t.Log("[STEP 4] Actualizar curso via POST /courses/:id")
	updatedName := "HTTP Course CRUD Updated"
	updatedDescription := "course updated by HTTP test"
	updateBody := map[string]any{
		"name":        updatedName,
		"description": updatedDescription,
	}

	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/courses/"+courseID, updateBody, teacherHeaders)
	if err != nil {
		t.Fatalf("update course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected update status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	updated := decodeMap(t, body, "update course")
	if mapString(t, updated, "Name", "update course") != updatedName {
		t.Fatalf("expected updated name=%q, got body=%s", updatedName, string(body))
	}
	t.Logf("[OK] Curso actualizado. name=%q", updatedName)

	t.Log("[STEP 5] Consultar curso via GET /courses/:id")
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, "/courses/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("get course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected get-by-id status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	reloaded := decodeMap(t, body, "get course by id")
	if mapString(t, reloaded, "ID", "get course by id") != courseID {
		t.Fatalf("expected course id=%s, got body=%s", courseID, string(body))
	}
	if mapString(t, reloaded, "Name", "get course by id") != updatedName {
		t.Fatalf("expected updated name=%q, got body=%s", updatedName, string(body))
	}
	t.Log("[OK] Detalles del curso validados")

	t.Log("[STEP 6] Consultar cursos own via GET /courses?scope=owned")
	listPath := fmt.Sprintf("/courses?scope=owned&teacherId=%s", teacherAccess.UserData.ID)
	status, body, err = httputils.DoJSONRequest(app, http.MethodGet, listPath, nil, teacherHeaders)
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
	status, body, err = httputils.DoJSONRequest(app, http.MethodDelete, "/courses/"+courseID, nil, teacherHeaders)
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	deleted := decodeMap(t, body, "delete course")
	removed, ok := deleted["removed"].(bool)
	if !ok || !removed {
		t.Fatalf("expected delete response with removed=true, got body=%s", string(body))
	}
	t.Log("[OK] Curso eliminado")
}

func ensureCourseHTTPAuthUserAccess(t *testing.T, app *fiber.App, email, password, name string) *user_dtos.UserAccess {
	t.Helper()

	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodPost, "/auth/login", loginBody, nil)
	if err != nil {
		t.Fatalf("login request failed for %s: %v", email, err)
	}

	if status == http.StatusOK {
		access := decodeUserAccess(t, body, "login")
		validateUserAccess(t, access, email)
		return access
	}

	t.Logf("login failed for %s with status=%d, trying register", email, status)
	registerBody := map[string]string{
		"email":    email,
		"name":     name,
		"password": password,
	}

	status, body, err = httputils.DoJSONRequest(app, http.MethodPost, "/auth/register", registerBody, nil)
	if err != nil {
		t.Fatalf("register request failed for %s: %v", email, err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected register status=%d for %s, got=%d body=%s", http.StatusCreated, email, status, string(body))
	}

	access := decodeUserAccess(t, body, "register")
	validateUserAccess(t, access, email)
	return access
}

func decodeUserAccess(t *testing.T, raw []byte, source string) *user_dtos.UserAccess {
	t.Helper()

	var access user_dtos.UserAccess
	if err := json.Unmarshal(raw, &access); err != nil {
		t.Fatalf("decode %s response failed: %v body=%s", source, err, string(raw))
	}

	return &access
}

func validateUserAccess(t *testing.T, access *user_dtos.UserAccess, expectedEmail string) {
	t.Helper()

	if access == nil || access.UserData == nil || access.Token == nil {
		t.Fatalf("expected valid access payload for %s", expectedEmail)
	}
	if access.UserData.ID == "" {
		t.Fatalf("expected user ID in access payload for %s", expectedEmail)
	}
	if access.Token.AccessToken == "" {
		t.Fatalf("expected access token in payload for %s", expectedEmail)
	}
	if access.UserData.Email != expectedEmail {
		t.Fatalf("expected email=%s in payload, got=%s", expectedEmail, access.UserData.Email)
	}
}

func authHeaders(access *user_dtos.UserAccess) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + access.Token.AccessToken,
		"X-User-Email":  access.UserData.Email,
	}
}

func decodeMap(t *testing.T, raw []byte, source string) map[string]any {
	t.Helper()

	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode %s map response failed: %v body=%s", source, err, string(raw))
	}

	return out
}

func mapString(t *testing.T, m map[string]any, key, source string) string {
	t.Helper()

	v, ok := m[key]
	if !ok {
		t.Fatalf("missing key=%s in %s response", key, source)
	}
	s, ok := v.(string)
	if !ok {
		t.Fatalf("key=%s in %s response is not string (type=%T)", key, source, v)
	}
	return s
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
