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

type GetActiveSessionUseCase struct {
	sessionRepository submissionRepository.SessionRepository
	userRepository userRepository.UserRepository
}

func NewGetActiveSessionUseCase(sessionRepository submissionRepository.SessionRepository, userRepository userRepository.UserRepository) *GetActiveSessionUseCase {
	return &GetActiveSessionUseCase{
		sessionRepository: sessionRepository,
		userRepository: userRepository,
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
	if role == user_entities.UserRoleProfessor {
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

	// [STEP 3] Verify existing student active session
	sessions, err := uc.sessionRepository.GetSessionsByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var active_session *Entity.Session
	session := getExistingSession(sessions)
	if session != nil && session.Status == Entity.SessionStatusActive {
		active_session = session
	}

	// [STEP 4] If there is an active session, return it. Else, throw error
	if active_session == nil {
		return nil, fmt.Errorf("no active session found for student")
	}

	return active_session, nil
}