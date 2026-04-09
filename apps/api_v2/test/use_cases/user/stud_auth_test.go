package auth_test

import (
	"context"
	"fmt"
	"testing"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestStudentAuth(t *testing.T) {
	process := test.StartTestWithApp(t, "Student Auth Use Cases")
	email := "stud@test.com"
	password := "Password123!"
	name := "Student Test"

	var studentID string

	// [STEP 1] Login Student user
	process.StartStep("Iniciar sesión con Cuenta de Estudiante")
	access, err := process.Application.UserModule.Login.Execute(email, password)
	if err != nil {
		process.Fail("student login", err)
	}
	if access == nil || access.UserData == nil || access.UserData.ID == "" {
		process.Fail("student login", fmt.Errorf("expected student user data with valid ID"))
	}
	if access.Token == nil || access.Token.AccessToken == "" {
		process.Fail("student login", fmt.Errorf("expected non-empty access token"))
	}
	studentID = access.UserData.ID
	process.Log(fmt.Sprintf("studentID=%s", studentID))

	process.EndStep()

	// [STEP 2] Try to register student user and validate response
	process.StartStep("Intentar registrar usuario estudiante y validar respuesta")
	registeredAccess, registerErr := process.Application.UserModule.Register.Execute(email, name, password)
	if registerErr != nil {
		process.Log(fmt.Sprintf("Register devolvió error (válido si ya existe): %v", registerErr))
	} else {
		if registeredAccess == nil || registeredAccess.UserData == nil || registeredAccess.UserData.ID == "" {
			process.Fail("student register", fmt.Errorf("expected registered student with valid ID"))
		}
		if registeredAccess.Token == nil || registeredAccess.Token.AccessToken == "" {
			process.Fail("student register", fmt.Errorf("expected access token in register response"))
		}
		studentID = registeredAccess.UserData.ID
		process.Log(fmt.Sprintf("Registro exitoso. studentID=%s", studentID))
	}

	process.EndStep()

	// [STEP 3] Get student data and validate
	process.StartStep("Obtener datos del estudiante y validar")
	ctx := services.WithAccessToken(services.WithUserEmail(context.Background(), email), access.Token.AccessToken)
	studentData, err := process.Application.UserModule.GetData.Execute(ctx, email)
	if err != nil {
		process.Fail("student get data", err)
	}
	if studentData == nil || studentData.ID == "" {
		process.Fail("student get data", fmt.Errorf("expected student data with valid ID"))
	}
	if studentData.Email != email {
		process.Fail("student get data", fmt.Errorf("expected student email %s, got %s", email, studentData.Email))
	}
	if studentID != "" && studentData.ID != studentID {
		process.Fail("student get data", fmt.Errorf("expected student ID %s, got %s", studentID, studentData.ID))
	}
	if access.UserData.Role != user_entities.UserRoleStudent {
		process.Fail("student role", fmt.Errorf("expected role %s, got %s", user_entities.UserRoleStudent, access.UserData.Role))
	}
	process.Log(fmt.Sprintf("Datos validados. email=%s role=%s", studentData.Email, studentData.Role))

	process.EndStep()
	process.End()
}
