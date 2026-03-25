package exam_usecases

import (
	"context"
	"fmt"
	
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examItemRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
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
	examItemRepository examItemRepository.ExamItemRepository
}

func NewGetTestCasesByChallengeUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, courseRepository courseRepository.CourseRepository, examItemRepository examItemRepository.ExamItemRepository) *GetTestCasesByChallengeUseCase {
	return &GetTestCasesByChallengeUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
		challengeRepository: challengeRepository,
		testCaseRepository: testCaseRepository,
		courseRepository: courseRepository,
		examItemRepository: examItemRepository,
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

	// [STEP 4] If user is professor, validate that exam belongs to the teacher or is public/teacher
	if role == user_entities.UserRoleProfessor && challenge.UserID != user.ID {
		if challenge.Status != exam_entities.ChallengeStatusPublished {
			return nil, fmt.Errorf("user is not the owner of the challenge with id %q", challenge.ID)
		}
	}

	// [STEP 5] If user is student, make validations to check if he has access
	if role == user_entities.UserRoleStudent {
		// [STEP 5.1] If challenge is not published, student cannot access
		if challenge.Status != exam_entities.ChallengeStatusPublished {
			return nil, fmt.Errorf("challenge with id %q is not published yet", challenge.ID)
		}

		// [STEP 5.2] If challenge is published, get challenge exams
		if input.ExamID == nil {
			return nil, fmt.Errorf("exam_id is required for students to access test cases")
		}

		examItems, err := uc.examItemRepository.GetExamItem(ctx, input.ExamID, &challenge.ID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving exam items for exam with id %q and challenge with id %q: %v", *input.ExamID, challenge.ID, err)
		}

		if len(examItems) == 0 {
			return nil, fmt.Errorf("no exam items found for exam with id %q and challenge with id %q", *input.ExamID, challenge.ID)
		}
		examItem := examItems[0]

		exam , err := uc.examRepository.GetExamByID(ctx, examItem.ExamID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving exam with id %q: %v", examItem.ExamID, err)
		}

		if exam == nil {
			return nil, fmt.Errorf("exam with id %q does not exist", examItem.ExamID)
		}

		// [STEP 5.3] Validate if exam is public for course and student is in the course associated with the exam
		if exam.Visibility == exam_entities.VisibilityCourse {
			courses, err := uc.courseRepository.GetCoursesByStudentID(ctx, user.ID)
			if err != nil {
				return nil, fmt.Errorf("error retrieving courses for student with id %q: %v", user.ID, err)
			}

			if !studentInCourse(exam.CourseID, courses) {
				return nil, fmt.Errorf("user does not have permissions to view test cases for challenge with id %q", input.ChallengeID)
			}
		// [STEP 5.4] Validate if exam is not public
		} else if exam.Visibility != exam_entities.VisibilityPublic {
			return nil, fmt.Errorf("exam with id %q is not public", exam.ID)
		}
	}

	// [STEP 6] Retrieve test cases for the specified challenge
	testCases, err := uc.testCaseRepository.GetTestCasesByChallengeID(ctx, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving test cases for challenge with id %q: %v", input.ChallengeID, err)
	}

	// [STEP 7] If user is student, filter test cases to return only public ones
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
