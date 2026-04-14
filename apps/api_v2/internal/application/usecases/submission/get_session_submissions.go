package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type GetSessionSubmissionsUseCase struct {
	userRepository userRepository.UserRepository
	challengeRepository examRepository.ChallengeRepository
	examRepository examRepository.ExamRepository
	submissionRepository submissionRepository.SubmissionRepository
	resultsRepository submissionRepository.SubmissionResultRepository
	sessionRepository submissionRepository.SessionRepository
}

func NewGetSessionSubmissionsUseCase(userRepository userRepository.UserRepository, challengeRepository examRepository.ChallengeRepository, examRepository examRepository.ExamRepository, submissionRepository submissionRepository.SubmissionRepository, resultsRepository submissionRepository.SubmissionResultRepository, sessionRepository submissionRepository.SessionRepository) *GetSessionSubmissionsUseCase {
	return &GetSessionSubmissionsUseCase{
		userRepository: userRepository,
		challengeRepository: challengeRepository,
		examRepository: examRepository,
		submissionRepository: submissionRepository,
		resultsRepository: resultsRepository,
		sessionRepository: sessionRepository,
	}
}

func (uc *GetSessionSubmissionsUseCase) Execute(ctx context.Context, input dtos.GetSessionSubmissionsInput) ([]*dtos.SubmissionOutputDTO, error) {
	// [STEP 1] Verify user and his role
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

	// [STEP 2] Verify that session exists
	session, err := uc.sessionRepository.GetSessionByID(ctx, input.SessionID)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, fmt.Errorf("session with id %q does not exist", input.SessionID)
	}

	// [STEP 3] If user is a student, verify that session belongs to him
	if role == user_constants.UserRoleStudent && session.StudentID != user.ID {
		return nil, fmt.Errorf("session with id %q does not belong to user with email %q", input.SessionID, userEmail)
	}

	// [STEP 4] Get submissions for the session (filtering by status, challengeID and testID if provided)
	submissions, err := uc.submissionRepository.GetSubmissionsBySessionID(ctx, input.SessionID, input.Status, input.TestID, input.ChallengeID)
		
	if err != nil {
		return nil, err
	} 

	// [STEP 5] If user is a teacher, query for all submissions for the session (only if he is the exam owner)
	exam, err := uc.examRepository.GetExamByID(ctx, session.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", session.ExamID)
	}

	if exam.ProfessorID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to view submissions for this session")
	}

	return uc.createSubmissionsOutputDTO(ctx, submissions)
}

func (uc *GetSessionSubmissionsUseCase) getSubmissionResults (ctx context.Context, submission *Entities.Submission) ([]Entities.SubmissionResult, error) {
	results, err := uc.resultsRepository.GetResultsBySubmissionID(ctx, submission.ID)
	if err != nil {
		return nil, err
	}

	derefResults := make([]Entities.SubmissionResult, len(results))
	for i, result := range results {
		derefResults[i] = *result
	}

	return derefResults, err
}

func (uc *GetSessionSubmissionsUseCase) createSubmissionsOutputDTO (ctx context.Context, submission []*Entities.Submission) ([]*dtos.SubmissionOutputDTO, error) {
	var dtos []*dtos.SubmissionOutputDTO
	for _, submission := range submission {
		results, err := uc.getSubmissionResults(ctx, submission)
		if err != nil {
			return nil, err
		}

		dto := mapper.MapSubmissionOutputDTO(submission, results)
		if dto == nil {
			return nil, fmt.Errorf("failed to map submission output DTO")
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}