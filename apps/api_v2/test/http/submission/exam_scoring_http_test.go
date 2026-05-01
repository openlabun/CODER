package submission_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamScoringHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Exam Scoring HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var examID string
	var firstChallengeID string
	var secondChallengeID string
	var firstTestCaseOneID string
	var firstTestCaseTwoID string
	var firstTestCaseThreeID string
	var secondTestCaseOneID string
	var secondTestCaseTwoID string
	var firstExamItemID string
	var secondExamItemID string
	var sessionID string
	var submissionID string

	defer func() {
		if sessionID != "" && teacherAccess != nil {
			_ = httputils.PostSessionClose(t, app, teacherAccess, sessionID)
		}
		if firstExamItemID != "" && teacherAccess != nil {
			_ = httputils.DeleteExamItemByID(t, app, teacherAccess, firstExamItemID)
		}
		if secondExamItemID != "" && teacherAccess != nil {
			_ = httputils.DeleteExamItemByID(t, app, teacherAccess, secondExamItemID)
		}
		if firstTestCaseThreeID != "" && teacherAccess != nil {
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, firstTestCaseThreeID)
		}
		if firstTestCaseTwoID != "" && teacherAccess != nil {
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, firstTestCaseTwoID)
		}
		if firstTestCaseOneID != "" && teacherAccess != nil {
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, firstTestCaseOneID)
		}
		if secondTestCaseTwoID != "" && teacherAccess != nil {
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, secondTestCaseTwoID)
		}
		if secondTestCaseOneID != "" && teacherAccess != nil {
			_ = httputils.DeleteTestCaseByID(t, app, teacherAccess, secondTestCaseOneID)
		}
		if firstChallengeID != "" && teacherAccess != nil {
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, firstChallengeID)
		}
		if secondChallengeID != "" && teacherAccess != nil {
			_ = httputils.DeleteChallengeByID(t, app, teacherAccess, secondChallengeID)
		}
		if examID != "" && teacherAccess != nil {
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examID)
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	process.EndStep()

	// [STEP 2] Create an exam
	process.StartStep("Crear examen publico (visibilidad public y sin curso)")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Exam Scoring Exam",
		"description":            "Examen para validar puntaje del examen",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             120,
		"try_limit":              3,
		"professor_id":           teacherAccess.UserID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create two challenges
	process.StartStep("Crear dos retos")
	resp = httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Exam Scoring Challenge 1",
		"description":         "Primer reto para validar score del examen",
		"tags":                []string{"exam", "scoring", "one"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve():\n    pass",
		},
		"input_variables": []map[string]any{{"name": "n", "type": "int", "value": "2"}},
		"output_variable": map[string]any{"name": "out", "type": "int", "value": "4"},
		"constraints":     "1 <= n <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create first challenge")
	firstChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	resp = httputils.PostChallengeCreate(t, app, teacherAccess, map[string]any{
		"title":               "Exam Scoring Challenge 2",
		"description":         "Segundo reto para validar score del examen",
		"tags":                []string{"exam", "scoring", "two"},
		"status":              "published",
		"difficulty":          "easy",
		"worker_time_limit":   1200,
		"worker_memory_limit": 256,
		"code_templates": map[string]any{
			"python": "def solve():\n    pass",
		},
		"input_variables": []map[string]any{{"name": "n", "type": "int", "value": "5"}},
		"output_variable": map[string]any{"name": "out", "type": "int", "value": "10"},
		"constraints":     "1 <= n <= 1000",
	})
	httputils.RequireStatus(t, resp, 201, "create second challenge")
	secondChallengeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 4] Create test cases for first challenge
	process.StartStep("Crear 2 casos de prueba con valor de 3 puntos para el primer reto")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "exam_score_case_1_1",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "2"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "4"},
		"is_sample":       false,
		"points":          3,
		"challenge_id":    firstChallengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create first test case 1")
	firstTestCaseOneID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "exam_score_case_1_2",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "5"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "10"},
		"is_sample":       false,
		"points":          3,
		"challenge_id":    firstChallengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create first test case 2")
	firstTestCaseTwoID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 5] Create impossible test case for first challenge
	process.StartStep("Crear un caso de prueba con valor de 6 puntos para el primer reto (debe ser imposible de cumplir)")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "exam_score_case_1_impossible",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "7"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "999"},
		"is_sample":       false,
		"points":          6,
		"challenge_id":    firstChallengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create first impossible test case")
	firstTestCaseThreeID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 6] Create exam items
	process.StartStep("Crear un punto de examen con valor de 30")
	resp = httputils.PostExamItemCreate(t, app, teacherAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": firstChallengeID,
		"order":        1,
		"points":       30,
	})
	httputils.RequireStatus(t, resp, 201, "create first exam item")
	firstExamItemID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	process.StartStep("Crear 2 casos de prueba con valor de 5 puntos para el segundo reto")
	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "exam_score_case_2_1",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "2"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "4"},
		"is_sample":       false,
		"points":          5,
		"challenge_id":    secondChallengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create second test case 1")
	secondTestCaseOneID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	resp = httputils.PostTestCaseCreate(t, app, teacherAccess, map[string]any{
		"name":            "exam_score_case_2_2",
		"input":           []map[string]any{{"name": "n", "type": "int", "value": "5"}},
		"expected_output": map[string]any{"name": "out", "type": "int", "value": "10"},
		"is_sample":       false,
		"points":          5,
		"challenge_id":    secondChallengeID,
	})
	httputils.RequireStatus(t, resp, 201, "create second test case 2")
	secondTestCaseTwoID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	process.StartStep("Crear un punto de examen con valor de 70")
	resp = httputils.PostExamItemCreate(t, app, teacherAccess, map[string]any{
		"exam_id":      examID,
		"challenge_id": secondChallengeID,
		"order":        2,
		"points":       70,
	})
	httputils.RequireStatus(t, resp, 201, "create second exam item")
	secondExamItemID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 7] Login as student
	process.StartStep("Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	process.EndStep()

	// [STEP 8] Create a session
	process.StartStep("Crear una sesion en el examen")
	resp = httputils.PostSessionCreate(t, app, studentAccess, map[string]any{"user_id": studentAccess.UserID, "exam_id": examID})
	httputils.RequireStatus(t, resp, 201, "create session")
	sessionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 9] Create a submission
	process.StartStep("Crear una revision")
	resp = httputils.PostSubmissionCreate(t, app, studentAccess, map[string]any{
		"code":         "import sys\nprint(int(sys.stdin.read().strip()) * 2)",
		"language":     "python",
		"challenge_id": firstChallengeID,
		"session_id":   sessionID,
	})
	httputils.RequireStatus(t, resp, 201, "create submission")
	submissionID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 10] Poll submission until terminal
	process.StartStep("Obtener el status de la revision hasta que su estado sea accepted o wrong_answer")
	deadline := time.Now().Add(2 * time.Minute)
	var lastStatus map[string]any
	timedOut := true
	for time.Now().Before(deadline) {
		resp = httputils.GetSubmissionByID(t, app, studentAccess, submissionID)
		if resp.StatusCode != 200 {
			process.Log(fmt.Sprintf("Polling status code: %d", resp.StatusCode))
			time.Sleep(2 * time.Second)
			continue
		}

		status := httputils.MustJSONMap(t, resp)
		resultsRaw, ok := status["results"].([]any)
		if !ok || len(resultsRaw) != 3 {
			time.Sleep(2 * time.Second)
			continue
		}

		allTerminal := true
		for _, item := range resultsRaw {
			result, ok := item.(map[string]any)
			if !ok {
				allTerminal = false
				break
			}
			s := httputils.StringField(result, "status")
			if s != "accepted" && s != "wrong_answer" {
				allTerminal = false
				break
			}
		}

		lastStatus = status
		if allTerminal {
			timedOut = false
			break
		}

		time.Sleep(2 * time.Second)
	}
	if timedOut {
		process.Fail("wait scoring terminal status", fmt.Errorf("submission %s did not reach terminal accepted/wrong_answer status before timeout", submissionID))
	}
	process.EndStep()

	// [STEP 11] Verify submission score equals expected (50)
	process.StartStep("Confirmar valor del atributo Score de la revision corresponde a 50")
	if lastStatus == nil {
		process.Fail("verify submission score", fmt.Errorf("expected submission status output"))
	}
	submission, ok := lastStatus["submission"].(map[string]any)
	if !ok {
		process.Fail("verify submission score", fmt.Errorf("expected submission field in status output"))
	}
	score, ok := submission["score"].(float64)
	if !ok || int(score) != 50 {
		process.Fail("verify submission score", fmt.Errorf("expected submission score 50, got %v", submission["score"]))
	}
	process.EndStep()

	// [STEP 12] Close session
	process.StartStep("Cerrar la sesion del examen")
	_ = httputils.PostSessionClose(t, app, studentAccess, sessionID)
	process.EndStep()

	// [STEP 13] Get exam scoring as student
	process.StartStep("Obtener resultados de su examen como estudiante")
	req := httputils.GetExamScoring(t, app, studentAccess, examID, studentAccess.UserID)
	httputils.RequireStatus(t, req, 200, "get exam scoring")
	examScoring := httputils.MustJSONMap(t, req)
	if httputils.StringField(examScoring, "exam_id") != examID {
		process.Fail("verify exam scoring", fmt.Errorf("expected exam id %s, got %s", examID, httputils.StringField(examScoring, "exam_id")))
	}
	bestScoreRaw, ok := examScoring["best_score"].(float64)
	if !ok || int(bestScoreRaw) != 15 {
		process.Fail("verify exam scoring", fmt.Errorf("expected best score 15, got %v", examScoring["best_score"]))
	}
	examScoresRaw, ok := examScoring["exam_scores"].([]any)
	if !ok || len(examScoresRaw) != 1 {
		process.Fail("verify exam scoring", fmt.Errorf("expected 1 exam score, got %d", len(examScoresRaw)))
	}
	examScoreObj, ok := examScoresRaw[0].(map[string]any)
	if !ok {
		process.Fail("verify exam scoring", fmt.Errorf("unexpected exam score object"))
	}
	examScoreVal, ok := examScoreObj["score"].(float64)
	if !ok || int(examScoreVal) != 15 {
		process.Fail("verify exam scoring", fmt.Errorf("expected stored exam score 15, got %v", examScoreObj["score"]))
	}
	process.EndStep()

	// [STEP 14] Get item scoring details
	process.StartStep("Obtener resultados de los punto del examen como estudiante")
	examScoreID := httputils.StringField(examScoreObj, "id")
	req = httputils.GetExamScoreDetails(t, app, studentAccess, examScoreID)
	httputils.RequireStatus(t, req, 200, "get item scoring")
	examScoreDetails := httputils.MustJSONMap(t, req)
	if httputils.StringField(examScoreDetails, "exam_score_id") != examScoreID {
		process.Fail("get item scoring", fmt.Errorf("expected exam score id %s, got %s", examScoreID, httputils.StringField(examScoreDetails, "exam_score_id")))
	}
	examItemsRaw, ok := examScoreDetails["exam_items"].([]any)
	if !ok || len(examItemsRaw) != 2 {
		process.Fail("get item scoring", fmt.Errorf("expected 2 exam items, got %d", len(examItemsRaw)))
	}
	firstFound := false
	secondFound := false
	for _, it := range examItemsRaw {
		item, ok := it.(map[string]any)
		if !ok {
			continue
		}
		examItemObj, ok := item["exam_item"].(map[string]any)
		if !ok {
			continue
		}
		examItemID := httputils.StringField(examItemObj, "id")
		if examItemID == firstExamItemID {
			firstFound = true
			itemScoreObj, _ := item["exam_item_score"].(map[string]any)
			scoreF, _ := itemScoreObj["score"].(float64)
			if int(scoreF) != 50 {
				process.Fail("get item scoring", fmt.Errorf("expected first exam item score 50, got %v", itemScoreObj["score"]))
			}
		}
		if examItemID == secondExamItemID {
			secondFound = true
			itemScoreObj, _ := item["exam_item_score"].(map[string]any)
			scoreF, _ := itemScoreObj["score"].(float64)
			if int(scoreF) != 0 {
				process.Fail("get item scoring", fmt.Errorf("expected second exam item score 0, got %v", itemScoreObj["score"]))
			}
		}
	}
	if !firstFound {
		process.Fail("get item scoring", fmt.Errorf("expected first exam item %s", firstExamItemID))
	}
	if !secondFound {
		process.Fail("get item scoring", fmt.Errorf("expected second exam item %s", secondExamItemID))
	}
	process.EndStep()

	// [STEP 15] Get users scores by exam as teacher
	process.StartStep("Obtener los resultados del examen para cada estudiante como profesor")
	req = httputils.GetUsersScoresByExam(t, app, teacherAccess, examID)
	httputils.RequireStatus(t, req, 200, "get users scores by exam")
	usersScores := httputils.MustJSONMap(t, req)
	if httputils.StringField(usersScores, "exam_id") != examID {
		process.Fail("get users scores by exam", fmt.Errorf("expected exam id %s, got %s", examID, httputils.StringField(usersScores, "exam_id")))
	}
	usersRaw, ok := usersScores["users"].([]any)
	if !ok || len(usersRaw) != 1 {
		process.Fail("get users scores by exam", fmt.Errorf("expected 1 user score entry, got %d", len(usersRaw)))
	}
	userObj, _ := usersRaw[0].(map[string]any)
	if httputils.StringField(userObj, "user_id") != studentAccess.UserID {
		process.Fail("get users scores by exam", fmt.Errorf("expected user id %s, got %s", studentAccess.UserID, httputils.StringField(userObj, "user_id")))
	}
	examScoresObj, _ := userObj["exam_scores"].(map[string]any)
	best, _ := examScoresObj["best_score"].(float64)
	if int(best) != 15 {
		process.Fail("get users scores by exam", fmt.Errorf("expected best score 15, got %v", best))
	}
	process.EndStep()

	// [STEP 16] Get exams scores by user as teacher
	process.StartStep("Obtener los resultados de todos los exámenes del estudiante como profesor")
	req = httputils.GetExamsScoresByUser(t, app, teacherAccess, studentAccess.UserID)
	httputils.RequireStatus(t, req, 200, "get exams scores by user")
	studentExams := httputils.MustJSONMap(t, req)
	if httputils.StringField(studentExams, "user_id") != studentAccess.UserID {
		process.Fail("get exams scores by user", fmt.Errorf("expected user id %s, got %s", studentAccess.UserID, httputils.StringField(studentExams, "user_id")))
	}
	examsRaw, _ := studentExams["exam_scores"].([]any)
	if len(examsRaw) != 1 {
		process.Fail("get exams scores by user", fmt.Errorf("expected 1 exam score entry, got %d", len(examsRaw)))
	}
	examObj, _ := examsRaw[0].(map[string]any)
	if httputils.StringField(examObj, "exam_id") != examID {
		process.Fail("get exams scores by user", fmt.Errorf("expected exam id %s, got %s", examID, httputils.StringField(examObj, "exam_id")))
	}
	promF, _ := studentExams["prom_score"].(float64)
	if int(promF) != 15 {
		process.Fail("get exams scores by user", fmt.Errorf("expected prom score 15, got %v", promF))
	}
	process.EndStep()

	process.End()
}
