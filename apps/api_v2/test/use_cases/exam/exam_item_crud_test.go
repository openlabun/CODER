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

func TestExamItemCRUD(t *testing.T) {
	process := test.StartTestWithApp(t, "ExamItem CRUD")
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var examID string
	var challengeID string
	var examItemID string
	var deletedExamItemID string

	defer func() {
		if examItemID != "" {
			t.Logf("[CLEANUP] Eliminando exam item %s", examItemID)
			_ = process.Application.ExamItemModule.DeleteExamItem.Execute(teacherCtx, exam_dtos.DeleteExamItemInput{ID: examItemID})
		}
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando exam %s", examID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
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

	// [STEP 2] Create an exam
	process.StartStep("Crear un examen")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "ExamItem CRUD Exam",
		Description:          "Exam auxiliar para CRUD de exam item",
		Visibility:           string(exam_entities.VisibilityPrivate),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: false,
		TimeLimit:            3600,
		TryLimit:             1,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	examID = createdExam.ID
	process.EndStep()

	// [STEP 3] Create a challenge
	process.StartStep("Crear un reto")
	createdChallenge, err := process.Application.ChallengeModule.CreateChallenge.Execute(teacherCtx, exam_dtos.CreateChallengeInput{
		Title:             "ExamItem CRUD Challenge",
		Description:       "Challenge auxiliar para exam item",
		Tags:              []string{"exam-item", "crud"},
		Status:            string(exam_entities.ChallengeStatusPublished),
		Difficulty:        string(exam_entities.ChallengeDifficultyEasy),
		WorkerTimeLimit:   1200,
		WorkerMemoryLimit: 256,
		InputVariables: []exam_dtos.IOVariableDTO{
			{Name: "n", Type: string(exam_entities.VariableFormatInt), Value: "10"},
		},
		OutputVariable: exam_dtos.IOVariableDTO{Name: "out", Type: string(exam_entities.VariableFormatInt), Value: "10"},
		Constraints:    "1 <= n <= 1000",
	})
	if err != nil {
		process.Fail("create challenge", err)
	}
	challengeID = createdChallenge.ID
	process.EndStep()

	// [STEP 4] Create an exam item with the challenge
	process.StartStep("Crear un punto de examen")
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
	process.EndStep()

	// [STEP 5] Try to create another exam item with the same challenge (expect error)
	process.StartStep("Crear otro punto de examen con mismo reto (espera error)")
	_, err = process.Application.ExamItemModule.CreateExamItem.Execute(teacherCtx, exam_dtos.CreateExamItemInput{
		ExamID:      examID,
		ChallengeID: challengeID,
		Order:       2,
		Points:      50,
	})
	if err == nil {
		process.Fail("duplicate exam item", fmt.Errorf("expected error when adding same challenge twice"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 6] Create another exam item with a different challenge
	process.StartStep("Actualizar punto de examen")
	updatedOrder := 3
	updatedPoints := 150
	updatedExamItem, err := process.Application.ExamItemModule.UpdateExamItem.Execute(teacherCtx, exam_dtos.UpdateExamItemInput{
		ID:     examItemID,
		Order:  &updatedOrder,
		Points: &updatedPoints,
	})
	if err != nil {
		process.Fail("update exam item", err)
	}
	if updatedExamItem == nil || updatedExamItem.Order != updatedOrder || updatedExamItem.Points != updatedPoints {
		process.Fail("update exam item", fmt.Errorf("expected updated exam item values"))
	}
	process.EndStep()

	// [STEP 7] Get exam items and validate the updated item
	process.StartStep("Obtener punto de examen y validar datos")
	items, err := process.Application.ExamModule.GetExamItems.Execute(teacherCtx, exam_dtos.GetExamItemsInput{ExamID: examID})
	if err != nil {
		process.Fail("get exam items", err)
	}
	found := false
	for _, item := range items {
		if item.ID == examItemID {
			found = true
			if item.Order != updatedOrder || item.Points != updatedPoints {
				process.Fail("get exam items", fmt.Errorf("unexpected exam item values"))
			}
		}
	}
	if !found {
		process.Fail("get exam items", fmt.Errorf("expected exam item %s in exam items list", examItemID))
	}
	process.EndStep()

	// [STEP 8] Delete the exam item
	process.StartStep("Eliminar punto de examen")
	err = process.Application.ExamItemModule.DeleteExamItem.Execute(teacherCtx, exam_dtos.DeleteExamItemInput{ID: examItemID})
	if err != nil {
		process.Fail("delete exam item", err)
	}
	deletedExamItemID = examItemID
	examItemID = ""
	process.EndStep()

	// [STEP 9] Verify deletion by trying to get exam items and checking the deleted item is not there
	process.StartStep("Verificar eliminación")
	itemsAfterDelete, err := process.Application.ExamModule.GetExamItems.Execute(teacherCtx, exam_dtos.GetExamItemsInput{ExamID: examID})
	if err != nil {
		process.Fail("get exam items after delete", err)
	}
	for _, item := range itemsAfterDelete {
		if item.ID == deletedExamItemID {
			process.Fail("verify exam item deletion", fmt.Errorf("exam item %s should not exist after deletion", deletedExamItemID))
		}
	}
	process.EndStep()

	process.End()
}
