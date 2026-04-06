package user_test

import (
	"testing"

	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestStudentAuth (t *testing.T) {
	process := test.StartTestWithApp(t, "Course Creation and Persistence")
	email := "stud@test.com"
	password := "Password123!"

	// [STEP 1] Login Student user
	process.StartStep("Iniciar sesión con Cuenta de Estudiante")

	process.EndStep()


	// [STEP 2] Try to register student user and validate response
	process.StartStep("Intentar registrar usuario estudiante y validar respuesta")

	process.EndStep()


	// [STEP 3] Get student data and validate
	process.StartStep("Obtener datos del estudiante y validar")

	process.EndStep()
	process.End()

	_ = email
	_ = password
}