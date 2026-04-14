package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
)

type BlockSessionUseCase struct {
	sessionRepository submissionRepository.SessionRepository
	userRepository userRepository.UserRepository
}

func NewBlockSessionUseCase(sessionRepository submissionRepository.SessionRepository, userRepository userRepository.UserRepository) *BlockSessionUseCase {
	return &BlockSessionUseCase{
		sessionRepository: sessionRepository,
		userRepository: userRepository,
	}
}

func (uc *BlockSessionUseCase) Execute(ctx context.Context, input dtos.BlockSessionInput) (*Entity.Session, error) {
	// [STEP 1] Verify user is student and has permissions to submit
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

	// [STEP 2] Verify if user is student
	if user.Role == user_constants.UserRoleStudent {
		return nil, fmt.Errorf("students are not allowed to block sessions")
	}
	
	// [STEP 3] Verify existing student session
	session, err := uc.sessionRepository.GetSessionByID(ctx, input.SessionID)
	if err != nil {
		return nil, err
	}

	// [STEP 4] Block Session
	err = state_machine.ApplyTranstion(session, constants.SessionStatusBlocked)
	if err != nil {
		return nil, fmt.Errorf("session is not active and cannot be blocked, got error: %w", err)
	}

	// [STEP 5] Update session in repository
	session, err = uc.sessionRepository.UpdateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}