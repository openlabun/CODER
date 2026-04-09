package submission_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSessionFreezeAndBlockHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Session Freeze and Block HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var studentID string
	var examID string
	var firstSessionID string
	var secondSessionID string

	defer func() {
		if secondSessionID != "" && teacherAccess != nil {
			_ = httputils.PostSessionClose(t, app, teacherAccess, secondSessionID)
		}
		if firstSessionID != "" && teacherAccess != nil {
			_ = httputils.PostSessionClose(t, app, teacherAccess, firstSessionID)
		}
		if examID != "" && teacherAccess != nil {
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examID)
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	process.EndStep()

	// [STEP 2] Create an exam
	process.StartStep("Crear examen publico (visibilidad public y sin curso)")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Session Freeze Block Exam",
		"description":            "Examen para bloqueo y congelamiento",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3600,
		"try_limit":              3,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Login as student
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	studentID = studentAccess.UserID
	process.EndStep()

	// [STEP 4] Create a session with the exam
	process.StartStep("Crear sesion con examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create first session")
	firstSessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Block the session from teacher view
	process.StartStep("Bloquear sesion desde cuenta de docente")
	resp = httputils.PostSessionBlock(t, app, teacherAccess, firstSessionID)
	httputils.RequireStatus(t, resp, 200, "block session")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "status") != "blocked" {
		process.Fail("block session", fmt.Errorf("expected blocked session status"))
	}
	firstSessionID = ""
	process.EndStep()

	// [STEP 6] Try to get the session and expect it to be blocked
	process.StartStep("Obtener la sesion")
	resp = httputils.GetActiveSession(t, app, studentAccess, "")
	if resp.StatusCode == 200 {
		process.Fail("get active session after block", fmt.Errorf("expected no active session after block"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Create another session with the same exam
	process.StartStep("Crear sesion con examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create second session")
	secondSessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 8] Wait for freeze time to elapse and check session is frozen
	process.StartStep("Esperar tiempo para congelamiento de examen")
	freezeTime, err := strconv.Atoi(os.Getenv("SESSION_FREEZE_TIME"))
	if err != nil || freezeTime <= 0 {
		freezeTime = 60
	}
	process.Log(fmt.Sprintf("Tiempo de congelamiento configurado: %d segundos", freezeTime))
	time.Sleep(time.Duration(freezeTime) * time.Second)
	process.EndStep()

	// [STEP 9] Get the session and check it is frozen
	process.StartStep("Obtener la sesion y comprobar que esta congelada")
	resp = httputils.GetActiveSession(t, app, studentAccess, studentAccess.UserID)
	httputils.RequireStatus(t, resp, 200, "get active session")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "status") != "frozen" {
		process.Fail("verify frozen session", fmt.Errorf("expected frozen session status"))
	}
	process.EndStep()

	// [STEP 10] Send heartbeat to reactivate session
	process.StartStep("Hacer heartbeat")
	resp = httputils.PostSessionHeartbeat(t, app, studentAccess, secondSessionID)
	httputils.RequireStatus(t, resp, 200, "second heartbeat")
	process.EndStep()

	// [STEP 11] Get the session and check it is active again
	process.StartStep("Comprobar que se volvio a activar")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "status") != "active" {
		process.Fail("verify reactivation", fmt.Errorf("expected session to be active after heartbeat"))
	}
	process.EndStep()

	// [STEP 12] Block the session again from teacher view
	process.StartStep("Bloquear sesion desde vista de docente")
	resp = httputils.PostSessionBlock(t, app, teacherAccess, secondSessionID)
	httputils.RequireStatus(t, resp, 200, "block reactivated session")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "status") != "blocked" {
		process.Fail("block reactivated session", fmt.Errorf("expected blocked status after teacher block"))
	}
	process.EndStep()

	// [STEP 13] Try to get the session and expect it to be blocked
	process.StartStep("Obtener la sesion y comprobar que esta bloqueada")
	resp = httputils.GetActiveSession(t, app, studentAccess, "")
	if resp.StatusCode == 200 {
		process.Fail("verify blocked session", fmt.Errorf("expected no active session after blocking from teacher view"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	secondSessionID = ""
	process.EndStep()

	process.End()
}
