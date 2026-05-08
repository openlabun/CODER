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
	submission_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	course_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	exam_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
	submission_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/submission"
	cache_infra "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	cache_course "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble/course"
	cache_exam "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble/exam"
	cache_submission "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble/submission"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_exam_infra "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	roble_submission_infra "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	hasher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestCacheSubmissionCRUD(t *testing.T) {
	process := test.StartTest(t, "Cache Submission Creation and Persistence")

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
	courseRepository := cache_course.NewCourseRepository(robleAdapter, cacheAdapter)
	examRepository := cache_exam.NewExamRepository(robleAdapter, cacheAdapter)
	examItemRepository := cache_exam.NewExamItemRepository(robleAdapter, cacheAdapter)
	examScoreRepository := cache_exam.NewExamScoreRepository(robleAdapter, cacheAdapter)
	examItemScoreRepository := cache_exam.NewExamItemScoreRepository(robleAdapter, cacheAdapter)
	ioVariableRepository := cache_exam.NewIOVariableRepository(robleAdapter, cacheAdapter)
	challengeRepository := cache_exam.NewChallengeRepository(robleAdapter, cacheAdapter)
	testCaseRepository := cache_exam.NewTestCaseRepository(robleAdapter, cacheAdapter)
	sessionRepository := cache_submission.NewSessionRepository(robleAdapter, cacheAdapter)
	submissionRepository := cache_submission.NewSubmissionRepository(robleAdapter, cacheAdapter)
	resultRepository := cache_submission.NewSubmissionResultRepository(robleAdapter, cacheAdapter)
	robleSubIoVarRepo := roble_exam_infra.NewIOVariableRepository(robleAdapter)
	robleDirectSession := roble_submission_infra.NewSessionRepository(robleAdapter)
	robleDirectSubmission := roble_submission_infra.NewSubmissionRepository(robleAdapter)
	robleDirectResult := roble_submission_infra.NewSubmissionResultRepository(robleAdapter, robleSubIoVarRepo)
	process.Log("Clientes y repositories inicializados")
	process.EndStep()

	// [STEP 2] Hash password
	process.StartStep("Hashear password")
	secAdapter := hasher.NewSecurityAdapter()
	hashedPassword, err := secAdapter.Hash("Password123!")
	if err != nil {
		process.Fail("hash password", err)
	}
	process.EndStep()

	// [STEP 3] Login teacher
	process.StartStep("Login docente")
	access, err := authAdapter.LoginUser("test@test.com", hashedPassword)
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
	process.Log(fmt.Sprintf("Login exitoso. userID=%s", teacherID))
	process.EndStep()

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("CACHE-SUB-%d", now.UnixNano())
	courseCode := fmt.Sprintf("CSB-%d", now.Unix()%100000)

	// [STEP 4] Create auxiliary course
	process.StartStep("Crear curso auxiliar")
	course, err := course_factory.NewCourse(
		"Cache Submission Integration Course",
		"Course for cache submission integration",
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
	courseID := createdCourse.ID
	process.Log(fmt.Sprintf("Curso creado. courseID=%s", courseID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando curso %s", courseID)
		_ = courseRepository.DeleteCourse(ctx, courseID)
	}()

	// [STEP 5] Create auxiliary exam
	process.StartStep("Crear examen auxiliar")
	startTime := now
	endTime := startTime.Add(2 * time.Hour)
	exam, err := exam_factory.NewExam(
		"Cache Submission Exam",
		"Exam for cache submission integration",
		exam_consts.VisibilityCourse,
		startTime,
		&endTime,
		false,
		5400,
		3,
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
	examID := createdExam.ID
	process.Log(fmt.Sprintf("Examen creado. examID=%s", examID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando examen %s", examID)
		_ = examRepository.DeleteExam(ctx, examID)
	}()

	// [STEP 6] Create challenge, exam item and test case
	process.StartStep("Crear challenge, exam item y test case auxiliares")
	inputA, err := exam_factory.NewIOVariable("a", exam_consts.VariableFormatInt, "2")
	if err != nil {
		process.Fail("create challenge input variable", err)
	}
	output, err := exam_factory.NewIOVariable("sum", exam_consts.VariableFormatInt, "5")
	if err != nil {
		process.Fail("create challenge output variable", err)
	}
	inputA, err = ioVariableRepository.CreateIOVariable(ctx, inputA)
	if err != nil {
		process.Fail("persist challenge input variable", err)
	}
	output, err = ioVariableRepository.CreateIOVariable(ctx, output)
	if err != nil {
		process.Fail("persist challenge output variable", err)
	}

	challenge, err := exam_factory.NewChallenge(
		"Cache Sum Challenge Submission",
		"Challenge for cache submission integration",
		[]string{"cache", "submission"},
		exam_consts.ChallengeStatusDraft,
		exam_consts.ChallengeDifficultyEasy,
		1500,
		256,
		[]exam_entities.CodeTemplate{
			{Language: "python", Template: "def solve(a):\n    return a + 3\n"},
		},
		[]exam_entities.IOVariable{*inputA},
		*output,
		"1 <= a <= 1000",
		examID,
	)
	if err != nil {
		process.Fail("build challenge with factory", err)
	}
	createdChallenge, err := challengeRepository.CreateChallenge(ctx, challenge)
	if err != nil {
		process.Fail("create challenge", err)
	}
	challengeID := createdChallenge.ID

	defer func() {
		t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
		_ = challengeRepository.DeleteChallenge(ctx, challengeID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, inputA.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, output.ID)
	}()

	examItem, err := exam_factory.NewExamItem(examID, challengeID, 1, 100, nil)
	if err != nil {
		process.Fail("build exam item", err)
	}
	createdExamItem, err := examItemRepository.CreateExamItem(ctx, examItem)
	if err != nil {
		process.Fail("create exam item", err)
	}
	defer func() {
		t.Logf("[CLEANUP] Eliminando exam item %s", createdExamItem.ID)
		_ = examItemRepository.DeleteExamItem(ctx, createdExamItem.ID)
	}()

	tcInput, err := exam_factory.NewIOVariable("a", exam_consts.VariableFormatInt, "10")
	if err != nil {
		process.Fail("create test case input", err)
	}
	tcOutput, err := exam_factory.NewIOVariable("sum", exam_consts.VariableFormatInt, "13")
	if err != nil {
		process.Fail("create test case output", err)
	}
	tcInput, err = ioVariableRepository.CreateIOVariable(ctx, tcInput)
	if err != nil {
		process.Fail("persist test case input", err)
	}
	tcOutput, err = ioVariableRepository.CreateIOVariable(ctx, tcOutput)
	if err != nil {
		process.Fail("persist test case output", err)
	}
	testCase, err := exam_factory.NewTestCase(
		"cache_sub_sample",
		[]exam_entities.IOVariable{*tcInput},
		*tcOutput,
		true,
		10,
		false,
		challengeID,
	)
	if err != nil {
		process.Fail("build test case", err)
	}
	createdTestCase, err := testCaseRepository.CreateTestCase(ctx, testCase)
	if err != nil {
		process.Fail("create test case", err)
	}
	testCaseID := createdTestCase.ID
	process.Log(fmt.Sprintf("Challenge=%s ExamItem=%s TestCase=%s", challengeID, createdExamItem.ID, testCaseID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando test case %s", testCaseID)
		_ = testCaseRepository.DeleteTestCase(ctx, testCaseID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, tcInput.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, tcOutput.ID)
	}()

	// [STEP 7] Create session
	process.StartStep("Crear sesion (write-through)")
	session, err := submission_factory.NewSession(teacherID, createdExam)
	if err != nil {
		process.Fail("build session with factory", err)
	}
	session.TimeLeft = 5400

	createdSession, err := sessionRepository.CreateSession(ctx, session)
	if err != nil {
		process.Fail("create session", err)
	}
	if createdSession == nil || createdSession.ID == "" {
		process.Fail("create session", fmt.Errorf("expected created session with ID"))
	}
	sessionID := createdSession.ID
	process.Log(fmt.Sprintf("Session creada. sessionID=%s", sessionID))

	robleSession, err := robleDirectSession.GetSessionByID(ctx, sessionID)
	if err != nil || robleSession == nil {
		process.Fail("roble verify create session", fmt.Errorf("expected session persisted in Roble, err=%v", err))
	}
	process.Log(fmt.Sprintf("Roble verify: session persisted id=%s", robleSession.ID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando session %s", sessionID)
		_ = sessionRepository.DeleteSession(ctx, sessionID)
	}()

	// [STEP 8] Session cache hit and update
	process.StartStep("Leer sesion desde cache y actualizar")
	cachedSession, err := sessionRepository.GetSessionByID(ctx, sessionID)
	if err != nil {
		process.Fail("get session by id", err)
	}
	if cachedSession == nil || cachedSession.ID != sessionID {
		process.Fail("session cache hit", fmt.Errorf("expected session from cache"))
	}
	process.Log(fmt.Sprintf("Session cache hit: status=%s", cachedSession.Status))

	createdSession.Attempts = 1
	createdSession.TimeLeft = 5000
	updatedSession, err := sessionRepository.UpdateSession(ctx, createdSession)
	if err != nil {
		process.Fail("update session", err)
	}
	if updatedSession == nil {
		process.Fail("update session", fmt.Errorf("expected updated session"))
	}
	process.Log(fmt.Sprintf("Session actualizada. attempts=%d timeLeft=%d", updatedSession.Attempts, updatedSession.TimeLeft))

	robleUpdatedSession, err := robleDirectSession.GetSessionByID(ctx, sessionID)
	if err != nil || robleUpdatedSession == nil || robleUpdatedSession.Attempts != 1 {
		process.Fail("roble verify update session", fmt.Errorf("expected updated session in Roble"))
	}
	process.Log(fmt.Sprintf("Roble verify: session updated attempts=%d", robleUpdatedSession.Attempts))

	sessionsByExam, err := sessionRepository.GetSessionsByExamID(ctx, examID)
	if err != nil {
		process.Fail("get sessions by exam", err)
	}
	sessionsByStudent, err := sessionRepository.GetSessionsByStudentID(ctx, teacherID)
	if err != nil {
		process.Fail("get sessions by student", err)
	}
	if len(sessionsByExam) == 0 || len(sessionsByStudent) == 0 {
		process.Fail("session list queries", fmt.Errorf("expected sessions in query results"))
	}
	process.Log(fmt.Sprintf("Session queries: byExam=%d byStudent=%d", len(sessionsByExam), len(sessionsByStudent)))
	process.EndStep()

	// [STEP 9] ExamScore + ExamItemScore
	process.StartStep("Crear ExamScore y ExamItemScore")
	examScore, err := exam_factory.NewExamScore(examID, sessionID, teacherID)
	if err != nil {
		process.Fail("build exam score", err)
	}
	examScore, err = examScoreRepository.CreateExamScore(ctx, examScore)
	if err != nil {
		process.Fail("create exam score", err)
	}
	defer func() {
		t.Logf("[CLEANUP] Eliminando exam score %s", examScore.ID)
		_ = examScoreRepository.DeleteExamScore(ctx, examScore.ID)
	}()

	examItemScore, err := exam_factory.NewExamItemScore(createdExamItem.ID, examScore.ID)
	if err != nil {
		process.Fail("build exam item score", err)
	}
	examItemScore, err = examItemScoreRepository.CreateExamItemScore(ctx, examItemScore)
	if err != nil {
		process.Fail("create exam item score", err)
	}
	defer func() {
		t.Logf("[CLEANUP] Eliminando exam item score %s", examItemScore.ID)
		_ = examItemScoreRepository.DeleteExamItemScore(ctx, examItemScore.ID)
	}()
	process.Log(fmt.Sprintf("ExamScore=%s ExamItemScore=%s", examScore.ID, examItemScore.ID))
	process.EndStep()

	// [STEP 10] Create submission
	process.StartStep("Crear submission (write-through)")
	submission, err := submission_factory.NewSubmission(
		"print(1+2)",
		submission_consts.LanguagePython,
		true,
		challengeID,
		sessionID,
		teacherID,
		&examItemScore.ID,
	)
	if err != nil {
		process.Fail("build submission with factory", err)
	}
	createdSubmission, err := submissionRepository.CreateSubmission(ctx, submission)
	if err != nil {
		process.Fail("create submission", err)
	}
	if createdSubmission == nil || createdSubmission.ID == "" {
		process.Fail("create submission", fmt.Errorf("expected created submission with ID"))
	}
	submissionID := createdSubmission.ID
	process.Log(fmt.Sprintf("Submission creada. submissionID=%s scorable=%v", submissionID, createdSubmission.Scorable))

	robleSubmission, err := robleDirectSubmission.GetSubmissionByID(ctx, submissionID)
	if err != nil || robleSubmission == nil {
		process.Fail("roble verify create submission", fmt.Errorf("expected submission persisted in Roble, err=%v", err))
	}
	process.Log(fmt.Sprintf("Roble verify: submission persisted id=%s", robleSubmission.ID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando submission %s", submissionID)
		_ = submissionRepository.DeleteSubmission(ctx, submissionID)
	}()

	// [STEP 11] GetSubmissionByID — verify NormalizeBools for scorable
	process.StartStep("Leer submission desde cache — verificar campo bool scorable")
	cachedSubmission, err := submissionRepository.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		process.Fail("get submission by id", err)
	}
	if cachedSubmission == nil {
		process.Fail("get submission by id", fmt.Errorf("expected submission from cache"))
	}
	if !cachedSubmission.Scorable {
		process.Fail("bool normalization", fmt.Errorf("expected Scorable=true after cache deserialization, got false"))
	}
	process.Log(fmt.Sprintf("Submission cache hit: language=%s scorable=%v", cachedSubmission.Language, cachedSubmission.Scorable))
	process.EndStep()

	// [STEP 12] Update submission
	process.StartStep("Actualizar submission y verificar cache")
	createdSubmission.Code = "print(2+3)"
	createdSubmission.Score = 100
	createdSubmission.TimeMsTotal = 45
	updatedSubmission, err := submissionRepository.UpdateSubmission(ctx, createdSubmission)
	if err != nil {
		process.Fail("update submission", err)
	}
	if updatedSubmission == nil || updatedSubmission.Score != 100 {
		process.Fail("update submission", fmt.Errorf("expected updated submission score"))
	}
	process.Log(fmt.Sprintf("Submission actualizada. score=%d timeMsTotal=%d", updatedSubmission.Score, updatedSubmission.TimeMsTotal))

	afterUpdate, err := submissionRepository.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		process.Fail("get submission after update", err)
	}
	if afterUpdate == nil || afterUpdate.Score != 100 {
		process.Fail("submission cache after update", fmt.Errorf("expected updated score in cache"))
	}

	robleUpdatedSub, err := robleDirectSubmission.GetSubmissionByID(ctx, submissionID)
	if err != nil || robleUpdatedSub == nil || robleUpdatedSub.Score != 100 {
		process.Fail("roble verify update submission", fmt.Errorf("expected updated submission in Roble"))
	}
	process.Log(fmt.Sprintf("Roble verify: submission updated score=%d", robleUpdatedSub.Score))
	process.EndStep()

	// [STEP 13] Submission list queries
	process.StartStep("Verificar queries de submission por sesion")
	bySession, err := submissionRepository.GetSubmissionsBySessionID(ctx, sessionID, nil, nil, nil)
	if err != nil {
		process.Fail("get submissions by session", err)
	}
	if len(bySession) == 0 {
		process.Fail("submission list", fmt.Errorf("expected submissions by session"))
	}
	process.Log(fmt.Sprintf("bySession=%d", len(bySession)))
	process.EndStep()

	// [STEP 14] Create SubmissionResult (with nested ActualOutput IOVariable)
	process.StartStep("Crear SubmissionResult con ActualOutput anidado")
	result, err := submission_factory.NewSubmissionResult(submissionID, testCaseID)
	if err != nil {
		process.Fail("build submission result", err)
	}
	actualOutput, err := exam_factory.NewIOVariable("sum", exam_consts.VariableFormatInt, "5")
	if err != nil {
		process.Fail("build actual output io variable", err)
	}
	actualOutput, err = ioVariableRepository.CreateIOVariable(ctx, actualOutput)
	if err != nil {
		process.Fail("persist actual output io variable", err)
	}
	result.ActualOutput = actualOutput
	result.Status = submission_consts.SubmissionStatusAccepted

	createdResult, err := resultRepository.CreateResult(ctx, result)
	if err != nil {
		process.Fail("create result", err)
	}
	if createdResult == nil || createdResult.ID == "" {
		process.Fail("create result", fmt.Errorf("expected created result with ID"))
	}
	resultID := createdResult.ID
	process.Log(fmt.Sprintf("SubmissionResult creado. resultID=%s status=%s", resultID, createdResult.Status))

	robleResult, err := robleDirectResult.GetResultByID(ctx, resultID)
	if err != nil || robleResult == nil {
		process.Fail("roble verify create result", fmt.Errorf("expected result persisted in Roble, err=%v", err))
	}
	process.Log(fmt.Sprintf("Roble verify: result persisted id=%s", robleResult.ID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando result %s", resultID)
		_ = resultRepository.DeleteResult(ctx, resultID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, actualOutput.ID)
	}()

	// [STEP 15] GetResultByID — verify nested ActualOutput round-trip through cache
	process.StartStep("Leer SubmissionResult desde cache — verificar ActualOutput anidado")
	cachedResult, err := resultRepository.GetResultByID(ctx, resultID)
	if err != nil {
		process.Fail("get result by id", err)
	}
	if cachedResult == nil {
		process.Fail("get result by id", fmt.Errorf("expected result from cache"))
	}
	if cachedResult.ActualOutput == nil || cachedResult.ActualOutput.ID == "" {
		process.Fail("result actualOutput round-trip", fmt.Errorf("expected ActualOutput with ID after cache deserialization"))
	}
	process.Log(fmt.Sprintf("Result cache hit: status=%s actualOutputID=%s", cachedResult.Status, cachedResult.ActualOutput.ID))
	process.EndStep()

	// [STEP 16] Update result
	process.StartStep("Actualizar SubmissionResult")
	createdResult.Status = submission_consts.SubmissionStatusWrongAnswer
	updatedResult, err := resultRepository.UpdateResult(ctx, createdResult)
	if err != nil {
		process.Fail("update result", err)
	}
	if updatedResult == nil || updatedResult.Status != submission_consts.SubmissionStatusWrongAnswer {
		process.Fail("update result", fmt.Errorf("expected updated result status"))
	}
	process.Log(fmt.Sprintf("Result actualizado. status=%s", updatedResult.Status))

	robleUpdatedResult, err := robleDirectResult.GetResultByID(ctx, resultID)
	if err != nil || robleUpdatedResult == nil || robleUpdatedResult.Status != submission_consts.SubmissionStatusWrongAnswer {
		process.Fail("roble verify update result", fmt.Errorf("expected updated result in Roble"))
	}
	process.Log(fmt.Sprintf("Roble verify: result updated status=%s", robleUpdatedResult.Status))
	process.EndStep()

	// [STEP 17] Result list queries
	process.StartStep("Verificar queries de result por submission y test case")
	bySubmission, err := resultRepository.GetResultsBySubmissionID(ctx, submissionID)
	if err != nil {
		process.Fail("get results by submission", err)
	}
	byTestCase, err := resultRepository.GetResultByTestCase(ctx, testCaseID)
	if err != nil {
		process.Fail("get results by test case", err)
	}
	if len(bySubmission) == 0 || len(byTestCase) == 0 {
		process.Fail("result queries", fmt.Errorf("expected results in query outputs"))
	}
	process.Log(fmt.Sprintf("bySubmission=%d byTestCase=%d", len(bySubmission), len(byTestCase)))
	process.EndStep()

	// [STEP 18] Delete result, submission and session
	process.StartStep("Eliminar Result, Submission y Session — verificar cache invalidation")
	if err := resultRepository.DeleteResult(ctx, resultID); err != nil {
		process.Fail("delete result", err)
	}
	deletedResult, err := resultRepository.GetResultByID(ctx, resultID)
	if err != nil {
		process.Fail("get result after delete", err)
	}
	if deletedResult != nil {
		process.Fail("result cache invalidation", fmt.Errorf("expected result deleted"))
	}
	robleDeletedResult, err := robleDirectResult.GetResultByID(ctx, resultID)
	if err != nil {
		process.Fail("roble verify delete result", err)
	}
	if robleDeletedResult != nil {
		process.Fail("roble verify delete result", fmt.Errorf("expected result %s deleted from Roble", resultID))
	}
	process.Log("Roble verify: result deleted from Roble")
	resultID = ""

	if err := submissionRepository.DeleteSubmission(ctx, submissionID); err != nil {
		process.Fail("delete submission", err)
	}
	deletedSubmission, err := submissionRepository.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		process.Fail("get submission after delete", err)
	}
	if deletedSubmission != nil {
		process.Fail("submission cache invalidation", fmt.Errorf("expected submission deleted"))
	}
	robleDeletedSub, err := robleDirectSubmission.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		process.Fail("roble verify delete submission", err)
	}
	if robleDeletedSub != nil {
		process.Fail("roble verify delete submission", fmt.Errorf("expected submission %s deleted from Roble", submissionID))
	}
	process.Log("Roble verify: submission deleted from Roble")
	submissionID = ""

	if err := sessionRepository.DeleteSession(ctx, sessionID); err != nil {
		process.Fail("delete session", err)
	}
	deletedSession, err := sessionRepository.GetSessionByID(ctx, sessionID)
	if err != nil {
		process.Fail("get session after delete", err)
	}
	if deletedSession != nil {
		process.Fail("session cache invalidation", fmt.Errorf("expected session deleted"))
	}
	robleDeletedSession, err := robleDirectSession.GetSessionByID(ctx, sessionID)
	if err != nil {
		process.Fail("roble verify delete session", err)
	}
	if robleDeletedSession != nil {
		process.Fail("roble verify delete session", fmt.Errorf("expected session %s deleted from Roble", sessionID))
	}
	process.Log("Roble verify: session deleted from Roble")
	sessionID = ""
	process.Log("Result, Submission y Session eliminados")
	process.EndStep()

	process.End()
}
