package user_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestAuthProcessHTTP(t *testing.T) {
	t.Log("[STEP 1] Initialize app for HTTP flow")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}
	t.Log("[OK] App initialized")

	t.Log("[STEP 2] Ensure teacher access via /auth/login or /auth/register")
	teacherAccess := ensureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	if teacherAccess.UserData == nil || teacherAccess.UserData.Email != "test@test.com" {
		t.Fatalf("expected teacher email=%s, got=%v", "test@test.com", teacherAccess.UserData)
	}
	t.Logf("[OK] Teacher access resolved. teacherID=%s", teacherAccess.UserData.ID)

	t.Log("[STEP 3] Ensure student access via /auth/login or /auth/register")
	studentAccess := ensureHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Test")
	if studentAccess.UserData == nil || studentAccess.UserData.Email != "stud@test.com" {
		t.Fatalf("expected student email=%s, got=%v", "stud@test.com", studentAccess.UserData)
	}
	t.Logf("[OK] Student access resolved. studentID=%s", studentAccess.UserData.ID)

	t.Log("[STEP 4] Validate teacher data using /auth/me")
	teacherData := getUserDataHTTP(t, app, teacherAccess)
	if teacherData.ID != teacherAccess.UserData.ID {
		t.Fatalf("expected teacher ID=%s from auth, got=%s from /auth/me", teacherAccess.UserData.ID, teacherData.ID)
	}
	if teacherData.Email != "test@test.com" {
		t.Fatalf("expected teacher email=%s, got=%s", "test@test.com", teacherData.Email)
	}
	t.Log("[OK] Teacher /auth/me validated")

	t.Log("[STEP 5] Validate student data using /auth/me")
	studentData := getUserDataHTTP(t, app, studentAccess)
	if studentData.ID != studentAccess.UserData.ID {
		t.Fatalf("expected student ID=%s from auth, got=%s from /auth/me", studentAccess.UserData.ID, studentData.ID)
	}
	if studentData.Email != "stud@test.com" {
		t.Fatalf("expected student email=%s, got=%s", "stud@test.com", studentData.Email)
	}
	t.Log("[OK] Student /auth/me validated")

	t.Log("[STEP 6] Validate teacher and student are distinct users")
	if teacherAccess.UserData.ID == studentAccess.UserData.ID {
		t.Fatal("expected teacher and student to be different users")
	}
	t.Log("[OK] HTTP auth process completed successfully")
}

func ensureHTTPAuthUserAccess(t *testing.T, app *fiber.App, email, password, name string) *dtos.UserAccess {
	t.Helper()

	t.Logf("[STEP] Attempt HTTP login for %s", email)
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

func getUserDataHTTP(t *testing.T, app *fiber.App, access *dtos.UserAccess) *Entities.User {
	t.Helper()

	headers := map[string]string{
		"Authorization": "Bearer " + access.Token.AccessToken,
		"X-User-Email":  access.UserData.Email,
	}

	status, body, err := httputils.DoJSONRequest(app, http.MethodGet, "/auth/me", nil, headers)
	if err != nil {
		t.Fatalf("/auth/me request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected /auth/me status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	var user Entities.User
	if err := json.Unmarshal(body, &user); err != nil {
		t.Fatalf("decode /auth/me response failed: %v body=%s", err, string(body))
	}

	if user.ID == "" {
		t.Fatal("expected /auth/me response with user ID")
	}

	return &user
}

func decodeUserAccess(t *testing.T, raw []byte, source string) *dtos.UserAccess {
	t.Helper()

	var access dtos.UserAccess
	if err := json.Unmarshal(raw, &access); err != nil {
		t.Fatalf("decode %s response failed: %v body=%s", source, err, string(raw))
	}

	return &access
}

func validateUserAccess(t *testing.T, access *dtos.UserAccess, expectedEmail string) {
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
