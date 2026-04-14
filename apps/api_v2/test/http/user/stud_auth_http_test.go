package user_test

import (
	"fmt"
	"testing"

	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestStudentAuthHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Student Auth HTTP")
	email := "stud@test.com"
	password := "Password123!"
	name := "Student Test"

	var studentID string
	var access *httputils.HTTPAccess

	// [STEP 1] Login Student user
	process.StartStep("Iniciar sesión con Cuenta de Estudiante")
	resp := httputils.PostAuthLogin(t, app, email, password)
	httputils.RequireStatus(t, resp, 200, "student login")

	access = httputils.ParseAccessResponse(t, resp, email)
	studentID = access.UserID
	if studentID == "" {
		process.Fail("student login", fmt.Errorf("expected student ID"))
	}
	process.Log(fmt.Sprintf("studentID=%s", studentID))
	process.EndStep()

	// [STEP 2] Try to register student user and validate response
	process.StartStep("Intentar registrar usuario estudiante y validar respuesta")
	resp = httputils.PostAuthRegister(t, app, email, name, password)
	if resp.StatusCode == 201 {
		registeredAccess := httputils.ParseAccessResponse(t, resp, email)
		studentID = registeredAccess.UserID
		process.Log(fmt.Sprintf("Registro exitoso. studentID=%s", studentID))
	} else if resp.StatusCode == 400 {
		process.Log("Register devolvio error (valido si ya existe)")
	} else {
		process.Fail("student register", fmt.Errorf("unexpected register status=%d body=%s", resp.StatusCode, string(resp.Body)))
	}
	process.EndStep()

	// [STEP 3] Get student data and validate
	process.StartStep("Obtener datos del estudiante y validar")
	resp = httputils.GetAuthMe(t, app, access, email)
	httputils.RequireStatus(t, resp, 200, "student get data")

	body := httputils.MustJSONMap(t, resp)
	if httputils.StringField(body, "id") == "" {
		process.Fail("student get data", fmt.Errorf("expected student id in response"))
	}
	if httputils.StringField(body, "email") != email {
		process.Fail("student get data", fmt.Errorf("expected student email %s, got %s", email, httputils.StringField(body, "email")))
	}
	if studentID != "" && httputils.StringField(body, "id") != studentID {
		process.Fail("student get data", fmt.Errorf("expected student ID %s, got %s", studentID, httputils.StringField(body, "id")))
	}
	if httputils.StringField(body, "role") != string(user_constants.UserRoleStudent) {
		process.Fail("student role", fmt.Errorf("expected role %s, got %s", user_constants.UserRoleStudent, httputils.StringField(body, "role")))
	}
	process.Log(fmt.Sprintf("Datos validados. email=%s role=%s", httputils.StringField(body, "email"), httputils.StringField(body, "role")))
	process.EndStep()

	process.End()
}
