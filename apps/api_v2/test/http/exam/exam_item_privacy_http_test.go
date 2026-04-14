package exam_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamItemPrivacyHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "ExamItem Privacy HTTP")
	ownerEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var ownerAccess *httputils.HTTPAccess
	var observerAccess *httputils.HTTPAccess
	var examID string
	var challengeID string
	var examItemID string

	defer func() {
		if examItemID != "" && ownerAccess != nil {
			t.Logf("[CLEANUP] Eliminando punto de examen %s", examItemID)
			_ = httputils.DeleteExamItemByID(t, app, ownerAccess, examItemID)
		}
		if examID != "" && ownerAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_ = httputils.DeleteExamByID(t, app, ownerAccess, examID)
		}
		if challengeID != "" && ownerAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto %s", challengeID)
			_ = httputils.DeleteChallengeByID(t, app, ownerAccess, challengeID)
		}
	}()

	// [STEP 1] Login as teacher owner
	process.StartStep("Iniciar sesion con docente dueno")
	ownerAccess = httputils.EnsureAuthUserAccess(t, app, ownerEmail, password, "Teacher Owner")
	process.EndStep()

	// [STEP 2] Create an exam and a challenge as owner
	process.StartStep("Crear exam y challenge del dueno")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, ownerAccess, map[string]any{
		"title":                  "ExamItem Privacy Exam",
		"description":            "Exam para privacidad",
		"visibility":             "private",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": false,
		"time_limit":             3600,
		"try_limit":              1,
		"professor_id":           ownerAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create owner exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	resp = httputils.PostChallengeCreate(t, app, ownerAccess, map[string]any{
		"title":               "ExamItem Privacy Challenge",
		"description":         "Challenge para privacidad",
		"tags":                []string{"exam-item", "privacy"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables":     []map[string]any{{"name": "x", "type": "int", "value": "7"}},
		"output_variable":     map[string]any{"name": "y", "type": "int", "value": "7"},
		"constraints":         "x >= 0",
	})
	httputils.RequireStatus(t, resp, 201, "create owner challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create an exam item in the owner's exam
	process.StartStep("Crear exam item del dueno")
	resp = httputils.PostExamItemCreate(t, app, ownerAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        1,
		"points":       100,
	})
	httputils.RequireStatus(t, resp, 201, "create owner exam item")
	examItemID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 4] Try to create an exam item in the observer's exam using the owner's private challenge (expect error)
	process.StartStep("Iniciar sesion con docente observador")
	observerAccess = httputils.EnsureAuthUserAccess(t, app, observerEmail, password, "Teacher Observer")
	process.EndStep()

	// [STEP 5] Create an exam for the observer teacher
	process.StartStep("Crear exam item en exam de otro docente (espera error)")
	resp = httputils.PostExamItemCreate(t, app, observerAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        2,
		"points":       100,
	})
	if resp.StatusCode == 201 {
		process.Fail("observer create exam item", fmt.Errorf("expected error when observer creates exam item in foreign exam"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 6] Try to update the owner's exam item as observer (expect error)
	process.StartStep("Actualizar exam item de otro docente (espera error)")
	resp = httputils.PatchExamItemUpdate(t, app, observerAccess, examItemID, map[string]any{"points": 200})
	if resp.StatusCode == 200 {
		process.Fail("observer update exam item", fmt.Errorf("expected error when observer updates foreign exam item"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Try to delete the owner's exam item as observer (expect error)
	process.StartStep("Eliminar exam item de otro docente (espera error)")
	resp = httputils.DeleteExamItemByID(t, app, observerAccess, examItemID)
	if resp.StatusCode == 200 {
		process.Fail("observer delete exam item", fmt.Errorf("expected error when observer deletes foreign exam item"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
