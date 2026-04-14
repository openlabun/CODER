package services

import (
	"context"

	courseRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
)

func RemoveCourse (ctx context.Context, 
	courseID string,
	courseRepository courseRepository.CourseRepository, 
	examRepository examRepository.ExamRepository,
	examItemRepository examRepository.ExamItemRepository,
	examScoreRepository examRepository.ExamScoreRepository,
	examItemScoreRepository examRepository.ExamItemScoreRepository,
	) error {
		// [STEP 1] Get all enrolled students for the course
		enrolledStudents, err := courseRepository.GetStudentsByCourseID(ctx, courseID)
		if err != nil {
			return err
		}

		// [STEP 2] Delete all existing enrollments for the course
		for _, student := range enrolledStudents {
			if student != nil {
				err = courseRepository.RemoveStudentFromCourse(ctx, courseID, student.ID)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 3] Get all exams for the course
		exams, err := examRepository.GetExamsByCourseID(ctx, courseID)
		if err != nil {
			return err
		}

		// [STEP 4] Delete all existing exams for the course
		for _, exam := range exams {
			if exam != nil {
				err = RemoveExam(ctx, exam.ID, examRepository, examItemRepository, examScoreRepository, examItemScoreRepository)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 5] Delete the course itself
		err = courseRepository.DeleteCourse(ctx, courseID)
		if err != nil {
			return err
		}

		return nil
	}
