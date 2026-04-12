package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestSubmissionCreateAndRead(t *testing.T) {
	process := test.StartTestWithApp(t, "Submission Create and Read")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var teacherID string
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

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	teacherID = teacherAccess.UserData.ID
	process.EndStep()

	// [STEP 2] Create a public exam (visibility public and no course)
	process.StartStep("Crear examen público (visibilidad public y sin curso)")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Submission Create Exam",
		Description:          "Examen para creación y consulta de revisiones",
		Visibility:           string(exam_consts.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            3600,
		TryLimit:             2,
		ProfessorID:          teacherID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	examID = createdExam.ID
	process.EndStep()

	// [STEP 3] Create a challenge
	process.StartStep("Crear un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Submission Create Challenge",
		Description:       "Challenge para pruebas de submissions",
		Tags:              []string{"submission", "create"},
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

	// [STEP 4] Create a test case for the challenge
	process.StartStep("Crear casos de prueba")
	createdTestCase, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "sample_submission_case",
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

	// [STEP 5] Create an exam item for the challenge
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

	// [STEP 6] Login as student
	process.StartStep("Iniciar sesión con usuario de estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, studentEmail, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	studentID = studentAccess.UserData.ID
	process.EndStep()

	// [STEP 7] Create a session with the exam
	process.StartStep("Crear sesión con examen")
	createdSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create student session", err)
	}
	sessionID = createdSession.ID
	process.EndStep()

	// [STEP 8] Create a submission for the challenge in the session
	process.StartStep("Crear una revisión")
	createdSubmission, err := process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "def solve(n):\n    return n",
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

	// [STEP 9] Retrieve the submission and verify its contents
	process.StartStep("Obtener revisiones a partir del ID del reto")
	challengeSubmissions, err := process.Application.SubmissionUseCases.GetChallengeSubmissions.Execute(teacherCtx, submission_dtos.GetChallengeSubmissionsInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("get challenge submissions", err)
	}
	if len(challengeSubmissions) == 0 {
		process.Fail("get challenge submissions", fmt.Errorf("expected at least one submission for challenge"))
	}
	process.EndStep()

	// [STEP 10] Retrieve submissions by session and by user
	process.StartStep("Obtener revisiones a partir del ID de la sesión")
	sessionSubmissions, err := process.Application.SubmissionUseCases.GetSessionSubmissions.Execute(teacherCtx, submission_dtos.GetSessionSubmissionsInput{SessionID: sessionID})
	if err != nil {
		process.Fail("get session submissions", err)
	}
	if len(sessionSubmissions) == 0 {
		process.Fail("get session submissions", fmt.Errorf("expected at least one submission for session"))
	}
	process.EndStep()

	// [STEP 11] Retrieve submissions by user ID
	process.StartStep("Obtener revisiones a partir del ID del usuario")
	userSubmissions, err := process.Application.SubmissionUseCases.GetUserSubmissions.Execute(teacherCtx, submission_dtos.GetUserSubmissionsInput{UserID: teacherID})
	if err != nil {
		process.Fail("get user submissions", err)
	}
	if userSubmissions == nil {
		process.Fail("get user submissions", fmt.Errorf("expected non-nil user submissions slice"))
	}
	process.EndStep()

	// [STEP 12] Retrieve the status of the submission
	process.StartStep("Obtener el status de la revisión")
	statusOutput, err := process.Application.SubmissionUseCases.GetSubmissionStatus.Execute(studentCtx, submission_dtos.GetSubmissionStatusInput{SubmissionID: submissionID})
	if err != nil {
		process.Fail("get submission status", err)
	}
	if statusOutput == nil || statusOutput.Submission.ID != submissionID {
		process.Fail("get submission status", fmt.Errorf("expected status for submission %s", submissionID))
	}
	process.EndStep()

	process.End()
}
