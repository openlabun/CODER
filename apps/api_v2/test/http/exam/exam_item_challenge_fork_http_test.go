package exam_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamItemChallengeForkHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "ExamItem Challenge Fork HTTP")
	creatorEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var creatorAccess *httputils.HTTPAccess
	var observerAccess *httputils.HTTPAccess
	var creatorChallengeID string
	var observerExamID string
	var forkedChallengeID string
	var observerExamItemID string

	defer func() {
		if observerExamItemID != "" && observerAccess != nil {
			t.Logf("[CLEANUP] Eliminando exam item %s", observerExamItemID)
			_ = httputils.DeleteExamItemByID(t, app, observerAccess, observerExamItemID)
		}
		if observerExamID != "" && observerAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", observerExamID)
			_ = httputils.DeleteExamByID(t, app, observerAccess, observerExamID)
		}
		if creatorChallengeID != "" && creatorAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto original %s", creatorChallengeID)
			_ = httputils.DeleteChallengeByID(t, app, creatorAccess, creatorChallengeID)
		}
		if forkedChallengeID != "" && observerAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto forkeado %s", forkedChallengeID)
			_ = httputils.DeleteChallengeByID(t, app, observerAccess, forkedChallengeID)
		}
	}()

	// [STEP 1] Login as creator teacher
	process.StartStep("Iniciar sesion con docente creador")
	creatorAccess = httputils.EnsureAuthUserAccess(t, app, creatorEmail, password, "Teacher Creator")
	process.EndStep()

	// [STEP 2] Create a published challenge as creator
	process.StartStep("Crear challenge publicado del creador")
	resp := httputils.PostChallengeCreate(t, app, creatorAccess, map[string]any{
		"title":               "ExamItem Fork Original",
		"description":         "Challenge original para prueba de fork en exam item",
		"tags":                []string{"exam-item", "fork"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1500,
		"worker_memory_limit": 256,
		"input_variables":     []map[string]any{{"name": "x", "type": "int", "value": "10"}},
		"output_variable":     map[string]any{"name": "y", "type": "int", "value": "10"},
		"constraints":         "1 <= x <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create creator challenge")
	creatorChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Login as observer teacher
	process.StartStep("Iniciar sesion con docente observador")
	observerAccess = httputils.EnsureAuthUserAccess(t, app, observerEmail, password, "Teacher Observer")
	process.EndStep()

	// [STEP 4] Create an exam for the observer teacher
	process.StartStep("Crear examen del observador")
	now := time.Now().UTC()
	resp = httputils.PostExamCreate(t, app, observerAccess, map[string]any{
		"title":                  "ExamItem Fork Exam",
		"description":            "Exam para prueba de fork en exam item",
		"visibility":             "private",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": false,
		"time_limit":             3600,
		"try_limit":              1,
		"professor_id":           observerAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create observer exam")
	observerExamID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Create an exam item in the observer's exam using the creator's challenge (expect the challenge to be forked)
	process.StartStep("Crear exam item con challenge del creador")
	resp = httputils.PostExamItemCreate(t, app, observerAccess, map[string]any{
		"exam_id":      observerExamID,
		"challenge_id": creatorChallengeID,
		"order":        1,
		"points":       100,
	})
	httputils.RequireStatus(t, resp, 201, "create exam item with foreign challenge")
	body := httputils.MustJSONMap(t, resp)
	observerExamItemID = httputils.StringField(body, "id")
	forkedChallengeID = httputils.StringField(body, "challenge_id")
	process.EndStep()

	// [STEP 6] Verify that the forked challenge has the same details as the original challenge
	process.StartStep("Modificar challenge original del creador")
	updatedOriginalTitle := "ExamItem Fork Original Updated"
	updatedOriginalDescription := "Challenge original actualizado tras crear exam item"
	resp = httputils.PatchChallengeUpdate(t, app, creatorAccess, creatorChallengeID, map[string]any{
		"title":       updatedOriginalTitle,
		"description": updatedOriginalDescription,
	})
	httputils.RequireStatus(t, resp, 200, "update creator challenge")
	process.EndStep()

	// [STEP 7] Get the exam item details and verify that the challenge details in the exam item are unchanged (indicating it was forked and not directly linked to original)
	process.StartStep("Obtener exam items y verificar que no cambie el challenge del exam item")
	resp = httputils.GetExamItems(t, app, observerAccess, observerExamID)
	httputils.RequireStatus(t, resp, 200, "get exam items")
	var items []map[string]any
	if err := json.Unmarshal(resp.Body, &items); err != nil {
		process.Fail("get exam items", fmt.Errorf("decode exam items: %w", err))
	}
	if len(items) != 1 {
		process.Fail("get exam items", fmt.Errorf("expected 1 exam item, got %d", len(items)))
	}
	item := items[0]
	if httputils.StringField(item, "id") != observerExamItemID {
		process.Fail("get exam items", fmt.Errorf("unexpected exam item id"))
	}
	challengeRaw, ok := item["challenge"].(map[string]any)
	if !ok || challengeRaw == nil {
		process.Fail("get exam items", fmt.Errorf("expected challenge details in exam item"))
	}
	if httputils.StringField(challengeRaw, "id") == creatorChallengeID {
		process.Fail("verify exam item challenge fork", fmt.Errorf("expected exam item challenge to be forked and have different ID from original"))
	}
	if httputils.StringField(challengeRaw, "title") == updatedOriginalTitle || httputils.StringField(challengeRaw, "description") == updatedOriginalDescription {
		process.Fail("verify exam item challenge fork", fmt.Errorf("expected exam item challenge to remain unchanged after original update"))
	}
	process.EndStep()

	process.End()
}
