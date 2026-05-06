package performance_tests

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestSessionsPerformance(t *testing.T) {
	process, app := httputils.StartConcurrentHTTPTestWithApp(t, "Sessions Performance")
	const n = 10
	teacherEmail := "test@test.com"
	password := "Password123!"

	var teacherAccess *httputils.HTTPAccess
	var examID string
	accesses := make([]*httputils.HTTPAccess, n)
	sessionIDs := make([]string, n)

	defer func() {
		for _, sid := range sessionIDs {
			if sid != "" && teacherAccess != nil {
				_ = httputils.PostSessionClose(t, app, teacherAccess, sid)
			}
		}
		if examID != "" && teacherAccess != nil {
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examID)
		}
	}()

	// [STEP 1] Login as teacher user
	t.Log("[STEP 1] Iniciar sesion con usuario de docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")

	// [STEP 2] Create public exam
	t.Log("[STEP 2] Crear examen publico")
	now := time.Now().UTC()
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Sessions Performance Exam",
		"description":            "Examen para prueba de rendimiento de sesiones",
		"visibility":             "public",
		"start_time":             now.Add(2 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3600,
		"try_limit":              10,
		"professor_id":           teacherAccess.UserID,
	})
	if resp.StatusCode != 201 {
		t.Fatalf("[STEP 2] create exam: status=%d body=%s", resp.StatusCode, string(resp.Body))
	}
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")

	// [STEP 3] Login with student users (critical)
	var loginCounter int64 = -1
	process.StartConcurrentStep(3, "Iniciar sesion con usuarios estudiantes", func() (int, []byte, error) {
		i := atomic.AddInt64(&loginCounter, 1)
		email := fmt.Sprintf("stud-%d@test.com", i)
		httputils.PostAuthRegister(t, app, email, fmt.Sprintf("Student %d", i), password)
		resp := httputils.PostAuthLogin(t, app, email, password)
		if resp.StatusCode == 200 && resp.JSON != nil {
			access := &httputils.HTTPAccess{Email: email}
			if tok, ok := resp.JSON["token"].(map[string]any); ok {
				access.AccessToken, _ = tok["access_token"].(string)
			}
			if user, ok := resp.JSON["user_data"].(map[string]any); ok {
				access.UserID, _ = user["id"].(string)
			}
			accesses[i] = access
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	// [STEP 4] Each student creates a session (critical)
	var sessionCounter int64 = -1
	process.StartConcurrentStep(4, "Cada estudiante crea una sesion", func() (int, []byte, error) {
		i := atomic.AddInt64(&sessionCounter, 1)
		if accesses[i] == nil {
			return 0, nil, fmt.Errorf("no access for student %d", i)
		}
		resp := httputils.PostSessionCreate(t, app, accesses[i], map[string]any{
			"user_id": accesses[i].UserID,
			"exam_id": examID,
		})
		if resp.StatusCode == 201 && resp.JSON != nil {
			sessionIDs[i], _ = resp.JSON["id"].(string)
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	// [STEP 5] Each student heartbeats their session (critical)
	var heartbeat1Counter int64 = -1
	process.StartConcurrentStep(5, "Cada estudiante hace heartbeat (antes del freeze)", func() (int, []byte, error) {
		i := atomic.AddInt64(&heartbeat1Counter, 1)
		if accesses[i] == nil || sessionIDs[i] == "" {
			return 0, nil, fmt.Errorf("no session for student %d", i)
		}
		resp := httputils.PostSessionHeartbeat(t, app, accesses[i], sessionIDs[i])
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	// [STEP 6] Wait for SESSION_FREEZE_TIME to elapse
	freezeTime, err := strconv.Atoi(os.Getenv("SESSION_FREEZE_TIME"))
	if err != nil || freezeTime <= 0 {
		freezeTime = 60
	}
	t.Logf("[STEP 6] Esperando tiempo de congelamiento: %d segundos", freezeTime)
	time.Sleep(time.Duration(freezeTime) * time.Second)

	// [STEP 7] Each student heartbeats again after freeze (critical)
	var heartbeat2Counter int64 = -1
	process.StartConcurrentStep(7, "Cada estudiante hace heartbeat (despues del freeze)", func() (int, []byte, error) {
		i := atomic.AddInt64(&heartbeat2Counter, 1)
		if accesses[i] == nil || sessionIDs[i] == "" {
			return 0, nil, fmt.Errorf("no session for student %d", i)
		}
		resp := httputils.PostSessionHeartbeat(t, app, accesses[i], sessionIDs[i])
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	// [STEP 8] Each student closes their session (critical)
	var closeCounter int64 = -1
	process.StartConcurrentStep(8, "Cada estudiante cierra su sesion", func() (int, []byte, error) {
		i := atomic.AddInt64(&closeCounter, 1)
		if accesses[i] == nil || sessionIDs[i] == "" {
			return 0, nil, fmt.Errorf("no session for student %d", i)
		}
		resp := httputils.PostSessionClose(t, app, accesses[i], sessionIDs[i])
		if resp.StatusCode == 200 {
			sessionIDs[i] = ""
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(n)

	process.End()
}
