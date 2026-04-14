package exam_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamItemCRUDHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "ExamItem CRUD HTTP")
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var examID string
	var challengeID string
	var examItemID string
	var deletedExamItemID string

	defer func() {
		if examItemID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando punto de examen %s", examItemID)
			_ = httputils.DeleteExamItemByID(t, app, teacherAccess, examItemID)
		}
		if examID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examID)
		}
		if challengeID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto %s", challengeID)
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	process.EndStep()

	// [STEP 2] Create an exam
	process.StartStep("Crear un examen")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "ExamItem CRUD Exam",
		"description":            "Exam auxiliar para CRUD de exam item",
		"visibility":             "private",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": false,
		"time_limit":             3600,
		"try_limit":              1,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create a challenge
	process.StartStep("Crear un reto")
	resp = httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "ExamItem CRUD Challenge",
		"description":         "Challenge auxiliar para exam item",
		"tags":                []string{"exam-item", "crud"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables":     []map[string]any{{"name": "n", "type": "int", "value": "10"}},
		"output_variable":     map[string]any{"name": "out", "type": "int", "value": "10"},
		"constraints":         "1 <= n <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 4] Create an exam item
	process.StartStep("Crear un punto de examen")
	resp = httputils.PostExamItemCreate(t, app, teacherAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        1,
		"points":       100,
	})
	httputils.RequireStatus(t, resp, 201, "create exam item")
	examItemID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Try to create another exam item with the same challenge (expect error)
	process.StartStep("Crear otro punto de examen con mismo reto (espera error)")
	resp = httputils.PostExamItemCreate(t, app, teacherAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": challengeID,
		"order":        2,
		"points":       50,
	})
	if resp.StatusCode == 201 {
		process.Fail("duplicate exam item", fmt.Errorf("expected error when adding same challenge twice"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 6] Update exam item
	process.StartStep("Actualizar punto de examen")
	resp = httputils.PatchExamItemUpdate(t, app, teacherAccess, examItemID, map[string]any{
		"order":  3,
		"points": 150,
	})
	httputils.RequireStatus(t, resp, 200, "update exam item")
	body := httputils.MustJSONMap(t, resp)
	if int(body["order"].(float64)) != 3 || int(body["points"].(float64)) != 150 {
		process.Fail("update exam item", fmt.Errorf("expected updated exam item values"))
	}
	process.EndStep()

	// [STEP 7] Get exam items and validate the updated item is correct
	process.StartStep("Obtener punto de examen y validar datos")
	resp = httputils.GetExamItems(t, app, teacherAccess, examID)
	httputils.RequireStatus(t, resp, 200, "get exam items")
	var items []map[string]any
	if err := json.Unmarshal(resp.Body, &items); err != nil {
		process.Fail("get exam items", fmt.Errorf("decode exam items: %w", err))
	}
	found := false
	for _, item := range items {
		if httputils.StringField(item, "id") == examItemID {
			found = true
			if int(item["order"].(float64)) != 3 || int(item["points"].(float64)) != 150 {
				process.Fail("get exam items", fmt.Errorf("unexpected exam item values"))
			}
		}
	}
	if !found {
		process.Fail("get exam items", fmt.Errorf("expected exam item %s in exam items list", examItemID))
	}
	process.EndStep()

	// [STEP 8] Delete exam item
	process.StartStep("Eliminar punto de examen")
	resp = httputils.DeleteExamItemByID(t, app, teacherAccess, examItemID)
	httputils.RequireStatus(t, resp, 200, "delete exam item")
	deletedExamItemID = examItemID
	examItemID = ""
	process.EndStep()

	// [STEP 9] Verify exam item deletion
	process.StartStep("Verificar eliminacion")
	resp = httputils.GetExamItems(t, app, teacherAccess, examID)
	httputils.RequireStatus(t, resp, 200, "get exam items after delete")
	items = nil
	if err := json.Unmarshal(resp.Body, &items); err != nil {
		process.Fail("get exam items after delete", fmt.Errorf("decode exam items after delete: %w", err))
	}
	for _, item := range items {
		if httputils.StringField(item, "id") == deletedExamItemID {
			process.Fail("verify exam item deletion", fmt.Errorf("exam item %s should not exist after deletion", deletedExamItemID))
		}
	}
	process.EndStep()

	process.End()
}
