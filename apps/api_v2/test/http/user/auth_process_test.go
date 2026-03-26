package user_test

import (
	"testing"

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
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	if teacherAccess.UserData == nil || teacherAccess.UserData.Email != "test@test.com" {
		t.Fatalf("expected teacher email=%s, got=%v", "test@test.com", teacherAccess.UserData)
	}
	t.Logf("[OK] Teacher access resolved. teacherID=%s", teacherAccess.UserData.ID)

	t.Log("[STEP 3] Ensure student access via /auth/login or /auth/register")
	studentAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Test")
	if studentAccess.UserData == nil || studentAccess.UserData.Email != "stud@test.com" {
		t.Fatalf("expected student email=%s, got=%v", "stud@test.com", studentAccess.UserData)
	}
	t.Logf("[OK] Student access resolved. studentID=%s", studentAccess.UserData.ID)

	t.Log("[STEP 4] Validate teacher data using /auth/me")
	teacherData := httputils.GetUserDataHTTP(t, app, teacherAccess)
	if teacherData.ID != teacherAccess.UserData.ID {
		t.Fatalf("expected teacher ID=%s from auth, got=%s from /auth/me", teacherAccess.UserData.ID, teacherData.ID)
	}
	if teacherData.Email != "test@test.com" {
		t.Fatalf("expected teacher email=%s, got=%s", "test@test.com", teacherData.Email)
	}
	t.Log("[OK] Teacher /auth/me validated")

	t.Log("[STEP 5] Validate student data using /auth/me")
	studentData := httputils.GetUserDataHTTP(t, app, studentAccess)
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

