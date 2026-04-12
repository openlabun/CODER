package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestTestCaseFromStudentView(t *testing.T) {
	process := test.StartTestWithApp(t, "TestCase Student View")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var challengeID string
	var sampleTestCaseID string
	var privateTestCaseID string
	var examID string
	var examItemID string

	defer func() {
		if examItemID != "" {
			t.Logf("[CLEANUP] Eliminando exam item %s", examItemID)
			_ = process.Application.ExamItemModule.DeleteExamItem.Execute(teacherCtx, exam_dtos.DeleteExamItemInput{ID: examItemID})
		}
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando exam %s", examID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
		if privateTestCaseID != "" {
			t.Logf("[CLEANUP] Eliminando test case privado %s", privateTestCaseID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: privateTestCaseID})
		}
		if sampleTestCaseID != "" {
			t.Logf("[CLEANUP] Eliminando test case sample %s", sampleTestCaseID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: sampleTestCaseID})
		}
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.EndStep()

	// [STEP 2] Login as student
	process.StartStep("Crear un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Challenge for TestCase Student View",
		Description:       "Challenge auxiliar para vista estudiante",
		Tags:              []string{"testcase", "student-view"},
		Status:            string(exam_consts.ChallengeStatusDraft),
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

	// [STEP 3] Create a sample test case and a private test case for the challenge
	process.StartStep("Crear un caso de prueba (isSample == true)")
	sample, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "sample_case",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		IsSample:       true,
		Points:         0,
		ChallengeID:    challengeID,
	})
	if err != nil {
		process.Fail("create sample test case", err)
	}
	sampleTestCaseID = sample.ID
	process.EndStep()

	// [STEP 4] Create a hidden test case (isSample == false)
	process.StartStep("Crear un caso de prueba (isSample == false)")
	hiddenCase, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "hidden_case",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "11"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "11"},
		IsSample:       false,
		Points:         10,
		ChallengeID:    challengeID,
	})
	if err != nil {
		process.Fail("create hidden test case", err)
	}
	privateTestCaseID = hiddenCase.ID
	process.EndStep()

	// [STEP 5] Get test cases by challenge as teacher and validate both test cases are returned
	process.StartStep("Obtener casos de prueba con vista de Docente")
	// Extra setup required by use case contract for student access.
	_, err = process.Application.ChallengeModule.PublishChallenge.Execute(teacherCtx, exam_dtos.PublishChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("publish challenge", err)
	}

	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "TestCase Student View Exam",
		Description:          "Exam público para acceder a test cases",
		Visibility:           string(exam_consts.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            3600,
		TryLimit:             1,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	examID = createdExam.ID

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

	teacherCases, err := process.Application.TestCaseModule.GetTestCasesByChallenge.Execute(teacherCtx, exam_dtos.GetTestCasesByChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("teacher get test cases", err)
	}
	if len(teacherCases) != 2 {
		process.Fail("teacher get test cases", fmt.Errorf("expected 2 test cases, got %d", len(teacherCases)))
	}
	process.EndStep()

	// [STEP 6] Get test cases by challenge as student and validate only the sample test case is returned
	process.StartStep("Iniciar sesión con usuario de estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, studentEmail, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	process.EndStep()

	// [STEP 7] Get test cases by challenge as student and validate only the sample test case is returned
	process.StartStep("Obtener casos de prueba con vista de Estudiante (espera solo 1)")
	studentCases, err := process.Application.TestCaseModule.GetTestCasesByChallenge.Execute(studentCtx, exam_dtos.GetTestCasesByChallengeInput{
		ChallengeID: challengeID,
		ExamID:      &examID,
	})
	if err != nil {
		process.Fail("student get test cases", err)
	}
	if len(studentCases) != 1 {
		process.Fail("student get test cases", fmt.Errorf("expected 1 public sample test case, got %d", len(studentCases)))
	}
	if studentCases[0] == nil || studentCases[0].ID != sampleTestCaseID {
		process.Fail("student get test cases", fmt.Errorf("expected only sample test case %s", sampleTestCaseID))
	}
	process.EndStep()

	process.End()
}
