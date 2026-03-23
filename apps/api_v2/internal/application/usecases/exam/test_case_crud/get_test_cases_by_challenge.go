package exam_usecases

import (
	"context"
	"fmt"
	
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	courseRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	exam_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type GetTestCasesByChallengeUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
	challengeRepository examRepository.ChallengeRepository
	testCaseRepository examRepository.TestCaseRepository
	courseRepository courseRepository.CourseRepository
}

func NewGetTestCasesByChallengeUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, courseRepository courseRepository.CourseRepository) *GetTestCasesByChallengeUseCase {
	return &GetTestCasesByChallengeUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
		challengeRepository: challengeRepository,
		testCaseRepository: testCaseRepository,
		courseRepository: courseRepository,
	}
}

func (uc *GetTestCasesByChallengeUseCase) Execute(ctx context.Context, input dtos.GetTestCasesByChallengeInput) ([]*Entities.TestCase, error) {
	// [STEP 1] Verify user and get role
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user with email %q does not exist", userEmail)
	}

	role := user.Role
	
	// [STEP 2] Validate that challenge exists
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	if challenge == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	// [STEP 3] Validate that exam exists
	exam, err := uc.examRepository.GetExamByID(ctx, challenge.ExamID)
	if err != nil {
		return nil, fmt.Errorf("exam with id %q does not exist", challenge.ExamID)
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", challenge.ExamID)
	}

	// [STEP 4] If user is professor, validate that exam belongs to the teacher or is public/teacher
	if role == user_entities.UserRoleProfessor && exam.ProfessorID != user.ID {
		if exam.Visibility == exam_entities.VisibilityPublic || exam.Visibility == exam_entities.VisibilityTeachers {
			return nil, fmt.Errorf("user is not the owner of the exam with id %q", challenge.ExamID)
		}
	}

	// [STEP 5] If user is student, validate that exam is public or student is in course
	if role == user_entities.UserRoleStudent {
		if exam.Visibility == exam_entities.VisibilityCourse {
			courses, err := uc.courseRepository.GetCoursesByStudentID(ctx, user.ID)
			if err != nil {
				return nil, fmt.Errorf("error retrieving courses for student with id %q: %v", user.ID, err)
			}
			if !studentInCourse(exam.CourseID, courses) {
				return nil, fmt.Errorf("user does not have permissions to view test cases for challenge with id %q", input.ChallengeID)
			}
		}
	}

	// [STEP 5] Retrieve test cases for the specified challenge
	testCases, err := uc.testCaseRepository.GetTestCasesByChallengeID(ctx, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving test cases for challenge with id %q: %v", input.ChallengeID, err)
	}

	// [STEP 6] If user is student, filter test cases to return only public ones
	if role == user_entities.UserRoleStudent {
		testCases = filterPublicTestCases(testCases)
	}

	return testCases, nil
}

func filterPublicTestCases(testCases []*Entities.TestCase) []*Entities.TestCase {
	publicTestCases := []*Entities.TestCase{}
	for _, testCase := range testCases {
		if testCase.IsSample {
			publicTestCases = append(publicTestCases, testCase)
		}
	}

	return publicTestCases
}



func studentInCourse (courseID string, courses []*course_entities.Course) bool {
    if courses == nil {
		return false
	}

	for _, course := range courses {
		if course == nil {
			continue
		}

		if course.ID == courseID {
			return true
		}
	}
	return false
}
