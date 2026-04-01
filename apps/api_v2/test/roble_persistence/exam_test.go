package roble_persistence_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	course_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	exam_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_course_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	roble_exam_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
)

func TestExamCRUD(t *testing.T) {
	t.Log("[STEP 1] Inicializando cliente Roble y repositories")

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
	examRepository := roble_exam_infrastructure.NewExamRepository(robleAdapter)
	ioVariableRepository := roble_exam_infrastructure.NewIOVariableRepository(robleAdapter)
	challengeRepository := roble_exam_infrastructure.NewChallengeRepository(robleAdapter, ioVariableRepository)
	testCaseRepository := roble_exam_infrastructure.NewTestCaseRepository(robleAdapter, ioVariableRepository)
	t.Log("[OK] Adapter y repositories inicializados")

	t.Log("[STEP 2] Login con usuario docente de pruebas")
	email := "test@test.com"
	password := "Testing123!"

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
	teacherID := access.UserData.ID
	t.Logf("[OK] Login exitoso. teacherID=%s", teacherID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-EXAM-%d", now.UnixNano())
	courseCode := fmt.Sprintf("EX-%d", now.Unix()%100000)

	t.Log("[STEP 3] Creando curso auxiliar para asociar el examen")
	course, err := course_factory.NewCourse(
		"Exam Integration Course",
		"Course for exam integration test",
		course_entities.CourseVisibilityPublic,
		course_entities.CourseColourBlue,
		courseCode,
		&course_entities.Period{Year: now.Year(), Semester: course_entities.AcademicFirstPeriod},
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
	t.Logf("[OK] Curso auxiliar creado. courseID=%s", courseID)

	defer func() {
		t.Logf("[CLEANUP] Eliminando curso auxiliar %s", courseID)
		_ = courseRepository.DeleteCourse(ctx, courseID)
	}()

	t.Log("[STEP 4] Crear examen (CRUD Exam)")
	startTime := time.Now().UTC()
	endTime := startTime.Add(90 * time.Minute)

	exam, err := exam_factory.NewExam(
		"Integration Exam",
		"Exam created by integration test",
		exam_entities.VisibilityCourse,
		startTime,
		&endTime,
		false,
		5400,
		2,
		teacherID,
		courseID,
	)
	if err != nil {
		t.Fatalf("build exam with factory failed: %v", err)
	}

	createdExam, err := examRepository.CreateExam(ctx, exam)
	if err != nil {
		t.Fatalf("create exam failed: %v", err)
	}
	if createdExam == nil || createdExam.ID == "" {
		t.Fatal("expected created exam with ID")
	}
	examID := createdExam.ID
	t.Logf("[OK] Examen creado. examID=%s", examID)

	defer func() {
		t.Logf("[CLEANUP] Eliminando examen %s", examID)
		_ = examRepository.DeleteExam(ctx, examID)
	}()

	t.Logf("[STEP 5] Leer examen por ID=%s", examID)
	reloadedExam, err := examRepository.GetExamByID(ctx, examID)
	if err != nil {
		t.Fatalf("get exam by id failed: %v", err)
	}
	if reloadedExam == nil {
		t.Fatal("expected reloaded exam")
	}
	t.Logf("[OK] Examen recargado. title=%q", reloadedExam.Title)

	t.Log("[STEP 6] Actualizar examen")
	createdExam.Title = "Integration Exam Updated"
	createdExam.Description = "Updated exam description"
	createdExam.TryLimit = 3
	updatedExam, err := examRepository.UpdateExam(ctx, createdExam)
	if err != nil {
		t.Fatalf("update exam failed: %v", err)
	}
	if updatedExam == nil || updatedExam.Title != "Integration Exam Updated" {
		t.Fatal("expected updated exam title")
	}
	t.Logf("[OK] Examen actualizado. title=%q", updatedExam.Title)

	t.Log("[STEP 7] Validar listados de examenes por courseID y teacherID")
	courseExams, err := examRepository.GetExamsByCourseID(ctx, courseID)
	if err != nil {
		t.Fatalf("get exams by course id failed: %v", err)
	}
	teacherExams, err := examRepository.GetExamsByTeacherID(ctx, teacherID)
	if err != nil {
		t.Fatalf("get exams by teacher id failed: %v", err)
	}
	if len(courseExams) == 0 {
		t.Fatal("expected at least one exam for course")
	}
	if len(teacherExams) == 0 {
		t.Fatal("expected at least one exam for teacher")
	}
	t.Logf("[OK] Listados validados. courseExams=%d teacherExams=%d", len(courseExams), len(teacherExams))

	t.Log("[STEP 8] Crear challenge con IOVariables (CRUD Challenge)")
	inputA, err := exam_factory.NewIOVariable("a", exam_entities.VariableFormatInt, "2")
	if err != nil {
		t.Fatalf("create input io variable failed: %v", err)
	}
	inputB, err := exam_factory.NewIOVariable("b", exam_entities.VariableFormatInt, "3")
	if err != nil {
		t.Fatalf("create input io variable failed: %v", err)
	}
	output, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "5")
	if err != nil {
		t.Fatalf("create output io variable failed: %v", err)
	}

	challenge, err := exam_factory.NewChallenge(
		"Sum Challenge",
		"Return the sum of two numbers",
		[]string{"math", "integration"},
		exam_entities.ChallengeStatusDraft,
		exam_entities.ChallengeDifficultyEasy,
		1500,
		256,
		[]exam_entities.IOVariable{*inputA, *inputB},
		*output,
		"1 <= a,b <= 1000",
		examID,
	)
	if err != nil {
		t.Fatalf("build challenge with factory failed: %v", err)
	}

	createdChallenge, err := challengeRepository.CreateChallenge(ctx, challenge)
	if err != nil {
		t.Fatalf("create challenge failed: %v", err)
	}
	if createdChallenge == nil || createdChallenge.ID == "" {
		t.Fatal("expected created challenge with ID")
	}
	challengeID := createdChallenge.ID
	t.Logf("[OK] Challenge creado. challengeID=%s", challengeID)

	defer func() {
		t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
		_ = challengeRepository.DeleteChallenge(ctx, challengeID)
	}()

	t.Logf("[STEP 9] Leer challenge por ID=%s e hidratar IOVariables", challengeID)
	reloadedChallenge, err := challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		t.Fatalf("get challenge by id failed: %v", err)
	}
	if reloadedChallenge == nil {
		t.Fatal("expected reloaded challenge")
	}
	if len(reloadedChallenge.InputVariables) == 0 || reloadedChallenge.OutputVariable.ID == "" {
		t.Fatal("expected challenge with hydrated input/output variables")
	}
	t.Logf("[OK] Challenge recargado con IO. inputCount=%d outputID=%s", len(reloadedChallenge.InputVariables), reloadedChallenge.OutputVariable.ID)

	t.Log("[STEP 10] Actualizar challenge")
	createdChallenge.Title = "Sum Challenge Updated"
	createdChallenge.Status = exam_entities.ChallengeStatusPublished
	createdChallenge.Tags = []string{"math", "updated"}
	updatedChallenge, err := challengeRepository.UpdateChallenge(ctx, createdChallenge)
	if err != nil {
		t.Fatalf("update challenge failed: %v", err)
	}
	if updatedChallenge == nil || updatedChallenge.Title != "Sum Challenge Updated" {
		t.Fatal("expected updated challenge title")
	}
	t.Logf("[OK] Challenge actualizado. title=%q status=%s", updatedChallenge.Title, updatedChallenge.Status)

	t.Log("[STEP 11] Validar listados y lecturas de IO de challenge")
	examChallenges, err := challengeRepository.GetChallengesByExamID(ctx, examID)
	if err != nil {
		t.Fatalf("get challenges by exam failed: %v", err)
	}
	if len(examChallenges) == 0 {
		t.Fatal("expected at least one challenge for exam")
	}
	challengeInputs, err := challengeRepository.GetInputVariablesByChallengeID(ctx, challengeID)
	if err != nil {
		t.Fatalf("get challenge input variables failed: %v", err)
	}
	challengeOutputs, err := challengeRepository.GetOutputVariablesByChallengeID(ctx, challengeID)
	if err != nil {
		t.Fatalf("get challenge output variables failed: %v", err)
	}
	if len(challengeInputs) == 0 || len(challengeOutputs) == 0 {
		t.Fatal("expected challenge input/output variable relations")
	}
	t.Logf("[OK] Relaciones de IO challenge validadas. inputs=%d outputs=%d", len(challengeInputs), len(challengeOutputs))

	t.Log("[STEP 12] Crear test case con IOVariables (CRUD TestCase)")
	tcInputA, err := exam_factory.NewIOVariable("a", exam_entities.VariableFormatInt, "10")
	if err != nil {
		t.Fatalf("create test case input failed: %v", err)
	}
	tcInputB, err := exam_factory.NewIOVariable("b", exam_entities.VariableFormatInt, "15")
	if err != nil {
		t.Fatalf("create test case input failed: %v", err)
	}
	tcOutput, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "25")
	if err != nil {
		t.Fatalf("create test case output failed: %v", err)
	}

	testCase, err := exam_factory.NewTestCase(
		"sample_1",
		[]exam_entities.IOVariable{*tcInputA, *tcInputB},
		*tcOutput,
		true,
		10,
		challengeID,
	)
	if err != nil {
		t.Fatalf("build test case with factory failed: %v", err)
	}

	createdTestCase, err := testCaseRepository.CreateTestCase(ctx, testCase)
	if err != nil {
		t.Fatalf("create test case failed: %v", err)
	}
	if createdTestCase == nil || createdTestCase.ID == "" {
		t.Fatal("expected created test case with ID")
	}
	testCaseID := createdTestCase.ID
	t.Logf("[OK] TestCase creado. testCaseID=%s", testCaseID)

	defer func() {
		t.Logf("[CLEANUP] Eliminando test case %s", testCaseID)
		_ = testCaseRepository.DeleteTestCase(ctx, testCaseID)
	}()

	t.Logf("[STEP 13] Leer test case por ID=%s e hidratar IOVariables", testCaseID)
	reloadedTestCase, err := testCaseRepository.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		t.Fatalf("get test case by id failed: %v", err)
	}
	if reloadedTestCase == nil {
		t.Fatal("expected reloaded test case")
	}
	if len(reloadedTestCase.Input) == 0 || reloadedTestCase.ExpectedOutput.ID == "" {
		t.Fatal("expected test case with hydrated input/output variables")
	}
	t.Logf("[OK] TestCase recargado con IO. inputCount=%d outputID=%s", len(reloadedTestCase.Input), reloadedTestCase.ExpectedOutput.ID)

	t.Log("[STEP 14] Actualizar test case")
	createdTestCase.Name = "sample_1_updated"
	createdTestCase.Points = 25
	updatedTestCase, err := testCaseRepository.UpdateTestCase(ctx, createdTestCase)
	if err != nil {
		t.Fatalf("update test case failed: %v", err)
	}
	if updatedTestCase == nil || updatedTestCase.Name != "sample_1_updated" {
		t.Fatal("expected updated test case name")
	}
	t.Logf("[OK] TestCase actualizado. name=%q points=%d", updatedTestCase.Name, updatedTestCase.Points)

	t.Log("[STEP 15] Validar listados y lecturas de IO de test case")
	challengeTestCases, err := testCaseRepository.GetTestCasesByChallengeID(ctx, challengeID)
	if err != nil {
		t.Fatalf("get test cases by challenge id failed: %v", err)
	}
	if len(challengeTestCases) == 0 {
		t.Fatal("expected at least one test case for challenge")
	}
	testCaseInputs, err := testCaseRepository.GetInputVariablesByTestCaseID(ctx, testCaseID)
	if err != nil {
		t.Fatalf("get test case input variables failed: %v", err)
	}
	testCaseOutputs, err := testCaseRepository.GetOutputVariablesByTestCaseID(ctx, testCaseID)
	if err != nil {
		t.Fatalf("get test case output variables failed: %v", err)
	}
	if len(testCaseInputs) == 0 || len(testCaseOutputs) == 0 {
		t.Fatal("expected test case input/output variable relations")
	}
	t.Logf("[OK] Relaciones de IO test case validadas. inputs=%d outputs=%d", len(testCaseInputs), len(testCaseOutputs))

	t.Log("[STEP 16] Eliminar test case y validar borrado")
	if err := testCaseRepository.DeleteTestCase(ctx, testCaseID); err != nil {
		t.Fatalf("delete test case failed: %v", err)
	}
	deletedTestCase, err := testCaseRepository.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		t.Fatalf("get test case after delete failed: %v", err)
	}
	if deletedTestCase != nil {
		t.Fatalf("expected test case %s to be deleted", testCaseID)
	}
	testCaseID = ""
	t.Log("[OK] TestCase eliminado")

	t.Log("[STEP 17] Eliminar challenge y validar borrado")
	if err := challengeRepository.DeleteChallenge(ctx, challengeID); err != nil {
		t.Fatalf("delete challenge failed: %v", err)
	}
	deletedChallenge, err := challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		t.Fatalf("get challenge after delete failed: %v", err)
	}
	if deletedChallenge != nil {
		t.Fatalf("expected challenge %s to be deleted", challengeID)
	}
	challengeID = ""
	t.Log("[OK] Challenge eliminado")

	t.Log("[STEP 18] Eliminar examen y validar borrado")
	if err := examRepository.DeleteExam(ctx, examID); err != nil {
		t.Fatalf("delete exam failed: %v", err)
	}
	deletedExam, err := examRepository.GetExamByID(ctx, examID)
	if err != nil {
		t.Fatalf("get exam after delete failed: %v", err)
	}
	if deletedExam != nil {
		t.Fatalf("expected exam %s to be deleted", examID)
	}
	examID = ""
	t.Log("[OK] Examen eliminado")
}
