package exam_test

import (
	"fmt"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestChallengeCRUDHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Challenge CRUD HTTP")
	email := "test@test.com"
	password := "Password123!"

	var teacherID string
	var teacherAccess *httputils.HTTPAccess
	var challengeID string
	var deletedChallengeID string

	defer func() {
		if challengeID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando reto %s", challengeID)
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
		}
	}()

	// [STEP 1] Login Teacher user
	process.StartStep("Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, email, password, "Teacher Test")
	teacherID = teacherAccess.UserID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))
	process.EndStep()

	// [STEP 2] Create Challenge
	process.StartStep("Crea un reto")
	resp := httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Challenge CRUD Test",
		"description":         "Challenge creado por test CRUD",
		"tags":                []string{"crud", "challenge"},
		"status":              "draft",
		"difficulty":          "easy",
		"worker_time_limit":   1500,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables": []map[string]any{
			{"name": "a", "type": "int", "value": "2"},
			{"name": "b", "type": "int", "value": "3"},
		},
		"output_variable": map[string]any{"name": "sum", "type": "int", "value": "5"},
		"constraints":     "1 <= a,b <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	if challengeID == "" {
		process.Fail("create challenge", fmt.Errorf("expected created challenge with ID"))
	}
	process.Log(fmt.Sprintf("challengeID=%s", challengeID))
	process.EndStep()

	// [STEP 3] Update Challenge
	process.StartStep("Actualiza el reto")
	updatedTitle := "Challenge CRUD Test Updated"
	updatedDescription := "Challenge actualizado por test CRUD"
	resp = httputils.PatchChallengeUpdate(t, app, teacherAccess, challengeID, map[string]any{
		"title":       updatedTitle,
		"description": updatedDescription,
		"difficulty":  "medium",
	})
	httputils.RequireStatus(t, resp, 200, "update challenge")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "title") != updatedTitle {
		process.Fail("update challenge", fmt.Errorf("expected updated challenge title"))
	}
	process.EndStep()

	// [STEP 4] Get Challenge details and validate
	process.StartStep("Obtiene los datos del reto y valida los cambios")
	resp = httputils.GetChallengeByID(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "get challenge details")
	body := httputils.MustJSONMap(t, resp)
	if httputils.StringField(body, "id") != challengeID {
		process.Fail("get challenge details", fmt.Errorf("expected challenge details for %s", challengeID))
	}
	if httputils.StringField(body, "title") != updatedTitle || httputils.StringField(body, "description") != updatedDescription {
		process.Fail("get challenge details", fmt.Errorf("challenge update not persisted"))
	}
	process.EndStep()

	// [STEP 5] Delete Challenge
	process.StartStep("Elimina el reto")
	resp = httputils.DeleteChallengeByID(t, app, teacherAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "delete challenge")
	deletedChallengeID = challengeID
	challengeID = ""
	process.EndStep()

	// [STEP 6] Verify deletion
	process.StartStep("Verifica eliminacion")
	resp = httputils.GetChallengeByID(t, app, teacherAccess, deletedChallengeID)
	if resp.StatusCode == 200 {
		process.Fail("verify challenge deletion", fmt.Errorf("expected error after deleting challenge"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
