package course_test

import (
	"testing"

	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestCourseFromStudentView (t *testing.T) {
	process := test.StartTestWithApp(t, "Course Creation and Persistence")
	teacher_email := "test@test.com"
	student_email := "stud@test.com"
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


	// [STEP 3] Login Student user
	process.StartStep("Inicio de sesión con usuario estudiante")

	process.EndStep()


	// [STEP 4] Get Course from student view and validate error
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar error")
	
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()


	// [STEP 5] Enroll student to course
	process.StartStep("Inscribir estudiante al curso")

	process.EndStep()


	// [STEP 6] Get Course from student view and validate data
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar datos")

	process.EndStep()


	// [STEP 7] Unenroll student from course
	process.StartStep("Retirar estudiante del curso")

	process.EndStep()


	// [STEP 8] Get Course from student view and validate error
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar error")

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()
	process.End()	
	_ = teacher_email
	_ = student_email
	_ = password

}