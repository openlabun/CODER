package auth_test

import (
	"testing"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestAuthProcess(t *testing.T) {
	t.Log("[STEP 1] Initialize application container with dependencies")
	app, err := container.BuildApplicationContainer()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application initialized")

	t.Log("[STEP 2] Attempt teacher login and validate access payload")
	teacherAccess := utils.EnsureAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	if teacherAccess.UserData == nil || teacherAccess.UserData.Email != "test@test.com" {
		t.Fatalf("expected teacher email=%s, got=%v", "test@test.com", teacherAccess.UserData)
	}
	t.Logf("[OK] Teacher access resolved. teacherID=%s", teacherAccess.UserData.ID)

	t.Log("[STEP 3] Attempt student login and validate access payload")
	studentAccess := utils.EnsureAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Test")
	if studentAccess.UserData == nil || studentAccess.UserData.Email != "stud@test.com" {
		t.Fatalf("expected student email=%s, got=%v", "stud@test.com", studentAccess.UserData)
	}
	t.Logf("[OK] Student access resolved. studentID=%s", studentAccess.UserData.ID)

	t.Log("[STEP 4] Get user data for teacher and verify consistency")
	ctx := utils.BuildUserCtx(teacherAccess)
	teacherData, err := app.UserModule.GetData.Execute(ctx, "test@test.com")
	if err != nil {
		t.Fatalf("teacher get data failed: %v", err)
	}
	if teacherData == nil || teacherData.ID == "" {
		t.Fatal("expected teacher data with ID")
	}
	if teacherData.Email != "test@test.com" {
		t.Fatalf("expected teacher email=%s, got=%s", "test@test.com", teacherData.Email)
	}
	if teacherData.ID != teacherAccess.UserData.ID {
		t.Fatalf("expected teacher ID=%s from login, got=%s from get-data", teacherAccess.UserData.ID, teacherData.ID)
	}
	t.Log("[OK] Teacher get-data validated")

	t.Log("[STEP 5] Get user data for student and verify consistency")
	ctx = utils.BuildUserCtx(studentAccess)
	studentData, err := app.UserModule.GetData.Execute(ctx, "stud@test.com")
	if err != nil {
		t.Fatalf("student get data failed: %v", err)
	}
	if studentData == nil || studentData.ID == "" {
		t.Fatal("expected student data with ID")
	}
	if studentData.Email != "stud@test.com" {
		t.Fatalf("expected student email=%s, got=%s", "stud@test.com", studentData.Email)
	}
	if studentData.ID != studentAccess.UserData.ID {
		t.Fatalf("expected student ID=%s from login, got=%s from get-data", studentAccess.UserData.ID, studentData.ID)
	}
	t.Log("[OK] Student get-data validated")

	t.Log("[STEP 6] Verify teacher and student are distinct users")
	if teacherAccess.UserData.ID == studentAccess.UserData.ID {
		t.Fatal("expected teacher and student to be different users")
	}
	t.Log("[OK] Auth process finished successfully")
}



