package exam_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamCRUDHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Exam CRUD HTTP")
	email := "test@test.com"
	password := "Password123!"

	var teacherID string
	var teacherAccess *httputils.HTTPAccess
	var courseID string
	var examID string
	var deletedExamID string

	defer func() {
		if examID != "" && teacherAccess != nil {
			_ = httputils.DeleteExamByID(t, app, teacherAccess, examID)
		}
		if courseID != "" && teacherAccess != nil {
			_ = httputils.DeleteCourseByID(t, app, teacherAccess, courseID)
		}
	}()

	// [STEP 1] Login Teacher user
	process.StartStep("Inicio de sesion con usuario docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, email, password, "Teacher Test")
	teacherID = teacherAccess.UserID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))
	process.EndStep()

	// [STEP 2] Create Course
	process.StartStep("Crear curso")
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-EXAM-%d", now.UnixNano())
	resp := httputils.PostCourseCreate(t, app, teacherAccess, map[string]any{
		"name":            "Exam CRUD Course",
		"description":     "Course for exam CRUD test",
		"visibility":      "public",
		"visual_identity": "blue",
		"code":            fmt.Sprintf("EXCRUD-%d", now.Unix()%100000),
		"year":            now.Year(),
		"semester":        "01",
		"enrollment_code": enrollmentCode,
		"enrollment_url":  "https://example.test/enroll/" + enrollmentCode,
		"teacher_id":      teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create course")
	courseID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	if courseID == "" {
		process.Fail("create course", fmt.Errorf("expected created course with ID"))
	}
	process.Log(fmt.Sprintf("courseID=%s", courseID))
	process.EndStep()

	// [STEP 3] Create Exam
	process.StartStep("Crea un Examen")
	startTime := now.Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)
	resp = httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"course_id":              courseID,
		"title":                  "Exam CRUD Test",
		"description":            "Created by HTTP test",
		"visibility":             "course",
		"start_time":             startTime.Format(time.RFC3339),
		"end_time":               endTimeStr,
		"allow_late_submissions": false,
		"time_limit":             5400,
		"try_limit":              2,
		"professor_id":           teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create exam")
	examID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	if examID == "" {
		process.Fail("create exam", fmt.Errorf("expected created exam with ID"))
	}
	process.Log(fmt.Sprintf("examID=%s", examID))
	process.EndStep()

	// [STEP 4] Attempt to create an exam with invalid course ID (expect error)
	process.StartStep("Actualizar Examen")
	updatedTitle := "Exam CRUD Test Updated"
	updatedDescription := "Updated by HTTP test"
	resp = httputils.PatchExamUpdate(t, app, teacherAccess, examID, map[string]any{
		"title":       updatedTitle,
		"description": updatedDescription,
		"try_limit":   3,
	})
	httputils.RequireStatus(t, resp, 200, "update exam")
	if httputils.StringField(httputils.MustJSONMap(t, resp), "title") != updatedTitle {
		process.Fail("update exam", fmt.Errorf("expected updated exam title"))
	}
	process.Log(fmt.Sprintf("Updated exam title=%q", updatedTitle))
	process.EndStep()

	// [STEP 5] Get Exam details and validate
	process.StartStep("Obtener datos del Examen y validarlos")
	resp = httputils.GetExamByID(t, app, teacherAccess, examID)
	httputils.RequireStatus(t, resp, 200, "get exam details")
	body := httputils.MustJSONMap(t, resp)
	v1 := httputils.StringField(body, "id")
	v2 := httputils.StringField(body, "title")
	if v1 != examID {
		process.Fail("get exam details", fmt.Errorf("expected exam details for %s", examID))
	}
	if v2 != updatedTitle {
		process.Fail("get exam details", fmt.Errorf("expected title %q, got %q", updatedTitle, v2))
	}
	process.Log(fmt.Sprintf("Exam details validated for %s", examID))
	process.EndStep()

	// [STEP 6] Attempt to update the exam with invalid end time (expect error)
	process.StartStep("Actualizar Examen con valores incorrectos (EndTime con valor de ayer)")
	yesterday := now.Add(-24 * time.Hour).Format(time.RFC3339)
	resp = httputils.PatchExamUpdate(t, app, teacherAccess, examID, map[string]any{"end_time": yesterday})
	if resp.StatusCode == 200 {
		process.Fail("invalid exam update", fmt.Errorf("expected error when updating exam with end time before start time"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Update the exam with valid data
	process.StartStep("Eliminar Examen")
	resp = httputils.DeleteExamByID(t, app, teacherAccess, examID)
	httputils.RequireStatus(t, resp, 200, "delete exam")
	deletedExamID = examID
	examID = ""
	process.EndStep()

	// [STEP 8] Attempt to get the deleted exam details (expect error)
	process.StartStep("Obtener datos del Examen y validar error")
	resp = httputils.GetExamByID(t, app, teacherAccess, deletedExamID)
	if resp.StatusCode == 200 {
		process.Fail("get exam details after delete", fmt.Errorf("expected error after deleting exam"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
