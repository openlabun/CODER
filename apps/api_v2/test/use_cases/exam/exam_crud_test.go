package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	course_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"
	exam_dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

	test "github.com/openlabun/CODER/apps/api_v2/test"
	utils "github.com/openlabun/CODER/apps/api_v2/test/use_cases"
)

func TestExamCRUD(t *testing.T) {
	process := test.StartTestWithApp(t, "Exam CRUD Use Cases")
	email := "test@test.com"
	password := "Password123!"

	var teacherID string
	var teacherCtx = context.Background()
	var courseID string
	var examID string
	var deletedExamID string

	defer func() {
		if examID != "" {
			t.Logf("[CLEANUP] Deleting created exam %s", examID)
			_, _ = process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
		}
		if courseID != "" {
			t.Logf("[CLEANUP] Deleting created course %s", courseID)
			_ = process.Application.CourseModule.DeleteCourse.Execute(teacherCtx, course_dtos.DeleteCourseInput{CourseID: courseID})
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
	enrollmentCode := fmt.Sprintf("ENR-UC-EXAM-%d", now.UnixNano())
	createdCourse, err := process.Application.CourseModule.CreateCourse.Execute(teacherCtx, course_dtos.CreateCourseInput{
		Name:           "Exam CRUD Course",
		Description:    "Course for exam CRUD test",
		Visibility:     string(course_entities.CourseVisibilityPublic),
		VisualIdentity: string(course_entities.CourseColourBlue),
		Code:           fmt.Sprintf("EXCRUD-%d", now.Unix()%100000),
		Year:           now.Year(),
		Semester:       string(course_entities.AcademicFirstPeriod),
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

	// [STEP 3] Create Exam
	process.StartStep("Crea un Examen")
	startTime := now.Add(2 * time.Hour)
	endTime := startTime.Add(90 * time.Minute)
	endTimeStr := endTime.Format(time.RFC3339)

	createdExam, err := process.Application.ExamModule.CreateExam.Execute(teacherCtx, exam_dtos.CreateExamInput{
		CourseID:             &courseID,
		Title:                "Exam CRUD Test",
		Description:          "Created by use case test",
		Visibility:           string(exam_entities.VisibilityCourse),
		StartTime:            startTime.Format(time.RFC3339),
		EndTime:              &endTimeStr,
		AllowLateSubmissions: false,
		TimeLimit:            5400,
		TryLimit:             2,
		ProfessorID:          teacherID,
	})
	if err != nil {
		process.Fail("create exam", err)
	}
	if createdExam == nil || createdExam.ID == "" {
		process.Fail("create exam", fmt.Errorf("expected created exam with ID"))
	}
	examID = createdExam.ID
	process.Log(fmt.Sprintf("examID=%s", examID))

	process.EndStep()

	// [STEP 4] Update Exam
	process.StartStep("Actualizar Examen")
	updatedTitle := "Exam CRUD Test Updated"
	updatedDescription := "Updated by use case test"
	updatedTryLimit := 3
	updatedExam, err := process.Application.ExamModule.UpdateExam.Execute(teacherCtx, exam_dtos.UpdateExamInput{
		ExamID:      examID,
		Title:       &updatedTitle,
		Description: &updatedDescription,
		TryLimit:    &updatedTryLimit,
	})
	if err != nil {
		process.Fail("update exam", err)
	}
	if updatedExam == nil || updatedExam.Title != updatedTitle {
		process.Fail("update exam", fmt.Errorf("expected updated exam title"))
	}
	process.Log(fmt.Sprintf("Updated exam title=%q", updatedExam.Title))

	process.EndStep()

	// [STEP 5] Get Exam and validate data
	process.StartStep("Obtener datos del Examen y validarlos")
	reloadedExam, err := process.Application.ExamModule.GetExamDetails.Execute(teacherCtx, exam_dtos.GetExamDetailsInput{ExamID: examID})
	if err != nil {
		process.Fail("get exam details", err)
	}
	if reloadedExam == nil || reloadedExam.ID != examID {
		process.Fail("get exam details", fmt.Errorf("expected exam details for %s", examID))
	}
	if reloadedExam.Title != updatedTitle {
		process.Fail("get exam details", fmt.Errorf("expected title %q, got %q", updatedTitle, reloadedExam.Title))
	}
	process.Log(fmt.Sprintf("Exam details validated for %s", examID))

	process.EndStep()

	// [STEP 6] Update wrong values in Exam (yesterday EndTime) - expects error
	process.StartStep("Actualizar Examen con valores incorrectos (EndTime con valor de ayer)")
	yesterday := now.Add(-24 * time.Hour).Format(time.RFC3339)
	_, err = process.Application.ExamModule.UpdateExam.Execute(teacherCtx, exam_dtos.UpdateExamInput{
		ExamID:  examID,
		EndTime: &yesterday,
	})
	if err == nil {
		process.Fail("invalid exam update", fmt.Errorf("expected error when updating exam with end time before start time"))
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()

	// [STEP 7] Delete Exam
	process.StartStep("Eliminar Examen")
	deletedExam, err := process.Application.ExamModule.DeleteExam.Execute(teacherCtx, exam_dtos.DeleteExamInput{ExamID: examID})
	if err != nil {
		process.Fail("delete exam", err)
	}
	if deletedExam == nil || deletedExam.ID == "" {
		process.Fail("delete exam", fmt.Errorf("expected deleted exam payload"))
	}
	deletedExamID = deletedExam.ID
	examID = ""

	process.EndStep()

	// [STEP 8] Get Exam and validate error
	process.StartStep("Obtener datos del Examen y validar error")
	_, err = process.Application.ExamModule.GetExamDetails.Execute(teacherCtx, exam_dtos.GetExamDetailsInput{ExamID: deletedExamID})
	if err == nil {
		process.Fail("get exam details after delete", fmt.Errorf("expected error after deleting exam"))
	}

	process.Log("Recibió ERROR, como se esperaba")
	process.EndStep()
	process.End()
}
