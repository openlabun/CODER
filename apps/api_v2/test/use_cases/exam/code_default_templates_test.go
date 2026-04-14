package usecases_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	sub_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestChallengeDefaultCodeTemplates(t *testing.T) {
	process := test.StartTestWithApp(t, "Challenge Default Code Templates")
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var challengeID string

	defer func() {
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesión con usuario de docente (creador)")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.Log(fmt.Sprintf("teacherID=%s", teacherAccess.UserData.ID))
	process.EndStep()

	// [STEP 2] Create challenge
	process.StartStep("Crear un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Challenge Templates Test",
		Description:       "Challenge para validar plantillas por defecto",
		Tags:              []string{"templates", "challenge"},
		Status:            string(exam_consts.ChallengeStatusDraft),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1500,
		WorkerMemoryLimit: 256,
		CodeTemplates: []exam_dtos.CodeTemplateDTO{
			{Language: "python", Template: "# custom template"},
		},
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "5"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "result", Type: string(exam_consts.VariableFormatInt), Value: "25"},
		Constraints:    "1 <= n <= 10^6",
	})
	if err != nil {
		process.Fail("create challenge", err)
	}
	if createdChallenge == nil || createdChallenge.ID == "" {
		process.Fail("create challenge", fmt.Errorf("expected created challenge with ID"))
	}
	challengeID = createdChallenge.ID
	process.Log(fmt.Sprintf("challengeID=%s", challengeID))
	process.EndStep()

	// [STEP 3] Create test case
	process.StartStep("Crear casos de prueba")
	createdTestCase, err := process.Application.TestCaseModule.CreateTestCase.Execute(teacherCtx, exam_dtos.CreateTestCaseInput{
		Name: "default_template_case",
		Input: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "5"},
		},
		ExpectedOutput: exam_dtos.IOVariableDTO{Name: "result", Type: string(exam_consts.VariableFormatInt), Value: "25"},
		IsSample:       true,
		Points:         10,
		ChallengeID:    challengeID,
	})
	if err != nil {
		process.Fail("create test case", err)
	}
	if createdTestCase == nil || createdTestCase.ID == "" {
		process.Fail("create test case", fmt.Errorf("expected created test case with ID"))
	}
	process.Log(fmt.Sprintf("testCaseID=%s", createdTestCase.ID))
	process.EndStep()

	// [STEP 4] Get default templates
	process.StartStep("Obtener plantillas por defecto para el reto")
	templates, err := process.Application.ExamModule.GetCodeDefaultTemplates.Execute(teacherCtx, exam_dtos.DefaultCodeTemplatesInput{
		Inputs: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "5"},
		},
		Output: exam_dtos.IOVariableDTO{Name: "result", Type: string(exam_consts.VariableFormatInt), Value: "25"},
	})
	if err != nil {
		process.Fail("get default code templates", err)
	}
	if len(templates) == 0 {
		process.Fail("get default code templates", fmt.Errorf("expected at least one template"))
	}
	process.EndStep()

	// [STEP 5] Validate expected variables and output print
	process.StartStep("Validar que se reciban todas las variables esperadas y el print con el output")
	templatesByLanguage := make(map[string]string, len(templates))
	for _, tpl := range templates {
		templatesByLanguage[string(tpl.Language)] = tpl.Template
	}

	for _, language := range sub_consts.SupportedProgrammingLanguages {
		template, ok := templatesByLanguage[string(language)]
		if !ok {
			process.Fail("validate default templates", fmt.Errorf("expected template for language %s", language))
		}
		if strings.TrimSpace(template) == "" {
			process.Fail("validate default templates", fmt.Errorf("template for language %s is empty", language))
		}
	}

	pythonTemplate, ok := templatesByLanguage[string(sub_consts.LanguagePython)]
	if !ok {
		process.Fail("validate default templates", fmt.Errorf("expected python template"))
	}
	if !strings.Contains(pythonTemplate, "n = int(input())") {
		process.Fail("validate default templates", fmt.Errorf("expected input variable assignment in python template"))
	}
	if !strings.Contains(pythonTemplate, "result = 0") {
		process.Fail("validate default templates", fmt.Errorf("expected output declaration in python template"))
	}
	if !strings.Contains(pythonTemplate, "print(result)") {
		process.Fail("validate default templates", fmt.Errorf("expected output print in python template"))
	}
	process.EndStep()

	process.End()
}