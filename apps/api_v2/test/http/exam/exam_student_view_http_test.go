package exam_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

func TestExamFromStudentViewHTTP(t *testing.T) {
	process, app := httputils.StartHTTPTestWithApp(t, "Exam Student View HTTP")
	teacherEmail := "test@test.com"
	studentEmail := "stud@test.com"
	password := "Password123!"

	var teacherID string
	var teacherAccess *httputils.HTTPAccess
	var studentAccess *httputils.HTTPAccess
	var courseID string
	var privateExamID string
	var courseExamID string
	var publicExamID string

	defer func() {
		if privateExamID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen privado %s", privateExamID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, privateExamID)
		}
		if courseExamID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen del curso %s", courseExamID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, courseExamID)
		}
		if publicExamID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando examen publico %s", publicExamID)
			_ = httputils.DeleteExamByID(t, app, teacherAccess, publicExamID)
		}
		if courseID != "" && teacherAccess != nil {
			t.Logf("[CLEANUP] Eliminando curso %s", courseID)
			_ = httputils.DeleteCourseByID(t, app, teacherAccess, courseID)
		}
	}()

	// [STEP 1] Login as teacher
	process.StartStep("Inicio de sesion con usuario docente")
	teacherAccess = httputils.EnsureAuthUserAccess(t, app, teacherEmail, password, "Teacher Test")
	teacherID = teacherAccess.UserID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))
	process.EndStep()

	// [STEP 2] Create a course
	process.StartStep("Crear curso")
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-STUD-EXAM-%d", now.UnixNano())
	resp := httputils.PostCourseCreate(t, app, teacherAccess, map[string]any{
		"name":            "Student View Course",
		"description":     "Course for student exam visibility test",
		"visibility":      "public",
		"visual_identity": "blue",
		"code":            fmt.Sprintf("STV-%d", now.Unix()%100000),
		"year":            now.Year(),
		"semester":        "10",
		"enrollment_code": enrollmentCode,
		"enrollment_url":  "https://example.test/enroll/" + enrollmentCode,
		"teacher_id":      teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create course")
	courseID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 3] Create a private exam for the course
	process.StartStep("Crear un Examen con visibilidad private para el curso")
	start := now.Add(2 * time.Hour)
	end := start.Add(60 * time.Minute)
	resp = httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"course_id":              courseID,
		"title":                  "Private Student View Exam",
		"description":            "Student should not access private exam",
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

	// [STEP 4] Login as student
	process.StartStep("Inicio de sesion con usuario estudiante")
	studentAccess = httputils.EnsureAuthUserAccess(t, app, studentEmail, password, "Student Test")
	process.Log(fmt.Sprintf("studentID=%s", studentAccess.UserID))
	process.EndStep()

	// [STEP 5] Try to get the private exam details as student (expect error)
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar error")
	resp = httputils.GetExamByID(t, app, studentAccess, privateExamID)
	if resp.StatusCode == 200 {
		process.Fail("student private exam view", fmt.Errorf("expected error when student views private exam"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 6] Get public exams as student and validate the private exam is not listed
	process.StartStep("Crear un Examen con visibilidad course para el curso")
	resp = httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"course_id":              courseID,
		"title":                  "Course Student View Exam",
		"description":            "Student can access after enrollment",
		"visibility":             "course",
		"start_time":             start.Add(24 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             5400,
		"try_limit":              2,
		"professor_id":           teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create course exam")
	courseExamID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 7] Try to get the course exam details as student before enrollment (expect error)
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar error")
	resp = httputils.GetExamByID(t, app, studentAccess, courseExamID)
	if resp.StatusCode == 200 {
		process.Fail("student course exam view before enroll", fmt.Errorf("expected error when student is not enrolled"))
	}
	process.Log("Recibio ERROR, como se esperaba")
	process.EndStep()

	// [STEP 8] Enroll student to course
	process.StartStep("Inscribir estudiante al curso")
	resp = httputils.PostCourseEnroll(t, app, teacherAccess, map[string]any{
		"course_id":     courseID,
		"student_email": studentEmail,
	})
	httputils.RequireStatus(t, resp, 200, "enroll student")
	process.EndStep()

	// [STEP 9] Get the course exam details as student after enrollment and validate data
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar datos")
	resp = httputils.GetExamsByCourseID(t, app, studentAccess, courseID)
	httputils.RequireStatus(t, resp, 200, "student get exams by course")
	var courseExams []map[string]any
	if err := json.Unmarshal(resp.Body, &courseExams); err != nil {
		process.Fail("student get exams by course", fmt.Errorf("decode exams by course: %w", err))
	}
	foundCourseExam := false
	foundPrivateExam := false
	for _, exam := range courseExams {
		if httputils.StringField(exam, "id") == courseExamID {
			foundCourseExam = true
		}
		if httputils.StringField(exam, "id") == privateExamID {
			foundPrivateExam = true
		}
	}
	if !foundCourseExam {
		process.Fail("student get exams by course", fmt.Errorf("expected to find course exam %s", courseExamID))
	}
	if foundPrivateExam {
		process.Fail("student get exams by course", fmt.Errorf("did not expect private exam %s in student list", privateExamID))
	}
	process.EndStep()

	// [STEP 10] Get public exams as student and validate the course exam is listed but private exam is not listed
	process.StartStep("Crear un Examen con visibilidad public")
	resp = httputils.PostExamCreate(t, app, teacherAccess, map[string]any{
		"title":                  "Public Student View Exam",
		"description":            "Public exam for student visibility",
		"visibility":             "public",
		"start_time":             start.Add(48 * time.Hour).Format(time.RFC3339),
		"allow_late_submissions": true,
		"time_limit":             2400,
		"try_limit":              1,
		"professor_id":           teacherID,
	})
	httputils.RequireStatus(t, resp, 201, "create public exam")
	publicExamID = httputils.StringField(httputils.MustJSONMap(t, resp), "id")
	process.EndStep()

	// [STEP 11] Get public exams as student and validate the public exam is listed but private exam is not listed
	process.StartStep("Obtener examenes publicos desde vista estudiante y validar datos")
	resp = httputils.GetPublicExams(t, app, studentAccess)
	httputils.RequireStatus(t, resp, 200, "student get public exams")
	var publicExams []map[string]any
	if err := json.Unmarshal(resp.Body, &publicExams); err != nil {
		process.Fail("student get public exams", fmt.Errorf("decode public exams: %w", err))
	}
	foundPublicExam := false
	for _, exam := range publicExams {
		if httputils.StringField(exam, "id") == publicExamID {
			foundPublicExam = true
			break
		}
	}
	if !foundPublicExam {
		process.Fail("student get public exams", fmt.Errorf("expected public exam %s in list", publicExamID))
	}
	process.EndStep()

	process.End()
}
