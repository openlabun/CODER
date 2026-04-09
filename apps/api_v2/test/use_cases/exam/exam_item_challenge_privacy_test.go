package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestExamItemChallengePrivacy(t *testing.T) {
	process := test.StartTestWithApp(t, "ExamItem Challenge Privacy")
	creatorEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var creatorCtx = context.Background()
	var observerCtx = context.Background()
	var examID string
	var privateChallengeID string
	var publishedChallengeID string
	var ownChallengeID string
	var forkedChallengeID string

	defer func() {
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando exam %s", examID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(observerCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
		if ownChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge propio %s", ownChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(observerCtx, exam_dtos.DeleteChallengeInput{ChallengeID: ownChallengeID})
		}
		if privateChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge privado %s", privateChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(creatorCtx, exam_dtos.DeleteChallengeInput{ChallengeID: privateChallengeID})
		}
		if publishedChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge publicado %s", publishedChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(creatorCtx, exam_dtos.DeleteChallengeInput{ChallengeID: publishedChallengeID})
		}
		if forkedChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge forkeado %s", forkedChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(observerCtx, exam_dtos.DeleteChallengeInput{ChallengeID: forkedChallengeID})
		}
	}()

	// [STEP 1] Login as teacher (creator)
	process.StartStep("Iniciar sesión con docente creador")
	creatorAccess := utils.EnsureAuthUserAccess(t, process.Application, creatorEmail, password, "Teacher Creator")
	creatorCtx = utils.BuildUserCtx(creatorAccess)
	process.EndStep()

	// [STEP 2] Create private challenge and published challenge as creator
	process.StartStep("Crear challenge privado y challenge publicado del creador")
	privateChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(creatorCtx, exam_dtos.CreateChallengeInput{
		Title:             "Creator Private Challenge",
		Description:       "Challenge privado para prueba de privacidad",
		Tags:              []string{"privacy", "private"},
		Status:            string(exam_entities.ChallengeStatusPrivate),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "x", Type: string(exam_entities.VariableFormatInt), Value: "1"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "y", Type: string(exam_entities.VariableFormatInt), Value: "1"},
		Constraints:    "x >= 0",
	})
	if err != nil {
		process.Fail("create private challenge", err)
	}
	privateChallengeID = privateChallenge.ID

	publishedChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(creatorCtx, exam_dtos.CreateChallengeInput{
		Title:             "Creator Published Challenge",
		Description:       "Challenge publicado para prueba de privacidad",
		Tags:              []string{"privacy", "published"},
		Status:            string(exam_entities.ChallengeStatusPublished),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "x", Type: string(exam_entities.VariableFormatInt), Value: "2"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "y", Type: string(exam_entities.VariableFormatInt), Value: "2"},
		Constraints:    "x >= 0",
	})
	if err != nil {
		process.Fail("create published challenge", err)
	}
	publishedChallengeID = publishedChallenge.ID
	process.EndStep()

	// [STEP 3] Login as teacher (observer)
	process.StartStep("Iniciar sesión con docente observador")
	observerAccess := utils.EnsureAuthUserAccess(t, process.Application, observerEmail, password, "Teacher Observer")
	observerCtx = utils.BuildUserCtx(observerAccess)
	process.EndStep()

	// [STEP 4] Create an exam for the observer
	process.StartStep("Crear examen del observador")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(observerCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "ExamItem Challenge Privacy Exam",
		Description:          "Exam de observador para pruebas de challenge privacy",
		Visibility:           string(exam_entities.VisibilityPrivate),
		StartTime:            now.Add(3 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: false,
		TimeLimit:            3600,
		TryLimit:             1,
		ProfessorID:          observerAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create observer exam", err)
	}
	examID = createdExam.ID
	process.EndStep()

	// [STEP 5] Create exam item with creator's private challenge (expect error)
	process.StartStep("Crear challenge propio del observador")
	ownChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(observerCtx, exam_dtos.CreateChallengeInput{
		Title:             "Observer Own Challenge",
		Description:       "Challenge propio del observador",
		Tags:              []string{"privacy", "own"},
		Status:            string(exam_entities.ChallengeStatusPrivate),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "x", Type: string(exam_entities.VariableFormatInt), Value: "3"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "y", Type: string(exam_entities.VariableFormatInt), Value: "3"},
		Constraints:    "x >= 0",
	})
	if err != nil {
		process.Fail("create own challenge", err)
	}
	ownChallengeID = ownChallenge.ID
	process.EndStep()

	// [STEP 6] Create exam item with creator's private challenge (expect error)
	process.StartStep("Crear exam item con challenge propio (ok)")
	_, err = process.Application.ExamItemModule.CreateExamItem.Execute(observerCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: ownChallengeID,
		Order:       1,
		Points:      100,
	})
	if err != nil {
		process.Fail("create exam item with own challenge", err)
	}
	process.EndStep()

	// [STEP 7] Create exam item with creator's private challenge (expect error)
	process.StartStep("Crear exam item con challenge privado de otro docente (espera error)")
	_, err = process.Application.ExamItemModule.CreateExamItem.Execute(observerCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: privateChallengeID,
		Order:       2,
		Points:      100,
	})
	if err == nil {
		process.Fail("create exam item with foreign private challenge", fmt.Errorf("expected error when using private challenge from another teacher"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 8] Create exam item with creator's published challenge (expect success)
	process.StartStep("Crear exam item con challenge publicado de otro docente (ok)")
	examItem, err := process.Application.ExamItemModule.CreateExamItem.Execute(observerCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: publishedChallengeID,
		Order:       3,
		Points:      120,
	})
	if err != nil {
		process.Fail("create exam item with foreign published challenge", err)
	}
	t.Logf("Original challenge ID: %s", publishedChallengeID)
	t.Logf("Forked challenge ID: %s", examItem.ChallengeID)
	forkedChallengeID = examItem.ChallengeID
	process.EndStep()

	process.End()
}
