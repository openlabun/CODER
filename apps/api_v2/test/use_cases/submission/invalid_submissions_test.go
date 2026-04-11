package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	exam_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam/exam_crud"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submission_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestInvalidSubmissions(t *testing.T) {
	process := test.StartTestWithApp(t, "Invalid Submissions")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var studentID string
	var examID string
	var challengeID string
	var testCaseID string
	var examItemID string
	var firstSessionID string
	var secondSessionID string
	var thirdSessionID string

	defer func() {
		if thirdSessionID != "" {
			t.Logf("[CLEANUP] Cerrando tercera sesión %s", thirdSessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: thirdSessionID})
		}
		if secondSessionID != "" {
			t.Logf("[CLEANUP] Cerrando segunda sesión %s", secondSessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: secondSessionID})
		}
		if firstSessionID != "" {
			t.Logf("[CLEANUP] Cerrando primera sesión %s", firstSessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: firstSessionID})
		}
		if examItemID != "" {
			t.Logf("[CLEANUP] Eliminando exam item %s", examItemID)
			_ = process.Application.ExamItemModule.DeleteExamItem.Execute(teacherCtx, exam_dtos.DeleteExamItemInput{ID: examItemID})
		}
		if testCaseID != "" {
			t.Logf("[CLEANUP] Eliminando test case %s", testCaseID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: testCaseID})
		}
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
	}()

	closeExamUC := exam_usecases.NewCloseExamUseCase(
		process.Application.Dependencies.UserRepository,
		process.Application.Dependencies.ExamRepository,
	)

	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.EndStep()

	process.StartStep("Crear examen público (visibilidad public, sin curso y 60 segundos de tiempo para resolver)")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Invalid Submissions Exam",
		Description:          "Examen para validar revisiones inválidas",
		Visibility:           string(exam_entities.VisibilityPublic),
		StartTime:            now.Add(-2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: false,
		TimeLimit:            60,
		TryLimit:             5,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	examID = createdExam.ID
	process.EndStep()

	process.StartStep("Crear un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Invalid Submissions Challenge",
		Description:       "Challenge para escenarios de revisión inválida",
		Tags:              []string{"submission", "invalid"},
		Status:            string(exam_entities.ChallengeStatusPublished),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_entities.VariableFormatInt), Value: "10"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_entities.VariableFormatInt), Value: "10"},
		Constraints:    "1 <= n <= 1000",
	})
	if err != nil {
		process.Fail("create challenge", err)
	}
	challengeID = createdChallenge.ID
	process.EndStep()

	process.StartStep("Crear casos de prueba")
	createdTestCase, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "invalid_submission_case",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_entities.VariableFormatInt), Value: "10"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_entities.VariableFormatInt), Value: "10"},
		IsSample:       true,
		Points:         10,
		ChallengeID:    challengeID,
	})
	if err != nil {
		process.Fail("create test case", err)
	}
	testCaseID = createdTestCase.ID
	process.EndStep()

	process.StartStep("Crear un punto de examen")
	createdExamItem, err := process.Application.ExamItemModule.CreateExamItem.Execute(teacherCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: challengeID,
		Order:       1,
		Points:      100,
	})
	if err != nil {
		process.Fail("create exam item", err)
	}
	examItemID = createdExamItem.ID
	process.EndStep()

	process.StartStep("Iniciar sesión con usuario de estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, studentEmail, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	studentID = studentAccess.UserData.ID
	process.EndStep()

	process.StartStep("Crear una revisión sin sesión (espera error)")
	_, err = process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(sys.stdin.read().strip())",
		Language:    "python",
		Score:       0,
		ChallengeID: challengeID,
		SessionID:   "session-not-found",
	})
	if err == nil {
		process.Fail("create submission without session", fmt.Errorf("expected error when creating submission without valid session"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.StartStep("Crear una sesión en el examen")
	sessionOne, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create first session", err)
	}
	firstSessionID = sessionOne.ID
	process.EndStep()

	process.StartStep("Cerrar el examen desde la vista de docente")
	_, err = closeExamUC.Execute(teacherCtx, exam_dtos.CloseExamInput{ExamID: examID})
	if err != nil {
		process.Fail("close exam", err)
	}
	process.EndStep()

	process.StartStep("Crear una revisión (espera error)")
	_, err = process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(sys.stdin.read().strip())",
		Language:    "python",
		Score:       0,
		ChallengeID: challengeID,
		SessionID:   firstSessionID,
	})
	if err == nil {
		process.Fail("create submission after exam close", fmt.Errorf("expected error after exam close"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.StartStep("Esperar 61 segundos y crear un revisión (espera error)")
	time.Sleep(61 * time.Second)
	_, err = process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(sys.stdin.read().strip())",
		Language:    "python",
		Score:       0,
		ChallengeID: challengeID,
		SessionID:   firstSessionID,
	})
	if err == nil {
		process.Fail("create submission after timeout", fmt.Errorf("expected error after timeout"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.StartStep("Obtener datos de sesión")
	expiredSession, err := process.Application.Dependencies.SessionRepository.GetSessionByID(teacherCtx, firstSessionID)
	if err != nil {
		process.Fail("get session by id", err)
	}
	if expiredSession == nil {
		process.Fail("get session by id", fmt.Errorf("expected session %s", firstSessionID))
	}
	process.EndStep()

	process.StartStep("Confirmar que la sesión tiene estado expired")
	if expiredSession.Status != submission_entities.SessionStatusExpired {
		process.Fail("verify expired session", fmt.Errorf("expected session status expired, got %s", expiredSession.Status))
	}
	firstSessionID = ""
	process.EndStep()

	process.StartStep("Crear una sesión en el examen")
	sessionTwo, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create second session", err)
	}
	secondSessionID = sessionTwo.ID
	process.EndStep()

	process.StartStep("Bloquear la sesión desde la vista de docente")
	blockedSession, err := process.Application.SessionModule.BlockSession.Execute(teacherCtx, submission_dtos.BlockSessionInput{SessionID: secondSessionID})
	if err != nil {
		process.Fail("block second session", err)
	}
	if blockedSession == nil || blockedSession.Status != submission_entities.SessionStatusBlocked {
		process.Fail("block second session", fmt.Errorf("expected blocked status"))
	}
	secondSessionID = ""
	process.EndStep()

	process.StartStep("Crear una revisión (espera error)")
	_, err = process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(sys.stdin.read().strip())",
		Language:    "python",
		Score:       0,
		ChallengeID: challengeID,
		SessionID:   blockedSession.ID,
	})
	if err == nil {
		process.Fail("create submission with blocked session", fmt.Errorf("expected error for blocked session"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.StartStep("Crear una sesión en el examen")
	sessionThree, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create third session", err)
	}
	thirdSessionID = sessionThree.ID
	process.EndStep()

	process.StartStep("Cerrar el examen")
	_, err = closeExamUC.Execute(teacherCtx, exam_dtos.CloseExamInput{ExamID: examID})
	if err != nil {
		process.Fail("close exam again", err)
	}
	process.EndStep()

	process.StartStep("Crear una revisión (espera error)")
	_, err = process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(sys.stdin.read().strip())",
		Language:    "python",
		Score:       0,
		ChallengeID: challengeID,
		SessionID:   thirdSessionID,
	})
	if err == nil {
		process.Fail("create submission after second close", fmt.Errorf("expected error after closing exam again"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	thirdSessionID = ""
	process.EndStep()

	process.End()
}
