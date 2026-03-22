package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"

	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type CreateSessionUseCase struct {
	userRepository userRepository.UserRepository
	sessionRepository submissionRepository.SessionRepository
	examRepository examRepository.ExamRepository
}

func NewCreateSessionUseCase(userRepository userRepository.UserRepository, sessionRepository submissionRepository.SessionRepository, examRepository examRepository.ExamRepository) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		userRepository: userRepository,
		sessionRepository: sessionRepository,
		examRepository: examRepository,
	}
}

func (uc *CreateSessionUseCase) Execute(ctx context.Context, input dtos.CreateSessionInput) (*Entity.Session, error) {
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

	// [STEP 3] If there is an active session, throw error
	if active_sesion != nil {
		return nil, fmt.Errorf("student already has an active session")
	}

	// [STEP 4] Retrieve exam
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", input.ExamID)
	}

	// [STEP 3] If no existing session, create new session for the student and return it
	sessionEntity, err := mapper.MapCreateSessionInputToSessionRecord(input, exam)
	if err != nil {
		return nil, err
	}

	// [STEP 4] Save in database and return session entity
	session, err := uc.sessionRepository.CreateSession(ctx, sessionEntity)
	if err != nil {
		return nil, err
	}
	
	return session, nil
}

func getExistingSession (sessions []*Entity.Session) (*Entity.Session) {
	for _, session := range sessions {
		if session == nil {
			continue
		}

		if session.Status == Entity.SessionStatusActive {
			return session
		}
	}

	return nil
}