package cache_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	course_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	exam_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
	cache_infra "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	cache_course "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble/course"
	cache_exam "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble/exam"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_course_infra "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	roble_exam_infra "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	hasher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestCacheExamCRUD(t *testing.T) {
	process := test.StartTest(t, "Cache Exam Creation and Persistence")
	email := "test@test.com"
	password := "Password123!"

	// [STEP 1] Initialize Roble + MariaDB clients, schema and repositories
	process.StartStep("Inicializando Roble, MariaDB y cache repositories")

	httpClient := &http.Client{Timeout: 15 * time.Second}
	robleClient, err := roble_infrastructure.NewRobleClient(httpClient)
	if err != nil {
		process.Fail("initialize roble client", err)
	}

	mariaClient, err := cache_infra.NewMariaDBClient()
	if err != nil {
		t.Skipf("cache DB unavailable, skipping: %v", err)
	}
	defer mariaClient.Close()

	cacheAdapter := cache_infra.NewCacheAdapter(mariaClient)
	initializer := cache_infra.NewSchemaInitializer(mariaClient, cache_infra.DefaultEntityConfigs())
	if err := initializer.Initialize(); err != nil {
		process.Fail("initialize cache schema", err)
	}

	robleAdapter := roble_infrastructure.NewRobleDatabaseAdapter(robleClient)
	userRepository := roble_user_infrastructure.NewUserRepository(robleAdapter)
	authAdapter := roble_user_infrastructure.NewRobleAuthAdapter(robleAdapter, userRepository)
	robleIoVarRepo := roble_exam_infra.NewIOVariableRepository(robleAdapter)
	robleDirectExam := roble_exam_infra.NewExamRepository(robleAdapter)
	robleDirectChallenge := roble_exam_infra.NewChallengeRepository(robleAdapter, robleIoVarRepo)
	robleDirectTestCase := roble_exam_infra.NewTestCaseRepository(robleAdapter, robleIoVarRepo)
	robleDirectCourse := roble_course_infra.NewCourseRepository(robleAdapter)
	robleDirectExamItem := roble_exam_infra.NewExamItemRepository(robleAdapter)
	courseRepository := cache_course.NewCourseRepository(robleDirectCourse, cacheAdapter)
	examRepository := cache_exam.NewExamRepository(robleDirectExam, cacheAdapter)
	examItemRepository := cache_exam.NewExamItemRepository(robleDirectExamItem, cacheAdapter)
	ioVariableRepository := cache_exam.NewIOVariableRepository(robleIoVarRepo, cacheAdapter)
	challengeRepository := cache_exam.NewChallengeRepository(robleDirectChallenge, cacheAdapter)
	testCaseRepository := cache_exam.NewTestCaseRepository(robleDirectTestCase, cacheAdapter)
	process.Log("Clientes y repositories inicializados")
	process.EndStep()

	// [STEP 2] Hash password
	process.StartStep("Hashear password")
	secAdapter := hasher.NewSecurityAdapter()
	hashedPassword, err := secAdapter.Hash(password)
	if err != nil {
		process.Fail("hash password", err)
	}
	process.EndStep()

	// [STEP 3] Login teacher
	process.StartStep("Login docente")
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
	enrollmentCode := fmt.Sprintf("CACHE-EXAM-%d", now.UnixNano())
	courseCode := fmt.Sprintf("CEX-%d", now.Unix()%100000)

	// [STEP 4] Create auxiliary course
	process.StartStep("Crear curso auxiliar")
	course, err := course_factory.NewCourse(
		"Cache Exam Integration Course",
		"Course for cache exam integration test",
		consts.CourseVisibilityPublic,
		consts.CourseColourBlue,
		courseCode,
		&course_entities.Period{Year: now.Year(), Semester: consts.AcademicFirstPeriod},
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

	// [STEP 5] Create exam
	process.StartStep("Crear examen (write-through)")
	startTime := now
	endTime := startTime.Add(90 * time.Minute)
	exam, err := exam_factory.NewExam(
		"Cache Integration Exam",
		"Exam created by cache integration test",
		exam_consts.VisibilityCourse,
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

	robleExam, err := robleDirectExam.GetExamByID(ctx, examID)
	if err != nil || robleExam == nil {
		process.Fail("roble verify create exam", fmt.Errorf("expected exam persisted in Roble, err=%v", err))
	}
	process.Log(fmt.Sprintf("Roble verify: exam persisted id=%s", robleExam.ID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando examen %s", examID)
		_ = examRepository.DeleteExam(ctx, examID)
	}()

	// [STEP 6] GetExamByID — cache hit
	process.StartStep("Leer examen por ID — cache hit")
	cachedExam, err := examRepository.GetExamByID(ctx, examID)
	if err != nil {
		process.Fail("get exam by id", err)
	}
	if cachedExam == nil || cachedExam.ID != examID {
		process.Fail("get exam by id", fmt.Errorf("expected exam from cache"))
	}
	process.Log(fmt.Sprintf("Cache hit: title=%q", cachedExam.Title))
	process.EndStep()

	// [STEP 7] Update exam
	process.StartStep("Actualizar examen (write-through update)")
	createdExam.Title = "Cache Integration Exam Updated"
	createdExam.TryLimit = 3
	updatedExam, err := examRepository.UpdateExam(ctx, createdExam)
	if err != nil {
		process.Fail("update exam", err)
	}
	if updatedExam == nil || updatedExam.Title != "Cache Integration Exam Updated" {
		process.Fail("update exam", fmt.Errorf("expected updated exam title"))
	}
	process.Log(fmt.Sprintf("Examen actualizado. title=%q tryLimit=%d", updatedExam.Title, updatedExam.TryLimit))

	robleUpdatedExam, err := robleDirectExam.GetExamByID(ctx, examID)
	if err != nil || robleUpdatedExam == nil || robleUpdatedExam.Title != "Cache Integration Exam Updated" {
		process.Fail("roble verify update exam", fmt.Errorf("expected updated exam in Roble"))
	}
	process.Log(fmt.Sprintf("Roble verify: updated title=%q", robleUpdatedExam.Title))
	process.EndStep()

	// [STEP 8] Validate exam lists
	process.StartStep("Validar listados de examenes por courseID y teacherID")
	courseExams, err := examRepository.GetExamsByCourseID(ctx, courseID)
	if err != nil {
		process.Fail("get exams by course id", err)
	}
	teacherExams, err := examRepository.GetExamsByTeacherID(ctx, teacherID)
	if err != nil {
		process.Fail("get exams by teacher id", err)
	}
	if len(courseExams) == 0 || len(teacherExams) == 0 {
		process.Fail("exam lists", fmt.Errorf("expected exams in list queries"))
	}
	process.Log(fmt.Sprintf("Listados: courseExams=%d teacherExams=%d", len(courseExams), len(teacherExams)))
	process.EndStep()

	// [STEP 9] Create IOVariables for challenge
	process.StartStep("Crear IOVariables para challenge")
	inputA, err := exam_factory.NewIOVariable("a", exam_consts.VariableFormatInt, "2")
	if err != nil {
		process.Fail("create input io variable", err)
	}
	inputB, err := exam_factory.NewIOVariable("b", exam_consts.VariableFormatInt, "3")
	if err != nil {
		process.Fail("create input io variable", err)
	}
	output, err := exam_factory.NewIOVariable("sum", exam_consts.VariableFormatInt, "5")
	if err != nil {
		process.Fail("create output io variable", err)
	}
	inputA, err = ioVariableRepository.CreateIOVariable(ctx, inputA)
	if inputA == nil || err != nil {
		process.Fail("persist input io variable a", err)
	}
	inputB, err = ioVariableRepository.CreateIOVariable(ctx, inputB)
	if inputB == nil || err != nil {
		process.Fail("persist input io variable b", err)
	}
	output, err = ioVariableRepository.CreateIOVariable(ctx, output)
	if output == nil || err != nil {
		process.Fail("persist output io variable", err)
	}
	process.Log(fmt.Sprintf("IOVariables creadas. inputA=%s inputB=%s output=%s", inputA.ID, inputB.ID, output.ID))
	process.EndStep()

	// [STEP 10] Create challenge
	process.StartStep("Crear challenge")
	challenge, err := exam_factory.NewChallenge(
		"Cache Sum Challenge",
		"Return the sum of two numbers",
		[]string{"math", "cache"},
		exam_consts.ChallengeStatusDraft,
		exam_consts.ChallengeDifficultyEasy,
		1500,
		256,
		[]exam_entities.CodeTemplate{
			{Language: "python", Template: "def solve(a, b):\n    return a + b\n"},
		},
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

	robleChallenge, err := robleDirectChallenge.GetChallengeByID(ctx, challengeID)
	if err != nil || robleChallenge == nil {
		process.Fail("roble verify create challenge", fmt.Errorf("expected challenge persisted in Roble, err=%v", err))
	}
	process.Log(fmt.Sprintf("Roble verify: challenge persisted id=%s", robleChallenge.ID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
		_ = challengeRepository.DeleteChallenge(ctx, challengeID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, inputA.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, inputB.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, output.ID)
	}()

	// [STEP 11] GetChallengeByID — verify nested IOVariables round-trip through cache
	process.StartStep("Leer challenge — verificar que IOVariables deserializan correctamente desde cache")
	cachedChallenge, err := challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		process.Fail("get challenge by id", err)
	}
	if cachedChallenge == nil {
		process.Fail("get challenge by id", fmt.Errorf("expected challenge from cache"))
	}
	if len(cachedChallenge.InputVariables) == 0 || cachedChallenge.OutputVariable.ID == "" {
		process.Fail("challenge io round-trip", fmt.Errorf("expected hydrated IOVariables after cache deserialization"))
	}
	process.Log(fmt.Sprintf("Challenge desde cache: inputs=%d outputID=%s", len(cachedChallenge.InputVariables), cachedChallenge.OutputVariable.ID))

	// Second read to confirm cache consistency for nested structs
	cachedChallenge2, err := challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		process.Fail("second get challenge by id", err)
	}
	if cachedChallenge2 == nil || len(cachedChallenge2.InputVariables) != len(cachedChallenge.InputVariables) {
		process.Fail("repeated challenge cache hit", fmt.Errorf("inconsistent IOVariables on repeated cache read"))
	}
	process.Log("IOVariables consistentes en lecturas repetidas desde cache")
	process.EndStep()

	// [STEP 12] GetInputVariablesByChallengeID / GetOutputVariablesByChallengeID
	process.StartStep("Verificar GetInputVariables y GetOutputVariables desde cache")
	inputs, err := challengeRepository.GetInputVariablesByChallengeID(ctx, challengeID)
	if err != nil {
		process.Fail("get input variables by challenge id", err)
	}
	outputs, err := challengeRepository.GetOutputVariablesByChallengeID(ctx, challengeID)
	if err != nil {
		process.Fail("get output variables by challenge id", err)
	}
	if len(inputs) == 0 || len(outputs) == 0 {
		process.Fail("challenge io relations", fmt.Errorf("expected challenge input/output variables"))
	}
	process.Log(fmt.Sprintf("IO variables: inputs=%d outputs=%d", len(inputs), len(outputs)))
	process.EndStep()

	// [STEP 13] Create ExamItem
	process.StartStep("Crear ExamItem")
	examItem, err := exam_factory.NewExamItem(challengeID, examID, 1, 100, nil)
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
	process.Log(fmt.Sprintf("ExamItem creado. id=%s", createdExamItem.ID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando exam item %s", createdExamItem.ID)
		_ = examItemRepository.DeleteExamItem(ctx, createdExamItem.ID)
	}()

	// [STEP 14] Create IOVariables and TestCase
	process.StartStep("Crear IOVariables y TestCase")
	tcInputA, err := exam_factory.NewIOVariable("a", exam_consts.VariableFormatInt, "10")
	if err != nil {
		process.Fail("create test case input", err)
	}
	tcOutput, err := exam_factory.NewIOVariable("sum", exam_consts.VariableFormatInt, "13")
	if err != nil {
		process.Fail("create test case output", err)
	}
	tcInputA, err = ioVariableRepository.CreateIOVariable(ctx, tcInputA)
	if tcInputA == nil || err != nil {
		process.Fail("persist test case input", err)
	}
	tcOutput, err = ioVariableRepository.CreateIOVariable(ctx, tcOutput)
	if tcOutput == nil || err != nil {
		process.Fail("persist test case output", err)
	}

	testCase, err := exam_factory.NewTestCase(
		"cache_sample_1",
		[]exam_entities.IOVariable{*tcInputA},
		*tcOutput,
		true,
		10,
		false,
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
	process.Log(fmt.Sprintf("TestCase creado. id=%s isSample=%v", testCaseID, createdTestCase.IsSample))

	robleTC, err := robleDirectTestCase.GetTestCaseByID(ctx, testCaseID)
	if err != nil || robleTC == nil {
		process.Fail("roble verify create test case", fmt.Errorf("expected test case persisted in Roble, err=%v", err))
	}
	process.Log(fmt.Sprintf("Roble verify: test case persisted id=%s isSample=%v", robleTC.ID, robleTC.IsSample))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando test case %s", testCaseID)
		_ = testCaseRepository.DeleteTestCase(ctx, testCaseID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, tcInputA.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, tcOutput.ID)
	}()

	// [STEP 15] GetTestCaseByID — verify NormalizeBools for is_sample and custom
	process.StartStep("Leer TestCase desde cache — verificar campos bool (is_sample, custom)")
	cachedTC, err := testCaseRepository.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		process.Fail("get test case by id", err)
	}
	if cachedTC == nil {
		process.Fail("get test case by id", fmt.Errorf("expected test case from cache"))
	}
	if !cachedTC.IsSample {
		process.Fail("bool normalization", fmt.Errorf("expected IsSample=true after cache deserialization, got false"))
	}
	if len(cachedTC.Input) == 0 || cachedTC.ExpectedOutput.ID == "" {
		process.Fail("test case io round-trip", fmt.Errorf("expected hydrated IOVariables in cached test case"))
	}
	process.Log(fmt.Sprintf("TestCase cache hit: isSample=%v custom=%v inputs=%d", cachedTC.IsSample, cachedTC.Custom, len(cachedTC.Input)))
	process.EndStep()

	// [STEP 16] GetTestCasesByChallengeID — list cache
	process.StartStep("Listar TestCases por challengeID")
	challengeTestCases, err := testCaseRepository.GetTestCasesByChallengeID(ctx, challengeID)
	if err != nil {
		process.Fail("get test cases by challenge id", err)
	}
	if len(challengeTestCases) == 0 {
		process.Fail("test case list", fmt.Errorf("expected at least one test case for challenge"))
	}
	process.Log(fmt.Sprintf("TestCases en challenge: %d", len(challengeTestCases)))
	process.EndStep()

	// [STEP 17] Delete test case and challenge
	process.StartStep("Eliminar TestCase y Challenge y verificar cache invalidation")
	if err := testCaseRepository.DeleteTestCase(ctx, testCaseID); err != nil {
		process.Fail("delete test case", err)
	}
	deletedTC, err := testCaseRepository.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		process.Fail("get test case after delete", err)
	}
	if deletedTC != nil {
		process.Fail("test case cache invalidation", fmt.Errorf("expected test case deleted"))
	}
	robleDeletedTC, err := robleDirectTestCase.GetTestCaseByID(ctx, testCaseID)
	if err != nil {
		process.Fail("roble verify delete test case", err)
	}
	if robleDeletedTC != nil {
		process.Fail("roble verify delete test case", fmt.Errorf("expected test case %s deleted from Roble", testCaseID))
	}
	process.Log("Roble verify: test case deleted from Roble")
	testCaseID = ""

	if err := challengeRepository.DeleteChallenge(ctx, challengeID); err != nil {
		process.Fail("delete challenge", err)
	}
	deletedChallenge, err := challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		process.Fail("get challenge after delete", err)
	}
	if deletedChallenge != nil {
		process.Fail("challenge cache invalidation", fmt.Errorf("expected challenge deleted"))
	}
	robleDeletedChallenge, err := robleDirectChallenge.GetChallengeByID(ctx, challengeID)
	if err != nil {
		process.Fail("roble verify delete challenge", err)
	}
	if robleDeletedChallenge != nil {
		process.Fail("roble verify delete challenge", fmt.Errorf("expected challenge %s deleted from Roble", challengeID))
	}
	process.Log("Roble verify: challenge deleted from Roble")
	challengeID = ""
	process.Log("TestCase y Challenge eliminados")
	process.EndStep()

	// [STEP 18] Delete exam
	process.StartStep("Eliminar examen y verificar cache invalidation")
	if err := examRepository.DeleteExam(ctx, examID); err != nil {
		process.Fail("delete exam", err)
	}
	deletedExam, err := examRepository.GetExamByID(ctx, examID)
	if err != nil {
		process.Fail("get exam after delete", err)
	}
	if deletedExam != nil {
		process.Fail("exam cache invalidation", fmt.Errorf("expected exam deleted"))
	}
	robleDeletedExam, err := robleDirectExam.GetExamByID(ctx, examID)
	if err != nil {
		process.Fail("roble verify delete exam", err)
	}
	if robleDeletedExam != nil {
		process.Fail("roble verify delete exam", fmt.Errorf("expected exam %s deleted from Roble", examID))
	}
	process.Log("Roble verify: exam deleted from Roble")
	examID = ""
	process.Log("Examen eliminado")
	process.EndStep()

	process.End()
}
