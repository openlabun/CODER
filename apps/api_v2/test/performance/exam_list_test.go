package performance_tests

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamListPerformance(t *testing.T) {
	process, app := httputils.StartConcurrentHTTPTestWithApp(t, "Exam List Performance")
	const n = 10
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	examIDs := make([]string, n)

	defer func() {
		for _, id := range examIDs {
			if id != "" && teacherAccess != nil {
				_ = httputils.DeleteExamByID(t, app, teacherAccess, id)
			}
		}
	}()

	// [STEP 1] Login as teacher user
	t.Log("[STEP 1] Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")

	// [STEP 2] Create multiple public exams concurrently (critical)
	now := time.Now().UTC()
	var createCounter int64 = -1
	process.StartConcurrentStep(2, "Crear examenes publicos concurrentemente", func() (int, []byte, error) {
		i := atomic.AddInt64(&createCounter, 1)
		resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
			"title":                  fmt.Sprintf("Perf Exam %d", i),
			"description":            "Examen para prueba de rendimiento",
			"visibility":             "public",
			"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
			"allow_late_submissions": true,
			"time_limit":             3600,
			"try_limit":              2,
			"professor_id":           teacherAccess.UserID,
		})
		if resp.StatusCode == 201 && resp.JSON != nil {
			examIDs[i], _ = resp.JSON["id"].(string)
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	// [STEP 3] Login as student user
	t.Log("[STEP 3] Iniciar sesion con usuario de estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")

	// [STEP 4] List public exams concurrently (critical)
	process.StartConcurrentStep(4, "Listar examenes publicos concurrentemente", func() (int, []byte, error) {
		resp := httputils.GetPublicExams(t, app, studentAccess)
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	process.End()
}
