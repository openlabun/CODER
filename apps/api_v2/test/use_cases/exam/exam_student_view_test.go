package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
	exam_consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"

	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestExamFromStudentView(t *testing.T) {
	process := test.StartTestWithApp(t, "Exam Student View Use Cases")
	teacher_email := "test@test.com"
	student_email := "stud@test.com"
	password := "Password123!"

	var teacherID string
	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var courseID string
	var privateExamID string
	var courseExamID string
	var publicExamID string

	defer func() {
		if privateExamID != "" {
			t.Logf("[CLEANUP] Deleting created 'private' exam %s", privateExamID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: privateExamID})
		}
		if courseExamID != "" {
			t.Logf("[CLEANUP] Deleting created 'course' exam %s", courseExamID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: courseExamID})
		}
		if publicExamID != "" {
			t.Logf("[CLEANUP] Deleting created 'public' exam %s", publicExamID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: publicExamID})
		}
		if courseID != "" {
			t.Logf("[CLEANUP] Deleting created course %s", courseID)
			_ = process.Application.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
		}
	}()

	// [STEP 1] Login Teacher user
	process.StartStep("Inicio de sesión con usuario docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacher_email, password, "Teacher Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	teacherID = teacherAccess.UserData.ID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))

	process.EndStep()

	// [STEP 2] Create Course
	process.StartStep("Crear curso")
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-STUD-EXAM-%d", now.UnixNano())
	createdCourse, err := process.Application.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "Student View Course",
		Description:    "Course for student exam visibility test",
		Visibility:     string(consts.CourseVisibilityPublic),
		VisualIdentity: string(consts.CourseColourBlue),
		Code:           fmt.Sprintf("STV-%d", now.Unix()%100000),
		Year:           now.Year(),
		Semester:       string(consts.AcademicFirstPeriod),
		EnrollmentCode: enrollmentCode,
		EnrollmentURL:  "https://example.test/enroll/" + enrollmentCode,
		TeacherID:      teacherID,
	})
	if err != nil {
		process.Fail("create course", err)
	}
	if createdCourse == nil || createdCourse.ID == "" {
		process.Fail("create course", fmt.Errorf("expected created course with ID"))
	}
	courseID = createdCourse.ID
	process.Log(fmt.Sprintf("courseID=%s", courseID))

	process.EndStep()

	// [STEP 3] Create "private" Exam
	process.StartStep("Crear un Examen con visibilidad 'private' para el curso")
	start := now.Add(2 * time.Hour)
	end := start.Add(60 * time.Minute)
	endStr := end.Format(time.RFC3339)
	privateExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             &courseID,
		Title:                "Private Student View Exam",
		Description:          "Student should not access private exam",
		Visibility:           string(exam_consts.VisibilityPrivate),
		StartTime:            start.Format(time.RFC3339),
		EndTime:              &endStr,
		AllowLateSubmissions: false,
		TimeLimit:            3600,
		TryLimit:             1,
		ProfessorID:          teacherID,
	})
	if err != nil {
		process.Fail("create private exam", err)
	}
	if privateExam == nil || privateExam.ID == "" {
		process.Fail("create private exam", fmt.Errorf("expected created private exam with ID"))
	}
	privateExamID = privateExam.ID

	process.EndStep()

	// [STEP 4] Login Student user
	process.StartStep("Inicio de sesión con usuario estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, student_email, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	process.Log(fmt.Sprintf("studentID=%s", studentAccess.UserData.ID))

	process.EndStep()

	// [STEP 5] Get Exam from student view and validate error
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar error")
	_, err = process.Application.ExamModule.GetExamDetails.Execute(studentCtx, exam_dtos.GetExamDetailsInput{ExamID: privateExamID})
	if err == nil {
		process.Fail("student private exam view", fmt.Errorf("expected error when student views private exam"))
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 6] Create "course" Exam
	process.StartStep("Crear un Examen con visibilidad 'course' para el curso")
	courseExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             &courseID,
		Title:                "Course Student View Exam",
		Description:          "Student can access after enrollment",
		Visibility:           string(exam_consts.VisibilityCourse),
		StartTime:            start.Add(24 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            5400,
		TryLimit:             2,
		ProfessorID:          teacherID,
	})
	if err != nil {
		process.Fail("create course exam", err)
	}
	if courseExam == nil || courseExam.ID == "" {
		process.Fail("create course exam", fmt.Errorf("expected created course exam with ID"))
	}
	courseExamID = courseExam.ID

	process.EndStep()

	// [STEP 7] Get Exam from student view and validate error
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar error")
	_, err = process.Application.ExamModule.GetExamDetails.Execute(studentCtx, exam_dtos.GetExamDetailsInput{ExamID: courseExamID})
	if err == nil {
		process.Fail("student course exam view before enroll", fmt.Errorf("expected error when student is not enrolled"))
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 8] Enroll student to course
	process.StartStep("Inscribir estudiante al curso")
	_, err = process.Application.CourseModule.EnrollInCourse.Execute(teacherCtx, course_dtos.EnrolledInCourseInput{
		CourseID:     courseID,
		StudentEmail: &student_email,
	})
	if err != nil {
		process.Fail("enroll student", err)
	}

	process.EndStep()

	// [STEP 9] Get Exam from student view and validate data
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar datos")
	courseExams, err := process.Application.ExamModule.GetExamsByCourse.Execute(studentCtx, exam_dtos.GetExamsByCourseInput{CourseID: courseID})
	if err != nil {
		process.Fail("student get exams by course", err)
	}
	foundCourseExam := false
	foundPrivateExam := false
	for _, exam := range courseExams {
		if exam == nil {
			continue
		}
		if exam.ID == courseExamID {
			foundCourseExam = true
		}
		if exam.ID == privateExamID {
			foundPrivateExam = true
		}
	}
	if !foundCourseExam {
		process.Fail("student get exams by course", fmt.Errorf("expected to find course exam %s", courseExamID))
	}
	if foundPrivateExam {
		process.Fail("student get exams by course", fmt.Errorf("did not expect private exam %s in student list", privateExamID))
	}
	process.Log(fmt.Sprintf("Student exam list validated. total=%d", len(courseExams)))

	process.EndStep()

	// [STEP 10] Create "public" Exam, and its not associated with a course
	process.StartStep("Crear un Examen con visibilidad 'public' para el curso")
	publicExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Public Student View Exam",
		Description:          "Public exam for student visibility",
		Visibility:           string(exam_consts.VisibilityPublic),
		StartTime:            start.Add(48 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            2400,
		TryLimit:             1,
		ProfessorID:          teacherID,
	})
	if err != nil {
		process.Fail("create public exam", err)
	}
	if publicExam == nil || publicExam.ID == "" {
		process.Fail("create public exam", fmt.Errorf("expected created public exam with ID"))
	}
	publicExamID = publicExam.ID

	process.EndStep()

	// [STEP 11] Get Exam from student view and validate data
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar datos")
	publicExams, err := process.Application.ExamModule.GetPublicExams.Execute(studentCtx)
	if err != nil {
		process.Fail("student get public exams", err)
	}

	process.Log(fmt.Sprintf("Student got %d public exams", len(publicExams)))

	foundPublicExam := false
	for _, exam := range publicExams {
		if exam != nil && exam.ID == publicExamID {
			foundPublicExam = true
			break
		}
	}
	if !foundPublicExam {
		process.Fail("student get public exams", fmt.Errorf("expected public exam %s in public exam list", publicExamID))
	}
	process.Log(fmt.Sprintf("Public exam list validated. total=%d", len(publicExams)))

	process.EndStep()
	process.End()
}
