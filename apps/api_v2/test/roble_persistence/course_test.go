package roble_persistence_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	test "github.com/openlabun/CODER/apps/api_v2/test"

	hasher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_course_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
)

func TestCourseCRUD (t *testing.T) {
	process := test.StartTest(t, "Course Creation and Persistence")
	email := "test@test.com"
	password := "Password123!"

	// [STEP 1] Initialize Roble client and repositories
	process.StartStep("Inicializando cliente Roble y repositorios")
	// Initialize Roble client and repositories
	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		process.Fail("initialize roble client", err)
	}
	process.Log("Cliente Roble inicializado")

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)
	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	courseRepository := roble_course_infrastructure.NewCourseRepository(robleAdapter)
	process.Log("Adapter y repositories inicializados")
	process.EndStep()

	// [STEP 2] Initialize hasher and hash password
	process.StartStep("Inicializar hasher y hashear password")
	adapter := hasher.NewSecurityAdapter()
	hashedPassword, err := adapter.Hash(password)
	if err != nil {
		process.Fail("hash password", err)
	}
	process.EndStep()

	// [STEP 3] Login Teacher user
	process.StartStep("Login docente")
	access, err := authAdapter.LoginUser(email, hashedPassword)
	if err != nil {
		process.Fail("login user", err)
	}
	if access == nil || access.UserData == nil || access.UserData.ID == "" {
		process.Fail("login user", fmt.Errorf("expected logged user data with valid ID"))
	}
	if access.Token == nil || access.Token.AccessToken == "" {
		process.Fail("login user", fmt.Errorf("expected access token in login response"))
	}
	ctx := services.WithAccessToken(context.Background(), access.Token.AccessToken)
	process.Log(fmt.Sprintf("teacherID=%s", access.UserData.ID))
	process.EndStep()

	// [STEP 4] Create a new course
	teacherID := access.UserData.ID
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-%d", now.UnixNano())

	process.StartStep("Creando curso con factory")
	process.Log(fmt.Sprintf("enrollmentCode=%s", enrollmentCode))

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
	process.Log(fmt.Sprintf("createdCourseID=%s", createdCourse.ID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando curso temporal %s", courseID)
		_ = courseRepository.DeleteCourse(ctx, courseID)
	}()

	// [STEP 5] Get all courses for the teacher and verify the new course is present
	process.StartStep("Listando cursos del docente")
	teacherCourses, err := courseRepository.GetCoursesByTeacherID(ctx, teacherID)
	if err != nil {
		process.Fail("courses list", err)
	}
	process.Log(fmt.Sprintf("Se recuperaron %d cursos del docente", len(teacherCourses)))

	found := false
	for _, c := range teacherCourses {
		if c != nil && c.ID == courseID {
			found = true
			break
		}
	}
	if !found {
		process.Fail("courses list", fmt.Errorf("expected created course %s in teacher course list", courseID))
	}
	process.Log(fmt.Sprintf("Validacion lista docente: curso %s encontrado", courseID))
	process.EndStep()

	// [STEP 6] Update course details and verify the changes
	process.StartStep("Actualizando nombre y descripcion del curso")
	createdCourse.Name = "Integration Course Updated"
	createdCourse.Description = "Updated by integration test"
	updatedCourse, err := courseRepository.UpdateCourse(ctx, createdCourse)
	if err != nil {
		process.Fail("update course", err)
	}
	if updatedCourse == nil || updatedCourse.Name != "Integration Course Updated" {
		process.Fail("update course", fmt.Errorf("expected updated course name"))
	}
	process.Log(fmt.Sprintf("Update aplicado: Nuevo nombre=%q", updatedCourse.Name))
	process.EndStep()

	// [STEP 7] Get course by ID and verify details
	process.StartStep("Recargando curso por ID")
	reloadedCourse, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		process.Fail("get course by id", err)
	}
	if reloadedCourse == nil {
		process.Fail("get course by id", fmt.Errorf("expected reloaded course"))
	}
	if reloadedCourse.Name != "Integration Course Updated" {
		process.Fail("get course by id", fmt.Errorf("expected persisted updated name, got %q", reloadedCourse.Name))
	}
	process.Log(fmt.Sprintf("Persistencia validada: Nombre recargado=%q", reloadedCourse.Name))

	// [STEP 8] Delete course
	process.StartStep("Eliminando curso")
	if err := courseRepository.DeleteCourse(ctx, courseID); err != nil {
		process.Fail("delete course", err)
	}
	process.EndStep()

	// [STEP 9] Verify course deletion
	process.StartStep("Verificando que el curso ya no existe")
	deletedCourse, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		process.Fail("get course after delete", err)
	}
	if deletedCourse != nil {
		process.Fail("get course after delete", fmt.Errorf("expected course %s to be deleted", courseID))
	}
	process.EndStep()

	process.End()
}
