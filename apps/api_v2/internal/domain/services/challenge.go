package services

import (
	"context"
	"fmt"

	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	examRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submissionRepository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
)

func CreateChallenge(
	ctx context.Context,
	challenge *examEntities.Challenge,
	challengeRepository examRepository.ChallengeRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) (*examEntities.Challenge, error) {
	if challenge == nil {
		return nil, fmt.Errorf("challenge is nil")
	}

	hydrated, err := hydrateChallengeIOVariables(ctx, challenge, ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return challengeRepository.CreateChallenge(ctx, hydrated)
}

func UpdateChallenge(
	ctx context.Context,
	challenge *examEntities.Challenge,
	challengeRepository examRepository.ChallengeRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) (*examEntities.Challenge, error) {
	if challenge == nil {
		return nil, fmt.Errorf("challenge is nil")
	}

	hydrated, err := hydrateChallengeIOVariables(ctx, challenge, ioVariableRepository)
	if err != nil {
		return nil, err
	}

	return challengeRepository.UpdateChallenge(ctx, hydrated)
}

func RemoveChallenge(ctx context.Context,
	challengeID string,
	challengeRepository examRepository.ChallengeRepository,
	testCaseRepository examRepository.TestCaseRepository,
	examItemRepository examRepository.ExamItemRepository,
	submissionRepository submissionRepository.SubmissionRepository,
	resultsRepository submissionRepository.SubmissionResultRepository,
	ioVariableRepository examRepository.IOVariableRepository,
) error {
	// [STEP 1] Get all exam items for the challenge
	examItems, err := examItemRepository.GetExamItem(ctx, nil, &challengeID)
	if err != nil {
		return err
	}

	// [STEP 2] Delete all existing exam items for the challenge
	for _, item := range examItems {
		if item != nil {
			err = examItemRepository.DeleteExamItem(ctx, item.ID)
			if err != nil {
				return err
			}
		}
	}

	// [STEP 3] Get all test cases for the challenge
	test_cases, err := testCaseRepository.GetTestCasesByChallengeID(ctx, challengeID)
	if err != nil {
		return err
	}

	// [STEP 4] Delete all existing test cases for the challenge
	for _, test_case := range test_cases {
		if test_case != nil {
			err = RemoveTestCase(ctx, test_case.ID, testCaseRepository, ioVariableRepository)
			if err != nil {
				return err
			}
		}
	}

	// [STEP 5] Get all existing submissions for the challenge
	submissions, err := submissionRepository.GetSubmissionsByChallengeID(ctx, challengeID, nil, nil)
	if err != nil {
		return err
	}

	// [STEP 6] Delete all existing submissions for the challenge
	for _, submission := range submissions {
		if submission != nil {
			err = RemoveSubmission(ctx, submission.ID, submissionRepository, resultsRepository, ioVariableRepository)
			if err != nil {
				return err
			}
		}
	}

	// [STEP 7] Delete all Challenge's IOVariables
	challenge, err := challengeRepository.GetChallengeByID(ctx, challengeID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return fmt.Errorf("challenge with id %q does not exist", challengeID)
	}

	for _, ioVariable := range challenge.InputVariables {
		err = ioVariableRepository.DeleteIOVariable(ctx, ioVariable.ID)
		if err != nil {
			return err
		}
	}

	err = ioVariableRepository.DeleteIOVariable(ctx, challenge.OutputVariable.ID)
	if err != nil {
		return err
	}

	// [STEP 8] Delete the challenge itself
	err = challengeRepository.DeleteChallenge(ctx, challengeID)
	if err != nil {
		return err
	}

	return nil
}
