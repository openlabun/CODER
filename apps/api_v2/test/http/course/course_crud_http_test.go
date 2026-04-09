package course_test

import (
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestCourseCRUDHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Course CRUD HTTP")
	email := "test@test.com"
	password := "Password123!"

	var teacherID string
	var teacherAccess *httputils.HTTPAccess
	var courseID string
	var deletedCourseID string

	defer func() {
		if courseID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando curso %s", courseID)
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
	enrollmentCode := fmt.Sprintf("ENR-HTTP-COURSE-%d", now.UnixNano())

	resp := httputils.PostCourseCreate(t, app, teacherAccess, map[string]any{
		"name":            "Course CRUD Test",
		"description":     "Course created by HTTP test",
		"visibility":      "public",
		"visual_identity": "blue",
		"code":            fmt.Sprintf("CCRUD-%d", now.Unix()%100000),
		"year":            now.Year(),
		"semester":        "01",
		"enrollment_code": enrollmentCode,
		"enrollment_url":  "https://example.test/enroll/" + enrollmentCode,
		"teacher_id":      teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create course")

	body := httputils.MustJSONMap(t, resp)
	v1 := httputils.StringField(body, "id")
	if v1 == "" {
		process.Fail("create course", fmt.Errorf("expected created course with ID"))
	}
	courseID = v1
	process.Log(fmt.Sprintf("courseID=%s", courseID))
	process.EndStep()

	// [STEP 3] Update Course
	process.StartStep("Actualizar curso")
	updatedName := "Course CRUD Test Updated"
	updatedDescription := "Course updated by HTTP test"

	resp = httputils.PostCourseUpdate(t, app, teacherAccess, courseID, map[string]any{
		"name":        updatedName,
		"description": updatedDescription,
	})
	httputils.RequireStatus(t, resp, 200, "update course")

	body = httputils.MustJSONMap(t, resp)
	v1 = httputils.StringField(body, "name")
	if v1 != updatedName {
		process.Fail("update course", fmt.Errorf("expected updated course name"))
	}
	process.Log(fmt.Sprintf("Updated course name=%q", v1))
	process.EndStep()

	// [STEP 4] Get Course details and validate
	process.StartStep("Obtener datos del curso y validarlos")
	resp = httputils.GetCourseByID(t, app, teacherAccess, courseID)
	httputils.RequireStatus(t, resp, 200, "get course details")

	body = httputils.MustJSONMap(t, resp)
	v1 = httputils.StringField(body, "id")
	v2 := httputils.StringField(body, "name")
	if v1 != courseID {
		process.Fail("get course details", fmt.Errorf("expected course details for %s", courseID))
	}
	if v2 != updatedName {
		process.Fail("get course details", fmt.Errorf("expected name %q, got %q", updatedName, v2))
	}
	process.Log(fmt.Sprintf("Course details validated for %s", courseID))
	process.EndStep()

	// [STEP 5] Create an Exam with teachers visibility for the course
	process.StartStep("Eliminar curso")
	resp = httputils.DeleteCourseByID(t, app, teacherAccess, courseID)
	httputils.RequireStatus(t, resp, 200, "delete course")

	deletedCourseID = courseID
	courseID = ""
	process.EndStep()

	// [STEP 6] Get Course details and validate error
	process.StartStep("Obtener datos del curso y validar error")
	resp = httputils.GetCourseByID(t, app, teacherAccess, deletedCourseID)
	if resp.StatusCode == 200 {
		process.Fail("get course details after delete", fmt.Errorf("expected error after deleting course"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
