package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"

	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
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
		return nil, err
	}

	if challenge == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	// [STEP 3] Get all submissions for the challenge
	submissions, err := uc.submissionRepository.GetSubmissionsByChallengeID(ctx, input.ChallengeID, input.Status, input.TestID)
		
	if err != nil {
		return nil, err
	}

	// [STEP 4] If user is a student, only query for his own submissions if challenge is published
	if role == user_entities.UserRoleStudent {

		if challenge.Status != examEntities.ChallengeStatusPublished {
			return nil, fmt.Errorf("challenge with id %q is not published yet or it was archived", input.ChallengeID)
		}
		
		userSubmissions, err := uc.filterUserSubmissions(user.ID, submissions)
		if err != nil {
			return nil, err
		}

		return uc.createSubmissionsOutputDTO(ctx, userSubmissions)
	} 

	// [STEP 5] If user is a teacher, query for all submissions for the challenge (only if he is owner)
	exam, err := uc.examRepository.GetExamByID(ctx, challenge.ExamID)
	if err != nil {
		return nil, err
	}

	if exam == nil {
		return nil, fmt.Errorf("exam with id %q does not exist", challenge.ExamID)
	}

	if exam.ProfessorID != user.ID {
		return nil, fmt.Errorf("user does not have permissions to view submissions for this challenge")
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
		return nil, err
	}

	derefResults := make([]Entities.SubmissionResult, len(results))
	for i, result := range results {
		derefResults[i] = *result
	}

	return derefResults, err
}

func (uc *GetChallengeSubmissionsUseCase) createSubmissionsOutputDTO (ctx context.Context, submission []*Entities.Submission) ([]*dtos.SubmissionOutputDTO, error) {
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