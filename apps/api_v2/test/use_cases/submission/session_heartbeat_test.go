package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestSessionHeartbeat(t *testing.T) {
	process := test.StartTestWithApp(t, "Session Heartbeat")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var studentID string
	var examID string
	var sessionID string
	var heartbeatBefore time.Time

	defer func() {
		if sessionID != "" {
			t.Logf("[CLEANUP] Cerrando sesión %s", sessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: sessionID})
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

	// [STEP 2] Create public exam (visibility public and no course)
	process.StartStep("Crear examen público (visibilidad public y sin curso)")
	now := time.Now().UTC()
	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Session Heartbeat Exam",
		Description:          "Examen para prueba de heartbeat",
		Visibility:           string(exam_consts.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            3600,
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

	// [STEP 4] Create session with exam
	process.StartStep("Crear sesión con examen")
	session, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create session", err)
	}
	sessionID = session.ID
	process.EndStep()

	// [STEP 5] Get the session
	process.StartStep("Obtener la sesión")
	loadedSession, err := process.Application.SessionModule.GetActiveSession.Execute(studentCtx, submission_dtos.GetActiveSessionInput{})
	if err != nil {
		process.Fail("get active session before heartbeat", err)
	}
	heartbeatBefore = loadedSession.LastHeartbeat
	process.EndStep()

	// [STEP 6] Heartbeat the session
	process.StartStep("Hacer heartbeat a la sesión")
	heartbeatedSession, err := process.Application.SessionModule.HeartBeatSession.Execute(studentCtx, submission_dtos.HeartbeatSessionInput{SessionID: sessionID})
	if err != nil {
		process.Fail("heartbeat session", err)
	}
	if heartbeatedSession == nil || heartbeatedSession.ID != sessionID {
		process.Fail("heartbeat session", fmt.Errorf("expected heartbeat response for session %s", sessionID))
	}
	process.EndStep()

	// [STEP 7] Get the session again and verify heartbeat timestamp is updated
	process.StartStep("Obtener la sesión y validar que esté activa")
	loadedAfterHeartbeat, err := process.Application.SessionModule.GetActiveSession.Execute(studentCtx, submission_dtos.GetActiveSessionInput{})
	if err != nil {
		process.Fail("get active session after heartbeat", err)
	}
	if loadedAfterHeartbeat.LastHeartbeat.Before(heartbeatBefore) {
		process.Fail("verify heartbeat update", fmt.Errorf("expected heartbeat timestamp to be updated"))
	}
	process.EndStep()

	process.End()
}
