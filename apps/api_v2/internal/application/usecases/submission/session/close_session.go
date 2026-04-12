package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
)

type CloseSessionUseCase struct {
	sessionRepository submissionRepository.SessionRepository
	userRepository userRepository.UserRepository
}

func NewCloseSessionUseCase(sessionRepository submissionRepository.SessionRepository, userRepository userRepository.UserRepository) *CloseSessionUseCase {
	return &CloseSessionUseCase{
		sessionRepository: sessionRepository,
		userRepository: userRepository,
	}
}

func (uc *CloseSessionUseCase) Execute(ctx context.Context, input dtos.CloseSessionInput) (*Entity.Session, error) {
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
	
	// [STEP 2] Verify existing student session
	session, err := uc.sessionRepository.GetSessionByID(ctx, input.SessionID)
	if err != nil {
		return nil, err
	}

	// [STEP 3] Verify if student is owner of the session
	if user.Role == user_entities.UserRoleStudent && session.StudentID != user.ID {
		return nil, fmt.Errorf("user is not the owner of the session")
	}

	// [STEP 4] Close Session
	err = state_machine.ApplyTranstion(session, constants.SessionStatusCompleted)
	if err != nil {
		return nil, fmt.Errorf("session is not active and cannot be closed, got error: %w", err)
	}

	// [STEP 5] Update session in repository
	session, err = uc.sessionRepository.UpdateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return session, nil
}