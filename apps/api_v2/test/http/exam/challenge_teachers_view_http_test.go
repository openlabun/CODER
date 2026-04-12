package exam_test

import (
	"encoding/json"
	"fmt"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestChallengeFromTeacherViewHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Challenge From Teachers View HTTP")
	creatorEmail := "test@test.com"
	observerEmail := "test2@test.com"
	password := "Password123!"

	var creatorAccess *httputils.HTTPAccess
	var observerAccess *httputils.HTTPAccess
	var challengeID string

	defer func() {
		if challengeID != "" && creatorAccess != nil {
			t.Logf("[CLEANUP] Eliminando challenge %s", challengeID)
			_ = httputils.DeleteChallengeByID(t, app, creatorAccess, challengeID)
		}
	}()

	// [STEP 1] Login as creator teacher
	process.StartStep("Iniciar sesion con usuario de docente (creador)")
	creatorAccess = httputils.EnsureAuthUserAccess(t, app, creatorEmail, password, "Teacher Creator")
	process.EndStep()

	// [STEP 2] Create a challenge with private visibility
	process.StartStep("Crea un reto (private)")
	resp := httputils.PostChallengeCreate(t, app, creatorAccess, map[string]any{
		"title":               "Challenge Teachers View",
		"description":         "Challenge para validar visibilidad entre docentes",
		"tags":                []string{"teachers-view", "challenge"},
		"status":              "private",
		"difficulty":          "easy",
		"worker_time_limit":   1400,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables":     []map[string]any{{"name": "n", "type": "int", "value": "7"}},
		"output_variable":     map[string]any{"name": "out", "type": "int", "value": "7"},
		"constraints":         "1 <= n <= 100",
	})
	httputils.RequireStatus(t, resp, 201, "create private challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Login as observer teacher
	process.StartStep("Iniciar sesion con usuario de docente (observador)")
	observerAccess = httputils.EnsureAuthUserAccess(t, app, observerEmail, password, "Teacher Observer")
	process.EndStep()

	// [STEP 4] Get the challenge details with observer teacher (expect error)
	process.StartStep("Obtiene datos del reto con docente observador (espera error)")
	resp = httputils.GetChallengeByID(t, app, observerAccess, challengeID)
	if resp.StatusCode == 200 {
		process.Fail("observer private challenge access", fmt.Errorf("expected error for observer on private challenge"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 5] Update the challenge to published
	process.StartStep("Actualiza el reto a visibilidad published")
	resp = httputils.PostChallengePublish(t, app, creatorAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "publish challenge")
	process.EndStep()

	// [STEP 6] Get the challenge details with observer teacher (expect success)
	process.StartStep("Obtiene datos del reto con docente observador")
	resp = httputils.GetPublicChallenges(t, app, observerAccess)
	httputils.RequireStatus(t, resp, 200, "observer get public challenges")
	var publishedList []map[string]any
	if err := json.Unmarshal(resp.Body, &publishedList); err != nil {
		process.Fail("observer get public challenges", fmt.Errorf("decode public challenges: %w", err))
	}
	foundPublished := false
	for _, c := range publishedList {
		if httputils.StringField(c, "id") == challengeID {
			foundPublished = true
			break
		}
	}
	if !foundPublished {
		process.Fail("observer get public challenges", fmt.Errorf("expected published challenge %s to be visible", challengeID))
	}
	process.EndStep()

	// [STEP 7] Update the challenge to archived
	process.StartStep("Actualiza el reto a visibilidad archived")
	resp = httputils.PostChallengeArchive(t, app, creatorAccess, challengeID)
	httputils.RequireStatus(t, resp, 200, "archive challenge")
	process.EndStep()

	// [STEP 8] Get the challenge details with observer teacher (expect error)
	process.StartStep("Obtiene datos del reto con docente observador (espera error)")
	resp = httputils.GetPublicChallenges(t, app, observerAccess)
	httputils.RequireStatus(t, resp, 200, "observer get public challenges after archive")
	var archivedList []map[string]any
	if err := json.Unmarshal(resp.Body, &archivedList); err != nil {
		process.Fail("observer get public challenges after archive", fmt.Errorf("decode public challenges archived: %w", err))
	}
	for _, c := range archivedList {
		if httputils.StringField(c, "id") == challengeID {
			process.Fail("observer archived challenge visibility", fmt.Errorf("archived challenge %s should not be visible", challengeID))
		}
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
