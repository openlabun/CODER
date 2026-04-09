package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestSessionCRUD(t *testing.T) {
	process := test.StartTestWithApp(t, "Session CRUD")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var studentID string
	var examOneID string
	var examTwoID string
	var sessionID string

	defer func() {
		if sessionID != "" {
			t.Logf("[CLEANUP] Cerrando sesión %s", sessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionID})
		}
		if examTwoID != "" {
			t.Logf("[CLEANUP] Eliminando examen %s", examTwoID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examTwoID})
		}
		if examOneID != "" {
			t.Logf("[CLEANUP] Eliminando examen %s", examOneID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examOneID})
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesión con usuario de docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacherEmail, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	process.EndStep()

	// [STEP 2] Create a public exam without course
	process.StartStep("Crear examen público (visibilidad public y sin curso)")
	now := time.Now().UTC()
	examOne, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Session CRUD Exam One",
		Description:          "Primer examen para CRUD de sesión",
		Visibility:           string(exam_entities.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            3600,
		TryLimit:             2,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create first exam", err)
	}
	examOneID = examOne.ID
	process.EndStep()

	// [STEP 3] Create another public exam
	process.StartStep("Crear otro examen público")
	examTwo, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Session CRUD Exam Two",
		Description:          "Segundo examen para CRUD de sesión",
		Visibility:           string(exam_entities.VisibilityPublic),
		StartTime:            now.Add(3 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            3600,
		TryLimit:             2,
		ProfessorID:          teacherAccess.UserData.ID,
	})
	if err != nil {
		process.Fail("create second exam", err)
	}
	examTwoID = examTwo.ID
	process.EndStep()

	// [STEP 4] Login as student
	process.StartStep("Iniciar sesión con usuario de estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, studentEmail, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	studentID = studentAccess.UserData.ID
	process.EndStep()

	// [STEP 5] Create a session with the first exam
	process.StartStep("Crear sesión con examen")
	createdSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examOneID,
	})
	if err != nil {
		process.Fail("create session", err)
	}
	sessionID = createdSession.ID
	process.EndStep()

	// [STEP 6] Attempt to create a session with the second exam (expect error)
	process.StartStep("Crear sesión con el otro examen (espera error)")
	_, err = process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examTwoID,
	})
	if err == nil {
		process.Fail("create second active session", fmt.Errorf("expected error when student already has an active session"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Get the active session
	process.StartStep("Obtener la sesión")
	activeSession, err := process.Application.SessionModule.GetActiveSession.Execute(studentCtx, submission_dtos.GetActiveSessionInput{})
	if err != nil {
		process.Fail("get active session", err)
	}
	if activeSession == nil || activeSession.ID != sessionID {
		process.Fail("get active session", fmt.Errorf("expected active session %s", sessionID))
	}
	process.EndStep()

	// [STEP 8] Heartbeat the session
	process.StartStep("Cerrar la sesión")
	closedSession, err := process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionID})
	if err != nil {
		process.Fail("close session", err)
	}
	if closedSession == nil {
		process.Fail("close session", fmt.Errorf("expected closed session response"))
	}
	sessionID = ""
	process.EndStep()

	// [STEP 9] Attempt to get active session after closing (expect error)
	process.StartStep("Obtener la sesión y confirmar cierre")
	_, err = process.Application.SessionModule.GetActiveSession.Execute(studentCtx, submission_dtos.GetActiveSessionInput{})
	if err == nil {
		process.Fail("verify session close", fmt.Errorf("expected no active session after closing"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
