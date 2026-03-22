package usecases_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	user_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	course_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	exam_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	submission_repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	security_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
)

func TestExamCRUD(t *testing.T) {
	t.Log("[STEP 1] Inicializando application container con dependencias")
	app, err := buildExamApplication()
	if err != nil {
		t.Fatalf("failed to build application: %v", err)
	}
	t.Log("[OK] Application container inicializado")

	t.Log("[STEP 2] Login de profesor")
	teacherAccess := mustLoginExamTeacher(t, app)
	teacherCtx := teacherExamCtx(teacherAccess)
	teacherID := teacherAccess.UserData.ID
	t.Logf("[OK] Login profesor exitoso. teacherID=%s", teacherID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-EXAM-%d", now.UnixNano())
	courseCode := fmt.Sprintf("UC-EXAM-%d", now.Unix()%100000)

	t.Logf("[STEP 3] Creando curso asociado para examenes code=%s enrollmentCode=%s", courseCode, enrollmentCode)
	createdCourse, err := app.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "UC Exam Course",
		Description:    "Course for exam CRUD use case test",
		Visibility:     string(course_entities.CourseVisibilityPublic),
		VisualIdentity: string(course_entities.CourseColourBlue),
		Code:           courseCode,
		Year:           now.Year(),
		Semester:       string(course_entities.AcademicFirstPeriod),
		EnrollmentCode: enrollmentCode,
		EnrollmentURL:  "https://example.test/enroll/" + enrollmentCode,
		TeacherID:      teacherID,
	})
	if err != nil {
		t.Fatalf("create course failed: %v", err)
	}
	if createdCourse == nil || createdCourse.ID == "" {
		t.Fatal("expected created course with ID")
	}
	courseID := createdCourse.ID
	t.Logf("[OK] Curso creado. courseID=%s", courseID)

	var examID1 string
	var examID2 string
	defer func() {
		if examID1 != "" {
			t.Logf("[CLEANUP] Eliminando examen 1 pendiente %s", examID1)
			_, _ = app.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID1})
		}
		if examID2 != "" {
			t.Logf("[CLEANUP] Eliminando examen 2 pendiente %s", examID2)
			_, _ = app.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID2})
		}
		if courseID != "" {
			t.Logf("[CLEANUP] Eliminando curso %s", courseID)
			_ = app.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
		}
	}()

	startTime := time.Now().UTC().Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)

	t.Log("[STEP 4] Creando primer examen con credenciales validas de profesor")
	createdExam1, err := app.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             courseID,
		Title:                "UC Exam 1",
		Description:          "First exam created from use case",
		Visibility:           string(exam_entities.VisibilityCourse),
		StartTime:            startTime.Format(time.RFC3339),
		EndTime:              &endTimeStr,
		AllowLateSubmissions: false,
		TimeLimit:            5400,
		TryLimit:             2,
		ProfessorID:          teacherID,
	})
	if err != nil {
		t.Fatalf("create exam 1 failed: %v", err)
	}
	if createdExam1 == nil || createdExam1.ID == "" {
		t.Fatal("expected created exam 1 with ID")
	}
	examID1 = createdExam1.ID
	t.Logf("[OK] Examen 1 creado. examID=%s", examID1)

	t.Log("[STEP 5] Actualizando examen 1 con nuevos valores")
	updatedTitle := "UC Exam 1 Updated"
	updatedDescription := "First exam updated from use case"
	updatedTryLimit := 3
	updatedExam1, err := app.ExamModule.UpdateExam.Execute(teacherCtx, exam_dtos.UpdateExamInput{
		ExamID:      examID1,
		Title:       &updatedTitle,
		Description: &updatedDescription,
		TryLimit:    &updatedTryLimit,
	})
	if err != nil {
		t.Fatalf("update exam 1 failed: %v", err)
	}
	if updatedExam1 == nil || updatedExam1.Title != updatedTitle {
		t.Fatal("expected updated exam 1 title")
	}
	t.Logf("[OK] Examen 1 actualizado. title=%q tryLimit=%d", updatedExam1.Title, updatedExam1.TryLimit)

	t.Logf("[STEP 6] Consultando detalle de examen 1 examID=%s", examID1)
	reloadedExam1, err := app.ExamModule.GetExamDetails.Execute(teacherCtx, exam_dtos.GetExamDetailsInput{ExamID: examID1})
	if err != nil {
		t.Fatalf("get exam details failed: %v", err)
	}
	if reloadedExam1 == nil {
		t.Fatal("expected exam details for exam 1")
	}
	if reloadedExam1.Title != updatedTitle || reloadedExam1.Description != updatedDescription || reloadedExam1.TryLimit != updatedTryLimit {
		t.Fatalf("expected updated exam values, got title=%q description=%q tryLimit=%d", reloadedExam1.Title, reloadedExam1.Description, reloadedExam1.TryLimit)
	}
	t.Logf("[OK] Detalle de examen 1 validado. title=%q", reloadedExam1.Title)

	t.Log("[STEP 7] Creando segundo examen")
	createdExam2, err := app.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             courseID,
		Title:                "UC Exam 2",
		Description:          "Second exam created from use case",
		Visibility:           string(exam_entities.VisibilityCourse),
		StartTime:            startTime.Add(24 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            3600,
		TryLimit:             1,
		ProfessorID:          teacherID,
	})
	if err != nil {
		t.Fatalf("create exam 2 failed: %v", err)
	}
	if createdExam2 == nil || createdExam2.ID == "" {
		t.Fatal("expected created exam 2 with ID")
	}
	examID2 = createdExam2.ID
	t.Logf("[OK] Examen 2 creado. examID=%s", examID2)

	t.Logf("[STEP 8] Listando examenes por curso courseID=%s y validando que esten ambos", courseID)
	examsByCourse, err := app.ExamModule.GetExamsByCourse.Execute(teacherCtx, exam_dtos.GetExamsByCourseInput{CourseID: courseID})
	if err != nil {
		t.Fatalf("get exams by course failed: %v", err)
	}
	if len(examsByCourse) < 2 {
		t.Fatalf("expected at least 2 exams for course, got %d", len(examsByCourse))
	}

	foundExam1 := false
	foundExam2 := false
	for _, exam := range examsByCourse {
		if exam == nil {
			continue
		}
		if exam.ID == examID1 {
			foundExam1 = true
		}
		if exam.ID == examID2 {
			foundExam2 = true
		}
	}
	if !foundExam1 || !foundExam2 {
		t.Fatal("expected both exams in course exam list")
	}
	t.Logf("[OK] Listado validado. totalExams=%d", len(examsByCourse))

	t.Logf("[STEP 9] Eliminando primer examen examID=%s", examID1)
	deletedExam1, err := app.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID1})
	if err != nil {
		t.Fatalf("delete exam 1 failed: %v", err)
	}
	if deletedExam1 == nil || deletedExam1.ID != examID1 {
		t.Fatal("expected deleted exam 1 payload")
	}
	examID1 = ""
	t.Log("[OK] Examen 1 eliminado")

	t.Log("[STEP 10] Listando examenes por curso nuevamente y validando que solo quede el segundo")
	examsAfterDelete, err := app.ExamModule.GetExamsByCourse.Execute(teacherCtx, exam_dtos.GetExamsByCourseInput{CourseID: courseID})
	if err != nil {
		t.Fatalf("get exams by course after delete failed: %v", err)
	}

	stillHasExam1 := false
	stillHasExam2 := false
	for _, exam := range examsAfterDelete {
		if exam == nil {
			continue
		}
		if exam.ID == createdExam1.ID {
			stillHasExam1 = true
		}
		if exam.ID == examID2 {
			stillHasExam2 = true
		}
	}
	if stillHasExam1 {
		t.Fatal("expected exam 1 to be absent after delete")
	}
	if !stillHasExam2 {
		t.Fatal("expected exam 2 to remain after deleting exam 1")
	}
	t.Logf("[OK] Listado post-borrado validado. totalExams=%d", len(examsAfterDelete))

	t.Logf("[STEP 11] Eliminando segundo examen examID=%s", examID2)
	deletedExam2, err := app.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID2})
	if err != nil {
		t.Fatalf("delete exam 2 failed: %v", err)
	}
	if deletedExam2 == nil || deletedExam2.ID != examID2 {
		t.Fatal("expected deleted exam 2 payload")
	}
	examID2 = ""
	t.Log("[OK] Examen 2 eliminado")
}


