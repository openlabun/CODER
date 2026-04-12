package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
)

type GetChallengeSubmissionsUseCase struct {
	userRepository userRepository.UserRepository
	challengeRepository examRepository.ChallengeRepository
	examRepository examRepository.ExamRepository
	submissionRepository submissionRepository.SubmissionRepository
	resultsRepository submissionRepository.SubmissionResultRepository
}

func NewGetChallengeSubmissionsUseCase(userRepository userRepository.UserRepository, challengeRepository examRepository.ChallengeRepository, examRepository examRepository.ExamRepository, submissionRepository submissionRepository.SubmissionRepository, resultsRepository submissionRepository.SubmissionResultRepository) *GetChallengeSubmissionsUseCase {
	return &GetChallengeSubmissionsUseCase{
		userRepository: userRepository,
		challengeRepository: challengeRepository,
		examRepository: examRepository,
		submissionRepository: submissionRepository,
		resultsRepository: resultsRepository,
	}
}

func (uc *GetChallengeSubmissionsUseCase) Execute(ctx context.Context, input dtos.GetChallengeSubmissionsInput) ([]*dtos.SubmissionOutputDTO, error) {
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

	// [STEP 2] Verify that challenge exists
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, fmt.Errorf("get challenge from repository: %w", err)
	}

	if challenge == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	// [STEP 3] If user is a student, check if the challenge is published
	if role == user_constants.UserRoleStudent {
		if challenge.Status != constants.ChallengeStatusPublished {
			return nil, fmt.Errorf("challenge with id %q is not published yet or it was archived", input.ChallengeID)
		}
	}

	// [STEP 4] If user is student only retrieve his own submissions
	var submissions []*Entities.Submission
	if role == user_constants.UserRoleStudent {
		submissions, err = uc.submissionRepository.GetSubmissionsByUserID(ctx, user.ID, input.Status, input.TestID, &input.ChallengeID)
		if err != nil {
			return nil, fmt.Errorf("get submissions by user: %w", err)
		}

		return uc.createSubmissionsOutputDTO(ctx, submissions)
	}

	// [STEP 5] If user is a teacher, query for all submissions for the challenge (only if he is owner)
	if role == user_constants.UserRoleProfessor {
		if challenge.UserID != user.ID {
			return nil, fmt.Errorf("user with email %q is not the owner of the challenge with id %q", userEmail, input.ChallengeID)
		}
	}

	// [STEP 6] Get all submissions for the challenge
	submissions, err = uc.submissionRepository.GetSubmissionsByChallengeID(ctx, input.ChallengeID, input.Status, input.TestID)
		
	if err != nil {
		return nil, fmt.Errorf("get submissions by challenge: %w", err)
	}

	return uc.createSubmissionsOutputDTO(ctx, submissions)
}

func (uc *GetChallengeSubmissionsUseCase) filterUserSubmissions (userID string, submissions []*Entities.Submission) ([]*Entities.Submission, error) {
	var userSubmissions []*Entities.Submission
	for _, submission := range submissions {
		if submission.UserID == userID {
			userSubmissions = append(userSubmissions, submission)
		}
	}

	return userSubmissions, nil
}

func (uc *GetChallengeSubmissionsUseCase) getSubmissionResults (ctx context.Context, submission *Entities.Submission) ([]Entities.SubmissionResult, error) {
	results, err := uc.resultsRepository.GetResultsBySubmissionID(ctx, submission.ID)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	derefResults := make([]Entities.SubmissionResult, len(results))
	for i, result := range results {
		derefResults[i] = *result
	}

	return derefResults, nil
}

func (uc *GetChallengeSubmissionsUseCase) createSubmissionsOutputDTO (ctx context.Context, submission []*Entities.Submission) ([]*dtos.SubmissionOutputDTO, error) {
	var subdtos []*dtos.SubmissionOutputDTO
	for _, submission := range submission {
		results, err := uc.getSubmissionResults(ctx, submission)
		if err != nil {
			return nil, fmt.Errorf("get submission results: %w", err)
		}

		dto := mapper.MapSubmissionOutputDTO(submission, results)
		if dto == nil {
			return nil, fmt.Errorf("failed to map submission output DTO")
		}

		subdtos = append(subdtos, dto)
	}

	if len(subdtos) == 0 {
		return []*dtos.SubmissionOutputDTO{}, nil
	}

	return subdtos, nil
}