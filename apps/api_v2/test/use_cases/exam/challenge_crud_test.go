package usecases_test

import (
	"context"
	"fmt"
	"testing"

	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestChallengeCRUD(t *testing.T) {
	process := test.StartTestWithApp(t, "Challenge CRUD")
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var challengeID string
	var deletedChallengeID string

	defer func() {
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
	}()
	
	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.Log(fmt.Sprintf("teacherID=%s", teacherAccess.UserData.ID))
	process.EndStep()

	// [STEP 2] Create a challenge
	process.StartStep("Crea un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Challenge CRUD Test",
		Description:       "Challenge creado por test CRUD",
		Tags:              []string{"crud", "challenge"},
		Status:            string(exam_consts.ChallengeStatusDraft),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1500,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "a", Type: string(exam_consts.VariableFormatInt), Value: "2"},
			{Name: "b", Type: string(exam_consts.VariableFormatInt), Value: "3"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "sum", Type: string(exam_consts.VariableFormatInt), Value: "5"},
		Constraints:    "1 <= a,b <= 1000",
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

	// [STEP 3] Update the challenge
	process.StartStep("Actualiza el reto")
	updatedTitle := "Challenge CRUD Test Updated"
	updatedDescription := "Challenge actualizado por test CRUD"
	updatedDifficulty := string(exam_consts.ChallengeDifficultyMedium)
	updatedChallenge, err := process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{
		ChallengeID: challengeID,
		Title:       &updatedTitle,
		Description: &updatedDescription,
		Difficulty:  &updatedDifficulty,
	})
	if err != nil {
		process.Fail("update challenge", err)
	}
	if updatedChallenge == nil || updatedChallenge.Title != updatedTitle {
		process.Fail("update challenge", fmt.Errorf("expected updated challenge title"))
	}
	process.EndStep()

	// [STEP 4] Get challenge details and validate updates
	process.StartStep("Obtiene los datos del reto y valida los cambios")
	reloadedChallenge, err := process.Application.ChallengeModule.GetChallengeDetails.Execute(teacherCtx, exam_dtos.GetChallengeDetailsInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("get challenge details", err)
	}
	if reloadedChallenge == nil || reloadedChallenge.ID != challengeID {
		process.Fail("get challenge details", fmt.Errorf("expected challenge details for %s", challengeID))
	}
	if reloadedChallenge.Title != updatedTitle || reloadedChallenge.Description != updatedDescription {
		process.Fail("get challenge details", fmt.Errorf("challenge update not persisted"))
	}
	process.EndStep()

	// [STEP 5] Delete the challenge
	process.StartStep("Elimina el reto")
	err = process.Application.ChallengeModule.DeleteChallenge.Execute(teacherCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("delete challenge", err)
	}
	deletedChallengeID = challengeID
	challengeID = ""
	process.EndStep()

	// [STEP 6] Verify deletion
	process.StartStep("Verifica eliminación")
	_, err = process.Application.ChallengeModule.GetChallengeDetails.Execute(teacherCtx, exam_dtos.GetChallengeDetailsInput{ChallengeID: deletedChallengeID})
	if err == nil {
		process.Fail("verify challenge deletion", fmt.Errorf("expected error after deleting challenge"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
