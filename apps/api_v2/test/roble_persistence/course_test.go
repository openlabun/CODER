package roble_persistence_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_course_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
)

func TestCourseCreation(t *testing.T) {
	t.Log("[STEP 1] Inicializando cliente Roble y repositorios")

	// Initialize Roble client and repositories
	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		t.Fatalf("initialize roble client: %v", err)
	}
	t.Log("[OK] Cliente Roble inicializado")

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)
	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	courseRepository := roble_course_infrastructure.NewCourseRepository(robleAdapter)
	t.Log("[OK] Adapter y repositories inicializados")

	// Login Teacher user
	email := "test@test.com"
	password := "Testing123!"
	t.Logf("[STEP 2] Login docente con email=%s", email)

	access, err := authAdapter.LoginUser(email, password)
	if err != nil {
		t.Fatalf("teacher login failed: %v", err)
	}
	if access == nil || access.UserData == nil || access.UserData.ID == "" {
		t.Fatal("expected logged user data with valid ID")
	}
	if access.Token == nil || access.Token.AccessToken == "" {
		t.Fatal("expected access token in login response")
	}
	ctx := services.WithAccessToken(context.Background(), access.Token.AccessToken)
	t.Logf("[OK] Login exitoso. teacherID=%s", access.UserData.ID)

	// Create a new course
	teacherID := access.UserData.ID
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-%d", now.UnixNano())
	t.Logf("[STEP 3] Creando curso con factory enrollmentCode=%s", enrollmentCode)

	course, err := factory.NewCourse(
		"Integration Course",
		"Created by integration test",
		course_entities.CourseVisibilityPublic,
		course_entities.CourseColourBlue,
		fmt.Sprintf("IT-%d", now.Unix()%100000),
		&course_entities.Period{
			Year:     now.Year(),
			Semester: course_entities.AcademicFirstPeriod,
		},
		enrollmentCode,
		"https://example.test/enroll/"+enrollmentCode,
		teacherID,
	)
	if err != nil {
		t.Fatalf("build course with factory failed: %v", err)
	}

	createdCourse, err := courseRepository.CreateCourse(ctx, course)
	if err != nil {
		t.Fatalf("create course failed: %v", err)
	}
	if createdCourse == nil || createdCourse.ID == "" {
		t.Fatal("expected created course with ID")
	}
	courseID := createdCourse.ID
	t.Logf("[OK] Curso creado. createdCourseID=%s", createdCourse.ID)

	defer func() {
		t.Logf("[CLEANUP] Eliminando curso temporal %s", courseID)
		_ = courseRepository.DeleteCourse(ctx, courseID)
	}()

	// Get all courses for the teacher and verify the new course is present
	t.Logf("[STEP 4] Listando cursos del docente teacherID=%s", teacherID)
	teacherCourses, err := courseRepository.GetCoursesByTeacherID(ctx, teacherID)
	if err != nil {
		t.Fatalf("get courses by teacher failed: %v", err)
	}
	t.Logf("[OK] Se recuperaron %d cursos del docente", len(teacherCourses))

	found := false
	for _, c := range teacherCourses {
		if c != nil && c.ID == courseID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected created course %s in teacher course list", courseID)
	}
	t.Logf("[OK] Validacion lista docente: curso %s encontrado", courseID)

	// Update course details and verify the changes
	t.Log("[STEP 5] Actualizando nombre y descripcion del curso")
	createdCourse.Name = "Integration Course Updated"
	createdCourse.Description = "Updated by integration test"
	updatedCourse, err := courseRepository.UpdateCourse(ctx, createdCourse)
	if err != nil {
		t.Fatalf("update course failed: %v", err)
	}
	if updatedCourse == nil || updatedCourse.Name != "Integration Course Updated" {
		t.Fatal("expected updated course name")
	}
	t.Logf("[OK] Update aplicado. Nuevo nombre=%q", updatedCourse.Name)

	t.Logf("[STEP 6] Recargando curso por ID=%s para validar persistencia", courseID)
	reloadedCourse, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		t.Fatalf("get course by id failed: %v", err)
	}
	if reloadedCourse == nil {
		t.Fatal("expected reloaded course")
	}
	if reloadedCourse.Name != "Integration Course Updated" {
		t.Fatalf("expected persisted updated name, got %q", reloadedCourse.Name)
	}
	t.Logf("[OK] Persistencia validada. Nombre recargado=%q", reloadedCourse.Name)

	// Delete course
	t.Logf("[STEP 7] Eliminando curso ID=%s", courseID)
	if err := courseRepository.DeleteCourse(ctx, courseID); err != nil {
		t.Fatalf("delete course failed: %v", err)
	}
	t.Log("[OK] Curso eliminado")

	t.Logf("[STEP 8] Verificando que el curso %s ya no existe", courseID)
	deletedCourse, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		t.Fatalf("get course after delete failed: %v", err)
	}
	if deletedCourse != nil {
		t.Fatalf("expected course %s to be deleted", courseID)
	}
	t.Log("[OK] Validacion final: curso eliminado correctamente")
}
