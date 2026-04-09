package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	session_states "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
)

type HeartBeatSessionUseCase struct {
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
	sessionRepository submissionRepository.SessionRepository
}

func NewHeartBeatSessionUseCase(userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository, sessionRepository submissionRepository.SessionRepository) *HeartBeatSessionUseCase {
	return &HeartBeatSessionUseCase{
		userRepository:  userRepository,
		examRepository:  examRepository,
		sessionRepository: sessionRepository,
	}
}

func (uc *HeartBeatSessionUseCase) Execute(ctx context.Context, input dtos.HeartbeatSessionInput) (*Entity.Session, error) {
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

	// [STEP 3] If there is not an active session, throw error 
	active_session := getExistingSession(sessions)
	if active_session == nil {
		return nil, fmt.Errorf("no active session found for student")
	}

	// [STEP 4] Get exam for the session
	exam, err := uc.examRepository.GetExamByID(ctx, active_session.ExamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam for session: %w", err)
	}

	// [STEP 4] Update its last_heartbeat timestamp
	err = session_states.UpdateSessionStatus(active_session, exam, services.Now(), true)
	if err != nil {
		return nil, fmt.Errorf("failed to update session status: %w", err)
	}

	// [STEP 5] Save in database and return session entity
	session, err := uc.sessionRepository.UpdateSession(ctx, active_session)
	if err != nil {
		return nil, err
	}

	return session, nil
}