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

func TestCourseEnrollment(t *testing.T) {
	process := test.StartTestWithApp(t, "Course Enrollment")
	email := "test@test.com"
	studentEmail := "stud@test.com"
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
	teacherAccess := utils.EnsureAuthUserAccess(t, process.Application, email, password, "Teacher Test")
	studentAccess := utils.EnsureAuthUserAccess(t, process.Application, studentEmail, password, "Student Test")
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	studentCtx = utils.BuildUserCtx(studentAccess)
	teacherID = teacherAccess.UserData.ID
	studentID = studentAccess.UserData.ID
	process.Log(fmt.Sprintf("teacherID=%s studentID=%s", teacherID, studentID))

	process.EndStep()

	// [STEP 2] Create Course
	process.StartStep("Crear curso")
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-ENROLL-%d", now.UnixNano())
	createdCourse, err := process.Application.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "Course Enrollment Test",
		Description:    "Course created for enrollment test",
		Visibility:     string(consts.CourseVisibilityPublic),
		VisualIdentity: string(consts.CourseColourBlue),
		Code:           fmt.Sprintf("CENR-%d", now.Unix()%100000),
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

	// [STEP 3] Enroll student to course
	process.StartStep("Inscribir estudiante al curso")
	_, err = process.Application.CourseModule.EnrollInCourse.Execute(teacherCtx, course_dtos.EnrolledInCourseInput{
		CourseID:     courseID,
		StudentEmail: &studentEmail,
	})
	if err != nil {
		process.Fail("enroll student", err)
	}

	process.EndStep()

	// [STEP 4] Get Course and validate student enrollment
	process.StartStep("Obtener datos del curso y validar inscripción del estudiante")
	students, err := process.Application.CourseModule.GetCourseStudents.Execute(teacherCtx, course_dtos.GetCourseStudentsInput{CourseID: courseID})
	if err != nil {
		process.Fail("get course students", err)
	}
	foundStudent := false
	for _, st := range students {
		if st != nil && st.ID == studentID {
			foundStudent = true
			break
		}
	}
	if !foundStudent {
		process.Fail("get course students", fmt.Errorf("expected enrolled student %s in course", studentID))
	}

	enrolledCourses, err := process.Application.CourseModule.GetEnrolledCourses.Execute(studentCtx)
	if err != nil {
		process.Fail("get enrolled courses", err)
	}
	foundCourseInStudentView := false
	for _, c := range enrolledCourses {
		if c != nil && c.ID == courseID {
			foundCourseInStudentView = true
			break
		}
	}
	if !foundCourseInStudentView {
		process.Fail("get enrolled courses", fmt.Errorf("expected course %s in student enrolled courses", courseID))
	}

	process.EndStep()

	// [STEP 5] Unenroll student from course
	process.StartStep("Retirar estudiante del curso")
	err = process.Application.CourseModule.RemoveStudentFromCourse.Execute(teacherCtx, course_dtos.RemoveStudentFromCourseInput{
		CourseID:  courseID,
		StudentID: studentID,
	})
	if err != nil {
		process.Fail("remove student from course", err)
	}

	process.EndStep()

	// [STEP 6] Get Course and validate student unenrollment
	process.StartStep("Obtener datos del curso y validar desinscripción del estudiante")
	studentsAfterRemoval, err := process.Application.CourseModule.GetCourseStudents.Execute(teacherCtx, course_dtos.GetCourseStudentsInput{CourseID: courseID})
	if err != nil {
		process.Fail("get course students after removal", err)
	}
	for _, st := range studentsAfterRemoval {
		if st != nil && st.ID == studentID {
			process.Fail("get course students after removal", fmt.Errorf("student %s should not remain enrolled", studentID))
		}
	}

	enrolledCoursesAfterRemoval, err := process.Application.CourseModule.GetEnrolledCourses.Execute(studentCtx)
	if err != nil {
		process.Fail("get enrolled courses after removal", err)
	}
	for _, c := range enrolledCoursesAfterRemoval {
		if c != nil && c.ID == courseID {
			process.Fail("get enrolled courses after removal", fmt.Errorf("course %s should not remain in student list", courseID))
		}
	}

	process.EndStep()

	process.End()
}
