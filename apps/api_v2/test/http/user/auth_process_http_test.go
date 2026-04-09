package user_test

import (
	"fmt"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestAuthProcessHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Auth Process HTTP")

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess

	// [STEP 1] Login Teacher user
	process.StartStep("Autenticar usuario docente")
	teacherResp := httputils.PostAuthLogin(t, app, "test@test.com", "Password123!")
	httputils.RequireStatus(t, teacherResp, 200, "teacher login")
	teacherAccess = httputils.ParseAccessResponse(t, teacherResp, "test@test.com")
	process.Log(fmt.Sprintf("teacherID=%s", teacherAccess.UserID))
	process.EndStep()

	// [STEP 2] Login Student user
	process.StartStep("Autenticar usuario estudiante")
	studentResp := httputils.PostAuthLogin(t, app, "stud@test.com", "Password123!")
	httputils.RequireStatus(t, studentResp, 200, "student login")
	studentAccess = httputils.ParseAccessResponse(t, studentResp, "stud@test.com")
	process.Log(fmt.Sprintf("studentID=%s", studentAccess.UserID))
	process.EndStep()

	// [STEP 3] Get teacher data and validate
	process.StartStep("Obtener datos de docente por /auth/me")
	teacherMeResp := httputils.GetAuthMe(t, app, teacherAccess, teacherAccess.Email)
	httputils.RequireStatus(t, teacherMeResp, 200, "get teacher me")
	teacherBody := httputils.MustJSONMap(t, teacherMeResp)
	v1 := httputils.StringField(teacherBody, "id")
	v2 := httputils.StringField(teacherBody, "email")
	if v1 != teacherAccess.UserID {
		process.Fail("get teacher me", fmt.Errorf("expected id=%s, got=%s", teacherAccess.UserID, v1))
	}
	if v2 != "test@test.com" {
		process.Fail("get teacher me", fmt.Errorf("expected email=%s, got=%s", "test@test.com", v2))
	}
	process.EndStep()

	// [STEP 4] Get student data and validate
	process.StartStep("Obtener datos de estudiante por /auth/me")
	studentMeResp := httputils.GetAuthMe(t, app, studentAccess, studentAccess.Email)
	httputils.RequireStatus(t, studentMeResp, 200, "get student me")
	studentBody := httputils.MustJSONMap(t, studentMeResp)
	v1 = httputils.StringField(studentBody, "id")
	v2 = httputils.StringField(studentBody, "email")
	if v1 != studentAccess.UserID {
		process.Fail("get student me", fmt.Errorf("expected id=%s, got=%s", studentAccess.UserID, v1))
	}
	if v2 != "stud@test.com" {
		process.Fail("get student me", fmt.Errorf("expected email=%s, got=%s", "stud@test.com", v2))
	}
	process.EndStep()

	// [STEP 5] Validate that teacher and student are different users
	process.StartStep("Validar que docente y estudiante son usuarios distintos")
	if teacherAccess.UserID == studentAccess.UserID {
		process.Fail("distinct users", fmt.Errorf("expected teacher and student to be different users"))
	}
	process.EndStep()

	process.End()
}
