package roble_persistence_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

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
)

func TestSubmissionCRUD(t *testing.T) {
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
	challengeRepository := roble_exam_infrastructure.NewChallengeRepository(robleAdapter)
	testCaseRepository := roble_exam_infrastructure.NewTestCaseRepository(robleAdapter)
	sessionRepository := roble_submission_infrastructure.NewSessionRepository(robleAdapter)
	submissionRepository := roble_submission_infrastructure.NewSubmissionRepository(robleAdapter)
	resultRepository := roble_submission_infrastructure.NewSubmissionResultRepository(robleAdapter)
	t.Log("[OK] Repositories inicializados")

	t.Log("[STEP 2] Login docente de pruebas")
	access, err := authAdapter.LoginUser("test@test.com", "Testing123!")
	if err != nil {
		t.Fatalf("teacher login failed: %v", err)
	}
	if access == nil || access.UserData == nil || access.UserData.ID == "" {
		t.Fatal("expected logged user data with valid ID")
	}
	if access.Token == nil || access.Token.AccessToken == "" {
		t.Fatal("expected access token in login response")
	}
	ctx := roble_infrastructure.WithAccessToken(context.Background(), access.Token.AccessToken)
	teacherID := access.UserData.ID
	t.Logf("[OK] Login exitoso. userID=%s", teacherID)

	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-SUB-%d", now.UnixNano())
	courseCode := fmt.Sprintf("SUB-%d", now.Unix()%100000)

	t.Log("[STEP 3] Crear curso auxiliar")
	course, err := course_factory.NewCourse(
		"Submission Integration Course",
		"Course for submission/session/result integration",
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
	t.Logf("[OK] Curso creado. courseID=%s", courseID)
	defer func() {
		if courseID != "" {
			t.Logf("[CLEANUP] Eliminando curso %s", courseID)
			_ = courseRepository.DeleteCourse(ctx, courseID)
		}
	}()

	t.Log("[STEP 4] Crear examen auxiliar")
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
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_ = examRepository.DeleteExam(ctx, examID)
		}
	}()

	t.Log("[STEP 5] Crear challenge y test case auxiliares")
	inputA, err := exam_factory.NewIOVariable("a", exam_entities.VariableFormatInt, "2")
	if err != nil {
		t.Fatalf("create challenge input variable failed: %v", err)
	}
	output, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "5")
	if err != nil {
		t.Fatalf("create challenge output variable failed: %v", err)
	}

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
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = challengeRepository.DeleteChallenge(ctx, challengeID)
		}
	}()

	tcInput, err := exam_factory.NewIOVariable("a", exam_entities.VariableFormatInt, "10")
	if err != nil {
		t.Fatalf("create test case input variable failed: %v", err)
	}
	tcOutput, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "10")
	if err != nil {
		t.Fatalf("create test case output variable failed: %v", err)
	}

	testCase, err := exam_factory.NewTestCase(
		"sample_submission",
		[]exam_entities.IOVariable{*tcInput},
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
	t.Logf("[OK] Test case creado. testCaseID=%s", testCaseID)
	defer func() {
		if testCaseID != "" {
			t.Logf("[CLEANUP] Eliminando test case %s", testCaseID)
			_ = testCaseRepository.DeleteTestCase(ctx, testCaseID)
		}
	}()

	t.Log("[STEP 6] CRUD Session")
	session, err := submission_factory.NewSession(teacherID, createdExam)
	if err != nil {
		t.Fatalf("build session with factory failed: %v", err)
	}
	session.TimeLeft = 5400

	createdSession, err := sessionRepository.CreateSession(ctx, session)
	if err != nil {
		t.Fatalf("create session failed: %v", err)
	}
	if createdSession == nil || createdSession.ID == "" {
		t.Fatal("expected created session with ID")
	}
	sessionID := createdSession.ID
	t.Logf("[OK] Session creada. sessionID=%s", sessionID)

	defer func() {
		if sessionID != "" {
			t.Logf("[CLEANUP] Eliminando session %s", sessionID)
			_ = sessionRepository.DeleteSession(ctx, sessionID)
		}
	}()

	reloadedSession, err := sessionRepository.GetSessionByID(ctx, sessionID)
	if err != nil {
		t.Fatalf("get session by id failed: %v", err)
	}
	if reloadedSession == nil {
		t.Fatal("expected reloaded session")
	}
	t.Logf("[OK] Session recargada. status=%s", reloadedSession.Status)

	createdSession.Status = submission_entities.SessionStatusFrozen
	createdSession.Attempts = 1
	createdSession.TimeLeft = 5000
	updatedSession, err := sessionRepository.UpdateSession(ctx, createdSession)
	if err != nil {
		t.Fatalf("update session failed: %v", err)
	}
	if updatedSession == nil || updatedSession.Status != submission_entities.SessionStatusFrozen {
		t.Fatal("expected updated session status")
	}
	t.Logf("[OK] Session actualizada. attempts=%d timeLeft=%d", updatedSession.Attempts, updatedSession.TimeLeft)

	sessionsByExam, err := sessionRepository.GetSessionsByExamID(ctx, examID)
	if err != nil {
		t.Fatalf("get sessions by exam failed: %v", err)
	}
	sessionsByStudent, err := sessionRepository.GetSessionsByStudentID(ctx, teacherID)
	if err != nil {
		t.Fatalf("get sessions by student failed: %v", err)
	}
	if len(sessionsByExam) == 0 || len(sessionsByStudent) == 0 {
		t.Fatal("expected sessions in query results")
	}
	t.Logf("[OK] Queries de session validadas. byExam=%d byStudent=%d", len(sessionsByExam), len(sessionsByStudent))

	t.Log("[STEP 7] CRUD Submission")
	submission, err := submission_factory.NewSubmission(
		"print(2+3)",
		submission_entities.LanguagePython,
		challengeID,
		sessionID,
		teacherID,
	)
	if err != nil {
		t.Fatalf("build submission with factory failed: %v", err)
	}

	createdSubmission, err := submissionRepository.CreateSubmission(ctx, submission)
	if err != nil {
		t.Fatalf("create submission failed: %v", err)
	}
	if createdSubmission == nil || createdSubmission.ID == "" {
		t.Fatal("expected created submission with ID")
	}
	submissionID := createdSubmission.ID
	t.Logf("[OK] Submission creada. submissionID=%s", submissionID)
	defer func() {
		if submissionID != "" {
			t.Logf("[CLEANUP] Eliminando submission %s", submissionID)
			_ = submissionRepository.DeleteSubmission(ctx, submissionID)
		}
	}()

	reloadedSubmission, err := submissionRepository.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		t.Fatalf("get submission by id failed: %v", err)
	}
	if reloadedSubmission == nil {
		t.Fatal("expected reloaded submission")
	}
	t.Logf("[OK] Submission recargada. language=%s", reloadedSubmission.Language)

	createdSubmission.Code = "print(2+3+0)"
	createdSubmission.Score = 100
	createdSubmission.TimeMsTotal = 45
	updatedSubmission, err := submissionRepository.UpdateSubmission(ctx, createdSubmission)
	if err != nil {
		t.Fatalf("update submission failed: %v", err)
	}
	if updatedSubmission == nil || updatedSubmission.Score != 100 {
		t.Fatal("expected updated submission score")
	}
	t.Logf("[OK] Submission actualizada. score=%d timeMsTotal=%d", updatedSubmission.Score, updatedSubmission.TimeMsTotal)

	bySession, err := submissionRepository.GetSubmissionsBySessionID(ctx, sessionID)
	if err != nil {
		t.Fatalf("get submissions by session failed: %v", err)
	}
	byUser, err := submissionRepository.GetSubmissionsByUserID(ctx, teacherID)
	if err != nil {
		t.Fatalf("get submissions by user failed: %v", err)
	}
	byChallenge, err := submissionRepository.GetSubmissionsByChallengeID(ctx, challengeID)
	if err != nil {
		t.Fatalf("get submissions by challenge failed: %v", err)
	}
	if len(bySession) == 0 || len(byUser) == 0 || len(byChallenge) == 0 {
		t.Fatal("expected submissions in query results")
	}
	t.Logf("[OK] Queries de submission validadas. bySession=%d byUser=%d byChallenge=%d", len(bySession), len(byUser), len(byChallenge))

	t.Log("[STEP 8] CRUD SubmissionResult")
	result, err := submission_factory.NewSubmissionResult(submissionID, testCaseID)
	if err != nil {
		t.Fatalf("build submission result with factory failed: %v", err)
	}

	actualOutput, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "5")
	if err != nil {
		t.Fatalf("build actual output io variable failed: %v", err)
	}
	result.Status = submission_entities.SubmissionStatusAccepted
	result.ActualOutput = actualOutput

	createdResult, err := resultRepository.CreateResult(ctx, result)
	if err != nil {
		t.Fatalf("create result failed: %v", err)
	}
	if createdResult == nil || createdResult.ID == "" {
		t.Fatal("expected created result with ID")
	}
	resultID := createdResult.ID
	t.Logf("[OK] SubmissionResult creado. resultID=%s", resultID)
	defer func() {
		if resultID != "" {
			t.Logf("[CLEANUP] Eliminando result %s", resultID)
			_ = resultRepository.DeleteResult(ctx, resultID)
		}
	}()

	reloadedResult, err := resultRepository.GetResultByID(ctx, resultID)
	if err != nil {
		t.Fatalf("get result by id failed: %v", err)
	}
	if reloadedResult == nil || reloadedResult.ActualOutput == nil {
		t.Fatal("expected reloaded result with actual output")
	}
	t.Logf("[OK] Result recargado. status=%s ActualOutput=%s", reloadedResult.Status, reloadedResult.ActualOutput.ID)

	updatedActualOutput, err := exam_factory.NewIOVariable("sum", exam_entities.VariableFormatInt, "4")
	if err != nil {
		t.Fatalf("build updated actual output failed: %v", err)
	}
	createdResult.Status = submission_entities.SubmissionStatusWrongAnswer
	createdResult.ActualOutput = updatedActualOutput
	createdResult.ErrorMessage = nil

	updatedResult, err := resultRepository.UpdateResult(ctx, createdResult)
	if err != nil {
		t.Fatalf("update result failed: %v", err)
	}
	if updatedResult == nil || updatedResult.Status != submission_entities.SubmissionStatusWrongAnswer {
		t.Fatal("expected updated result status")
	}
	t.Logf("[OK] Result actualizado. status=%s", updatedResult.Status)

	resultsBySubmission, err := resultRepository.GetResultsBySubmissionID(ctx, submissionID)
	if err != nil {
		t.Fatalf("get results by submission failed: %v", err)
	}
	resultsByTestCase, err := resultRepository.GetResultByTestCase(ctx, testCaseID)
	if err != nil {
		t.Fatalf("get results by test case failed: %v", err)
	}
	if len(resultsBySubmission) == 0 || len(resultsByTestCase) == 0 {
		t.Fatal("expected results in query outputs")
	}
	t.Logf("[OK] Queries de result validadas. bySubmission=%d byTestCase=%d", len(resultsBySubmission), len(resultsByTestCase))

	t.Log("[STEP 9] Eliminar Result, Submission y Session y validar borrado")
	if err := resultRepository.DeleteResult(ctx, resultID); err != nil {
		t.Fatalf("delete result failed: %v", err)
	}
	deletedResult, err := resultRepository.GetResultByID(ctx, resultID)
	if err != nil {
		t.Fatalf("get result after delete failed: %v", err)
	}
	if deletedResult != nil {
		t.Fatal("expected result deleted")
	}
	resultID = ""
	t.Log("[OK] Result eliminado")

	if err := submissionRepository.DeleteSubmission(ctx, submissionID); err != nil {
		t.Fatalf("delete submission failed: %v", err)
	}
	deletedSubmission, err := submissionRepository.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		t.Fatalf("get submission after delete failed: %v", err)
	}
	if deletedSubmission != nil {
		t.Fatal("expected submission deleted")
	}
	submissionID = ""
	t.Log("[OK] Submission eliminada")

	if err := sessionRepository.DeleteSession(ctx, sessionID); err != nil {
		t.Fatalf("delete session failed: %v", err)
	}
	deletedSession, err := sessionRepository.GetSessionByID(ctx, sessionID)
	if err != nil {
		t.Fatalf("get session after delete failed: %v", err)
	}
	if deletedSession != nil {
		t.Fatal("expected session deleted")
	}
	sessionID = ""
	t.Log("[OK] Session eliminada")
}
