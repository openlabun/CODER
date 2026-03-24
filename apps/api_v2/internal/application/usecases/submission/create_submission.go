package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

	submissionPorts "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/submission"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"
)

type CreateSubmissionUseCase struct {
	userRepository userRepository.UserRepository
	submissionRepository submissionRepository.SubmissionRepository
	sessionRepository submissionRepository.SessionRepository
	challengeRepository examRepository.ChallengeRepository
	testCaseRepository examRepository.TestCaseRepository
	resultRepository submissionRepository.SubmissionResultRepository
	publisherPort submissionPorts.SubmissionPublisherPort
}

func NewCreateSubmissionUseCase(userRepository userRepository.UserRepository, submissionRepository submissionRepository.SubmissionRepository, sessionRepository submissionRepository.SessionRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, resultRepository submissionRepository.SubmissionResultRepository, publisherPort submissionPorts.SubmissionPublisherPort) *CreateSubmissionUseCase {
	return &CreateSubmissionUseCase{
		userRepository: userRepository,
		submissionRepository: submissionRepository,
		sessionRepository: sessionRepository,
		challengeRepository: challengeRepository,
		testCaseRepository: testCaseRepository,
		resultRepository: resultRepository,
		publisherPort: publisherPort,
	}
}

func (uc *CreateSubmissionUseCase) Execute(ctx context.Context, input dtos.CreateSubmissionInput) (*Entities.Submission, error) {
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

	if user.Role != user_entities.UserRoleStudent {
		return nil, fmt.Errorf("only students can create submissions")
	}

	if input.Language != string(Entities.LanguagePython) {
		return nil, fmt.Errorf("unsupported language %q. currently only python is enabled", input.Language)
	}

	// [STEP 2] Session is optional for challenge-only flow.
	// If it is provided (exam flow), validate it exists.
	if input.SessionID != "" {
		session, err := uc.sessionRepository.GetSessionByID(ctx, input.SessionID)
		if err != nil {
			return nil, err
		}
		if session == nil {
			return nil, fmt.Errorf("no active session found for student %q", user.Username)
		}
	}

	// [STEP 3] Create submission with user provided values
	submission, err := mapper.MapCreateSubmissionInputToSubmissionEntity(user.ID, input)
	if err != nil {
		return nil, err
	}

	// [STEP 4] Save submission in database
	createdSubmission, err := uc.submissionRepository.CreateSubmission(ctx, submission)
	if err != nil {
		return nil, err
	}

	// [STEP 5] Get challenge
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}

	if challenge == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	// [STEP 6] Get test cases of the challenge
	testCases, err := uc.testCaseRepository.GetTestCasesByChallengeID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}

	if len(testCases) == 0 {
		return nil, fmt.Errorf("no test cases found for challenge with id %q", input.ChallengeID)
	}

	// [STEP 7] Create submission results for each test case of the challenge
	var publishedResults []dtos.SubmissionResultPublishedDTO
	for _, testCase := range testCases {
		result, err := uc.createSubmissionResultsForTestCases(ctx, createdSubmission.ID, *testCase)
		if err != nil {
			return nil, err
		}

		if result != nil {
			// [STEP 8] Create DTO for publishing
			publishedResult := mapper.MapSubmissionResultToPublishedDTO(*createdSubmission, *result, *testCase, *challenge)
			if publishedResult != nil {
				publishedResults = append(publishedResults, *publishedResult)
			}
		}
	}
	
	// [STEP 8] Publish submission created event to message broker for asynchronous processing of submission results
	for _, publishedResult := range publishedResults {
		err = uc.publisherPort.PublishSubmission(publishedResult)
		if err != nil {
			return nil, fmt.Errorf("failed to publish submission result: %w", err)
		}
	}

	// [STEP 8] Return created submission entity
	return createdSubmission, nil
}

func (uc *CreateSubmissionUseCase) createSubmissionResultsForTestCases(ctx context.Context, submissionID string, testCase examEntities.TestCase) (*Entities.SubmissionResult, error) {
	// [STEP 7.1] Create submission result entity with user provided values
	submissionResult, err := mapper.MapSubmissionResultEntity(submissionID, testCase.ID)
	if err != nil {
		return nil, err
	}

	// [STEP 7.2] Save submission result in database
	result, err := uc.resultRepository.CreateResult(ctx, submissionResult)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}
