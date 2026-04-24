package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestSubmissionTryLimit(t *testing.T) {
	process := test.StartTestWithApp(t, "Submission Try Limit")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var studentID string
	var examID string
	var challengeID string
	var testCaseOneID string
	var testCaseTwoID string
	var examItemID string
	var sessionID string

	defer func() {
		if sessionID != "" {
			t.Logf("[CLEANUP] Cerrando sesión %s", sessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionID})
		}
		if examItemID != "" {
			t.Logf("[CLEANUP] Eliminando exam item %s", examItemID)
			_ = process.Application.ExamItemModule.DeleteExamItem.Execute(teacherCtx, exam_dtos.DeleteExamItemInput{ID: examItemID})
		}
		if testCaseTwoID != "" {
			t.Logf("[CLEANUP] Eliminando test case 2 %s", testCaseTwoID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: testCaseTwoID})
		}
		if testCaseOneID != "" {
			t.Logf("[CLEANUP] Eliminando test case 1 %s", testCaseOneID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: testCaseOneID})
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
	process.EndStep()

	// [STEP 2] Create public exam with try_limit = 2
	process.StartStep("Crear examen público")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Submission Try Limit Exam",
		Description:          "Examen para validar límite de intentos por revisión",
		Visibility:           string(exam_consts.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            300,
		TryLimit:             3,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	examID = createdExam.ID
	process.EndStep()

	// [STEP 3] Create a challenge for the exam item
	process.StartStep("Crear un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Submission Try Limit Challenge",
		Description:       "Challenge para validar límite de revisiones",
		Tags:              []string{"submission", "try-limit"},
		Status:            string(exam_consts.ChallengeStatusPublished),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		CodeTemplates: []exam_dtos.CodeTemplateDTO{
			{Language: "python", Template: "def solve():\n    pass"},
		},
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "2"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "4"},
		Constraints:    "1 <= n <= 1000",
	})
	if err != nil {
		process.Fail("create challenge", err)
	}
	challengeID = createdChallenge.ID
	process.EndStep()

	// [STEP 4] Create test cases for the challenge
	process.StartStep("Crear casos de prueba")
	createdTestCaseOne, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "submission_try_limit_case_1",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "2"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "4"},
		IsSample:       false,
		Points:         5,
		ChallengeID:    challengeID,
	})
	if err != nil {
		process.Fail("create test case 1", err)
	}
	testCaseOneID = createdTestCaseOne.ID

	createdTestCaseTwo, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "submission_try_limit_case_2",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "5"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		IsSample:       false,
		Points:         5,
		ChallengeID:    challengeID,
	})
	if err != nil {
		process.Fail("create test case 2", err)
	}
	testCaseTwoID = createdTestCaseTwo.ID
	process.EndStep()

	// [STEP 5] Create an exam item with try_limit = 2 for the challenge
	process.StartStep("Crear un punto de examen (con TryLimit == 2)")
	itemTryLimit := 2
	createdExamItem, err := process.Application.ExamItemModule.CreateExamItem.Execute(teacherCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: challengeID,
		Order:       1,
		Points:      100,
		TryLimit:    &itemTryLimit,
	})
	if err != nil {
		process.Fail("create exam item", err)
	}
	examItemID = createdExamItem.ID
	process.EndStep()

	// [STEP 6] Update the exam to set TryLimit = 2 at exam level
	process.StartStep("Iniciar sesión con usuario de estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, studentEmail, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	studentID = studentAccess.UserData.ID
	process.EndStep()

	// [STEP 7] Create first session for the student
	process.StartStep("Crear una sesión en el examen")
	createdSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create session", err)
	}
	sessionID = createdSession.ID
	process.EndStep()

	// [STEP 8] Create first submission for the session - should succeed
	process.StartStep("Crear una revisión")
	_, err = process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		Language:    "python",
		ChallengeID: challengeID,
		SessionID:   sessionID,
	})
	if err != nil {
		process.Fail("create first submission", err)
	}
	process.EndStep()

	// [STEP 9] Create second submission for the session - should succeed
	process.StartStep("Crear una revisión")
	_, err = process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		Language:    "python",
		ChallengeID: challengeID,
		SessionID:   sessionID,
	})
	if err != nil {
		process.Fail("create second submission", err)
	}
	process.EndStep()

	// [STEP 10] Create third submission for the session - should fail with try limit exceeded error
	process.StartStep("Crear una revisión (espera error)")
	_, err = process.Application.SubmissionUseCases.CreateSubmission.Execute(studentCtx, submission_dtos.CreateSubmissionInput{
		Code:        "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		Language:    "python",
		ChallengeID: challengeID,
		SessionID:   sessionID,
	})
	if err == nil {
		process.Fail("create third submission", fmt.Errorf("expected error when submission try limit is exceeded"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
