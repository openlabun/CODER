package exam_test

import (
	"net/http"
	"testing"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestChallengeCRUDHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}

	t.Log("[STEP 2] Autenticando profesor")
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := httputils.AuthHeaders(teacherAccess)

	t.Log("[STEP 3] Crear curso para challenges")
	courseID := createCourseHTTP(t, app, teacherAccess, "challenge-crud")

	t.Log("[STEP 4] Crear challenge")
	challengeID := createChallengeHTTP(t, app, teacherAccess, "HTTP Challenge")
	t.Logf("[OK] Challenge creado. challengeID=%s", challengeID)

	t.Log("[STEP 5] Actualizar challenge")
	updatedTitle := "HTTP Challenge Updated"
	updateBody := map[string]any{
		"title":           updatedTitle,
		"difficulty":      "medium",
		"workerTimeLimit": 2000,
	}
	status, body, err := httputils.PatchChallengesById(teacherHeaders, map[string]any{"id": challengeID}, updateBody)
	if err != nil {
		t.Fatalf("update challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected update status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	updated := httputils.DecodeMap(t, body, "update challenge")
	if httputils.MapString(t, updated, "title", "update challenge") != updatedTitle {
		t.Fatalf("expected title=%q after update, got body=%s", updatedTitle, string(body))
	}
	t.Log("[OK] Challenge actualizado")

	t.Log("[STEP 6] Publicar challenge")
	status, body, err = httputils.PostChallengesPublish(teacherHeaders, map[string]any{"id": challengeID})
	if err != nil {
		t.Fatalf("publish challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected publish status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	published := httputils.DecodeMap(t, body, "publish challenge")
	if httputils.MapString(t, published, "status", "publish challenge") != "published" {
		t.Fatalf("expected status=published, got body=%s", string(body))
	}
	t.Log("[OK] Challenge publicado")

	t.Log("[STEP 7] Listar challenges por examen")
	status, body, err = httputils.GetChallenges(teacherHeaders, nil)
	if err != nil {
		t.Fatalf("list challenges request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	challenges := httputils.DecodeSliceMap(t, body, "list challenges")
	if !httputils.ContainsID(challenges, challengeID) {
		t.Fatalf("expected challengeID=%s in list, got body=%s", challengeID, string(body))
	}
	t.Logf("[OK] Challenge encontrado en listado. total=%d", len(challenges))

	t.Log("[STEP 8] Archivar challenge")
	status, body, err = httputils.PostChallengesArchive(teacherHeaders, map[string]any{"id": challengeID})
	if err != nil {
		t.Fatalf("archive challenge request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected archive status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	archived := httputils.DecodeMap(t, body, "archive challenge")
	if httputils.MapString(t, archived, "status", "archive challenge") != "archived" {
		t.Fatalf("expected status=archived, got body=%s", string(body))
	}
	t.Log("[OK] Challenge archivado")

	t.Log("[STEP 9] Cleanup via DELETE /courses/:id")
	status, body, err = httputils.DeleteCoursesById(teacherHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected course delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	if !httputils.MapBool(t, httputils.DecodeMap(t, body, "delete course"), "removed", "delete course") {
		t.Fatalf("expected removed=true for course delete, got body=%s", string(body))
	}
	t.Log("[OK] Cleanup completado")
}

func TestChallengeFromStudentViewHTTP(t *testing.T) {
	t.Log("[STEP 1] Inicializando app HTTP")
	app, err := httputils.InitApp()
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}

	t.Log("[STEP 2] Autenticando profesor")
	teacherAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "test@test.com", "Password123!", "Teacher Test")
	teacherHeaders := httputils.AuthHeaders(teacherAccess)

	t.Log("[STEP 3] Crear curso + examen")
	courseID := createCourseHTTP(t, app, teacherAccess, "challenge-student")
	examID := createExamHTTP(t, app, teacherAccess, courseID, "HTTP Challenge Student View Exam")

	t.Log("[STEP 4] Crear 3 challenges en draft")
	publishedID := createChallengeHTTP(t, app, teacherAccess, "Challenge Published")
	archivedID := createChallengeHTTP(t, app, teacherAccess, "Challenge Archived")
	draftID := createChallengeHTTP(t, app, teacherAccess, "Challenge Draft")

	t.Log("[STEP 5] Publicar uno y archivar otro")
	_, _, _ = httputils.PostChallengesPublish(teacherHeaders, map[string]any{"id": publishedID})
	_, _, _ = httputils.PostChallengesPublish(teacherHeaders, map[string]any{"id": archivedID})
	_, _, _ = httputils.PostChallengesArchive(teacherHeaders, map[string]any{"id": archivedID})

	t.Log("[STEP 5.1] Crear ExamItems para el challenge publicado")
	createExamItemHTTP(t, app, teacherAccess, examID, publishedID, 1)

	t.Log("[STEP 5.2] Validar que genera error al crear ExamItem para challenge en borrador")
	status, body, err := httputils.PostExamItems(httputils.AuthHeaders(teacherAccess), map[string]any{
		"exam_id":      examID,
		"challenge_id": draftID,
		"order":        2,
		"points":       100,
	})
	if err != nil {
		t.Fatalf("create exam item for draft challenge request failed: %v", err)
	}

	if status != http.StatusBadRequest {
		t.Fatalf("expected status=%d when creating exam item for draft challenge, got %d body=%s", http.StatusBadRequest, status, string(body))
	}

	t.Log("[STEP 5.3] Validar que genera error al crear ExamItem para challenge archivado")
	status, body, err = httputils.PostExamItems(httputils.AuthHeaders(teacherAccess), map[string]any{
		"exam_id":      examID,
		"challenge_id": archivedID,
		"order":        3,
		"points":       100,
	})
	if err != nil {
		t.Fatalf("create exam item for archived challenge request failed: %v", err)
	}

	if status != http.StatusBadRequest {
		t.Fatalf("expected status=%d when creating exam item for archived challenge, got %d body=%s", http.StatusBadRequest, status, string(body))
	}

	t.Log("[STEP 6] Autenticando estudiante y matriculando")
	studentAccess := httputils.EnsureHTTPAuthUserAccess(t, app, "stud@test.com", "Password123!", "Student Test")
	studentHeaders := httputils.AuthHeaders(studentAccess)
	enrollBody := map[string]string{"course_id": courseID, "student_id": studentAccess.UserData.ID}
	status, body, err = httputils.PostCoursesEnroll(studentHeaders, enrollBody)
	if err != nil {
		t.Fatalf("enroll student request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected enroll status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	t.Log("[STEP 7] Listar ExamItems y obtener los challenges como estudiante")
	status, body, err = httputils.GetExamItems(studentHeaders, map[string]any{"exam_id": examID})
	if err != nil {
		t.Fatalf("list exam items as student failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected list status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	examItems := httputils.DecodeSliceMap(t, body, "student exam items list")
	foundPublished := false
	foundArchived := false
	foundDraft := false
	for _, item := range examItems {
		challenge, ok := item["challenge"].(map[string]any)
		if !ok || challenge == nil {
			continue
		}
		id, _ := challenge["id"].(string)
		if id == publishedID {
			foundPublished = true
		}
		if id == archivedID {
			foundArchived = true
		}
		if id == draftID {
			foundDraft = true
		}
	}

	if !foundPublished {
		t.Fatal("expected published challenge in student view")
	}
	if foundArchived {
		t.Fatal("did not expect archived challenge in student view")
	}
	if foundDraft {
		t.Fatal("did not expect draft challenge in student view")
	}
	t.Logf("[OK] Vista estudiante validada. visibleChallenges=%d", len(examItems))

	t.Log("[STEP 8] Cleanup via DELETE /courses/:id")
	status, body, err = httputils.DeleteCoursesById(teacherHeaders, map[string]any{"id": courseID})
	if err != nil {
		t.Fatalf("delete course request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected course delete status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}
	if !httputils.MapBool(t, httputils.DecodeMap(t, body, "delete course"), "removed", "delete course") {
		t.Fatalf("expected removed=true for course delete, got body=%s", string(body))
	}
	t.Log("[OK] Cleanup completado")
}
