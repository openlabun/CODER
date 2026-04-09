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

func TestExamItemPrivacy(t *testing.T) {
	process := test.StartTestWithApp(t, "ExamItem Privacy")
	ownerEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var ownerCtx = context.Background()
	var observerCtx = context.Background()
	var examID string
	var challengeID string
	var examItemID string

	defer func() {
		if examItemID != "" {
			t.Logf("[CLEANUP] Eliminando exam item %s", examItemID)
			_ = process.Application.ExamItemModule.DeleteExamItem.Execute(ownerCtx, exam_dtos.DeleteExamItemInput{ID: examItemID})
		}
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando exam %s", examID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(ownerCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
		if challengeID != "" {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = process.Application.ChallengeModule.DeleteChallenge.Execute(ownerCtx, exam_dtos.DeleteChallengeInput{ChallengeID: challengeID})
		}
	}()

	// [STEP 1] Login as teacher (owner)
	process.StartStep("Iniciar sesión con docente dueño")
	ownerAccess := utils.EnsureAuthUserAccess(t, process.Application, ownerEmail, password, "Teacher Owner")
	ownerCtx = utils.BuildUserCtx(ownerAccess)
	process.EndStep()

	// [STEP 2] Create an exam and a challenge as owner
	process.StartStep("Crear exam y challenge del dueño")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(ownerCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "ExamItem Privacy Exam",
		Description:          "Exam para pruebas de privacidad de exam items",
		Visibility:           string(exam_entities.VisibilityPrivate),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: false,
		TimeLimit:            3600,
		TryLimit:             1,
		ProfessorID:          ownerAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create owner exam", err)
	}
	examID = createdExam.ID

	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(ownerCtx, exam_dtos.CreateChallengeInput{
		Title:             "ExamItem Privacy Challenge",
		Description:       "Challenge para pruebas de privacidad",
		Tags:              []string{"exam-item", "privacy"},
		Status:            string(exam_entities.ChallengeStatusPublished),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "x", Type: string(exam_entities.VariableFormatInt), Value: "7"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "y", Type: string(exam_entities.VariableFormatInt), Value: "7"},
		Constraints:    "x >= 0",
	})
	if err != nil {
		process.Fail("create owner challenge", err)
	}
	challengeID = createdChallenge.ID
	process.EndStep()

	// [STEP 3] Create an exam item with the challenge as owner
	process.StartStep("Crear exam item del dueño")
	createdExamItem, err := process.Application.ExamItemModule.CreateExamItem.Execute(ownerCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: challengeID,
		Order:       1,
		Points:      100,
	})
	if err != nil {
		process.Fail("create owner exam item", err)
	}
	examItemID = createdExamItem.ID
	process.EndStep()

	// [STEP 4] Login as teacher (observer)
	process.StartStep("Iniciar sesión con docente observador")
	observerAccess := utils.EnsureAuthUserAccess(t, process.Application, observerEmail, password, "Teacher Observer")
	observerCtx = utils.BuildUserCtx(observerAccess)
	process.EndStep()

	// [STEP 5] Try to create an exam item in the owner's exam with the owner's challenge (expect error)
	process.StartStep("Crear exam item en exam de otro docente (espera error)")
	_, err = process.Application.ExamItemModule.CreateExamItem.Execute(observerCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: challengeID,
		Order:       2,
		Points:      100,
	})
	if err == nil {
		process.Fail("observer create exam item", fmt.Errorf("expected error when observer creates exam item in foreign exam"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 6] Try to update the owner's exam item as observer (expect error)
	process.StartStep("Actualizar exam item de otro docente (espera error)")
	newPoints := 200
	_, err = process.Application.ExamItemModule.UpdateExamItem.Execute(observerCtx, exam_dtos.UpdateExamItemInput{
		ID:     examItemID,
		Points: &newPoints,
	})
	if err == nil {
		process.Fail("observer update exam item", fmt.Errorf("expected error when observer updates foreign exam item"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Try to delete the owner's exam item as observer (expect error)
	process.StartStep("Eliminar exam item de otro docente (espera error)")
	err = process.Application.ExamItemModule.DeleteExamItem.Execute(observerCtx, exam_dtos.DeleteExamItemInput{ID: examItemID})
	if err == nil {
		process.Fail("observer delete exam item", fmt.Errorf("expected error when observer deletes foreign exam item"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
