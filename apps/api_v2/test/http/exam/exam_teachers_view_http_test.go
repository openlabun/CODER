package exam_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamFromTeacherViewHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Exam Teacher View HTTP")
	teacherEmail := "test@test.com"
	teacher2Email := "test2@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherID string
	var teacherAccess *httputils.HTTPAccess
	var teacher2Access *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var privateExamID string
	var teachersExamID string
	var teachersCourseExamID string
	var courseID string

	defer func() {
		if privateExamID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen privado %s", privateExamID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, privateExamID)
		}
		if teachersExamID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen para docentes %s", teachersExamID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, teachersExamID)
		}
		if teachersCourseExamID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen para docentes con curso %s", teachersCourseExamID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, teachersCourseExamID)
		}
		if courseID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando curso %s", courseID)
			_ = httputils.DeleteCourseByID(t, app, teacherAccess, courseID)
		}
	}()

	// [STEP 1] Login as first teacher
	process.StartStep("Inicio de sesion con usuario docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher One")
	teacherID = teacherAccess.UserID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))
	process.EndStep()

	// [STEP 2] Create a private exam
	process.StartStep("Crear un Examen con visibilidad private")
	now := time.Now().UTC()
	start := now.Add(2 * time.Hour)
	end := start.Add(75 * time.Minute)
	resp := httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Private Teachers View Exam",
		"description":            "Only owner should access",
		"visibility":             "private",
		"start_time":             start.Format(time.RFC3339),
		"end_time":               end.Format(time.RFC3339),
		"allow_late_submissions": false,
		"time_limit":             3600,
		"try_limit":              1,
		"professor_id":           teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create private exam")
	privateExamID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create a teachers-only exam without course relation
	process.StartStep("Inicio de sesion con segundo docente")
	teacher2Access = httputils.EnsureAuthUserAccess(t, app, teacher2Email, password, "Teacher Two")
	process.EndStep()

	// [STEP 4] Try to get the private exam details as another teacher (expect error)
	process.StartStep("Obtener datos del Examen desde vista del otro docente y validar error")
	resp = httputils.GetExamByID(t, app, teacher2Access, privateExamID)
	if resp.StatusCode == 200 {
		process.Fail("other teacher private exam view", fmt.Errorf("expected error when another teacher views private exam"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 5] Create a teachers-only exam without course relation
	process.StartStep("Crear un Examen con visibilidad teachers sin curso")
	resp = httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Teachers Visibility Exam",
		"description":            "Visible for teachers",
		"visibility":             "teachers",
		"start_time":             start.Add(24 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             2700,
		"try_limit":              2,
		"professor_id":           teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create teachers exam")
	teachersExamID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 6] Get the teachers-only exam details as another teacher and validate data
	process.StartStep("Obtener datos del Examen desde vista del otro docente y validar datos")
	resp = httputils.GetExamByID(t, app, teacher2Access, teachersExamID)
	httputils.RequireStatus(t, resp, 200, "other teacher teachers exam view")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "id") != teachersExamID {
		process.Fail("other teacher teachers exam view", fmt.Errorf("expected exam details for %s", teachersExamID))
	}
	process.EndStep()

	// [STEP 7] Try to get the teachers-only exam details as student (expect error)
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar error")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	resp = httputils.GetExamByID(t, app, studentAccess, teachersExamID)
	if resp.StatusCode == 200 {
		process.Fail("student teachers exam view", fmt.Errorf("expected error when student views teachers exam"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 8] Create a course and an exam with teachers visibility related to that course
	process.StartStep("Crear curso")
	enrollmentCode := fmt.Sprintf("ENR-TEA-EXAM-%d", now.UnixNano())
	resp = httputils.PostCourseCreate(t, app, teacherAccess, map[string]any{
		"name":            "Teachers View Course",
		"description":     "Course for teachers visibility test",
		"visibility":      "public",
		"visual_identity": "blue",
		"code":            fmt.Sprintf("TV-%d", now.Unix()%100000),
		"year":            now.Year(),
		"semester":        "10",
		"enrollment_code": enrollmentCode,
		"enrollment_url":  "https://example.test/enroll/" + enrollmentCode,
		"teacher_id":      teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create course")
	courseID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 9] Create an exam with teachers visibility related to the course
	process.StartStep("Crear un Examen con visibilidad teachers para el curso")
	resp = httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"course_id":              courseID,
		"title":                  "Teachers Course Visibility Exam",
		"description":            "Teachers-only exam with course relation",
		"visibility":             "teachers",
		"start_time":             start.Add(48 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             3000,
		"try_limit":              1,
		"professor_id":           teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create teachers course exam")
	teachersCourseExamID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 10] Try to get the teachers-only course exam details as student (expect error)
	process.StartStep("Obtener datos del Examen desde vista de docente y validar datos")
	resp = httputils.GetExamByID(t, app, teacher2Access, teachersCourseExamID)
	httputils.RequireStatus(t, resp, 200, "other teacher teachers course exam view")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "id") != teachersCourseExamID {
		process.Fail("other teacher teachers course exam view", fmt.Errorf("expected exam details for %s", teachersCourseExamID))
	}
	process.EndStep()

	process.End()
}
