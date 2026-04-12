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

func TestChallengeFromTeacherView(t *testing.T) {
	process := test.StartTestWithApp(t, "Challenge From Teachers View")
	creatorEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var creatorCtx = context.Background()
	var observerCtx = context.Background()
	var challengeID string

	defer func() {
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(creatorCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
	}()

	// [STEP 1] Login as teacher (creator)
	process.StartStep("Iniciar sesión con usuario de docente (creador)")
	creatorAccess := utils.EnsureAuthUserAccess(t, process.Application, creatorEmail, password, "Teacher Creator")
	creatorCtx = utils.BuildUserCtx(creatorAccess)
	process.EndStep()

	// [STEP 2] Create a challenge with private visibility
	process.StartStep("Crea un reto (private)")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(creatorCtx, exam_dtos.CreateChallengeInput{
		Title:             "Challenge Teachers View",
		Description:       "Challenge para validar visibilidad entre docentes",
		Tags:              []string{"teachers-view", "challenge"},
		Status:            string(exam_consts.ChallengeStatusPrivate),
		Difficulty:        string(exam_consts.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1400,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_consts.VariableFormatInt), Value: "7"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_consts.VariableFormatInt), Value: "7"},
		Constraints:    "1 <= n <= 100",
	})
	if err != nil {
		process.Fail("create private challenge", err)
	}
	if createdChallenge == nil || createdChallenge.ID == "" {
		process.Fail("create private challenge", fmt.Errorf("expected challenge with ID"))
	}
	challengeID = createdChallenge.ID
	process.EndStep()

	// [STEP 3] Login as teacher (observer)
	process.StartStep("Iniciar sesión con usuario de docente (observador)")
	observerAccess := utils.EnsureAuthUserAccess(t, process.Application, observerEmail, password, "Teacher Observer")
	observerCtx = utils.BuildUserCtx(observerAccess)
	process.EndStep()

	// [STEP 4] Try to get challenge details with observer (expect error)
	process.StartStep("Obtiene datos del reto con docente observador (espera error)")
	_, err = process.Application.ChallengeModule.GetChallengeDetails.Execute(observerCtx, exam_dtos.GetChallengeDetailsInput{ChallengeID: challengeID})
	if err == nil {
		process.Fail("observer private challenge access", fmt.Errorf("expected error for observer on private challenge"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 5] Update the challenge to published visibility
	process.StartStep("Actualiza el reto a visibilidad published")
	_, err = process.Application.ChallengeModule.PublishChallenge.Execute(creatorCtx, exam_dtos.PublishChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("publish challenge", err)
	}
	process.EndStep()

	// [STEP 6] Try to get challenge details with observer (expect success)
	process.StartStep("Obtiene datos del reto con docente observador")
	publishedList, err := process.Application.ChallengeModule.GetPublicChallenges.Execute(observerCtx, exam_dtos.GetPublicChallengesInput{})
	if err != nil {
		process.Fail("observer get public challenges", err)
	}
	foundPublished := false
	for _, c := range publishedList {
		if c != nil && c.ID == challengeID {
			foundPublished = true
			break
		}
	}
	if !foundPublished {
		process.Fail("observer get public challenges", fmt.Errorf("expected published challenge %s to be visible", challengeID))
	}
	process.EndStep()

	// [STEP 7] Create an exam item with the published challenge as observer (expect success)
	process.StartStep("Actualiza el reto a visibilidad archived")
	_, err = process.Application.ChallengeModule.ArchiveChallenge.Execute(creatorCtx, exam_dtos.ArchiveChallengeInput{ChallengeID: challengeID})
	if err != nil {
		process.Fail("archive challenge", err)
	}
	process.EndStep()

	// [STEP 8] Try to get challenge details with observer (expect error)
	process.StartStep("Obtiene datos del reto con docente observador (espera error)")
	archivedList, err := process.Application.ChallengeModule.GetPublicChallenges.Execute(observerCtx, exam_dtos.GetPublicChallengesInput{})
	if err != nil {
		process.Fail("observer get public challenges after archive", err)
	}
	for _, c := range archivedList {
		if c != nil && c.ID == challengeID {
			process.Fail("observer archived challenge visibility", fmt.Errorf("archived challenge %s should not be visible", challengeID))
		}
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
