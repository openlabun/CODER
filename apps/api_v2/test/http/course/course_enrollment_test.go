package course_test

import (
	"testing"

	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestCourseEnrollment (t *testing.T) {
	process := test.StartTestWithApp(t, "Course Enrollment")
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

	
	// [STEP 3] Enroll student to course
	process.StartStep("Inscribir estudiante al curso")

	process.EndStep()


	// [STEP 4] Get Course and validate student enrollment
	process.StartStep("Obtener datos del curso y validar inscripción del estudiante")

	process.EndStep()


	// [STEP 5] Unenroll student from course
	process.StartStep("Retirar estudiante del curso")

	process.EndStep()


	// [STEP 6] Get Course and validate student unenrollment
	process.StartStep("Obtener datos del curso y validar desinscripción del estudiante")

	process.EndStep()


	process.End()	
	_ = email
	_ = password

}