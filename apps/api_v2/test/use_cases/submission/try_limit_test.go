package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestTryLimit(t *testing.T) {
	process := test.StartTestWithApp(t, "Try Limit")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var studentID string
	var examID string
	var sessionOneID string
	var sessionTwoID string

	defer func() {
		if sessionTwoID != "" {
			t.Logf("[CLEANUP] Cerrando segunda sesión %s", sessionTwoID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionTwoID})
		}
		if sessionOneID != "" {
			t.Logf("[CLEANUP] Cerrando primera sesión %s", sessionOneID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionOneID})
		}
		if examID != "" {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.EndStep()

	// [STEP 2] Create public exam with try_limit = 2
	process.StartStep("Crear examen público (con solo 2 intentos try_limit)")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Try Limit Exam",
		Description:          "Examen para validar límite de intentos",
		Visibility:           string(exam_consts.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            120,
		TryLimit:             2,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	examID = createdExam.ID
	process.EndStep()

	// [STEP 3] Login as student
	process.StartStep("Iniciar sesión con usuario de estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, studentEmail, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	studentID = studentAccess.UserData.ID
	process.EndStep()

	// [STEP 4] Create first session for the student
	process.StartStep("Crear una sesión en el examen")
	firstSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create first session", err)
	}
	sessionOneID = firstSession.ID
	process.EndStep()

	// [STEP 5] Create second session for the student
	process.StartStep("Cerrar la sesión")
	_, err = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionOneID})
	if err != nil {
		process.Fail("close first session", err)
	}
	sessionOneID = ""
	process.EndStep()

	// [STEP 6] Create second session for the student
	process.StartStep("Crear una sesión en el examen")
	secondSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create second session", err)
	}
	sessionTwoID = secondSession.ID
	process.EndStep()

	// [STEP 7] Close second session
	process.StartStep("Cerrar la sesión")
	_, err = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionTwoID})
	if err != nil {
		process.Fail("close second session", err)
	}
	sessionTwoID = ""
	process.EndStep()

	// [STEP 8] Try to create third session and expect error due to try_limit
	process.StartStep("Crear una sesión en el examen (espera error)")
	_, err = process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err == nil {
		process.Fail("create third session", fmt.Errorf("expected error when try_limit is exceeded"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
