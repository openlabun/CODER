package usecases_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	user_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	user_repo "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	security_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	rabbitmq_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/publisher/rabbitMQ"

	course_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	exam_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	submission_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
)


func TestCoursesCRUD(t *testing.T) {
	t.Log("[STEP 1] Inicializando application container")
	// Initialize application
	h, err := buildApplication()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application container inicializado")

	t.Log("[STEP 2] Login de profesor")
	// Login Teacher user
	teacherAccess := mustLoginTeacher(t, h)
	ctx := teacherCourseCtx(teacherAccess)
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



func buildApplication() (*container.Application, error) {
	// verbose logs are intentionally omitted here to keep the helper reusable in tests.
	// Start clients
	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		return nil, fmt.Errorf("initialize roble client: %w", err)
	}

	// Start adapters and repositories
	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)
	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	passwordHasher := security_infrastructure.NewSecurityAdapter()

	courseRepository := course_repository.NewCourseRepository(robleAdapter)
	examRepository := exam_repository.NewExamRepository(robleAdapter)
	challengeRepository := exam_repository.NewChallengeRepository(robleAdapter)
	testCaseRepository := exam_repository.NewTestCaseRepository(robleAdapter)
	
	submissionRepository := submission_repository.NewSubmissionRepository(robleAdapter)
	sessionRepository := submission_repository.NewSessionRepository(robleAdapter)
	submissionResRepository := submission_repository.NewSubmissionResultRepository(robleAdapter)
	publisherPort, err := rabbitmq_infrastructure.NewRabbitMQAdapter()
	if err != nil {
		return nil, fmt.Errorf("initialize publisher adapter: %w", err)
	}


	deps := container.NewApplicationDependencies(
		authAdapter,
		authAdapter,
		userRepository,
		authAdapter,

		passwordHasher,

		userRepository,
		courseRepository,

		examRepository,
		challengeRepository,
		testCaseRepository,

		submissionRepository,
		sessionRepository,
		submissionResRepository,
		publisherPort,
	)

	appContainer, err := container.NewApplication(deps)
	if err != nil {
		return nil, fmt.Errorf("initialize application container: %w", err)
	}

	return appContainer, nil
}

func mustLoginTeacher(t *testing.T, app *container.Application) *user_dtos.UserAccess {
	t.Helper()
	t.Log("[AUTH] Intentando login de profesor")

	email := "test@test.com"
	password := "Testing123!"

	access, err := app.Dependencies.LoginService.LoginUser(email, password)
	if err != nil {
		t.Fatalf("teacher login failed: %v", err)
	}
	if access == nil || access.UserData == nil || access.UserData.ID == "" {
		t.Fatal("expected teacher user data with ID")
	}
	if access.Token == nil || access.Token.AccessToken == "" {
		t.Fatal("expected teacher access token")
	}
	t.Logf("[AUTH][OK] Login profesor completado. teacherID=%s", access.UserData.ID)

	return access
}

func ensureStudentAccess(t *testing.T, auth *user_repo.RobleAuthAdapter) *user_dtos.UserAccess {
	t.Helper()
	t.Log("[AUTH] Intentando login de estudiante")

	email := "stud@test.com"
	password := "Testing123!"

	access, err := auth.LoginUser(email, password)
	if err == nil && access != nil && access.UserData != nil && access.UserData.ID != "" && access.Token != nil && access.Token.AccessToken != "" {
		t.Logf("[AUTH][OK] Login estudiante exitoso. studentID=%s", access.UserData.ID)
		return access
	}
	t.Log("[AUTH] Estudiante no existe o login fallo, registrando usuario")

	registered, registerErr := auth.RegisterUserDirect(email, password, "Student Test")
	if registerErr != nil {
		t.Fatalf("register student failed: %v", registerErr)
	}
	if registered == nil || registered.UserData == nil || registered.UserData.ID == "" || registered.Token == nil || registered.Token.AccessToken == "" {
		t.Fatal("expected registered student with ID and access token")
	}
	t.Logf("[AUTH][OK] Registro de estudiante exitoso. studentID=%s", registered.UserData.ID)

	return registered
}

func teacherCourseCtx(teacherAccess *user_dtos.UserAccess) context.Context {
	ctx := services.WithAccessToken(context.Background(), teacherAccess.Token.AccessToken)
	ctx = services.WithUserEmail(ctx, teacherAccess.UserData.Email)
	return ctx
}