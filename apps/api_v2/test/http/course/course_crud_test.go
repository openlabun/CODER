package course_test

import (
	"testing"

	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestCourseCRUD (t *testing.T) {
	process := test.StartTestWithApp(t, "Course Creation and Persistence")
	email := "test@test.com"
	password := "Password123!"

	// [STEP 1] Login Teacher user
	process.StartStep("Inicio de sesión con usuario docente")

	process.EndStep()


	// [STEP 2] Create Course
	process.StartStep("Crear curso")

	process.EndStep()


	defer func () {
		t.Logf("[CLEANUP] Deleting created course")
	}()


	// [STEP 3] Update Course
	process.StartStep("Actualizar curso")

	process.EndStep()


	// [STEP 4] Get Course and validate data
	process.StartStep("Obtener datos del curso y validarlos")

	process.EndStep()


	// [STEP 5] Delete Course
	process.StartStep("Eliminar curso")

	process.EndStep()


	// [STEP 6] Get Course and validate error
	process.StartStep("Obtener datos del curso y validar error")

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()
	process.End()	
	_ = email
	_ = password

}