package exam_usecases

import (
	"context"
	"fmt"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type CreateTestCaseUseCase struct {
	userRepository       userRepository.UserRepository
	examRepository       examRepository.ExamRepository
	challengeRepository  examRepository.ChallengeRepository
	testCaseRepository   examRepository.TestCaseRepository
	ioVariableRepository examRepository.IOVariableRepository
}

func NewCreateTestCaseUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, ioVariableRepository examRepository.IOVariableRepository) *CreateTestCaseUseCase {
	return &CreateTestCaseUseCase{
		userRepository:       userRepository,
		examRepository:       examRepository,
		challengeRepository:  challengeRepository,
		testCaseRepository:   testCaseRepository,
		ioVariableRepository: ioVariableRepository,
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

	if user.Role != user_constants.UserRoleProfessor {
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

	// [STEP 3] Validate that challenge belongs to the teacher
	if challenge.UserID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to create a test case for this challenge")
	}

	// [STEP 4] Create test case entity with user provided values
	testCase, err := mapper.MapCreateTestCaseInputToTestCaseEntity(input)
	if err != nil {
		return nil, err
	}

	// [STEP 5] Store test case in database
	testCase, err = domain_services.CreateTestCase(ctx, testCase, uc.testCaseRepository, uc.ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return testCase, nil
}
