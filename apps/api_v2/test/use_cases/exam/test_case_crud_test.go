package usecases_test

import (
	"context"
	"fmt"
	"testing"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestTestCaseCRUD(t *testing.T) {
	process := test.StartTestWithApp(t, "TestCase CRUD")
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var challengeID string
	var testCaseID string
	var deletedTestCaseID string

	defer func() {
		if testCaseID != "" {
			t.Logf("[CLEANUP] Eliminando test case %s", testCaseID)
			_ = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: testCaseID})
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

	// [STEP 2] Create a challenge to associate the test cases with
	process.StartStep("Crear un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Challenge for TestCase CRUD",
		Description:       "Challenge auxiliar para test case CRUD",
		Tags:              []string{"testcase", "crud"},
		Status:            string(exam_consts.ChallengeStatusDraft),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "x", Type: string(exam_consts.VariableFormatInt), Value: "1"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "1"},
		Constraints:    "1 <= x <= 1000",
	})
	if err != nil {
		process.Fail("create challenge", err)
	}
	if createdChallenge == nil || createdChallenge.ID == "" {
		process.Fail("create challenge", fmt.Errorf("expected challenge with ID"))
	}
	challengeID = createdChallenge.ID
	process.EndStep()

	// [STEP 3] Try to create a test case without input values (expect error)
	process.StartStep("Crear un caso de uso sin valores de entrada (espera error)")
	_, err = process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name:           "invalid_test_case",
		Input:          []exam_dtos.IOVariableDTO{},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "1"},
		IsSample:       true,
		Points:         0,
		ChallengeID:    challengeID,
	})
	if err == nil {
		process.Fail("create invalid test case", fmt.Errorf("expected error when creating test case without inputs"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 4] Try to create a test case without expected output (expect error)
	process.StartStep("Crear un caso de uso válido")
	createdTestCase, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "sample_valid",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "a", Type: string(exam_consts.VariableFormatInt), Value: "2"},
			{Name: "b", Type: string(exam_consts.VariableFormatInt), Value: "3"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "sum", Type: string(exam_consts.VariableFormatInt), Value: "5"},
		IsSample:       true,
		Points:         10,
		ChallengeID:    challengeID,
	})
	if err != nil {
		process.Fail("create valid test case", err)
	}
	if createdTestCase == nil || createdTestCase.ID == "" {
		process.Fail("create valid test case", fmt.Errorf("expected test case with ID"))
	}
	testCaseID = createdTestCase.ID
	process.EndStep()

	// [STEP 5] Update the test case
	process.StartStep("Actualizar el caso de uso")
	updatedName := "sample_valid_updated"
	updatedPoints := 25
	updatedTestCase, err := process.Application.TestCaseModule.UpdateTestCase.Execute(teacherCtx, exam_dtos.UpdateTestCaseInput{
		ID:     testCaseID,
		Name:   &updatedName,
		Points: &updatedPoints,
	})
	if err != nil {
		process.Fail("update test case", err)
	}
	if updatedTestCase == nil || updatedTestCase.Name != updatedName || updatedTestCase.Points != updatedPoints {
		process.Fail("update test case", fmt.Errorf("expected updated test case values"))
	}
	process.EndStep()

	// [STEP 6] Get test cases by challenge and validate the updated test case is there
	process.StartStep("Eliminar el caso de uso")
	err = process.Application.TestCaseModule.DeleteTestCase.Execute(teacherCtx, exam_dtos.DeleteTestCaseInput{TestCaseID: testCaseID})
	if err != nil {
		process.Fail("delete test case", err)
	}
	deletedTestCaseID = testCaseID
	testCaseID = ""
	process.EndStep()

	// [STEP 7] Verify deletion by trying to get test cases and checking the deleted one is not there
	process.StartStep("Verificar eliminación")
	remaining, err := process.Application.TestCaseModule.GetTestCasesByChallenge.Execute(teacherCtx, exam_dtos.GetTestCasesByChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("get test cases by challenge", err)
	}
	for _, tc := range remaining {
		if tc != nil && tc.ID == deletedTestCaseID {
			process.Fail("verify test case deletion", fmt.Errorf("test case %s should not exist after deletion", deletedTestCaseID))
		}
	}
	process.EndStep()

	process.End()
}
