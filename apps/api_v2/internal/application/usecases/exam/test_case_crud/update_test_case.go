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

type UpdateTestCaseUseCase struct {
	userRepository       userRepository.UserRepository
	examRepository       examRepository.ExamRepository
	challengeRepository  examRepository.ChallengeRepository
	testCaseRepository   examRepository.TestCaseRepository
	ioVariableRepository examRepository.IOVariableRepository
}

func NewUpdateTestCaseUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, ioVariableRepository examRepository.IOVariableRepository) *UpdateTestCaseUseCase {
	return &UpdateTestCaseUseCase{
		userRepository:       userRepository,
		examRepository:       examRepository,
		challengeRepository:  challengeRepository,
		testCaseRepository:   testCaseRepository,
		ioVariableRepository: ioVariableRepository,
	}
}

func (uc *UpdateTestCaseUseCase) Execute(ctx context.Context, input dtos.UpdateTestCaseInput) (*Entities.TestCase, error) {
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

	// [STEP 2] Verify that test case exists
	testCase, err := uc.testCaseRepository.GetTestCaseByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("test case with id %q does not exist", input.ID)
	}

	if testCase == nil {
		return nil, fmt.Errorf("test case with id %q does not exist", input.ID)
	}

	// [STEP 3] Validate that challenge exists
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, testCase.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", testCase.ChallengeID)
	}

	if challenge == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", testCase.ChallengeID)
	}

	// [STEP 4] Validate that challenge belongs to the teacher
	if challenge.UserID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to update a test case for this challenge")
	}

	// [STEP 5] Create update test case entity with user provided values and existing test case
	updatedTestCase, err := mapper.MapUpdateTestCaseInputToTestCaseEntity(testCase, input)
	if err != nil {
		return nil, err
	}

	// [STEP 6] Save test case changes in database
	updatedTestCase, err = domain_services.UpdateTestCase(ctx, updatedTestCase, uc.testCaseRepository, uc.ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return updatedTestCase, nil
}
