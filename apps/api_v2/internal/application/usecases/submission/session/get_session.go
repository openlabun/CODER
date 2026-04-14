package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
)

type GetActiveSessionUseCase struct {
	sessionRepository submissionRepository.SessionRepository
	userRepository userRepository.UserRepository
	examRepository examRepository.ExamRepository
}

func NewGetActiveSessionUseCase(sessionRepository submissionRepository.SessionRepository, userRepository userRepository.UserRepository, examRepository examRepository.ExamRepository) *GetActiveSessionUseCase {
	return &GetActiveSessionUseCase{
		sessionRepository: sessionRepository,
		userRepository:    userRepository,
		examRepository:    examRepository,
	}
}

func (uc *GetActiveSessionUseCase) Execute(ctx context.Context, input dtos.GetActiveSessionInput) (*Entity.Session, error) {
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

	role := user.Role

	// [STEP 2] If user is teacher get user id from input
	studentID := user.ID
	if role == user_constants.UserRoleProfessor {
		if input.UserID == nil {
			return nil, fmt.Errorf("student_id is required for teachers to get a student session")
		}

		studentID = *input.UserID

		// [STEP 2.1] Get student by id and verify it exists
		student, err := uc.userRepository.GetUserByID(ctx, studentID)
		if err != nil {
			return nil, err
		}
		if student == nil {
			return nil, fmt.Errorf("student with id %q does not exist", studentID)
		}
	}

	// [STEP 3] Get all sessions for a student
	sessions, err := uc.sessionRepository.GetSessionsByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	session := getExistingSession(sessions)
	if session == nil {
		return nil, fmt.Errorf("no active session found for student")
	}

	// [STEP 4] Update session status and persist it
	exam, err := uc.examRepository.GetExamByID(ctx, session.ExamID)
	if err != nil {
		return nil, err
	}

	if err := state_machine.UpdateSessionStatus(session, exam, services.Now(), false); err != nil {
		return nil, fmt.Errorf("failed to update session status: %w", err)
	}

	session, err = uc.sessionRepository.UpdateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session status: %w", err)
	}

	// [STEP 5] Get active session (updated)
	var active_session *Entity.Session
	
	if session != nil {
		if session.Status == constants.SessionStatusActive || session.Status == constants.SessionStatusFrozen {
			active_session = session
		}
	}

	// [STEP 6] If there is an active session, return it. Else, throw error
	if active_session == nil {
		return nil, fmt.Errorf("no active session found for student")
	}

	return active_session, nil
}