package submission_usecases

import (
	"context"
	"fmt"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

	submissionPorts "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/submission"
	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	userRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	domain_services "github.com/openlabun/CODER/apps/api_v2/internal/domain/services"
	session_states "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/session"

	mapper "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission/mapper"
	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"
)

type CreateCustomSubmissionUseCase struct {
	userRepository       userRepository.UserRepository
	submissionRepository submissionRepository.SubmissionRepository
	sessionRepository    submissionRepository.SessionRepository
	examRepository    	 examRepository.ExamRepository
	challengeRepository  examRepository.ChallengeRepository
	testCaseRepository   examRepository.TestCaseRepository
	resultRepository     submissionRepository.SubmissionResultRepository
	ioVariableRepository examRepository.IOVariableRepository
	publisherPort        submissionPorts.SubmissionPublisherPort
}

func NewCreateCustomSubmissionUseCase(userRepository userRepository.UserRepository, submissionRepository submissionRepository.SubmissionRepository, sessionRepository submissionRepository.SessionRepository, examRepository examRepository.ExamRepository, challengeRepository examRepository.ChallengeRepository, testCaseRepository examRepository.TestCaseRepository, resultRepository submissionRepository.SubmissionResultRepository, ioVariableRepository examRepository.IOVariableRepository, publisherPort submissionPorts.SubmissionPublisherPort) *CreateCustomSubmissionUseCase {
	return &CreateCustomSubmissionUseCase{
		userRepository:       userRepository,
		submissionRepository: submissionRepository,
		sessionRepository:    sessionRepository,
		examRepository:    	  examRepository,
		challengeRepository:  challengeRepository,
		testCaseRepository:   testCaseRepository,
		resultRepository:     resultRepository,
		ioVariableRepository: ioVariableRepository,
		publisherPort:        publisherPort,
	}
}

func (uc *CreateCustomSubmissionUseCase) Execute(ctx context.Context, input dtos.CreateCustomExecutionInput) (*Entities.Submission, error) {
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

	if user.Role == user_constants.UserRoleProfessor {
		return nil, fmt.Errorf("only students have permissions to make submissions")
	}

	// [STEP 2] Verify existing student session, it belongs to student and its active
	session, err := uc.sessionRepository.GetSessionByID(ctx, input.SessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("no active session found for student %q", user.Username)
	}
	if session.StudentID != user.ID {
		return nil, fmt.Errorf("session with id %q does not belong to student %q", input.SessionID, user.Username)
	}
	if session.Status != constants.SessionStatusActive {
		return nil, fmt.Errorf("session with id %q is not active", input.SessionID)
	}	

	// [STEP 3] Verify exam exists for the session
	exam, err := uc.examRepository.GetExamByID(ctx, session.ExamID)
	if err != nil {
		return nil, err
	}

	// [STEP 4] Update session status
	err = session_states.UpdateSessionStatus(session, exam, services.Now(), false)
	if err != nil {
		return nil, fmt.Errorf("failed to update session status: %w", err)
	}

	// [STEP 5] Save updated session status in database
	_, err = uc.sessionRepository.UpdateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session status: %w", err)
	}

	// [STEP 6] Validate session is still active after status update
	if session.Status != constants.SessionStatusActive {
		return nil, fmt.Errorf("session with id %q is not active", input.SessionID)
	}

	// [STEP 7] Create submission with user provided values
	submission, testCase, err := mapper.MapCreateCustomExecutionInputToEntities(user.ID, input)
	if err != nil {
		return nil, err
	}

	// [STEP 8] Save submission in database
	createdSubmission, err := uc.submissionRepository.CreateSubmission(ctx, submission)
	if err != nil {
		return nil, err
	}

	// [STEP 9] Get challenge and validate language is allowed
	challenge, err := uc.challengeRepository.GetChallengeByID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}

	if challenge == nil {
		return nil, fmt.Errorf("challenge with id %q does not exist", input.ChallengeID)
	}

	if challenge.GetLanguageTemplate(submission.Language) == nil {
		return nil, fmt.Errorf("language %q is not allowed for challenge with id %q", submission.Language, input.ChallengeID)
	}

	// [STEP 10] Get test cases of the challenge
	testCases, err := uc.testCaseRepository.GetTestCasesByChallengeID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}

	if len(testCases) == 0 {
		return nil, fmt.Errorf("no test cases found for challenge with id %q", input.ChallengeID)
	}

	// [STEP 11] Save custom test Case
	testCase, err = domain_services.CreateTestCase(ctx, testCase, uc.testCaseRepository, uc.ioVariableRepository)
	if err != nil {
		return nil, fmt.Errorf("failed to save custom test case: %w", err)
	}

	// [STEP 12] Create submission results for each test case of the challenge
	submissionResult, err := mapper.MapSubmissionResultEntity(submission.ID, testCase.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to map submission result entity: %w", err)
	}

	// [STEP 13] Save submission result in database
	result, err := domain_services.CreateSubmissionResult(ctx, submissionResult, uc.resultRepository, uc.ioVariableRepository)
	if err != nil {
		return nil, fmt.Errorf("failed to create submission result: %w", err)
	}

	// [STEP 14] Create DTO for publishing
	var publishedResults []dtos.SubmissionResultPublishedDTO
	if result != nil {
		publishedResult := mapper.MapSubmissionResultToPublishedDTO(*createdSubmission, *result, *testCase, *challenge)
		if publishedResult != nil {
			publishedResults = append(publishedResults, *publishedResult)
		}
	}

	// [STEP 15] Publish submission created event to message broker for asynchronous processing of submission results
	for _, publishedResult := range publishedResults {
		err = uc.publisherPort.PublishSubmission(publishedResult)
		if err != nil {
			return nil, fmt.Errorf("failed to publish submission result: %w", err)
		}
	}

	return createdSubmission, nil
}