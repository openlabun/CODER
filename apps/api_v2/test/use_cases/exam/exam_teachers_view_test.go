package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestExamFromTeacherView(t *testing.T) {
	process := test.StartTestWithApp(t, "Exam Teacher View Use Cases")
	teacher_email := "test@test.com"
	teacher_2_email := "test2@test.com"
	student_email := "stud@test.com"
	password := "Password123!"

	var teacherID string
	var teacherCtx = context.Background()
	var teacher2Ctx = context.Background()
	var studentCtx = context.Background()
	var privateExamID string
	var teachersExamID string
	var teachersCourseExamID string
	var courseID string

	defer func() {
		if privateExamID != "" {
			t.Logf("[CLEANUP] Deleting created 'private' exam %s", privateExamID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: privateExamID})
		}
		if teachersExamID != "" {
			t.Logf("[CLEANUP] Deleting created 'teachers' exam %s", teachersExamID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: teachersExamID})
		}
		if teachersCourseExamID != "" {
			t.Logf("[CLEANUP] Deleting created 'teachers' exam for course %s", teachersCourseExamID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: teachersCourseExamID})
		}
		if courseID != "" {
			t.Logf("[CLEANUP] Deleting created course %s", courseID)
			_ = process.Application.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
		}
	}()

	// [STEP 1] Login Teacher user
	process.StartStep("Inicio de sesión con usuario docente")
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, teacher_email, password, "Teacher One")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	teacherID = teacherAccess.UserData.ID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))

	process.EndStep()

	// [STEP 2] Create "private" Exam
	process.StartStep("Crear un Examen con visibilidad 'private' para el curso")
	now := time.Now().UTC()
	start := now.Add(2 * time.Hour)
	end := start.Add(75 * time.Minute)
	endStr := end.Format(time.RFC3339)

	privateExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Private Teachers View Exam",
		Description:          "Only owner should access",
		Visibility:           string(exam_entities.VisibilityPrivate),
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

	// [STEP 3] Login Teacher 2 user
	process.StartStep("Inicio de sesión con usuario docente")
	teacher2Access := utils.EnsureAuthUserAccess(t, process.Application, teacher_2_email, password, "Teacher Two")
	teacher2Ctx = utils.BuildUserCtx(teacher2Access)
	process.Log(fmt.Sprintf("teacher2ID=%s", teacher2Access.UserData.ID))

	process.EndStep()

	// [STEP 4] Get Exam from other teacher view and validate error
	process.StartStep("Obtener datos del Examen desde vista del otro docente y validar error")
	_, err = process.Application.ExamModule.GetExamDetails.Execute(teacher2Ctx, exam_dtos.GetExamDetailsInput{ExamID: privateExamID})
	if err == nil {
		process.Fail("other teacher private exam view", fmt.Errorf("expected error when another teacher views private exam"))
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 5] Create "teachers" Exam
	process.StartStep("Crear un Examen con visibilidad 'teachers' sin curso")
	teachersExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             nil,
		Title:                "Teachers Visibility Exam",
		Description:          "Visible for teachers",
		Visibility:           string(exam_entities.VisibilityTeachers),
		StartTime:            start.Add(24 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            2700,
		TryLimit:             2,
		ProfessorID:          teacherID,
	})
	if err != nil {
		process.Fail("create teachers exam", err)
	}
	if teachersExam == nil || teachersExam.ID == "" {
		process.Fail("create teachers exam", fmt.Errorf("expected created teachers exam with ID"))
	}
	teachersExamID = teachersExam.ID

	process.EndStep()

	// [STEP 6] Get Exam from other teacher view and validate data
	process.StartStep("Obtener datos del Examen desde vista del otro docente y validar datos")
	otherTeacherExam, err := process.Application.ExamModule.GetExamDetails.Execute(teacher2Ctx, exam_dtos.GetExamDetailsInput{ExamID: teachersExamID})
	if err != nil {
		process.Fail("other teacher teachers exam view", err)
	}
	if otherTeacherExam == nil || otherTeacherExam.ID != teachersExamID {
		process.Fail("other teacher teachers exam view", fmt.Errorf("expected exam details for %s", teachersExamID))
	}
	process.EndStep()

	// [STEP 7] Get Exam from student view and validate error
	process.StartStep("Obtener datos del Examen desde vista de estudiante y validar error")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, student_email, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	_, err = process.Application.ExamModule.GetExamDetails.Execute(studentCtx, exam_dtos.GetExamDetailsInput{ExamID: teachersExamID})
	if err == nil {
		process.Fail("student teachers exam view", fmt.Errorf("expected error when student views teachers exam"))
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 8] Create Course
	process.StartStep("Crear curso")
	enrollmentCode := fmt.Sprintf("ENR-TEA-EXAM-%d", now.UnixNano())
	createdCourse, err := process.Application.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "Teachers View Course",
		Description:    "Course for teachers visibility test",
		Visibility:     string(consts.CourseVisibilityPublic),
		VisualIdentity: string(consts.CourseColourBlue),
		Code:           fmt.Sprintf("TV-%d", now.Unix()%100000),
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

	process.EndStep()

	// [STEP 9] Create a "teachers" Exam for the course
	process.StartStep("Crear un Examen con visibilidad 'teachers' para el curso")
	teachersCourseExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             &courseID,
		Title:                "Teachers Course Visibility Exam",
		Description:          "Teachers-only exam with course relation",
		Visibility:           string(exam_entities.VisibilityTeachers),
		StartTime:            start.Add(48 * time.Hour).Format(time.RFC3339),
		EndTime:              nil,
		AllowLateSubmissions: true,
		TimeLimit:            3000,
		TryLimit:             1,
		ProfessorID:          teacherID,
	})
	if err != nil {
		process.Fail("create teachers course exam", err)
	}
	if teachersCourseExam == nil || teachersCourseExam.ID == "" {
		process.Fail("create teachers course exam", fmt.Errorf("expected created teachers course exam with ID"))
	}
	teachersCourseExamID = teachersCourseExam.ID

	process.EndStep()

	// [STEP 10] Get Exam from teacher view and validate data
	process.StartStep("Obtener datos del Examen desde vista de docente y validar datos")
	otherTeacherExam, err = process.Application.ExamModule.GetExamDetails.Execute(teacher2Ctx, exam_dtos.GetExamDetailsInput{ExamID: teachersCourseExamID})
	if err != nil {
		process.Fail("other teacher teachers course exam view", err)
	}
	if otherTeacherExam == nil || otherTeacherExam.ID != teachersCourseExamID {
		process.Fail("other teacher teachers course exam view", fmt.Errorf("expected exam details for %s", teachersCourseExamID))
	}
	process.EndStep()
	process.End()
}
