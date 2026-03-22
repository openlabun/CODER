package exam_usecases

import (
	"context"
	"fmt"
	
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
)

type CreateTestCaseUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
	challengeRepository examRepository.ChallengeRepository
	testCaseRepository examRepository.TestCaseRepository
}

func NewCreateTestCaseUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository) *CreateTestCaseUseCase {
	return &CreateTestCaseUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
		challengeRepository: challengeRepository,
		testCaseRepository: testCaseRepository,
	}
}

func (uc *CreateTestCaseUseCase) Execute(ctx context.Context, input dtos.CreateTestCaseInput) (*Entities.TestCase, error) {
	// [STEP 1] Verify user is teacher and has permissions to create an exam
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

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create an exam")
	}
	
	// [STEP 2] Validate that challenge exists
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	if challenge == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	// [STEP 3] Validate that challenge belongs to an exam owned by the teacher
	exam, err := uc.examRepository.GetExamByID(ctx, challenge.ExamID)
	if err != nil {
		return nil, fmt.Errorf("error fetching exam with id %q: %v", challenge.ExamID, err)
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", challenge.ExamID)
	}

	if exam.ProfessorID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to create a test case for this challenge")
	}

	// [STEP 4] Create test case entity with user provided values
	testCase, err := mapper.MapCreateTestCaseInputToTestCaseEntity(input)
	if err != nil {
		return nil, err
	}

	// [STEP 5] Store test case in database
	testCase, err = uc.testCaseRepository.CreateTestCase(ctx, testCase)
	if err != nil {
		return nil, err
	}

	return testCase, nil
}
