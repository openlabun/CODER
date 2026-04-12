package roble_persistence_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
	hasher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/security"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submission_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	course_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	exam_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
	submission_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/submission"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	roble_course_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/course"
	roble_exam_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
	roble_submission_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"
	roble_user_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
	test "github.com/openlabun/CODER/apps/api_v2/test"
)

func TestSubmissionCRUD(t *testing.T) {
	process := test.StartTest(t, "Submission Creation and Persistence")

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
	sessionRepository := roble_submission_infrastructure.NewSessionRepository(robleAdapter)
	submissionRepository := roble_submission_infrastructure.NewSubmissionRepository(robleAdapter)
	resultRepository := roble_submission_infrastructure.NewSubmissionResultRepository(robleAdapter, ioVariableRepository)
	process.Log("Repositories inicializados")
	process.EndStep()

	// [STEP 2] Initialize hasher and hash password
	process.StartStep("Inicializar hasher y hashear password")
	adapter := hasher.NewSecurityAdapter()
	hashedPassword, err := adapter.Hash("Password123!")
	if err != nil {
		process.Fail("hash password", err)
	}
	process.EndStep()

	// [STEP 3] Login Teacher user
	process.StartStep("Login docente de pruebas")
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
	enrollmentCode := fmt.Sprintf("ENR-SUB-%d", now.UnixNano())
	courseCode := fmt.Sprintf("SUB-%d", now.Unix()%100000)

	// [STEP 4] Create course
	process.StartStep("Crear curso auxiliar")
	course, err := course_factory.NewCourse(
		"Submission Integration Course",
		"Course for submission/session/result integration",
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
	process.Log(fmt.Sprintf("Curso creado. courseID=%s", courseID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando curso %s", courseID)
		_ = courseRepository.DeleteCourse(ctx, courseID)
	}()

	// [STEP 5] Create auxiliary exam
	process.StartStep("Crear examen auxiliar")
	startTime := time.Now().UTC()
	endTime := startTime.Add(2 * time.Hour)
	exam, err := exam_factory.NewExam(
		"Submission Integration Exam",
		"Exam for submission/session/result integration",
		exam_entities.VisibilityCourse,
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

	// [STEP 6] Create auxiliary challenge, examitem and test case
	process.StartStep("Crear challenge y test case auxiliares")
	inputA, err := exam_factory.NewIOVariable("a", exam_entities.VariableFormatInt, "2")
	if err != nil {
		process.Fail("create challenge input variable", err)
	}
	output, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "5")
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
	process.Log("Input y output IOVariables creadas")

	challenge, err := exam_factory.NewChallenge(
		"Sum Challenge Submission",
		"Challenge for submission integration",
		[]string{"submission", "integration"},
		exam_entities.ChallengeStatusDraft,
		exam_entities.ChallengeDifficultyEasy,
		1500,
		256,
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
	if createdChallenge == nil || createdChallenge.ID == "" {
		process.Fail("create challenge", fmt.Errorf("expected created challenge with ID"))
	}
	challengeID := createdChallenge.ID
	process.Log(fmt.Sprintf("Challenge creado: challengeID=%s", challengeID))

	defer func() {
		t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
		_ = challengeRepository.DeleteChallenge(ctx, challengeID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, inputA.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, output.ID)
	}()

	examItem, err := exam_factory.NewExamItem(examID, challengeID, 1, 100)
	if err != nil {
		process.Fail("create exam item", err)
	}
	createdExamItem, err := examItemRepository.CreateExamItem(ctx, examItem)
	if err != nil {
		process.Fail("persist exam item", err)
	}
	process.Log(fmt.Sprintf("ExamItem creado: examID=%s y challengeID=%s", createdExamItem.ExamID, createdExamItem.ChallengeID))

	defer func() {
		t.Logf("[CLEANUP] Eliminando exam item del challenge %s en el examen %s", challengeID, examID)
		_ = examItemRepository.DeleteExamItem(ctx, createdExamItem.ID)
	}()


	tcInput, err := exam_factory.NewIOVariable("a", exam_entities.VariableFormatInt, "10")
	if err != nil {
		process.Fail("create test case input variable", err)
	}
	tcOutput, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "10")
	if err != nil {
		process.Fail("create test case output variable", err)
	}
	tcInput, err = ioVariableRepository.CreateIOVariable(ctx, tcInput)
	if err != nil {
		process.Fail("persist test case input variable", err)
	}
	tcOutput, err = ioVariableRepository.CreateIOVariable(ctx, tcOutput)
	if err != nil {
		process.Fail("persist test case output variable", err)
	}
	process.Log("Input y output IOVariables para test case creadas")

	testCase, err := exam_factory.NewTestCase(
		"sample_submission",
		[]exam_entities.IOVariable{*tcInput},
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
	process.Log(fmt.Sprintf("Test case creado: testCaseID=%s", testCaseID))
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando test case %s", testCaseID)
		_ = testCaseRepository.DeleteTestCase(ctx, testCaseID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, tcInput.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, tcOutput.ID)
	}()

	// [STEP 7] Create session and confirm creation
	process.StartStep("Creación de la sesión para el examen")
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
	process.EndStep()

	defer func() {
		t.Logf("[CLEANUP] Eliminando session %s", sessionID)
		_ = sessionRepository.DeleteSession(ctx, sessionID)
	}()

	// [STEP 8] Get session by ID, update it and confirm its active
	process.StartStep("Confirmar sesión activa")
	reloadedSession, err := sessionRepository.GetSessionByID(ctx, sessionID)
	if err != nil {
		process.Fail("get session by id", err)
	}
	if reloadedSession == nil {
		process.Fail("get session by id", fmt.Errorf("expected reloaded session"))
	}
	process.Log(fmt.Sprintf("Session recargada. status=%s", reloadedSession.Status))

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

	sessionsByExam, err := sessionRepository.GetSessionsByExamID(ctx, examID)
	if err != nil {
		process.Fail("get sessions by exam", err)
	}
	sessionsByStudent, err := sessionRepository.GetSessionsByStudentID(ctx, teacherID)
	if err != nil {
		process.Fail("get sessions by student", err)
	}
	if len(sessionsByExam) == 0 || len(sessionsByStudent) == 0 {
		process.Fail("session queries", fmt.Errorf("expected sessions in query results"))
	}
	process.Log(fmt.Sprintf("Queries de session validadas. byExam=%d byStudent=%d", len(sessionsByExam), len(sessionsByStudent)))
	process.EndStep()

	// [STEP 9] Create submission
	process.StartStep("Crear Submission")
	submission, err := submission_factory.NewSubmission(
		"print(1+2)",
		submission_entities.LanguagePython,
		challengeID,
		sessionID,
		teacherID,
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
	process.Log(fmt.Sprintf("Submission creada. submissionID=%s", submissionID))
	process.EndStep()
	defer func() {
		t.Logf("[CLEANUP] Eliminando submission %s", submissionID)
		_ = submissionRepository.DeleteSubmission(ctx, submissionID)	
	}()
	

	// [STEP 10] Get submission by ID, update it and confirm changes
	process.StartStep("Cargar y Actualizar Submission")
	reloadedSubmission, err := submissionRepository.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		process.Fail("get submission by id", err)
	}
	if reloadedSubmission == nil {
		process.Fail("get submission by id", fmt.Errorf("expected reloaded submission"))
	}
	process.Log(fmt.Sprintf("Submission recargada. language=%s", reloadedSubmission.Language))

	createdSubmission.Code = "print(2+3+0)"
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
	process.EndStep()

	// [STEP 11] Get submissions by session
	process.StartStep("Obtener Submission por sesión")
	bySession, err := submissionRepository.GetSubmissionsBySessionID(ctx, sessionID, nil, nil, nil)
	if err != nil {
		process.Fail("get submissions by session", err)
	}
	process.EndStep()

	// [STEP 12] Get submissions by user
	process.StartStep("Obtener Submission por usuario")
	byUser, err := submissionRepository.GetSubmissionsByUserID(ctx, teacherID, nil, nil, nil)
	if err != nil {
		process.Fail("get submissions by user", err)
	}
	process.EndStep()

	// [STEP 13] Get submission by challenge
	process.StartStep("Obtener Submission por challenge")
	byChallenge, err := submissionRepository.GetSubmissionsByChallengeID(ctx, challengeID, nil, nil)
	if err != nil {
		process.Fail("get submissions by challenge", err)
	}
	process.EndStep()

	// [STEP 14] Validate submission queries
	process.StartStep("Validar queries de Submission")
	if len(bySession) == 0 || len(byUser) == 0 || len(byChallenge) == 0 {
		process.Fail("submission queries", fmt.Errorf("expected submissions in query results"))
	}
	process.Log(fmt.Sprintf("Queries de submission validadas. bySession=%d byUser=%d byChallenge=%d", len(bySession), len(byUser), len(byChallenge)))
	process.EndStep()

	// [STEP 15] CRUD submission result
	process.StartStep("CRUD de SubmissionResult")
	result, err := submission_factory.NewSubmissionResult(submissionID, testCaseID)
	if err != nil {
		process.Fail("build submission result with factory", err)
	}

	actualOutput, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "5")
	if err != nil {
		process.Fail("build actual output io variable", err)
	}
	actualOutput, err = ioVariableRepository.CreateIOVariable(ctx, actualOutput)
	if err != nil {
		process.Fail("persist actual output io variable", err)
	}
	process.Log("Actual output IOVariable creada y persistida")

	result.ActualOutput = actualOutput
	result.Status = submission_entities.SubmissionStatusAccepted
	
	createdResult, err := resultRepository.CreateResult(ctx, result)
	if err != nil {
		process.Fail("create result", err)
	}
	if createdResult == nil || createdResult.ID == "" {
		process.Fail("create result", fmt.Errorf("expected created result with ID"))
	}
	resultID := createdResult.ID
	process.Log(fmt.Sprintf("SubmissionResult creado. resultID=%s", resultID))
	defer func() {
		t.Logf("[CLEANUP] Eliminando result %s", resultID)
		_ = resultRepository.DeleteResult(ctx, resultID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, actualOutput.ID)
	}()

	reloadedResult, err := resultRepository.GetResultByID(ctx, resultID)
	if err != nil {
		process.Fail("get result by id", err)
	}
	if reloadedResult == nil || reloadedResult.ActualOutput == nil {
		process.Fail("get result by id", fmt.Errorf("expected reloaded result with actual output"))
	}
	process.Log(fmt.Sprintf("Result recargado. status=%s ActualOutput=%s", reloadedResult.Status, reloadedResult.ActualOutput.ID))

	updatedActualOutput, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "4")
	if err != nil {
		process.Fail("build updated actual output", err)
	}
	updatedActualOutput, err = ioVariableRepository.CreateIOVariable(ctx, updatedActualOutput)
	if err != nil {
		process.Fail("persist updated actual output", err)
	}
	process.Log("Updated actual output IOVariable creada y persistida")

	defer func () {
		t.Logf("[CLEANUP] Eliminando updated actual output io variable %s", updatedActualOutput.ID)
		_ = ioVariableRepository.DeleteIOVariable(ctx, updatedActualOutput.ID)
	} ()

	createdResult.Status = submission_entities.SubmissionStatusWrongAnswer
	createdResult.ActualOutput = updatedActualOutput
	createdResult.ErrorMessage = nil

	updatedResult, err := resultRepository.UpdateResult(ctx, createdResult)
	if err != nil {
		process.Fail("update result", err)
	}
	if updatedResult == nil || updatedResult.Status != submission_entities.SubmissionStatusWrongAnswer {
		process.Fail("update result", fmt.Errorf("expected updated result status"))
	}
	process.Log(fmt.Sprintf("Result actualizado. status=%s", updatedResult.Status))

	resultsBySubmission, err := resultRepository.GetResultsBySubmissionID(ctx, submissionID)
	if err != nil {
		process.Fail("get results by submission", err)
	}
	resultsByTestCase, err := resultRepository.GetResultByTestCase(ctx, testCaseID)
	if err != nil {
		process.Fail("get results by test case", err)
	}
	if len(resultsBySubmission) == 0 || len(resultsByTestCase) == 0 {
		process.Fail("result queries", fmt.Errorf("expected results in query outputs"))
	}
	process.Log(fmt.Sprintf("Queries de result validadas. bySubmission=%d byTestCase=%d", len(resultsBySubmission), len(resultsByTestCase)))
	process.EndStep()


	process.StartStep("Eliminar Result, Submission y Session y validar borrado")
	if err := resultRepository.DeleteResult(ctx, resultID); err != nil {
		process.Fail("delete result", err)
	}
	deletedResult, err := resultRepository.GetResultByID(ctx, resultID)
	if err != nil {
		process.Fail("get result after delete", err)
	}
	if deletedResult != nil {
		process.Fail("get result after delete", fmt.Errorf("expected result deleted"))
	}
	resultID = ""
	process.Log("Result eliminado")

	if err := submissionRepository.DeleteSubmission(ctx, submissionID); err != nil {
		process.Fail("delete submission", err)
	}
	deletedSubmission, err := submissionRepository.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		process.Fail("get submission after delete", err)
	}
	if deletedSubmission != nil {
		process.Fail("get submission after delete", fmt.Errorf("expected submission deleted"))
	}
	submissionID = ""
	process.Log("Submission eliminada")

	if err := sessionRepository.DeleteSession(ctx, sessionID); err != nil {
		process.Fail("delete session", err)
	}
	deletedSession, err := sessionRepository.GetSessionByID(ctx, sessionID)
	if err != nil {
		process.Fail("get session after delete", err)
	}
	if deletedSession != nil {
		process.Fail("get session after delete", fmt.Errorf("expected session deleted"))
	}
	sessionID = ""
	process.Log("Session eliminada")
	process.EndStep()

	process.End()
}
