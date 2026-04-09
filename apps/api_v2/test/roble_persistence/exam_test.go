package roble_persistence_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	test "github.com/openlabun/CODER/apps/api_v2/test"

	hasher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
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
	process := test.StartTest(t, "Exam Creation and Persistence")
	email := "test@test.com"
	password := "Password123!"

	// [STEP 1] Initialize Roble client and repositories
	process.StartStep("Inicializando cliente Roble y repositories")

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
	examRepository := roble_exam_infrastructure.NewExamRepository(robleAdapter)
	examItemRepository := roble_exam_infrastructure.NewExamItemRepository(robleAdapter)
	ioVariableRepository := roble_exam_infrastructure.NewIOVariableRepository(robleAdapter)
	challengeRepository := roble_exam_infrastructure.NewChallengeRepository(robleAdapter, ioVariableRepository)
	testCaseRepository := roble_exam_infrastructure.NewTestCaseRepository(robleAdapter, ioVariableRepository)
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

	// [STEP 3] Login teacher user
	process.StartStep("Login con usuario docente de pruebas")
	access, err := authAdapter.LoginUser(email, hashedPassword)
	if err != nil {
		process.Fail("teacher login", err)
	}
	if access == nil || access.UserData == nil || access.UserData.ID == "" {
		process.Fail("teacher login", fmt.Errorf("expected logged user data with valid ID"))
	}
	if access.Token == nil || access.Token.AccessToken == "" {
		process.Fail("teacher login", fmt.Errorf("expected access token in login response"))
	}

	ctx := services.WithAccessToken(context.Background(), access.Token.AccessToken)
	teacherID := access.UserData.ID
	process.Log(fmt.Sprintf("Login exitoso. teacherID=%s", teacherID))
	process.EndStep()

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-EXAM-%d", now.UnixNano())
	courseCode := fmt.Sprintf("EX-%d", now.Unix()%100000)

	// [STEP 4] Create an auxiliary course to associate the exam with
	process.StartStep("Creando curso auxiliar para asociar el examen")
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
		process.Fail("build course with factory", err)
	}

	createdCourse, err := courseRepository.CreateCourse(ctx, course)
	if err != nil {
		process.Fail("create course", err)
	}
	if createdCourse == nil || createdCourse.ID == "" {
		process.Fail("create course", fmt.Errorf("expected created course with ID"))
	}
	courseID := createdCourse.ID
	process.Log(fmt.Sprintf("Curso auxiliar creado. courseID=%s", courseID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando curso auxiliar %s", courseID)
		_ = courseRepository.DeleteCourse(ctx, courseID)
	}()

	startTime := time.Now().UTC()
	endTime := startTime.Add(90 * time.Minute)

	// [STEP 5] Create an exam associated with the course
	process.StartStep("Crear examen (CRUD Exam)")

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
		&courseID,
	)
	if err != nil {
		process.Fail("build exam with factory", err)
	}

	createdExam, err := examRepository.CreateExam(ctx, exam)
	if err != nil {
		process.Fail("create exam", err)
	}
	if createdExam == nil || createdExam.ID == "" {
		process.Fail("create exam", fmt.Errorf("expected created exam with ID"))
	}
	examID := createdExam.ID
	process.Log(fmt.Sprintf("Examen creado. examID=%s", examID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando examen %s", examID)
		_ = examRepository.DeleteExam(ctx, examID)
	}()

	// [STEP 6] Get exam by ID and verify details
	process.StartStep(fmt.Sprintf("Leer examen por ID=%s", examID))
	reloadedExam, err := examRepository.GetExamByID(ctx, examID)
	if err != nil {
		process.Fail("get exam by id", err)
	}
	if reloadedExam == nil {
		process.Fail("get exam by id", fmt.Errorf("expected reloaded exam"))
	}
	process.Log(fmt.Sprintf("Examen recargado. title=%q", reloadedExam.Title))
	process.EndStep()

	// [STEP 7] Update exam details and verify changes
	process.StartStep("Actualizar examen")
	createdExam.Title = "Integration Exam Updated"
	createdExam.Description = "Updated exam description"
	createdExam.TryLimit = 3
	updatedExam, err := examRepository.UpdateExam(ctx, createdExam)
	if err != nil {
		process.Fail("update exam", err)
	}
	if updatedExam == nil || updatedExam.Title != "Integration Exam Updated" {
		process.Fail("update exam", fmt.Errorf("expected updated exam title"))
	}
	process.Log(fmt.Sprintf("Examen actualizado. title=%q", updatedExam.Title))
	process.EndStep()

	// [STEP 8] Check exam details after update
	process.StartStep("Validar listados de examenes por courseID y teacherID")
	courseExams, err := examRepository.GetExamsByCourseID(ctx, courseID)
	if err != nil {
		process.Fail("get exams by course id", err)
	}
	teacherExams, err := examRepository.GetExamsByTeacherID(ctx, teacherID)
	if err != nil {
		process.Fail("get exams by teacher id", err)
	}
	if len(courseExams) == 0 {
		process.Fail("get exams by course id", fmt.Errorf("expected at least one exam for course"))
	}
	if len(teacherExams) == 0 {
		process.Fail("get exams by teacher id", fmt.Errorf("expected at least one exam for teacher"))
	}
	process.Log(fmt.Sprintf("Listados validados. courseExams=%d teacherExams=%d", len(courseExams), len(teacherExams)))
	process.EndStep()

	// [STEP 9] Create IOVariables for challenge
	process.StartStep("Crear IOVariables para Challenge")
	inputA, err := exam_factory.NewIOVariable("a", exam_entities.VariableFormatInt, "2")
	if err != nil {
		process.Fail("create input io variable", err)
	}
	inputB, err := exam_factory.NewIOVariable("b", exam_entities.VariableFormatInt, "3")
	if err != nil {
		process.Fail("create input io variable", err)
	}
	output, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "5")
	if err != nil {
		process.Fail("create output io variable", err)
	}
	inputA, err = ioVariableRepository.CreateIOVariable(ctx, inputA)
	if inputA == nil || err != nil {
		process.Fail("persist input io variable", err)
	}
	inputB, err = ioVariableRepository.CreateIOVariable(ctx, inputB)
	if inputB == nil || err != nil {
		process.Fail("persist input io variable", err)
	}
	output, err = ioVariableRepository.CreateIOVariable(ctx, output)
	if output == nil || err != nil {
		process.Fail("persist output io variable", err)
	}

	process.EndStep()

	// [STEP 10] Create Challenge
	process.StartStep("Crear Challenge")
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
		teacherID,
	)
	if err != nil {
		process.Fail("build challenge with factory", err)
	}

	createdChallenge, err := challengeRepository.CreateChallenge(ctx, challenge)
	if err != nil {
		process.Fail("create challenge", err)
	}
	if createdChallenge == nil || createdChallenge.ID == "" {
		process.Fail("create challenge", fmt.Errorf("expected created challenge with ID"))
	}
	challengeID := createdChallenge.ID
	process.Log(fmt.Sprintf("Challenge creado. challengeID=%s", challengeID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
		_ = challengeRepository.DeleteChallenge(ctx, challengeID)

		_ = ioVariableRepository.DeleteIOVariable(ctx, inputA.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, inputB.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, output.ID)
	}()

	// [STEP 11] Get challenge by ID and verify details
	process.StartStep(fmt.Sprintf("Leer challenge por ID=%s e hidratar IOVariables", challengeID))
	reloadedChallenge, err := challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		process.Fail("get challenge by id", err)
	}
	if reloadedChallenge == nil {
		process.Fail("get challenge by id", fmt.Errorf("expected reloaded challenge"))
	}
	if len(reloadedChallenge.InputVariables) == 0 || reloadedChallenge.OutputVariable.ID == "" {
		process.Fail("get challenge by id", fmt.Errorf("expected challenge with hydrated input/output variables"))
	}
	process.Log(fmt.Sprintf("Challenge recargado con IO. inputCount=%d outputID=%s", len(reloadedChallenge.InputVariables), reloadedChallenge.OutputVariable.ID))
	process.EndStep()

	// [STEP 12] Create ExamItem associating the challenge with the exam
	process.StartStep("Crear ExamItem para asociar el Challenge con el Examen")
	examItem, err := exam_factory.NewExamItem(
		challengeID,
		examID,
		1,
		100,
	)
	if err != nil {
		process.Fail("build exam item with factory", err)
	}

	createdExamItem, err := examItemRepository.CreateExamItem(ctx, examItem)
	if err != nil {
		process.Fail("create exam item", err)
	}
	if createdExamItem == nil || createdExamItem.ID == "" {
		process.Fail("create exam item", fmt.Errorf("expected created exam item with ID"))
	}

	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando exam item %s", createdExamItem.ID)
		_ = examItemRepository.DeleteExamItem(ctx, createdExamItem.ID)
	}()

	// [STEP 13] Update challenge details and verify changes
	process.StartStep("Actualizar challenge")
	createdChallenge.Title = "Sum Challenge Updated"
	createdChallenge.Status = exam_entities.ChallengeStatusPublished
	createdChallenge.Tags = []string{"math", "updated"}
	updatedChallenge, err := challengeRepository.UpdateChallenge(ctx, createdChallenge)
	if err != nil {
		process.Fail("update challenge", err)
	}
	if updatedChallenge == nil || updatedChallenge.Title != "Sum Challenge Updated" {
		process.Fail("update challenge", fmt.Errorf("expected updated challenge title"))
	}
	process.Log(fmt.Sprintf("Challenge actualizado. title=%q status=%s", updatedChallenge.Title, updatedChallenge.Status))
	process.EndStep()

	// [STEP 14] Validate challenge listings and IO variable relations
	process.StartStep("Validar listados y lecturas de IO de challenge")
	examChallenges, err := challengeRepository.GetChallengesByExamID(ctx, examID)
	if err != nil {
		process.Fail("get challenges by exam", err)
	}
	if len(examChallenges) == 0 {
		process.Fail("get challenges by exam", fmt.Errorf("expected at least one challenge for exam"))
	}
	challengeInputs, err := challengeRepository.GetInputVariablesByChallengeID(ctx, challengeID)
	if err != nil {
		process.Fail("get challenge input variables", err)
	}
	challengeOutputs, err := challengeRepository.GetOutputVariablesByChallengeID(ctx, challengeID)
	if err != nil {
		process.Fail("get challenge output variables", err)
	}
	if len(challengeInputs) == 0 || len(challengeOutputs) == 0 {
		process.Fail("challenge io relations", fmt.Errorf("expected challenge input/output variable relations"))
	}
	process.Log(fmt.Sprintf("Relaciones de IO challenge validadas. inputs=%d outputs=%d", len(challengeInputs), len(challengeOutputs)))
	process.EndStep()

	// [STEP 15] Create IOVariables for TestCase
	process.StartStep("Crear IOVariables for TestCase")
	tcInputA, err := exam_factory.NewIOVariable("a", exam_entities.VariableFormatInt, "10")
	if err != nil {
		process.Fail("create test case input", err)
	}
	tcInputB, err := exam_factory.NewIOVariable("b", exam_entities.VariableFormatInt, "15")
	if err != nil {
		process.Fail("create test case input", err)
	}
	tcOutput, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "25")
	if err != nil {
		process.Fail("create test case output", err)
	}
	tcInputA, err = ioVariableRepository.CreateIOVariable(ctx, tcInputA)
	if tcInputA == nil || err != nil {
		process.Fail("persist input io variable", err)
	}
	tcInputB, err = ioVariableRepository.CreateIOVariable(ctx, tcInputB)
	if tcInputB == nil || err != nil {
		process.Fail("persist input io variable", err)
	}
	tcOutput, err = ioVariableRepository.CreateIOVariable(ctx, tcOutput)
	if tcOutput == nil || err != nil {
		process.Fail("persist output io variable", err)
	}
	process.EndStep()

	// [STEP 16] Create a TestCase associated with the challenge
	process.StartStep("Crear TestCase")
	testCase, err := exam_factory.NewTestCase(
		"sample_1",
		[]exam_entities.IOVariable{*tcInputA, *tcInputB},
		*tcOutput,
		true,
		10,
		challengeID,
	)
	if err != nil {
		process.Fail("build test case with factory", err)
	}

	createdTestCase, err := testCaseRepository.CreateTestCase(ctx, testCase)
	if err != nil {
		process.Fail("create test case", err)
	}
	if createdTestCase == nil || createdTestCase.ID == "" {
		process.Fail("create test case", fmt.Errorf("expected created test case with ID"))
	}
	testCaseID := createdTestCase.ID
	process.Log(fmt.Sprintf("TestCase creado. testCaseID=%s", testCaseID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando test case %s", testCaseID)
		_ = testCaseRepository.DeleteTestCase(ctx, testCaseID)

		_ = ioVariableRepository.DeleteIOVariable(ctx, tcInputA.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, tcInputB.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, tcOutput.ID)
	}()

	// [STEP 17] Get test case by ID and verify details
	process.StartStep(fmt.Sprintf("Leer test case por ID=%s e hidratar IOVariables", testCaseID))
	reloadedTestCase, err := testCaseRepository.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		process.Fail("get test case by id", err)
	}
	if reloadedTestCase == nil {
		process.Fail("get test case by id", fmt.Errorf("expected reloaded test case"))
	}
	if len(reloadedTestCase.Input) == 0 || reloadedTestCase.ExpectedOutput.ID == "" {
		process.Fail("get test case by id", fmt.Errorf("expected test case with hydrated input/output variables"))
	}
	process.Log(fmt.Sprintf("TestCase recargado con IO. inputCount=%d outputID=%s", len(reloadedTestCase.Input), reloadedTestCase.ExpectedOutput.ID))
	process.EndStep()

	// [STEP 18] Update test case details
	process.StartStep("Actualizar test case")
	createdTestCase.Name = "sample_1_updated"
	createdTestCase.Points = 25
	updatedTestCase, err := testCaseRepository.UpdateTestCase(ctx, createdTestCase)
	if err != nil {
		process.Fail("update test case", err)
	}
	if updatedTestCase == nil || updatedTestCase.Name != "sample_1_updated" {
		process.Fail("update test case", fmt.Errorf("expected updated test case name"))
	}
	process.Log(fmt.Sprintf("TestCase actualizado. name=%q points=%d", updatedTestCase.Name, updatedTestCase.Points))
	process.EndStep()

	// [STEP 19] Validate test case listings and IO variable relations
	process.StartStep("Validar listados y lecturas de IO de test case")
	challengeTestCases, err := testCaseRepository.GetTestCasesByChallengeID(ctx, challengeID)
	if err != nil {
		process.Fail("get test cases by challenge id", err)
	}
	if len(challengeTestCases) == 0 {
		process.Fail("get test cases by challenge id", fmt.Errorf("expected at least one test case for challenge"))
	}
	testCaseInputs, err := testCaseRepository.GetInputVariablesByTestCaseID(ctx, testCaseID)
	if err != nil {
		process.Fail("get test case input variables", err)
	}
	testCaseOutputs, err := testCaseRepository.GetOutputVariablesByTestCaseID(ctx, testCaseID)
	if err != nil {
		process.Fail("get test case output variables", err)
	}
	if len(testCaseInputs) == 0 || len(testCaseOutputs) == 0 {
		process.Fail("test case io relations", fmt.Errorf("expected test case input/output variable relations"))
	}
	process.Log(fmt.Sprintf("Relaciones de IO test case validadas. inputs=%d outputs=%d", len(testCaseInputs), len(testCaseOutputs)))
	process.EndStep()

	// [STEP 20] Delete test case
	process.StartStep("Eliminar test case y validar borrado")
	if err := testCaseRepository.DeleteTestCase(ctx, testCaseID); err != nil {
		process.Fail("delete test case", err)
	}
	deletedTestCase, err := testCaseRepository.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		process.Fail("get test case after delete", err)
	}
	if deletedTestCase != nil {
		process.Fail("get test case after delete", fmt.Errorf("expected test case %s to be deleted", testCaseID))
	}
	testCaseID = ""
	process.Log("TestCase eliminado")
	process.EndStep()

	// [STEP 21] Delete challenge
	process.StartStep("Eliminar challenge y validar borrado")
	if err := challengeRepository.DeleteChallenge(ctx, challengeID); err != nil {
		process.Fail("delete challenge", err)
	}
	deletedChallenge, err := challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		process.Fail("get challenge after delete", err)
	}
	if deletedChallenge != nil {
		process.Fail("get challenge after delete", fmt.Errorf("expected challenge %s to be deleted", challengeID))
	}
	challengeID = ""
	process.Log("Challenge eliminado")
	process.EndStep()

	// [STEP 22] Delete exam
	process.StartStep("Eliminar examen y validar borrado")
	if err := examRepository.DeleteExam(ctx, examID); err != nil {
		process.Fail("delete exam", err)
	}
	deletedExam, err := examRepository.GetExamByID(ctx, examID)
	if err != nil {
		process.Fail("get exam after delete", err)
	}
	if deletedExam != nil {
		process.Fail("get exam after delete", fmt.Errorf("expected exam %s to be deleted", examID))
	}
	examID = ""
	process.Log("Examen eliminado")
	process.EndStep()

	// [STEP 23] Delete Course
	process.StartStep("Eliminar curso y validar borrado")
	if err := courseRepository.DeleteCourse(ctx, courseID); err != nil {
		process.Fail("delete course", err)
	}
	deletedCourse, err := courseRepository.GetCourseByID(ctx, courseID)
	if err != nil {
		process.Fail("get course after delete", err)
	}
	if deletedCourse != nil {
		process.Fail("get course after delete", fmt.Errorf("expected course %s to be deleted", courseID))
	}
	courseID = ""
	process.Log("Curso eliminado")
	process.EndStep()

	process.End()
}
