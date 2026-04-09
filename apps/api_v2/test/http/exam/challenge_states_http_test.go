package exam_test

import (
	"fmt"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestChallengeStatesHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Challenge States HTTP")
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var challengeID string

	defer func() {
		if challengeID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto %s", challengeID)
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	process.EndStep()

	// [STEP 2] Create a challenge in draft status
	process.StartStep("Crea un reto en estado draft")
	resp := httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Challenge States Test",
		"description":         "Challenge creado para validar transiciones",
		"tags":                []string{"states", "challenge"},
		"status":              "draft",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"input_variables":     []map[string]any{{"name": "n", "type": "int", "value": "10"}},
		"output_variable":     map[string]any{"name": "out", "type": "int", "value": "10"},
		"constraints":         "1 <= n <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create draft challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Update the challenge
	process.StartStep("Actualiza el reto")
	step3Title := "Challenge States Test Updated"
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{"title": step3Title})
	httputils.RequireStatus(t, resp, 200, "update draft challenge")
	process.EndStep()

	// [STEP 4] Get the challenge details and validate the changes
	process.StartStep("Obtiene los datos del reto y valida los cambios")
	resp = httputils.GetChallengeByID(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "get challenge details")
	body := httputils.MustJSONMap(t, resp)
	if httputils.StringField(body, "title") != step3Title || httputils.StringField(body, "status") != "draft" {
		process.Fail("get challenge details", fmt.Errorf("unexpected challenge state after step 4"))
	}
	process.EndStep()

	// [STEP 5] Update the challenge to published
	process.StartStep("Actualiza el reto a estado published")
	resp = httputils.PostChallengePublish(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "publish challenge")
	process.EndStep()

	// [STEP 6] Attempt to update the challenge (expect error)
	process.StartStep("Actualiza el reto (espera error)")
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{"status": "draft"})
	if resp.StatusCode == 200 {
		process.Fail("invalid transition published->draft", fmt.Errorf("expected error on invalid transition"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Update the challenge to private
	process.StartStep("Actualiza el reto a estado private")
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{"status": "private"})
	httputils.RequireStatus(t, resp, 200, "transition published->private")
	process.EndStep()

	// [STEP 8] Attempt to update the challenge to draft (expect error)
	process.StartStep("Actualiza el reto (espera error)")
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{"status": "draft"})
	if resp.StatusCode == 200 {
		process.Fail("invalid transition private->draft", fmt.Errorf("expected error on invalid transition"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 9] Update the challenge back to published
	process.StartStep("Actualiza el reto a estado archived")
	resp = httputils.PostChallengeArchive(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "transition to archived")
	process.EndStep()

	// [STEP 10] Attempt to update the challenge to private (expect error)
	process.StartStep("Actualiza el reto (espera error)")
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{"status": "private"})
	if resp.StatusCode == 200 {
		process.Fail("invalid transition archived->private", fmt.Errorf("expected error on invalid transition"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 11] Update the challenge back to published
	process.StartStep("Actualiza el reto a estado published")
	resp = httputils.PostChallengePublish(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "transition archived->published")
	process.EndStep()

	// [STEP 12] Attempt to update the challenge to draft (expect error)
	process.StartStep("Actualiza el reto a estado archived")
	resp = httputils.PostChallengeArchive(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "transition published->archived")
	process.EndStep()

	// [STEP 13] Attempt to update the challenge to published (expect error)
	process.StartStep("Actualiza el reto a estado private")
	resp = httputils.PostChallengePublish(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "bridge archived->published")
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{"status": "private"})
	httputils.RequireStatus(t, resp, 200, "transition published->private")
	process.EndStep()

	// [STEP 14] Attempt to update the challenge to draft (expect error)
	process.StartStep("Actualiza el reto a estado draft (espera error)")
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{"status": "draft"})
	if resp.StatusCode == 200 {
		process.Fail("invalid transition private->draft", fmt.Errorf("expected error on invalid transition"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 15] Update the challenge back to published
	process.StartStep("Actualiza el reto a estado published")
	resp = httputils.PostChallengePublish(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "transition private->published")
	process.EndStep()

	// [STEP 16] Attempt to update the challenge to draft (expect error)
	process.StartStep("Actualiza el reto a estado draft (espera error)")
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{"status": "draft"})
	if resp.StatusCode == 200 {
		process.Fail("invalid transition published->draft", fmt.Errorf("expected error on invalid transition"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
