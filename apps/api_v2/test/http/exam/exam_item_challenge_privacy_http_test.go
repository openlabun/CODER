package exam_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamItemChallengePrivacyHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "ExamItem Challenge Privacy HTTP")
	creatorEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var creatorAccess *httputils.HTTPAccess
	var observerAccess *httputils.HTTPAccess
	var examID string
	var privateChallengeID string
	var publishedChallengeID string
	var ownChallengeID string
	var forkedChallengeID string

	defer func() {
		if examID != "" && observerAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen %s", examID)
			_ = httputils.DeleteExamByID(t, app, observerAccess, examID)
		}
		if ownChallengeID != "" && observerAccess != nil {
			t.Logf("[CLEANUP] Eliminando challenge propio %s", ownChallengeID)
			_ = httputils.DeleteChallengeByID(t, app, observerAccess, ownChallengeID)
		}
		if privateChallengeID != "" && creatorAccess != nil {
			t.Logf("[CLEANUP] Eliminando challenge privado %s", privateChallengeID)
			_ = httputils.DeleteChallengeByID(t, app, creatorAccess, privateChallengeID)
		}
		if publishedChallengeID != "" && creatorAccess != nil {
			t.Logf("[CLEANUP] Eliminando challenge publicado %s", publishedChallengeID)
			_ = httputils.DeleteChallengeByID(t, app, creatorAccess, publishedChallengeID)
		}
		if forkedChallengeID != "" && observerAccess != nil {
			t.Logf("[CLEANUP] Eliminando challenge forkeado %s", forkedChallengeID)
			_ = httputils.DeleteChallengeByID(t, app, observerAccess, forkedChallengeID)
		}
	}()

	// [STEP 1] Login as creator teacher
	process.StartStep("Iniciar sesion con docente creador")
	creatorAccess = httputils.EnsureAuthUserAccess(t, app, creatorEmail, password, "Teacher Creator")
	process.EndStep()

	// [STEP 2] Create a private challenge and a published challenge as creator
	process.StartStep("Crear challenge privado y challenge publicado del creador")
	resp := httputils.PostChallengeCreate(t, app, creatorAccess, map[string]any{
		"title":               "Creator Private Challenge",
		"description":         "Challenge privado",
		"tags":                []string{"privacy", "private"},
		"status":              "private",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables":     []map[string]any{{"name": "x", "type": "int", "value": "1"}},
		"output_variable":     map[string]any{"name": "y", "type": "int", "value": "1"},
		"constraints":         "x >= 0",
	})
	httputils.RequireStatus(t, resp, 201, "create private challenge")
	privateChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	resp = httputils.PostChallengeCreate(t, app, creatorAccess, map[string]any{
		"title":               "Creator Published Challenge",
		"description":         "Challenge publicado",
		"tags":                []string{"privacy", "published"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables":     []map[string]any{{"name": "x", "type": "int", "value": "2"}},
		"output_variable":     map[string]any{"name": "y", "type": "int", "value": "2"},
		"constraints":         "x >= 0",
	})
	httputils.RequireStatus(t, resp, 201, "create published challenge")
	publishedChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Login as observer teacher
	process.StartStep("Iniciar sesion con docente observador")
	observerAccess = httputils.EnsureAuthUserAccess(t, app, observerEmail, password, "Teacher Observer")
	process.EndStep()

	// [STEP 4] Create an exam for the observer teacher
	process.StartStep("Crear examen del observador")
	now := time.Now().UTC()
	resp = httputils.PostExamCreate(t, app, observerAccess, map[string]any{
		"title":                  "ExamItem Challenge Privacy Exam",
		"description":            "Exam de observador",
		"visibility":             "private",
		"start_time":             now.Add(3 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": false,
		"time_limit":             3600,
		"try_limit":              1,
		"professor_id":           observerAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create observer exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Create an exam item in the observer's exam using the creator's private challenge (expect error)
	process.StartStep("Crear challenge propio del observador")
	resp = httputils.PostChallengeCreate(t, app, observerAccess, map[string]any{
		"title":               "Observer Own Challenge",
		"description":         "Challenge propio del observador",
		"tags":                []string{"privacy", "own"},
		"status":              "private",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables":     []map[string]any{{"name": "x", "type": "int", "value": "3"}},
		"output_variable":     map[string]any{"name": "y", "type": "int", "value": "3"},
		"constraints":         "x >= 0",
	})
	httputils.RequireStatus(t, resp, 201, "create own challenge")
	ownChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 6] Create an exam item in the observer's exam using the creator's private challenge (expect error)
	process.StartStep("Crear exam item con challenge propio (ok)")
	resp = httputils.PostExamItemCreate(t, app, observerAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": ownChallengeID,
		"order":        1,
		"points":       100,
	})
	httputils.RequireStatus(t, resp, 201, "create exam item with own challenge")
	process.EndStep()

	// [STEP 7] Create an exam item in the observer's exam using the creator's private challenge (expect error)
	process.StartStep("Crear exam item con challenge privado de otro docente (espera error)")
	resp = httputils.PostExamItemCreate(t, app, observerAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": privateChallengeID,
		"order":        2,
		"points":       100,
	})
	if resp.StatusCode == 201 {
		process.Fail("create exam item with foreign private challenge", fmt.Errorf("expected error when using private challenge from another teacher"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 8] Create an exam item in the observer's exam using the creator's published challenge (expect it to succeed and fork the challenge)
	process.StartStep("Crear exam item con challenge publicado de otro docente (ok)")
	resp = httputils.PostExamItemCreate(t, app, observerAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": publishedChallengeID,
		"order":        3,
		"points":       120,
	})
	httputils.RequireStatus(t, resp, 201, "create exam item with foreign published challenge")
	forkedChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "challenge_id")
	process.Log(fmt.Sprintf("Original challenge ID: %s", publishedChallengeID))
	process.Log(fmt.Sprintf("Forked challenge ID: %s", forkedChallengeID))
	process.EndStep()

	process.End()
}
