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

func TestExamItemChallengeFork(t *testing.T) {
	process := test.StartTestWithApp(t, "ExamItem Challenge Fork")
	creatorEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var creatorCtx = context.Background()
	var observerCtx = context.Background()
	var creatorChallengeID string
	var observerExamID string
	var forkedChallengeID string

	defer func() {
		if observerExamID != "" {
			t.Logf("[CLEANUP] Eliminando exam %s", observerExamID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(observerCtx, exam_dtos.DeleteExamInput{ExamID: observerExamID})
		}
		if creatorChallengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge creador %s", creatorChallengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(creatorCtx, exam_dtos.DeleteChallengeInput{ChallengeID: creatorChallengeID})
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
	process.StartStep("Crear challenge publicado del creador")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(creatorCtx, exam_dtos.CreateChallengeInput{
		Title:             "ExamItem Fork Original",
		Description:       "Challenge original para prueba de fork en exam item",
		Tags:              []string{"exam-item", "fork"},
		Status:            string(exam_entities.ChallengeStatusPublished),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1500,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "x", Type: string(exam_entities.VariableFormatInt), Value: "10"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "y", Type: string(exam_entities.VariableFormatInt), Value: "10"},
		Constraints:    "1 <= x <= 1000",
	})
	if err != nil {
		process.Fail("create creator challenge", err)
	}
	creatorChallengeID = createdChallenge.ID
	process.EndStep()

	// [STEP 3] Login as teacher (observer)
	process.StartStep("Iniciar sesión con docente observador")
	observerAccess := utils.EnsureAuthUserAccess(t, process.Application, observerEmail, password, "Teacher Observer")
	observerCtx = utils.BuildUserCtx(observerAccess)
	process.EndStep()

	// [STEP 4] Create an exam item with creator's challenge
	process.StartStep("Crear examen del observador")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(observerCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "ExamItem Fork Exam",
		Description:          "Exam para prueba de fork en exam item",
		Visibility:           string(exam_entities.VisibilityPrivate),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: false,
		TimeLimit:            3600,
		TryLimit:             1,
		ProfessorID:          observerAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create observer exam", err)
	}
	observerExamID = createdExam.ID
	process.EndStep()

	// [STEP 5] Create exam item with creator's challenge (expect fork)
	process.StartStep("Crear exam item con challenge del creador")
	observerExamItem, err := process.Application.ExamItemModule.CreateExamItem.Execute(observerCtx, exam_dtos.CreateExamItemInput{
		ExamID:      observerExamID,
		ChallengeID: creatorChallengeID,
		Order:       1,
		Points:      100,
	})
	if err != nil {
		process.Fail("create exam item with foreign challenge", err)
	}
	forkedChallengeID = observerExamItem.ChallengeID
	process.EndStep()

	// [STEP 6] Verify that the exam item challenge is a fork and not the original
	process.StartStep("Modificar challenge original del creador")
	updatedOriginalTitle := "ExamItem Fork Original Updated"
	updatedOriginalDescription := "Challenge original actualizado tras crear exam item"
	_, err = process.Application.ChallengeModule.UpdateChallenge.Execute(creatorCtx, exam_dtos.UpdateChallengeInput{
		ChallengeID: creatorChallengeID,
		Title:       &updatedOriginalTitle,
		Description: &updatedOriginalDescription,
	})
	if err != nil {
		process.Fail("update creator challenge", err)
	}
	process.EndStep()

	// [STEP 7] Get exam items and verify that the exam item challenge remains unchanged (forked) after original update
	process.StartStep("Obtener exam items y verificar que no cambie el challenge del exam item")
	items, err := process.Application.ExamModule.GetExamItems.Execute(observerCtx, exam_dtos.GetExamItemsInput{ExamID: observerExamID})
	if err != nil {
		process.Fail("get exam items", err)
	}
	if len(items) != 1 {
		process.Fail("get exam items", fmt.Errorf("expected 1 exam item, got %d", len(items)))
	}
	item := items[0]
	if item.ID != observerExamItem.ID {
		process.Fail("get exam items", fmt.Errorf("unexpected exam item id, expected %s got %s", observerExamItem.ID, item.ID))
	}
	if item.Challenge == nil {
		process.Fail("get exam items", fmt.Errorf("expected challenge details in exam item"))
	}
	if item.Challenge.ID == creatorChallengeID {
		process.Fail("verify exam item challenge fork", fmt.Errorf("expected exam item challenge to be forked and have different ID from original"))
	}
	if item.Challenge.Title == updatedOriginalTitle || item.Challenge.Description == updatedOriginalDescription {
		process.Fail("verify exam item challenge fork", fmt.Errorf("expected exam item challenge to remain unchanged after original update"))
	}
	process.EndStep()

	process.End()
}
