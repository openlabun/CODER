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

func TestChallengeFork(t *testing.T) {
	process := test.StartTestWithApp(t, "Challenge Fork")
	creatorEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var creatorCtx = context.Background()
	var observerCtx = context.Background()
	var originalChallengeID string
	var forkedChallengeID string

	defer func() {
		if forkedChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge fork %s", forkedChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(observerCtx, exam_dtos.DeleteChallengeInput{ChallengeID: forkedChallengeID})
		}
		if originalChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge original %s", originalChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(creatorCtx, exam_dtos.DeleteChallengeInput{ChallengeID: originalChallengeID})
		}
	}()

	// [STEP 1] Login as creator teacher
	process.StartStep("Iniciar sesión con usuario de docente (creador)")
	creatorAccess := utils.EnsureAuthUserAccess(t, process.Application, creatorEmail, password, "Teacher Creator")
	creatorCtx = utils.BuildUserCtx(creatorAccess)
	process.Log(fmt.Sprintf("creatorID=%s", creatorAccess.UserData.ID))
	process.EndStep()

	// [STEP 2] Create a challenge in private status
	process.StartStep("Crea un reto en estado private")
	originalChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(creatorCtx, exam_dtos.CreateChallengeInput{
		Title:             "Challenge Fork Original",
		Description:       "Challenge original para prueba de fork",
		Tags:              []string{"fork", "original"},
		Status:            string(exam_consts.ChallengeStatusPrivate),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1500,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "x", Type: string(exam_consts.VariableFormatInt), Value: "5"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "y", Type: string(exam_consts.VariableFormatInt), Value: "5"},
		Constraints:    "1 <= x <= 1000",
	})
	if err != nil {
		process.Fail("create original challenge", err)
	}
	if originalChallenge == nil || originalChallenge.ID == "" {
		process.Fail("create original challenge", fmt.Errorf("expected original challenge with ID"))
	}
	originalChallengeID = originalChallenge.ID
	process.EndStep()

	// [STEP 3] Login as observer teacher
	process.StartStep("Iniciar sesión con usuario de docente (observador)")
	observerAccess := utils.EnsureAuthUserAccess(t, process.Application, observerEmail, password, "Teacher Observer")
	observerCtx = utils.BuildUserCtx(observerAccess)
	process.Log(fmt.Sprintf("observerID=%s", observerAccess.UserData.ID))
	process.EndStep()

	// [STEP 4] Try to fork the original challenge (expect error because it's private)
	process.StartStep("Hace fork al reto (espera error)")
	_, err = process.Application.ChallengeModule.ForkChallenge.Execute(observerCtx, exam_dtos.ForkChallengeInput{ChallengeID: originalChallengeID})
	if err == nil {
		process.Fail("fork private challenge", fmt.Errorf("expected error when forking private challenge"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 5] Publish the original challenge
	process.StartStep("Actualiza reto original a estado published")
	_, err = process.Application.ChallengeModule.PublishChallenge.Execute(creatorCtx, exam_dtos.PublishChallengeInput{ChallengeID: originalChallengeID})
	if err != nil {
		process.Fail("publish original challenge", err)
	}
	process.EndStep()

	// [STEP 6] Fork the original challenge
	process.StartStep("Hace fork al reto")
	forkedChallenge, err := process.Application.ChallengeModule.ForkChallenge.Execute(observerCtx, exam_dtos.ForkChallengeInput{ChallengeID: originalChallengeID})
	if err != nil {
		process.Fail("fork published challenge", err)
	}
	if forkedChallenge == nil || forkedChallenge.ID == "" {
		process.Fail("fork published challenge", fmt.Errorf("expected forked challenge with ID"))
	}
	if forkedChallenge.ID == originalChallengeID {
		process.Fail("fork published challenge", fmt.Errorf("forked challenge must have a different ID"))
	}
	forkedChallengeID = forkedChallenge.ID
	process.EndStep()

	// [STEP 7] Create an exam item with the forked challenge to verify it's usable
	process.StartStep("Actualiza el reto copiado")
	updatedForkTitle := "Challenge Fork Copied Updated"
	updatedForkDescription := "Reto fork actualizado"
	updatedFork, err := process.Application.ChallengeModule.UpdateChallenge.Execute(observerCtx, exam_dtos.UpdateChallengeInput{
		ChallengeID: forkedChallengeID,
		Title:       &updatedForkTitle,
		Description: &updatedForkDescription,
	})
	if err != nil {
		process.Fail("update forked challenge", err)
	}
	if updatedFork == nil || updatedFork.Title != updatedForkTitle {
		process.Fail("update forked challenge", fmt.Errorf("expected updated forked challenge title"))
	}
	process.EndStep()

	// [STEP 8] Verify original challenge has not changed
	process.StartStep("Verifica que no haya cambios en el reto original")
	originalAfterForkUpdate, err := process.Application.ChallengeModule.GetChallengeDetails.Execute(creatorCtx, exam_dtos.GetChallengeDetailsInput{ChallengeID: originalChallengeID})
	if err != nil {
		process.Fail("get original challenge after fork update", err)
	}
	if originalAfterForkUpdate == nil {
		process.Fail("get original challenge after fork update", fmt.Errorf("expected original challenge details"))
	}
	if originalAfterForkUpdate.Title == updatedForkTitle || originalAfterForkUpdate.Description == updatedForkDescription {
		process.Fail("verify original unchanged", fmt.Errorf("original challenge should not be modified by fork update"))
	}
	process.EndStep()

	// [STEP 9] Verify forked challenge has the updated values
	process.StartStep("Verifica los cambios en el reto copiado")
	forkAfterUpdate, err := process.Application.ChallengeModule.GetChallengeDetails.Execute(observerCtx, exam_dtos.GetChallengeDetailsInput{ChallengeID: forkedChallengeID})
	if err != nil {
		process.Fail("get forked challenge after update", err)
	}
	if forkAfterUpdate == nil || forkAfterUpdate.Title != updatedForkTitle || forkAfterUpdate.Description != updatedForkDescription {
		process.Fail("verify fork updated", fmt.Errorf("expected forked challenge to keep updated values"))
	}
	process.EndStep()

	process.End()
}
