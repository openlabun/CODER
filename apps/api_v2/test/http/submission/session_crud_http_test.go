package submission_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSessionCRUDHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Session CRUD HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var studentID string
	var examOneID string
	var examTwoID string
	var sessionID string

	defer func() {
		if sessionID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando sesion %s", sessionID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, sessionID)
		}
		if examTwoID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", examTwoID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examTwoID)
		}
		if examOneID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", examOneID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examOneID)
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	process.EndStep()

	// [STEP 2] Create two exams
	process.StartStep("Crear examen publico (visibilidad public y sin curso)")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Session CRUD Exam One",
		"description":            "Primer examen para CRUD de sesion",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3600,
		"try_limit":              2,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create first exam")
	examOneID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create another public exam
	process.StartStep("Crear otro examen publico")
	resp = httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Session CRUD Exam Two",
		"description":            "Segundo examen para CRUD de sesion",
		"visibility":             "public",
		"start_time":             now.Add(3 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3600,
		"try_limit":              2,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create second exam")
	examTwoID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 4] Login as student
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	studentID = studentAccess.UserID
	process.EndStep()

	// [STEP 5] Create a session for the first exam
	process.StartStep("Crear sesion con examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentID, "exam_id": examOneID})
	httputils.RequireStatus(t, resp, 201, "create session")
	sessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 6] Try to create another session with the second exam (expect error since student already has an active session)
	process.StartStep("Crear sesion con el otro examen (espera error)")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentID, "exam_id": examTwoID})
	if resp.StatusCode == 201 {
		process.Fail("create second active session", fmt.Errorf("expected error when student already has an active session"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Get the active session and validate data
	process.StartStep("Obtener la sesion")
	resp = httputils.GetActiveSession(t, app, studentAccess, "")
	httputils.RequireStatus(t, resp, 200, "get active session")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "id") != sessionID {
		process.Fail("get active session", fmt.Errorf("expected active session %s", sessionID))
	}
	process.EndStep()

	// [STEP 8] Try to get the active session with another student (expect error)
	process.StartStep("Cerrar la sesion")
	resp = httputils.PostSessionClose(t, app, teacherAccess, sessionID)
	httputils.RequireStatus(t, resp, 200, "close session")
	sessionID = ""
	process.EndStep()

	// [STEP 9] Try to get the active session after closing (expect error)
	process.StartStep("Obtener la sesion y confirmar cierre")
	resp = httputils.GetActiveSession(t, app, studentAccess, "")
	if resp.StatusCode == 200 {
		process.Fail("verify session close", fmt.Errorf("expected no active session after closing"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
