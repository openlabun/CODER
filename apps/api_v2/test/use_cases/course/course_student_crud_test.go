package usecases_test

import (
	"fmt"
	"testing"
	"time"

	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"

	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestCoursesWithStudentsCRUD(t *testing.T) {
	t.Log("[STEP 1] Inicializando application container")
	// Initialize application
	h, err := container.BuildApplicationContainer()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application container inicializado")

	t.Log("[STEP 2] Login de profesor")
	// Login Teacher user
	teacherAccess := utils.EnsureTeacherAccess(t, h)
	teacherCtx := utils.TeacherCourseCtx(teacherAccess)
	t.Logf("[OK] Login profesor exitoso. teacherID=%s", teacherAccess.UserData.ID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-STU-%d", now.UnixNano())
	courseCode := fmt.Sprintf("UC-STU-%d", now.Unix()%100000)

	t.Logf("[STEP 3] Creando curso para flujo con estudiantes code=%s", courseCode)
	// Create a new course with valid teacher credentials
	createdCourse, err := h.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "UC Course CRUD",
		Description:    "course created by use case test",
		Visibility:     string(Entities.CourseVisibilityPublic),
		VisualIdentity: string(Entities.CourseColourBlue),
		Code:           courseCode,
		Year:           now.Year(),
		Semester:       string(Entities.AcademicFirstPeriod),
		EnrollmentCode: enrollmentCode,
		EnrollmentURL:  "https://example.test/enroll/" + enrollmentCode,
		TeacherID:      teacherAccess.UserData.ID,
	})
	if err != nil {
		t.Fatalf("create course failed: %v", err)
	}
	if createdCourse == nil || createdCourse.ID == "" {
		t.Fatal("expected created course with ID")
	}
	courseID := createdCourse.ID
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	defer func() {
		t.Logf("[CLEANUP] Eliminando curso %s", courseID)
		_ = h.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
	}()

	t.Log("[STEP 4] Actualizando curso")
	// Update course details with valid teacher credentials
	updatedName := "UC Course Students Updated"
	updatedCourse, err := h.CourseModule.UpdateCourse.Execute(teacherCtx, course_dtos.UpdateCourseInput{
		ID:   courseID,
		Name: &updatedName,
	})
	if err != nil {
		t.Fatalf("update course failed: %v", err)
	}
	if updatedCourse == nil || updatedCourse.Name != updatedName {
		t.Fatal("expected updated course name")
	}
	t.Logf("[OK] Curso actualizado. name=%q", updatedCourse.Name)

	t.Logf("[STEP 5] Consultando detalles del curso courseID=%s", courseID)
	// Get course details with valid teacher credentials
	_, err = h.CourseModule.GetCourseDetails.Execute(teacherCtx, course_dtos.GetCourseDetailsInput{CourseID: courseID})
	if err != nil {
		t.Fatalf("get course details failed: %v", err)
	}
	t.Log("[OK] Detalles del curso recuperados")

	t.Log("[STEP 6] Obtener/crear estudiante y construir contexto")
	// Create student user
	studentAccess := utils.EnsureStudentAccess(t, h)
	studentCtx := utils.BuildContext(studentAccess.Token.AccessToken, studentAccess.UserData.Email)
	t.Logf("[OK] Estudiante listo. studentID=%s", studentAccess.UserData.ID)

	t.Log("[STEP 7] Matricular estudiante en el curso")
	// Enroll a student to the course with valid student credentials
	_, err = h.CourseModule.EnrollInCourse.Execute(studentCtx, course_dtos.EnrolledInCourseInput{
		CourseID:  courseID,
		StudentID: &studentAccess.UserData.ID,
		StudentEmail: nil,
	})
	if err != nil {
		t.Fatalf("enroll student failed: %v", err)
	}
	t.Log("[OK] Estudiante matriculado")

	t.Log("[STEP 8] Validar que el estudiante aparece en la lista del curso")
	students, err := h.CourseModule.GetCourseStudents.Execute(studentCtx, course_dtos.GetCourseStudentsInput{CourseID: courseID})
	if err != nil {
		t.Fatalf("get course students failed: %v", err)
	}
	found := false
	for _, s := range students {
		if s != nil && s.ID == studentAccess.UserData.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected enrolled student in course")
	}
	t.Logf("[OK] Estudiante encontrado en lista. totalStudents=%d", len(students))

	t.Log("[STEP 9] Remover estudiante del curso")
	// Remove a student from the course with valid student credentials
	if err := h.CourseModule.RemoveStudentFromCourse.Execute(studentCtx, course_dtos.RemoveStudentFromCourseInput{
		CourseID:  courseID,
		StudentID: studentAccess.UserData.ID,
	}); err != nil {
		t.Fatalf("remove student failed: %v", err)
	}
	t.Log("[OK] Estudiante removido")

	t.Log("[STEP 10] Validar que el estudiante ya no aparece en la lista")
	afterRemoval, err := h.CourseModule.GetCourseStudents.Execute(studentCtx, course_dtos.GetCourseStudentsInput{CourseID: courseID})
	if err != nil {
		t.Fatalf("get course students after removal failed: %v", err)
	}
	for _, s := range afterRemoval {
		if s != nil && s.ID == studentAccess.UserData.ID {
			t.Fatal("expected removed student to be absent from course")
		}
	}
	t.Logf("[OK] Validacion post-remocion completada. totalStudents=%d", len(afterRemoval))

	t.Logf("[STEP 11] Eliminando curso courseID=%s", courseID)
	// Delete course with valid teacher credentials
	if err := h.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID}); err != nil {
		t.Fatalf("delete course failed: %v", err)
	}
	courseID = ""
	t.Log("[OK] Curso eliminado")
}