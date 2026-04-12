package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	submission_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestSubmissionExecution(t *testing.T) {
	process := test.StartTestWithApp(t, "Submission Execution")
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
	var sessionID string
	var submissionID string

	defer func() {
		if sessionID != "" {
			t.Logf("[CLEANUP] Cerrando sesión %s", sessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionID})
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

	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.EndStep()

	process.StartStep("Crear examen público (visibilidad public y sin curso)")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Submission Execution Exam",
		Description:          "Examen para validación de ejecución",
		Visibility:           string(exam_consts.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            120,
		TryLimit:             3,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	examID = createdExam.ID
	process.EndStep()

	process.StartStep("Crear un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Submission Execution Challenge",
		Description:       "Challenge para validar ejecución asíncrona",
		Tags:              []string{"submission", "execution"},
		Status:            string(exam_consts.ChallengeStatusPublished),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		Constraints:    "1 <= n <= 1000",
	})
	if err != nil {
		process.Fail("create challenge", err)
	}
	challengeID = createdChallenge.ID
	process.EndStep()

	process.StartStep("Crear casos de prueba")
	createdTestCase, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "execution_case",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "10"},
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

	process.StartStep("Crear una sesión en el examen")
	createdSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{UserID: studentID, ExamID: examID})
	if err != nil {
		process.Fail("create session", err)
	}
	sessionID = createdSession.ID
	process.EndStep()

	process.StartStep("Crear una revisión")
	createdSubmission, err := process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(int(sys.stdin.read().strip()))",
		Language:    "python",
		Score:       0,
		ChallengeID: challengeID,
		SessionID:   sessionID,
	})
	if err != nil {
		process.Fail("create submission", err)
	}
	submissionID = createdSubmission.ID
	process.EndStep()

	process.StartStep("Obtener el status de la revisión hasta que su estado sea accepted")
	deadline := time.Now().Add(2 * time.Minute)
	accepted := false
	for time.Now().Before(deadline) {
		statusOutput, statusErr := process.Application.SubmissionUseCases.GetSubmissionStatus.Execute(studentCtx, submission_dtos.GetSubmissionStatusInput{SubmissionID: submissionID})
		if statusErr != nil {
			process.Log(fmt.Sprintf("Polling status error: %v", statusErr))
			time.Sleep(2 * time.Second)
			continue
		}
		if statusOutput != nil && len(statusOutput.Results) > 0 {
			allAccepted := true
			for _, r := range statusOutput.Results {
				if r.Status != submission_entities.SubmissionStatusAccepted {
					allAccepted = false
					break
				}
			}
			if allAccepted {
				accepted = true
				break
			}
		}
		time.Sleep(2 * time.Second)
	}
	if !accepted {
		process.Fail("wait accepted submission status", fmt.Errorf("submission %s did not reach accepted status before timeout", submissionID))
	}
	process.EndStep()

	process.End()
}
