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

func TestChallengeStates(t *testing.T) {
	process := test.StartTestWithApp(t, "Challenge States Transitions and Valitadions")
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
	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.EndStep()

	// [STEP 2] Create a challenge in draft state
	process.StartStep("Crea un reto en estado draft")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "Challenge States Test",
		Description:       "Challenge creado para validar transiciones",
		Tags:              []string{"states", "challenge"},
		Status:            string(exam_consts.ChallengeStatusDraft),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		CodeTemplates: []exam_dtos.CodeTemplateDTO{
			{Language: "python", Template: "def solve() { return; }"},
		},
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "10"},
		Constraints:    "1 <= n <= 1000",
	})
	if err != nil {
		process.Fail("create draft challenge", err)
	}
	challengeID = createdChallenge.ID
	process.EndStep()

	// [STEP 3] Update the challenge to published
	process.StartStep("Actualiza el reto")
	step3Title := "Challenge States Test Updated"
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{
		ChallengeID: challengeID,
		Title:       &step3Title,
	})
	if err != nil {
		process.Fail("update draft challenge", err)
	}
	process.EndStep()

	// [STEP 4] Get challenge details and validate updates
	process.StartStep("Obtiene los datos del reto y valida los cambios")
	challenge, err := process.Application.ChallengeModule.GetChallengeDetails.Execute(teacherCtx, exam_dtos.GetChallengeDetailsInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("get challenge details", err)
	}
	if challenge == nil || challenge.Title != step3Title || challenge.Status != exam_consts.ChallengeStatusDraft {
		process.Fail("get challenge details", fmt.Errorf("unexpected challenge state after step 4"))
	}
	process.EndStep()

	// [STEP 5] Update the challenge to published
	process.StartStep("Actualiza el reto a estado published")
	_, err = process.Application.ChallengeModule.PublishChallenge.Execute(teacherCtx, exam_dtos.PublishChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("publish challenge", err)
	}
	process.EndStep()

	// [STEP 6] Updates challenge
	process.StartStep("Actualiza el reto")
	description := "Update 1"
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{ChallengeID: challengeID, Description: &description})
	if err != nil {
		process.Fail("error updating challenge", fmt.Errorf("unexpected error when updating published challenge: %v", err))
	}
	process.EndStep()

	// [STEP 7] Update the challenge to archived
	process.StartStep("Actualiza el reto a estado private")
	toPrivate := string(exam_consts.ChallengeStatusPrivate)
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{ChallengeID: challengeID, Status: &toPrivate})
	if err != nil {
		process.Fail("transition published->private", err)
	}
	process.EndStep()

	// [STEP 8] Updates Challenge
	process.StartStep("Actualiza el reto")
	description = "Update 2"
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{ChallengeID: challengeID, Description: &description})
	if err != nil {
		process.Fail("error updating challenge", fmt.Errorf("unexpected error when updating private challenge: %v", err))
	}
	process.EndStep()

	// [STEP 9] Update the challenge to archived
	process.StartStep("Actualiza el reto a estado archived")
	_, err = process.Application.ChallengeModule.ArchiveChallenge.Execute(teacherCtx, exam_dtos.ArchiveChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("transition to archived", err)
	}
	process.EndStep()

	// [STEP 10] Get challenge details and validate archived state
	process.StartStep("Actualiza el reto (espera error)")
	description = "Update 3"
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{ChallengeID: challengeID, Description: &description})
	if err == nil {
		process.Fail("invalid transition archived->private", fmt.Errorf("expected error on invalid transition"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 11] Update the challenge back to published
	process.StartStep("Actualiza el reto a estado published")
	_, err = process.Application.ChallengeModule.PublishChallenge.Execute(teacherCtx, exam_dtos.PublishChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("transition archived->published", err)
	}
	process.EndStep()

	// [STEP 12] Get challenge details and validate published state
	process.StartStep("Actualiza el reto a estado archived")
	_, err = process.Application.ChallengeModule.ArchiveChallenge.Execute(teacherCtx, exam_dtos.ArchiveChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("transition published->archived", err)
	}
	process.EndStep()

	// [STEP 13] Get challenge details and validate archived state
	process.StartStep("Actualiza el reto a estado private")
	// Bridge required by current state machine: archived -> published -> private.
	_, err = process.Application.ChallengeModule.PublishChallenge.Execute(teacherCtx, exam_dtos.PublishChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("bridge transition archived->published", err)
	}
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{ChallengeID: challengeID, Status: &toPrivate})
	if err != nil {
		process.Fail("transition published->private", err)
	}
	process.EndStep()

	// [STEP 14] Get challenge details and validate private state
	process.StartStep("Actualiza el reto a estado draft (espera error)")
	invalidDraft := string(exam_consts.ChallengeStatusDraft)
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{ChallengeID: challengeID, Status: &invalidDraft})
	if err == nil {
		process.Fail("invalid transition private->draft", fmt.Errorf("expected error on invalid transition"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 15] Update the challenge to published
	process.StartStep("Actualiza el reto a estado published")
	_, err = process.Application.ChallengeModule.PublishChallenge.Execute(teacherCtx, exam_dtos.PublishChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("transition private->published", err)
	}
	process.EndStep()

	// [STEP 16] Get challenge details and validate published state
	process.StartStep("Actualiza el reto a estado draft (espera error)")
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(teacherCtx, exam_dtos.UpdateChallengeInput{ChallengeID: challengeID, Status: &invalidDraft})
	if err == nil {
		process.Fail("invalid transition published->draft", fmt.Errorf("expected error on invalid transition"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
