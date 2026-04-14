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

func TestCourseCRUD(t *testing.T) {
	process := test.StartTestWithApp(t, "Course CRUD Use Cases")
	email := "test@test.com"
	password := "Password123!"

	var teacherID string
	var teacherCtx = context.Background()
	var courseID string
	var deletedCourseID string

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
	teacherCtx = utils.BuildUserCtx(teacherAccess)
	teacherID = teacherAccess.UserData.ID
	process.Log(fmt.Sprintf("teacherID=%s", teacherID))

	process.EndStep()

	// [STEP 2] Create Course
	process.StartStep("Crear curso")
	now := time.Now().UTC()
	enrollmentCode := fmt.Sprintf("ENR-UC-COURSE-%d", now.UnixNano())

	createdCourse, err := process.Application.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "Course CRUD Test",
		Description:    "Course created by use-case test",
		Visibility:     string(consts.CourseVisibilityPublic),
		VisualIdentity: string(consts.CourseColourBlue),
		Code:           fmt.Sprintf("CCRUD-%d", now.Unix()%100000),
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

	// [STEP 3] Update Course
	process.StartStep("Actualizar curso")
	updatedName := "Course CRUD Test Updated"
	updatedDescription := "Course updated by use-case test"
	updatedCourse, err := process.Application.CourseModule.UpdateCourse.Execute(teacherCtx, course_dtos.UpdateCourseInput{
		ID:          courseID,
		Name:        &updatedName,
		Description: &updatedDescription,
	})
	if err != nil {
		process.Fail("update course", err)
	}
	if updatedCourse == nil || updatedCourse.Name != updatedName {
		process.Fail("update course", fmt.Errorf("expected updated course name"))
	}
	process.Log(fmt.Sprintf("Updated course name=%q", updatedCourse.Name))

	process.EndStep()

	// [STEP 4] Get Course and validate data
	process.StartStep("Obtener datos del curso y validarlos")
	reloadedCourse, err := process.Application.CourseModule.GetCourseDetails.Execute(teacherCtx, course_dtos.GetCourseDetailsInput{CourseID: courseID})
	if err != nil {
		process.Fail("get course details", err)
	}
	if reloadedCourse == nil || reloadedCourse.ID != courseID {
		process.Fail("get course details", fmt.Errorf("expected course details for %s", courseID))
	}
	if reloadedCourse.Name != updatedName {
		process.Fail("get course details", fmt.Errorf("expected name %q, got %q", updatedName, reloadedCourse.Name))
	}
	process.Log(fmt.Sprintf("Course details validated for %s", courseID))

	process.EndStep()

	// [STEP 5] Delete Course
	process.StartStep("Eliminar curso")
	err = process.Application.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
	if err != nil {
		process.Fail("delete course", err)
	}
	deletedCourseID = courseID
	courseID = ""

	process.EndStep()

	// [STEP 6] Get Course and validate error
	process.StartStep("Obtener datos del curso y validar error")
	_, err = process.Application.CourseModule.GetCourseDetails.Execute(teacherCtx, course_dtos.GetCourseDetailsInput{CourseID: deletedCourseID})
	if err == nil {
		process.Fail("get course details after delete", fmt.Errorf("expected error after deleting course"))
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()
	process.End()
}
