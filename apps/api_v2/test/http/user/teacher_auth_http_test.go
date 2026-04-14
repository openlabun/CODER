package user_test

import (
	"fmt"
	"testing"

	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestTeacherAuthHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Teacher Auth HTTP")
	email := "test@test.com"
	password := "Password123!"
	name := "Teacher Test"

	var teacherID string
	var access *httputils.HTTPAccess

	// [STEP 1] Login Teacher user
	process.StartStep("Iniciar sesión con Cuenta de Docente")
	resp := httputils.PostAuthLogin(t, app, email, password)
	httputils.RequireStatus(t, resp, 200, "teacher login")

	access = httputils.ParseAccessResponse(t, resp, email)
	teacherID = access.UserID
	if teacherID == "" {
		process.Fail("teacher login", fmt.Errorf("expected teacher ID"))
	}
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))
	process.EndStep()

	// [STEP 2] Try to register teacher user and validate response
	process.StartStep("Intentar registrar usuario docente y validar respuesta")
	resp = httputils.PostAuthRegister(t, app, email, name, password)
	if resp.StatusCode == 201 {
		registeredAccess := httputils.ParseAccessResponse(t, resp, email)
		teacherID = registeredAccess.UserID
		process.Log(fmt.Sprintf("Registro exitoso. teacherID=%s", teacherID))
	} else if resp.StatusCode == 400 {
		process.Log("Register devolvio error (valido si ya existe)")
	} else {
		process.Fail("teacher register", fmt.Errorf("unexpected register status=%d body=%s", resp.StatusCode, string(resp.Body)))
	}
	process.EndStep()

	// [STEP 3] Get teacher data and validate
	process.StartStep("Obtener datos del docente y validar")
	resp = httputils.GetAuthMe(t, app, access, email)
	httputils.RequireStatus(t, resp, 200, "teacher get data")

	body := httputils.MustJSONMap(t, resp)
	if httputils.StringField(body, "id") == "" {
		process.Fail("teacher get data", fmt.Errorf("expected teacher id in response"))
	}
	if httputils.StringField(body, "email") != email {
		process.Fail("teacher get data", fmt.Errorf("expected teacher email %s, got %s", email, httputils.StringField(body, "email")))
	}
	if teacherID != "" && httputils.StringField(body, "id") != teacherID {
		process.Fail("teacher get data", fmt.Errorf("expected teacher ID %s, got %s", teacherID, httputils.StringField(body, "id")))
	}
	process.Log(fmt.Sprintf("Datos validados. email=%s role=%s", httputils.StringField(body, "email"), httputils.StringField(body, "role")))
	process.EndStep()

	// [STEP 4] Validate teacher role
	process.StartStep("Validar rol de Docente")
	if httputils.StringField(body, "role") != string(user_constants.UserRoleProfessor) {
		process.Fail("teacher role", fmt.Errorf("expected role %s, got %s", user_constants.UserRoleProfessor, httputils.StringField(body, "role")))
	}
	process.Log(fmt.Sprintf("Rol validado: %s", httputils.StringField(body, "role")))
	process.EndStep()

	process.End()
}
