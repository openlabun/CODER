package course_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestCourseFromStudentViewHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Course Student View HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherID string
	var studentID string
	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var courseID string

	defer func() {
		if courseID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando curso %s", courseID)
			_ = httputils.DeleteCourseByID(t, app, teacherAccess, courseID)
		}
	}()

	// [STEP 1] Login Teacher user 
	process.StartStep("Inicio de sesion con usuario docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	teacherID = teacherAccess.UserID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))
	process.EndStep()

	// [STEP 2] Create Course
	process.StartStep("Crear curso")
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-HTTP-STUDVIEW-%d", now.UnixNano())
	resp := httputils.PostCourseCreate(t, app, teacherAccess, map[string]any{
		"name":            "Course Student View Test",
		"description":     "Course created for student view HTTP test",
		"visibility":      "public",
		"visual_identity": "blue",
		"code":            fmt.Sprintf("CSV-%d", now.Unix()%100000),
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

	// [STEP 3] Login Student user
	process.StartStep("Inicio de sesion con usuario estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	studentID = studentAccess.UserID
	process.Log(fmt.Sprintf("studentID=%s", studentID))
	process.EndStep()

	// [STEP 4] Get Course details from student view and validate error (not enrolled yet)
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar error")
	resp = httputils.GetCourses(t, app, studentAccess, "enrolled", studentID, "")
	httputils.RequireStatus(t, resp, 200, "get enrolled courses before enrollment")

	var enrolledBefore []map[string]any
	if err := json.Unmarshal(resp.Body, &enrolledBefore); err != nil {
		process.Fail("get enrolled courses before enrollment", fmt.Errorf("decode enrolled courses before enrollment: %w", err))
	}
	for _, c := range enrolledBefore {
		if httputils.StringField(c, "id") == courseID {
			process.Fail("get enrolled courses before enrollment", fmt.Errorf("course %s should not be visible before enrollment", courseID))
		}
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 5] Enroll student to course
	process.StartStep("Inscribir estudiante al curso")
	resp = httputils.PostCourseEnroll(t, app, teacherAccess, map[string]any{
		"course_id":     courseID,
		"student_email": studentEmail,
	})
	httputils.RequireStatus(t, resp, 200, "enroll student")
	process.EndStep()

	// [STEP 6] Get Course details from student view and validate data
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar datos")
	resp = httputils.GetCourses(t, app, studentAccess, "enrolled", studentID, "")
	httputils.RequireStatus(t, resp, 200, "get enrolled courses after enrollment")

	var enrolledAfter []map[string]any
	if err := json.Unmarshal(resp.Body, &enrolledAfter); err != nil {
		process.Fail("get enrolled courses after enrollment", fmt.Errorf("decode enrolled courses after enrollment: %w", err))
	}
	foundCourse := false
	for _, c := range enrolledAfter {
		if httputils.StringField(c, "id") == courseID {
			foundCourse = true
			break
		}
	}
	if !foundCourse {
		process.Fail("get enrolled courses after enrollment", fmt.Errorf("expected course %s in student list", courseID))
	}
	process.EndStep()

	// [STEP 7] Remove student from course
	process.StartStep("Retirar estudiante del curso")
	resp = httputils.DeleteCourseStudent(t, app, teacherAccess, courseID, studentID)
	httputils.RequireStatus(t, resp, 200, "remove student from course")
	process.EndStep()

	// [STEP 8] Get Course details from student view and validate error (after removal)
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar error")
	resp = httputils.GetCourses(t, app, studentAccess, "enrolled", studentID, "")
	httputils.RequireStatus(t, resp, 200, "get enrolled courses after removal")

	var enrolledAfterRemoval []map[string]any
	if err := json.Unmarshal(resp.Body, &enrolledAfterRemoval); err != nil {
		process.Fail("get enrolled courses after removal", fmt.Errorf("decode enrolled courses after removal: %w", err))
	}
	for _, c := range enrolledAfterRemoval {
		if httputils.StringField(c, "id") == courseID {
			process.Fail("get enrolled courses after removal", fmt.Errorf("course %s should not be visible after unenrollment", courseID))
		}
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	process.End()
}