func buildExamApplication() (*container.Application, error) {
	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		return nil, fmt.Errorf("initialize roble client: %w", err)
	}

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)
	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	passwordHasher := security_infrastructure.NewSecurityAdapter()

	courseRepository := course_repository.NewCourseRepository(robleAdapter)
	examRepository := exam_repository.NewExamRepository(robleAdapter)
	challengeRepository := exam_repository.NewChallengeRepository(robleAdapter)
	testCaseRepository := exam_repository.NewTestCaseRepository(robleAdapter)

	subRepo := submission_repository.NewSubmissionRepository(robleAdapter)
	sessionRepo := submission_repository.NewSessionRepository(robleAdapter)
	resultRepo := submission_repository.NewSubmissionResultRepository(robleAdapter)

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
		subRepo,
		sessionRepo,
		resultRepo,
	)

	app, err := container.NewApplication(deps)
	if err != nil {
		return nil, fmt.Errorf("initialize application container: %w", err)
	}

	return app, nil
}

func mustLoginExamTeacher(t *testing.T, app *container.Application) *user_dtos.UserAccess {
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

func teacherExamCtx(teacherAccess *user_dtos.UserAccess) context.Context {
	ctx := roble_infrastructure.WithAccessToken(context.Background(), teacherAccess.Token.AccessToken)
	ctx = services.WithUserEmail(ctx, teacherAccess.UserData.Email)
	return ctx
}