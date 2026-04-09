package submission_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSessionHeartbeatHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Session Heartbeat HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var studentID string
	var examID string
	var sessionID string
	var heartbeatBefore string

	defer func() {
		if sessionID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando sesion %s", sessionID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, sessionID)
		}
		if examID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
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
		"title":                  "Session Heartbeat Exam",
		"description":            "Examen para prueba de heartbeat",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3600,
		"try_limit":              2,
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
	httputils.RequireStatus(t, resp, 201, "create session")
	sessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Get the session and check last heartbeat
	process.StartStep("Obtener la sesion")
	resp = httputils.GetActiveSession(t, app, studentAccess, "")
	httputils.RequireStatus(t, resp, 200, "get active session before heartbeat")
	heartbeatBefore = httputils.StringField(httputils.MustJSONMap(t, resp), "last_heartbeat")
	process.EndStep()

	// [STEP 6] Send heartbeat to the session
	process.StartStep("Hacer heartbeat a la sesion")
	resp = httputils.PostSessionHeartbeat(t, app, studentAccess, sessionID)
	httputils.RequireStatus(t, resp, 200, "heartbeat session")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "id") != sessionID {
		process.Fail("heartbeat session", fmt.Errorf("expected heartbeat response for session %s", sessionID))
	}
	process.EndStep()

	// [STEP 7] Get the session again and check last heartbeat was updated
	process.StartStep("Obtener la sesion y validar que este activa")
	resp = httputils.GetActiveSession(t, app, studentAccess, "")
	httputils.RequireStatus(t, resp, 200, "get active session after heartbeat")
	after := httputils.StringField(httputils.MustJSONMap(t, resp), "last_heartbeat")
	if after == "" || after < heartbeatBefore {
		process.Fail("verify heartbeat update", fmt.Errorf("expected heartbeat timestamp to be updated"))
	}
	process.EndStep()

	process.End()
}
