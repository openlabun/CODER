package exam_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type CreateExamItemUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
	testCaseRepository examRepository.TestCaseRepository
	examItemRepository examRepository.ExamItemRepository
	challengeRepository examRepository.ChallengeRepository
	ioVariableRepository examRepository.IOVariableRepository
}

func NewCreateExamItemUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, examItemRepository examRepository.ExamItemRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, ioVariableRepository examRepository.IOVariableRepository) *CreateExamItemUseCase {
	return &CreateExamItemUseCase{
		userRepository: userRepository,
		examRepository: examRepository,
		testCaseRepository: testCaseRepository,
		examItemRepository: examItemRepository,
		challengeRepository: challengeRepository,
		ioVariableRepository: ioVariableRepository,
	}
}

func (uc *CreateExamItemUseCase) Execute(ctx context.Context, input dtos.CreateExamItemInput) (*Entities.ExamItem, error) {
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

	// [STEP 2] Validate Exam exists
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam: %v", err)
	}
	if exam == nil {
		return nil, fmt.Errorf("exam with ID %q does not exist", input.ExamID)
	}

	// [STEP 3] Validate user owns the Exam
	if exam.ProfessorID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to modify this exam")
	}

	// [STEP 4] Validate Challenge exists
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %v", err)
	}
	if challenge == nil {
		return nil, fmt.Errorf("challenge with ID %q does not exist", input.ChallengeID)
	}

	// [STEP 5] If challenge is draft or its archived, throw error
	if challenge.Status == constants.ChallengeStatusDraft || challenge.Status == constants.ChallengeStatusArchived {
		return nil, fmt.Errorf("challenge with ID %q is not published and cannot be added to the exam", input.ChallengeID)
	}

	// [STEP 6] Validate Challenge is not already in the Exam
	existingItems, err := uc.examItemRepository.GetExamItem(ctx, &input.ExamID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing exam items: %v", err)
	}

	for _, item := range existingItems {
		if item.ChallengeID == input.ChallengeID {
			return nil, fmt.Errorf("challenge with ID %q is already in the exam", input.ChallengeID)
		}
	}

	// [STEP 7] Validate if challenge belongs to user, or if challenge is public
	if challenge.UserID != user.ID && challenge.Status != constants.ChallengeStatusPublished {
		return nil, fmt.Errorf("challenge with ID %q is not owned by the user and is not public", input.ChallengeID)
	}

	// [STEP 8] If challenge is not owned by the teacher, fork it and use the forked challenge in the exam item
	if challenge.UserID != user.ID {
		challenge, err = services.ForkChallenge(ctx, *challenge, user.ID, uc.challengeRepository, uc.testCaseRepository, uc.ioVariableRepository)
		if err != nil {
			return nil, fmt.Errorf("failed to fork challenge: %v", err)
		}
		
		input.ChallengeID = challenge.ID
	}

	// [STEP 9] Create ExamItem
	examItem, err := mapper.MapCreateExamItemInputToExamItemEntity(input)
	if err != nil {
		return nil, fmt.Errorf("failed to map exam item: %v", err)
	}

	// [STEP 10] Save ExamItem in database
	createdExamItem, err := uc.examItemRepository.CreateExamItem(ctx, examItem)
	if err != nil {
		return nil, fmt.Errorf("failed to create exam item: %v", err)
	}


	return createdExamItem, nil
}