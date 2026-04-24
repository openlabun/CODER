package submission_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamTryLimitHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Exam Try Limit HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var examID string
	var sessionOneID string
	var sessionTwoID string

	defer func() {
		if sessionTwoID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando segunda sesion %s", sessionTwoID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, sessionTwoID)
		}
		if sessionOneID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Cerrando primera sesion %s", sessionOneID)
			_ = httputils.PostSessionClose(t, app, teacherAccess, sessionOneID)
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

	// [STEP 2] Create a public exam with try_limit = 2
	process.StartStep("Crear examen publico (con solo 2 intentos try_limit)")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Try Limit Exam HTTP",
		"description":            "Examen para validar limite de intentos por sesiones",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             120,
		"try_limit":              2,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Login as student
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	process.EndStep()

	// [STEP 4] Create first session for the student
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create first session")
	sessionOneID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Close the first session
	process.StartStep("Cerrar la sesion")
	resp = httputils.PostSessionClose(t, app, teacherAccess, sessionOneID)
	httputils.RequireStatus(t, resp, 200, "close first session")
	sessionOneID = ""
	process.EndStep()

	// [STEP 6] Create second session for the student
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create second session")
	sessionTwoID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 7] Close the second session
	process.StartStep("Cerrar la sesion")
	resp = httputils.PostSessionClose(t, app, teacherAccess, sessionTwoID)
	httputils.RequireStatus(t, resp, 200, "close second session")
	sessionTwoID = ""
	process.EndStep()

	// [STEP 8] Create third session for the student - should fail with try limit exceeded error
	process.StartStep("Crear una sesion en el examen (espera error)")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	if resp.StatusCode == 201 {
		process.Fail("create third session", fmt.Errorf("expected error when try_limit is exceeded"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
