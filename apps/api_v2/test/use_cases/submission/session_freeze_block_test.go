package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"
	"os"
	"strconv"

	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	submission_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submission_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestSessionFreezeAndBlock(t *testing.T) {
	process := test.StartTestWithApp(t, "Session Freeze and Block")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var studentID string
	var examID string
	var firstSessionID string
	var secondSessionID string

	defer func() {
		if secondSessionID != "" {
			t.Logf("[CLEANUP] Cerrando segunda sesión %s", secondSessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: secondSessionID})
		}
		if firstSessionID != "" {
			t.Logf("[CLEANUP] Cerrando primera sesión %s", firstSessionID)
			_, _ = process.Application.SessionModule.CloseSession.Execute(teacherCtx, submission_dtos.CloseSessionInput{SessionID: firstSessionID})
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
		Title:                "Session Freeze Block Exam",
		Description:          "Examen para pruebas de bloqueo y congelamiento",
		Visibility:           string(exam_entities.VisibilityPublic),
		StartTime:            now.Add(2 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            3600,
		TryLimit:             3,
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
	firstSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create first session", err)
	}
	firstSessionID = firstSession.ID
	process.EndStep()

	// [STEP 5] Block session from teacher account
	process.StartStep("Bloquear sesión desde cuenta de docente")
	blockedSession, err := process.Application.SessionModule.BlockSession.Execute(teacherCtx, submission_dtos.BlockSessionInput{SessionID: firstSessionID})
	if err != nil {
		process.Fail("block session", err)
	}
	if blockedSession == nil || blockedSession.Status != submission_entities.SessionStatusBlocked {
		process.Fail("block session", fmt.Errorf("expected blocked session status"))
	}
	firstSessionID = ""
	process.EndStep()

	// [STEP 6] Try to get the session (expect error because it's blocked)
	process.StartStep("Obtener la sesión")
	_, err = process.Application.SessionModule.GetActiveSession.Execute(studentCtx, submission_dtos.GetActiveSessionInput{})
	if err == nil {
		process.Fail("get active session after block", fmt.Errorf("expected no active session after block"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Create another session to test freeze (should succeed because block is on previous session)
	process.StartStep("Crear sesión con examen")
	secondSession, err := process.Application.SessionModule.CreateSession.Execute(studentCtx, submission_dtos.CreateSessionInput{
		UserID: studentID,
		ExamID: examID,
	})
	if err != nil {
		process.Fail("create second session", err)
	}
	secondSessionID = secondSession.ID
	process.EndStep()

	// [STEP 8] Wait for time to trigger freeze
	process.StartStep("Esperar tiempo para congelamiento de examen")
	freeze_time, err :=  strconv.Atoi(os.Getenv("SESSION_FREEZE_TIME"))
	if err != nil {
		freeze_time = 60 // default freeze time if env variable is not set or invalid
	}
	process.Log(fmt.Sprintf("Tiempo de congelamiento configurado: %d segundos", freeze_time))
	freeze_time_duration := time.Duration(freeze_time) * time.Second
	time.Sleep(freeze_time_duration)
	process.EndStep()

	// [STEP 9] Get the session and verify it's frozen
	process.StartStep("Obtener la sesión y comprobar que está congelada")
	frozenSession, err := process.Application.SessionModule.HeartBeatSession.Execute(studentCtx, submission_dtos.HeartbeatSessionInput{SessionID: secondSessionID})
	if err != nil {
		process.Fail("heartbeat to freeze session", err)
	}
	if frozenSession == nil || frozenSession.Status != submission_entities.SessionStatusFrozen {
		process.Fail("verify frozen session", fmt.Errorf("expected frozen session status"))
	}
	process.EndStep()

	// [STEP 10] Try to heartbeat frozen session (expect it to reactivate)
	process.StartStep("Hacer heartbeat")
	reheartbeatedSession, err := process.Application.SessionModule.HeartBeatSession.Execute(studentCtx, submission_dtos.HeartbeatSessionInput{SessionID: secondSessionID})
	if err != nil {
		process.Fail("second heartbeat", err)
	}
	if reheartbeatedSession == nil {
		process.Fail("second heartbeat", fmt.Errorf("expected session response on second heartbeat"))
	}
	process.EndStep()

	// [STEP 11] Verify that the session is active again after heartbeat
	process.StartStep("Comprobar que se volvió a activar")
	if reheartbeatedSession.Status != submission_entities.SessionStatusActive {
		process.Log(fmt.Sprintf("Estado actual tras segundo heartbeat: %s", reheartbeatedSession.Status))
		process.Fail("verify reactivation", fmt.Errorf("expected session to be active after heartbeat"))
	}
	process.EndStep()

	// [STEP 12] Block session from teacher view
	process.StartStep("Bloquear sesión desde vista de docente")
	blockedAgainSession, err := process.Application.SessionModule.BlockSession.Execute(teacherCtx, submission_dtos.BlockSessionInput{SessionID: secondSessionID})
	if err != nil {
		process.Fail("block reactivated session", err)
	}
	if blockedAgainSession == nil || blockedAgainSession.Status != submission_entities.SessionStatusBlocked {
		process.Fail("block reactivated session", fmt.Errorf("expected blocked status after teacher block"))
	}
	process.EndStep()

	// [STEP 13] Get session and verify blocked state
	process.StartStep("Obtener la sesión y comprobar que está bloqueada")
	_, err = process.Application.SessionModule.GetActiveSession.Execute(studentCtx, submission_dtos.GetActiveSessionInput{})
	if err == nil {
		process.Fail("verify blocked session", fmt.Errorf("expected no active session after blocking from teacher view"))
	}
	process.Log("Recibió ERROR, como se esperaba")
	secondSessionID = ""
	process.EndStep()

	process.End()
}
