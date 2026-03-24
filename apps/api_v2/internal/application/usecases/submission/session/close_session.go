package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type CloseSessionUseCase struct {
	userRepository userRepository.UserRepository
	sessionRepository submissionRepository.SessionRepository
	examRepository examRepository.ExamRepository
}

func NewCloseSessionUseCase(userRepository userRepository.UserRepository, sessionRepository submissionRepository.SessionRepository, examRepository examRepository.ExamRepository) *CloseSessionUseCase {
	return &CloseSessionUseCase{
		userRepository: userRepository,
		sessionRepository: sessionRepository,
		examRepository: examRepository,
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
	sessions, err := uc.sessionRepository.GetSessionsByStudentID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	active_sesion := getExistingSession(sessions)

	// [STEP 3] If session is not active, return error
	if active_sesion == nil {
		return nil, fmt.Errorf("student does not have an active session")
	}

	// [STEP 4] Close session
	err = state_machine.ApplyTranstion(active_sesion, Entity.SessionStatusCompleted)
	if err != nil {
		return nil, err
	}

	// [STEP 5] Persist session changes
	session, err := uc.sessionRepository.UpdateSession(ctx, active_sesion)
	if err != nil {
		return nil, err
	}

	return session, nil
}
