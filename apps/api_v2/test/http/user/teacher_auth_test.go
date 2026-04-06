package user_test

import (
	"testing"

	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestTeacherAuth (t *testing.T) {
	process := test.StartTestWithApp(t, "Course Creation and Persistence")
	email := "test@test.com"
	password := "Password123!"

	// [STEP 1] Login Teacher user
	process.StartStep("Iniciar sesión con Cuenta de Docente")

	process.EndStep()


	// [STEP 2] Try to register teacher user and validate response
	process.StartStep("Intentar registrar usuario docente y validar respuesta")

	process.EndStep()


	// [STEP 3] Get teacher data and validate
	process.StartStep("Obtener datos del docente y validar")

	process.EndStep()


	// [STEP 4] Validate teacher role
	process.StartStep("Validar rol de Docente")

	process.EndStep()
	process.End()

	_ = email
	_ = password
}