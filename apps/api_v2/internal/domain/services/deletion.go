package services

import (
	"context"
	
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	courseRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
)

func RemoveCourse (ctx context.Context, 
	courseID string,
	courseRepository courseRepository.CourseRepository, 
	examRepository examRepository.ExamRepository,
	examItemRepository examRepository.ExamItemRepository,
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
				err = RemoveExam(ctx, exam.ID, examRepository, examItemRepository)
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

func RemoveExam (ctx context.Context,
	examID string,
	examRepository examRepository.ExamRepository,
	examItemRepository examRepository.ExamItemRepository,
	) error {
		// [STEP 1] Get all exam items for the exam
		examItems, err := examItemRepository.GetExamItem(ctx, &examID, nil)
		if err != nil {
			return err
		}

		// [STEP 2] Delete all existing exam items for the exam
		for _, item := range examItems {
			if item != nil {
				err = examItemRepository.DeleteExamItem(ctx, item.ID)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 3] Delete the exam itself
		err = examRepository.DeleteExam(ctx, examID)
		if err != nil {
			return err
		}

		return nil
	}

func RemoveChallenge (ctx context.Context,
	challengeID string,
	challengeRepository examRepository.ChallengeRepository,
	testCaseRepository examRepository.TestCaseRepository,
	examItemRepository examRepository.ExamItemRepository,
	submissionRepository submissionRepository.SubmissionRepository,
	resultsRepository submissionRepository.SubmissionResultRepository,
	) error {
		// [STEP 1] Get all exam items for the challenge
		examItems, err := examItemRepository.GetExamItem(ctx, nil, &challengeID)
		if err != nil {
			return err
		}

		// [STEP 2] Delete all existing exam items for the challenge
		for _, item := range examItems {
			if item != nil {
				err = examItemRepository.DeleteExamItem(ctx, item.ID)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 3] Get all test cases for the challenge
		test_cases , err := testCaseRepository.GetTestCasesByChallengeID(ctx, challengeID)
		if err != nil {
			return err
		}

		// [STEP 4] Delete all existing test cases for the challenge
		for _, test_case := range test_cases {
			if test_case != nil {
				err = testCaseRepository.DeleteTestCase(ctx, test_case.ID)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 5] Get all existing submissions for the challenge
		submissions, err := submissionRepository.GetSubmissionsByChallengeID(ctx, challengeID, nil, nil)
		if err != nil {
			return err
		}

		// [STEP 6] Delete all existing submissions for the challenge
		for _, submission := range submissions {
			if submission != nil {
				err = RemoveSubmission(ctx, submission.ID, submissionRepository, resultsRepository)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 8] Delete the challenge itself
		err = challengeRepository.DeleteChallenge(ctx, challengeID)
		if err != nil {
			return err
		}

		return nil
	}


func RemoveSubmission (ctx context.Context,
	submissionID string,
	submissionRepository submissionRepository.SubmissionRepository,
	resultsRepository submissionRepository.SubmissionResultRepository,
	) error {
		// [STEP 1] Get all Submission results for the submission
		results, err := resultsRepository.GetResultsBySubmissionID(ctx, submissionID)
		if err != nil {
			return err
		}

		// [STEP 2] Delete all existing Submission results for the submission
		for _, result := range results {
			if result != nil {
				err = resultsRepository.DeleteResult(ctx, result.ID)
				if err != nil {
					return err
				}
			}
		}

		// [STEP 3] Delete the submission itself
		err = submissionRepository.DeleteSubmission(ctx, submissionID)
		if err != nil {
			return err
		}

		return nil
	}