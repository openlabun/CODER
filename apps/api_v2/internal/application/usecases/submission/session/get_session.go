package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
)

type GetSessionUseCase struct {
	sessionRepository submissionRepository.SessionRepository
	userRepository userRepository.UserRepository
}

func NewGetSessionUseCase(sessionRepository submissionRepository.SessionRepository, userRepository userRepository.UserRepository) *GetSessionUseCase {
	return &GetSessionUseCase{
		sessionRepository: sessionRepository,
		userRepository: userRepository,
	}
}

func (uc *GetSessionUseCase) Execute(ctx context.Context, input dtos.GetSessionInput) (*Entity.Session, error) {
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

	if user.Role != user_entities.UserRoleProfessor {
		return nil, fmt.Errorf("user does not have permissions to create an exam")
	}
	
	// [STEP 2] Verify existing student session
	sessions, err := uc.sessionRepository.GetSessionsByStudentID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	active_session := getExistingSession(sessions)

	// [STEP 3] If there is an active session, return it. Else, throw error
	if active_session == nil {
		return nil, fmt.Errorf("no active session found for student")
	}

	return active_session, nil
}