package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"
	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
)

type CloseSessionUseCase struct {
	sessionRepository submissionRepository.SessionRepository
	examScoreRepository examRepository.ExamScoreRepository
	examItemRepository examRepository.ExamItemRepository
	examItemScoreRepository examRepository.ExamItemScoreRepository
	submissionRepository submissionRepository.SubmissionRepository
	userRepository userRepository.UserRepository
}

func NewCloseSessionUseCase(sessionRepository submissionRepository.SessionRepository, examScoreRepository examRepository.ExamScoreRepository, examItemRepository examRepository.ExamItemRepository, examItemScoreRepository examRepository.ExamItemScoreRepository, submissionRepository submissionRepository.SubmissionRepository, userRepository userRepository.UserRepository) *CloseSessionUseCase {
	return &CloseSessionUseCase{
		sessionRepository: sessionRepository,
		examScoreRepository: examScoreRepository,
		examItemRepository: examItemRepository,
		examItemScoreRepository: examItemScoreRepository,
		submissionRepository: submissionRepository,
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
	if user.Role == user_constants.UserRoleStudent && session.StudentID != user.ID {
		return nil, fmt.Errorf("user is not the owner of the session")
	}

	// [STEP 4] Get ExamScore associated with session
	examScore, err := uc.examScoreRepository.GetExamScoreBySessionID(ctx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve exam score for session: %w", err)
	}
	if examScore == nil {
		return nil, fmt.Errorf("no exam score found for session")
	}

	// [STEP 5] Calculate final score for the session
	_, err = domain_services.CalculateExamScore(ctx, examScore, uc.examScoreRepository, uc.examItemRepository, uc.examItemScoreRepository, uc.submissionRepository)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate exam score: %w", err)
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