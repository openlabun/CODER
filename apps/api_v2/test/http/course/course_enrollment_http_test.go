package course_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestCourseEnrollmentHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Course Enrollment HTTP")
	email := "test@test.com"
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

	// [STEP 1] Login Teacher and Student users
	process.StartStep("Inicio de sesion con usuario docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, email, password, "Teacher Test")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	teacherID = teacherAccess.UserID
	studentID = studentAccess.UserID
	process.Log(fmt.Sprintf("teacherID=%s studentID=%s", teacherID, studentID))
	process.EndStep()

	// [STEP 2] Create Course
	process.StartStep("Crear curso")
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-HTTP-ENROLL-%d", now.UnixNano())
	resp := httputils.PostCourseCreate(t, app, teacherAccess, map[string]any{
		"name":            "Course Enrollment Test",
		"description":     "Course created for enrollment test",
		"visibility":      "public",
		"visual_identity": "blue",
		"code":            fmt.Sprintf("CENR-%d", now.Unix()%100000),
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

	// [STEP 3] Enroll student to course
	process.StartStep("Inscribir estudiante al curso")
	resp = httputils.PostCourseEnroll(t, app, teacherAccess, map[string]any{
		"course_id":     courseID,
		"student_email": studentEmail,
	})
	httputils.RequireStatus(t, resp, 200, "enroll student")
	process.EndStep()

	// [STEP 4] Get Course details and validate student enrollment
	process.StartStep("Obtener datos del curso y validar inscripcion del estudiante")
	resp = httputils.GetCourseStudents(t, app, teacherAccess, courseID)
	httputils.RequireStatus(t, resp, 200, "get course students")

	var students []map[string]any
	if err := json.Unmarshal(resp.Body, &students); err != nil {
		process.Fail("get course students", fmt.Errorf("decode students list: %w", err))
	}
	foundStudent := false
	for _, st := range students {
		if httputils.StringField(st, "id") == studentID {
			foundStudent = true
			break
		}
	}
	if !foundStudent {
		process.Fail("get course students", fmt.Errorf("expected enrolled student %s in course", studentID))
	}

	resp = httputils.GetCourses(t, app, studentAccess, "enrolled", studentID, "")
	httputils.RequireStatus(t, resp, 200, "get enrolled courses")

	var enrolledCourses []map[string]any
	if err := json.Unmarshal(resp.Body, &enrolledCourses); err != nil {
		process.Fail("get enrolled courses", fmt.Errorf("decode enrolled courses list: %w", err))
	}
	foundCourseInStudentView := false
	for _, c := range enrolledCourses {
		if httputils.StringField(c, "id") == courseID {
			foundCourseInStudentView = true
			break
		}
	}
	if !foundCourseInStudentView {
		process.Fail("get enrolled courses", fmt.Errorf("expected course %s in student enrolled courses", courseID))
	}
	process.EndStep()

	// [STEP 5] Remove student from course
	process.StartStep("Retirar estudiante del curso")
	resp = httputils.DeleteCourseStudent(t, app, teacherAccess, courseID, studentID)
	httputils.RequireStatus(t, resp, 200, "remove student from course")
	process.EndStep()

	// [STEP 6] Get Course details and validate student removal
	process.StartStep("Obtener datos del curso y validar desinscripcion del estudiante")
	resp = httputils.GetCourseStudents(t, app, teacherAccess, courseID)
	httputils.RequireStatus(t, resp, 200, "get course students after removal")

	students = nil
	if err := json.Unmarshal(resp.Body, &students); err != nil {
		process.Fail("get course students after removal", fmt.Errorf("decode students list after removal: %w", err))
	}
	for _, st := range students {
		if httputils.StringField(st, "id") == studentID {
			process.Fail("get course students after removal", fmt.Errorf("student %s should not remain enrolled", studentID))
		}
	}

	resp = httputils.GetCourses(t, app, studentAccess, "enrolled", studentID, "")
	httputils.RequireStatus(t, resp, 200, "get enrolled courses after removal")

	enrolledCourses = nil
	if err := json.Unmarshal(resp.Body, &enrolledCourses); err != nil {
		process.Fail("get enrolled courses after removal", fmt.Errorf("decode enrolled courses after removal: %w", err))
	}
	for _, c := range enrolledCourses {
		if httputils.StringField(c, "id") == courseID {
			process.Fail("get enrolled courses after removal", fmt.Errorf("course %s should not remain in student list", courseID))
		}
	}
	process.EndStep()

	process.End()
}
