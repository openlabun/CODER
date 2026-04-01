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


func TestCoursesCRUD(t *testing.T) {
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
	ctx := utils.TeacherCourseCtx(teacherAccess)
	t.Logf("[OK] Login profesor exitoso. teacherID=%s", teacherAccess.UserData.ID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-CRUD-%d", now.UnixNano())
	courseCode := fmt.Sprintf("UC-CRUD-%d", now.Unix()%100000)

	t.Logf("[STEP 3] Creando curso con code=%s enrollmentCode=%s", courseCode, enrollmentCode)
	// Create a new course with valid teacher credentials
	createdCourse, err := h.CourseModule.CreateCourse.Execute(ctx, course_dtos.CreateCourseInput{
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
		_ = h.CourseModule.DeleteCourse.Execute(ctx, course_dtos.DeleteCourseInput{CourseID: courseID})
	}()

	t.Log("[STEP 4] Actualizando curso")
	// Update course details with valid teacher credentials
	updatedName := "UC Course CRUD Updated"
	updatedDescription := "course updated by use case test"
	updatedCourse, err := h.CourseModule.UpdateCourse.Execute(ctx, course_dtos.UpdateCourseInput{
		ID:          courseID,
		Name:        &updatedName,
		Description: &updatedDescription,
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
	reloadedCourse, err := h.CourseModule.GetCourseDetails.Execute(ctx, course_dtos.GetCourseDetailsInput{CourseID: courseID})
	if err != nil {
		t.Fatalf("get course details failed: %v", err)
	}
	if reloadedCourse == nil || reloadedCourse.ID != courseID {
		t.Fatal("expected course details for created course")
	}
	if reloadedCourse.Name != updatedName {
		t.Fatalf("expected updated name %q, got %q", updatedName, reloadedCourse.Name)
	}
	t.Logf("[OK] Detalles validados. name=%q", reloadedCourse.Name)

	t.Logf("[STEP 6] Eliminando curso courseID=%s", courseID)
	// Delete course with valid teacher credentials
	if err := h.CourseModule.DeleteCourse.Execute(ctx, course_dtos.DeleteCourseInput{CourseID: courseID}); err != nil {
		t.Fatalf("delete course failed: %v", err)
	}
	courseID = ""
	t.Log("[OK] Curso eliminado")
}
