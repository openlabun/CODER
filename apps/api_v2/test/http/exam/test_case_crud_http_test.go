package exam_test

import (
	"encoding/json"
	"fmt"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestTestCaseCRUDHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "TestCase CRUD HTTP")
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var challengeID string
	var testCaseID string
	var deletedTestCaseID string

	defer func() {
		if testCaseID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando caso de prueba %s", testCaseID)
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, testCaseID)
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

	// [STEP 2] Create a challenge to associate with test cases
	process.StartStep("Crear un reto")
	resp := httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Challenge for TestCase CRUD",
		"description":         "Challenge auxiliar para test case CRUD",
		"tags":                []string{"testcase", "crud"},
		"status":              "draft",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve() { return; }",
		},
		"input_variables":     []map[string]any{{"name": "x", "type": "int", "value": "1"}},
		"output_variable":     map[string]any{"name": "out", "type": "int", "value": "1"},
		"constraints":         "1 <= x <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create challenge")
	challengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Try to create a test case without input values (expect error)
	process.StartStep("Crear un caso de uso sin valores de entrada (espera error)")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "invalid_test_case",
		"input":           []map[string]any{},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "1"},
		"is_sample":       true,
		"points":          0,
		"challenge_id":    challengeID,
	})
	if resp.StatusCode == 201 {
		process.Fail("create invalid test case", fmt.Errorf("expected error when creating test case without inputs"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 4] Create a valid test case
	process.StartStep("Crear un caso de uso valido")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name": "sample_valid",
		"input": []map[string]any{
			{"name": "a", "type": "int", "value": "2"},
			{"name": "b", "type": "int", "value": "3"},
		},
		"expected_output": map[string]any{"name": "sum", "type": "int", "value": "5"},
		"is_sample":       true,
		"points":          10,
		"challenge_id":    challengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create valid test case")
	testCaseID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Update the test case
	process.StartStep("Actualizar el caso de uso")
	resp = httputils.PatchTestCaseUpdate(t, app, teacherAccess, testCaseID, map[string]any{
		"name":   "sample_valid_updated",
		"points": 25,
	})
	httputils.RequireStatus(t, resp, 200, "update test case")
	body := httputils.MustJSONMap(t, resp)
	if httputils.StringField(body, "name") != "sample_valid_updated" || int(body["points"].(float64)) != 25 {
		process.Fail("update test case", fmt.Errorf("expected updated test case values"))
	}
	process.EndStep()

	// [STEP 6] Get the test case details and validate updated data
	process.StartStep("Eliminar el caso de uso")
	resp = httputils.DeleteTestCaseByID(t, app, teacherAccess, testCaseID)
	httputils.RequireStatus(t, resp, 200, "delete test case")
	deletedTestCaseID = testCaseID
	testCaseID = ""
	process.EndStep()

	// [STEP 7] Try to get the deleted test case details (expect error)
	process.StartStep("Verificar eliminacion")
	resp = httputils.GetTestCasesByChallenge(t, app, teacherAccess, challengeID, "")
	httputils.RequireStatus(t, resp, 200, "get test cases by challenge")
	var remaining []map[string]any
	if err := json.Unmarshal(resp.Body, &remaining); err != nil {
		process.Fail("get test cases by challenge", fmt.Errorf("decode remaining test cases: %w", err))
	}
	for _, tc := range remaining {
		if httputils.StringField(tc, "id") == deletedTestCaseID {
			process.Fail("verify test case deletion", fmt.Errorf("test case %s should not exist after deletion", deletedTestCaseID))
		}
	}
	process.EndStep()

	process.End()
}
