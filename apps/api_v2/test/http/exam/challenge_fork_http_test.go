package exam_test

import (
	"fmt"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestChallengeForkHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Challenge Fork HTTP")
	creatorEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var creatorAccess *httputils.HTTPAccess
	var observerAccess *httputils.HTTPAccess
	var originalChallengeID string
	var forkedChallengeID string

	defer func() {
		if forkedChallengeID != "" && observerAccess != nil {
			t.Logf("[CLEANUP] Eliminando challenge fork %s", forkedChallengeID)
			_ = httputils.DeleteChallengeByID(t, app, observerAccess, forkedChallengeID)
		}
		if originalChallengeID != "" && creatorAccess != nil {
			t.Logf("[CLEANUP] Eliminando challenge original %s", originalChallengeID)
			_ = httputils.DeleteChallengeByID(t, app, creatorAccess, originalChallengeID)
		}
	}()

	// [STEP 1] Login as creator teacher
	process.StartStep("Iniciar sesion con usuario de docente (creador)")
	creatorAccess = httputils.EnsureAuthUserAccess(t, app, creatorEmail, password, "Teacher Creator")
	process.Log(fmt.Sprintf("creatorID=%s", creatorAccess.UserID))
	process.EndStep()

	// [STEP 2] Create a challenge in private status
	process.StartStep("Crea un reto en estado private")
	resp := httputils.PostChallengeCreate(t, app, creatorAccess, map[string]any{
		"title":               "Challenge Fork Original",
		"description":         "Challenge original para prueba de fork",
		"tags":                []string{"fork", "original"},
		"status":              "private",
		"difficulty":          "easy",
		"worker_time_limit":   1500,
		"worker_memory_limit": 256,
		"input_variables":     []map[string]any{{"name": "x", "type": "int", "value": "5"}},
		"output_variable":     map[string]any{"name": "y", "type": "int", "value": "5"},
		"constraints":         "1 <= x <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create original challenge")
	originalChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Login as observer teacher
	process.StartStep("Iniciar sesion con usuario de docente (observador)")
	observerAccess = httputils.EnsureAuthUserAccess(t, app, observerEmail, password, "Teacher Observer")
	process.Log(fmt.Sprintf("observerID=%s", observerAccess.UserID))
	process.EndStep()

	// [STEP 4] Attempt to fork the private challenge (expect error)
	process.StartStep("Hace fork al reto (espera error)")
	resp = httputils.PostChallengeFork(t, app, observerAccess, originalChallengeID)
	if resp.StatusCode == 201 {
		process.Fail("fork private challenge", fmt.Errorf("expected error when forking private challenge"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 5] Publish the original challenge
	process.StartStep("Actualiza reto original a estado published")
	resp = httputils.PostChallengePublish(t, app, creatorAccess, originalChallengeID)
	httputils.RequireStatus(t, resp, 200, "publish original challenge")
	process.EndStep()

	// [STEP 6] Fork the published challenge
	process.StartStep("Hace fork al reto")
	resp = httputils.PostChallengeFork(t, app, observerAccess, originalChallengeID)
	httputils.RequireStatus(t, resp, 201, "fork published challenge")
	forkedChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	if forkedChallengeID == "" || forkedChallengeID == originalChallengeID {
		process.Fail("fork published challenge", fmt.Errorf("forked challenge must have a different ID"))
	}
	process.EndStep()

	// [STEP 7] Update the forked challenge and verify original challenge is unchanged
	process.StartStep("Actualiza el reto copiado")
	updatedForkTitle := "Challenge Fork Copied Updated"
	updatedForkDescription := "Reto fork actualizado"
	resp = httputils.PatchChallengeUpdate(t, app, observerAccess, forkedChallengeID, map[string]any{
		"title":       updatedForkTitle,
		"description": updatedForkDescription,
	})
	httputils.RequireStatus(t, resp, 200, "update forked challenge")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "title") != updatedForkTitle {
		process.Fail("update forked challenge", fmt.Errorf("expected updated forked challenge title"))
	}
	process.EndStep()

	// [STEP 8] Verify original challenge is unchanged
	process.StartStep("Verifica que no haya cambios en el reto original")
	resp = httputils.GetChallengeByID(t, app, creatorAccess, originalChallengeID)
	httputils.RequireStatus(t, resp, 200, "get original challenge after fork update")
	body := httputils.MustJSONMap(t, resp)
	if httputils.StringField(body, "title") == updatedForkTitle || httputils.StringField(body, "description") == updatedForkDescription {
		process.Fail("verify original unchanged", fmt.Errorf("original challenge should not be modified by fork update"))
	}
	process.EndStep()

	// [STEP 9] Verify forked challenge has the updated values
	process.StartStep("Verifica los cambios en el reto copiado")
	resp = httputils.GetChallengeByID(t, app, observerAccess, forkedChallengeID)
	httputils.RequireStatus(t, resp, 200, "get forked challenge after update")
	body = httputils.MustJSONMap(t, resp)
	if httputils.StringField(body, "title") != updatedForkTitle || httputils.StringField(body, "description") != updatedForkDescription {
		process.Fail("verify fork updated", fmt.Errorf("expected forked challenge to keep updated values"))
	}
	process.EndStep()

	process.End()
}
