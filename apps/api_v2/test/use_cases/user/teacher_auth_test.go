package auth_test

import (
	"context"
	"fmt"
	"testing"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestTeacherAuth(t *testing.T) {
	process := test.StartTestWithApp(t, "Teacher Auth Use Cases")
	email := "test@test.com"
	password := "Password123!"
	name := "Teacher Test"

	var teacherID string

	// [STEP 1] Login Teacher user
	process.StartStep("Iniciar sesión con Cuenta de Docente")
	access, err := process.Application.UserModule.Login.Execute(email, password)
	if err != nil {
		process.Fail("teacher login", err)
	}
	if access == nil || access.UserData == nil || access.UserData.ID == "" {
		process.Fail("teacher login", fmt.Errorf("expected teacher user data with valid ID"))
	}
	if access.Token == nil || access.Token.AccessToken == "" {
		process.Fail("teacher login", fmt.Errorf("expected non-empty access token"))
	}
	teacherID = access.UserData.ID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))

	process.EndStep()

	// [STEP 2] Try to register teacher user and validate response
	process.StartStep("Intentar registrar usuario docente y validar respuesta")
	registeredAccess, registerErr := process.Application.UserModule.Register.Execute(email, name, password)
	if registerErr != nil {
		process.Log(fmt.Sprintf("Register devolvió error (válido si ya existe): %v", registerErr))
	} else {
		if registeredAccess == nil || registeredAccess.UserData == nil || registeredAccess.UserData.ID == "" {
			process.Fail("teacher register", fmt.Errorf("expected registered teacher with valid ID"))
		}
		if registeredAccess.Token == nil || registeredAccess.Token.AccessToken == "" {
			process.Fail("teacher register", fmt.Errorf("expected access token in register response"))
		}
		teacherID = registeredAccess.UserData.ID
		process.Log(fmt.Sprintf("Registro exitoso. teacherID=%s", teacherID))
	}

	process.EndStep()

	// [STEP 3] Get teacher data and validate
	process.StartStep("Obtener datos del docente y validar")
	ctx := services.WithAccessToken(services.WithUserEmail(context.Background(), email), access.Token.AccessToken)
	teacherData, err := process.Application.UserModule.GetData.Execute(ctx, email)
	if err != nil {
		process.Fail("teacher get data", err)
	}
	if teacherData == nil || teacherData.ID == "" {
		process.Fail("teacher get data", fmt.Errorf("expected teacher data with valid ID"))
	}
	if teacherData.Email != email {
		process.Fail("teacher get data", fmt.Errorf("expected teacher email %s, got %s", email, teacherData.Email))
	}
	if teacherID != "" && teacherData.ID != teacherID {
		process.Fail("teacher get data", fmt.Errorf("expected teacher ID %s, got %s", teacherID, teacherData.ID))
	}
	process.Log(fmt.Sprintf("Datos validados. email=%s role=%s", teacherData.Email, teacherData.Role))

	process.EndStep()

	// [STEP 4] Validate teacher role
	process.StartStep("Validar rol de Docente")
	if access.UserData.Role != user_entities.UserRoleProfessor {
		process.Fail("teacher role", fmt.Errorf("expected role %s, got %s", user_entities.UserRoleProfessor, access.UserData.Role))
	}
	process.Log(fmt.Sprintf("Rol validado: %s", access.UserData.Role))

	process.EndStep()
	process.End()
}
