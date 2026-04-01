package exam_usecases

import (
	"context"
	"fmt"

	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type DeleteTestCaseUseCase struct {
	userRepository       userRepository.UserRepository
	examRepository       examRepository.ExamRepository
	challengeRepository  examRepository.ChallengeRepository
	testCaseRepository   examRepository.TestCaseRepository
	ioVariableRepository examRepository.IOVariableRepository
}

func NewDeleteTestCaseUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, ioVariableRepository examRepository.IOVariableRepository) *DeleteTestCaseUseCase {
	return &DeleteTestCaseUseCase{
		userRepository:       userRepository,
		examRepository:       examRepository,
		challengeRepository:  challengeRepository,
		testCaseRepository:   testCaseRepository,
		ioVariableRepository: ioVariableRepository,
	}
}

func (uc *DeleteTestCaseUseCase) Execute(ctx context.Context, input dtos.DeleteTestCaseInput) error {
	// [STEP 1] Verify user is teacher and has permissions to create an exam
	userEmail, err := services.UserEmailFromContext(ctx)
	if err != nil {
		return err
	}

	user, err := uc.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	if user == nil {
		return fmt.Errorf("user with email %q does not exist", userEmail)
	}

	if user.Role != user_entities.UserRoleProfessor {
		return fmt.Errorf("user does not have permissions to create an exam")
	}

	// [STEP 2] Validate that test case exists
	test_case, err := uc.testCaseRepository.GetTestCaseByID(ctx, input.TestCaseID)
	if err != nil {
		return fmt.Errorf("test case with id %q does not exist", input.TestCaseID)
	}

	if test_case == nil {
		return fmt.Errorf("test case with id %q does not exist", input.TestCaseID)
	}

	// [STEP 3] Validate that challenge belongs to the teacher
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, test_case.ChallengeID)
	if err != nil {
		return fmt.Errorf("challenge with id %q does not exist", test_case.ChallengeID)
	}

	if challenge == nil {
		return fmt.Errorf("challenge with id %q does not exist", test_case.ChallengeID)
	}

	if challenge.UserID != user.ID {
		return fmt.Errorf("user does not have permissions to delete a test case for this challenge")
	}

	// [STEP 4] Delete test case entity
	if err := domain_services.RemoveTestCase(ctx, input.TestCaseID, uc.testCaseRepository, uc.ioVariableRepository); err != nil {
		return fmt.Errorf("failed to delete test case with id %q: %v", input.TestCaseID, err)
	}

	return nil
}
