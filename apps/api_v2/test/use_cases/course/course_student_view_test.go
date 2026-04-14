package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestCourseFromStudentView(t *testing.T) {
	process := test.StartTestWithApp(t, "Course Student View Use Cases")
	teacher_email := "test@test.com"
	student_email := "stud@test.com"
	password := "Password123!"

	var teacherID string
	var studentID string
	var teacherCtx = context.Background()
	var studentCtx = context.Background()
	var courseID string

	defer func() {
		if courseID != "" {
			process.Log(fmt.Sprintf("[CLEANUP] Deleting created course %s", courseID))
			if err := process.Application.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID}); err != nil {
				process.Log(fmt.Sprintf("[CLEANUP] Error deleting course %s: %v", courseID, err))
			}
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
	enrollmentCode := fmt.Sprintf("ENR-UC-STUDVIEW-%d", now.UnixNano())
	createdCourse, err := process.Application.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "Course Student View Test",
		Description:    "Course created for student view use case",
		Visibility:     string(consts.CourseVisibilityPublic),
		VisualIdentity: string(consts.CourseColourBlue),
		Code:           fmt.Sprintf("CSV-%d", now.Unix()%100000),
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

	// [STEP 3] Login Student user
	process.StartStep("Inicio de sesión con usuario estudiante")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, student_email, password, "Student Test")
	studentCtx = utils.BuildUserCtx(studentAccess)
	studentID = studentAccess.UserData.ID
	process.Log(fmt.Sprintf("studentID=%s", studentID))

	process.EndStep()

	// [STEP 4] Get Course from student view and validate error
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar error")
	enrolledBefore, err := process.Application.CourseModule.GetEnrolledCourses.Execute(studentCtx)
	if err != nil {
		process.Fail("get enrolled courses before enrollment", err)
	}
	for _, c := range enrolledBefore {
		if c != nil && c.ID == courseID {
			process.Fail("get enrolled courses before enrollment", fmt.Errorf("course %s should not be visible before enrollment", courseID))
		}
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 5] Enroll student to course
	process.StartStep("Inscribir estudiante al curso")
	_, err = process.Application.CourseModule.EnrollInCourse.Execute(teacherCtx, course_dtos.EnrolledInCourseInput{
		CourseID:     courseID,
		StudentEmail: &student_email,
	})
	if err != nil {
		process.Fail("enroll student", err)
	}

	process.EndStep()

	// [STEP 6] Get Course from student view and validate data
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar datos")
	enrolledAfter, err := process.Application.CourseModule.GetEnrolledCourses.Execute(studentCtx)
	if err != nil {
		process.Fail("get enrolled courses after enrollment", err)
	}
	foundCourse := false
	for _, c := range enrolledAfter {
		if c != nil && c.ID == courseID {
			foundCourse = true
			break
		}
	}
	if !foundCourse {
		process.Fail("get enrolled courses after enrollment", fmt.Errorf("expected course %s in student list", courseID))
	}

	process.EndStep()

	// [STEP 7] Unenroll student from course
	process.StartStep("Retirar estudiante del curso")
	err = process.Application.CourseModule.RemoveStudentFromCourse.Execute(teacherCtx, course_dtos.RemoveStudentFromCourseInput{
		CourseID:  courseID,
		StudentID: studentID,
	})
	if err != nil {
		process.Fail("remove student from course", err)
	}

	process.EndStep()

	// [STEP 8] Get Course from student view and validate error
	process.StartStep("Obtener datos del curso desde vista de estudiante y validar error")
	enrolledAfterRemoval, err := process.Application.CourseModule.GetEnrolledCourses.Execute(studentCtx)
	if err != nil {
		process.Fail("get enrolled courses after removal", err)
	}
	for _, c := range enrolledAfterRemoval {
		if c != nil && c.ID == courseID {
			process.Fail("get enrolled courses after removal", fmt.Errorf("course %s should not be visible after unenrollment", courseID))
		}
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()
	process.End()
}
