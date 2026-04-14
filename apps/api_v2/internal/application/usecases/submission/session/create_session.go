package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	exam_mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam/mapper"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	Entity "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"
)

type CreateSessionUseCase struct {
	userRepository userRepository.UserRepository
	sessionRepository submissionRepository.SessionRepository
	examRepository examRepository.ExamRepository
	examScoreRepository examRepository.ExamScoreRepository
	examItemRepository examRepository.ExamItemRepository
	examItemScoreRepository examRepository.ExamItemScoreRepository
}

func NewCreateSessionUseCase(userRepository userRepository.UserRepository, sessionRepository submissionRepository.SessionRepository, examRepository examRepository.ExamRepository, examScoreRepository examRepository.ExamScoreRepository, examItemRepository examRepository.ExamItemRepository, examItemScoreRepository examRepository.ExamItemScoreRepository) *CreateSessionUseCase {
	return &CreateSessionUseCase{
		userRepository: userRepository,
		sessionRepository: sessionRepository,
		examScoreRepository: examScoreRepository,
		examItemRepository: examItemRepository,
		examItemScoreRepository: examItemScoreRepository,
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

	session := getExistingSession(sessions)

	// [STEP 3] If session is frozen, close it, if its active keep it 
	var active_sesion *Entity.Session
	if session != nil && session.Status == constants.SessionStatusFrozen {
		if err := state_machine.ApplyTranstion(session, constants.SessionStatusExpired); err != nil {
			return nil, fmt.Errorf("failed to expire frozen session: %w", err)
		}
		
		_, err = uc.sessionRepository.UpdateSession(ctx, session)
		if err != nil {
			return nil, err
		}
	} else if session != nil && session.Status == constants.SessionStatusActive {
		active_sesion = session
	}

	// [STEP 4] If there is an active session, throw error
	if active_sesion != nil {
		return nil, fmt.Errorf("student already has an active session")
	}

	// [STEP 5] Retrieve exam
	exam, err := uc.examRepository.GetExamByID(ctx, input.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", input.ExamID)
	}

	// [STEP 6] If no existing session, create new session for the student and return it
	sessionEntity, err := mapper.MapCreateSessionInputToSessionRecord(input, exam)
	if err != nil {
		return nil, err
	}

	// [STEP 7] Save in database and return session entity
	session, err = uc.sessionRepository.CreateSession(ctx, sessionEntity)
	if err != nil {
		return nil, err
	}

	// [STEP 8] Create ExamScore for exam evaluation
	examScore, err := exam_mapper.MapExamScore(exam, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create exam score: %w", err)
	}

	examScore, err = uc.examScoreRepository.CreateExamScore(ctx, examScore)
	if err != nil {
		return nil, fmt.Errorf("failed to create exam score: %w", err)
	}

	// [STEP 9] Create ExamItemScores for exam evaluation
	examItems, err := uc.examItemRepository.GetExamItem(ctx, &exam.ID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get exam items: %w", err)
	}

	for _, examItem := range examItems {
		// [STEP 9.1] Create ExamItemScore for each exam item
		examItemScore, err := exam_mapper.MapExamItemScore(examItem, examScore)
		if err != nil {
			return nil, fmt.Errorf("failed to create exam item score: %w", err)
		}

		// [STEP 9.2] Save ExamItemScore in database
		examItemScore, err = uc.examItemScoreRepository.CreateExamItemScore(ctx, examItemScore)
		if err != nil {
			return nil, fmt.Errorf("failed to create exam item score: %w", err)
		}
	}
	
	return session, nil
}

func getExistingSession (sessions []*Entity.Session) (*Entity.Session) {
	for _, session := range sessions {
		if session == nil {
			continue
		}

		if session.Status == constants.SessionStatusActive || session.Status == constants.SessionStatusFrozen {
			return session
		}
	}

	return nil
}