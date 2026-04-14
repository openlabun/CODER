package submission_states

import (
	"fmt"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

// State machine for submission processing
// States: queued -> running -> executed|timeout|error -> accepted|wrong_answer
//   - queued: initial state when a submission is made
//   - running: when the submission is being evaluated
//   - executed: when the submission has been executed
//   - timeout: when the submission times out
//   - accepted: when the submission is accepted
//   - wrong_answer: when the submission is rejected due to wrong answer
//   - error: when there is an internal error during processing (e.g., runtime error, compilation error)

// If a submission is in the accepted/wrong_answer state: ExpectedOutput and ActualOutput may be populated for feedback purposes
// If a submission is in the error state: ErrorMessage may be populated with details about the error

var submissionAllowedTransitions = map[constants.SubmissionStatus]map[constants.SubmissionStatus]struct{}{
	constants.SubmissionStatusQueued: {
		constants.SubmissionStatusRunning: {},
	},
	constants.SubmissionStatusRunning: {
		constants.SubmissionStatusExecuted: {},
		constants.SubmissionStatusTimeout:  {},
		constants.SubmissionStatusError:    {},
	},
	constants.SubmissionStatusExecuted: {
		constants.SubmissionStatusAccepted:    {},
		constants.SubmissionStatusWrongAnswer: {},
	},
}

func IsValidState(state constants.SubmissionStatus) bool {
	switch state {
	case constants.SubmissionStatusQueued:
		return true
	case constants.SubmissionStatusRunning:
		return true
	case constants.SubmissionStatusTimeout:
		return true
	case constants.SubmissionStatusExecuted:
		return true
	case constants.SubmissionStatusAccepted:
		return true
	case constants.SubmissionStatusWrongAnswer:
		return true
	case constants.SubmissionStatusError:
		return true
	default:
		return false
	}
}

func canTransitionState(from constants.SubmissionStatus, to constants.SubmissionStatus) bool {
	if !IsValidState(from) || !IsValidState(to) {
		return false
	}

	nextStates, ok := submissionAllowedTransitions[from]
	if !ok {
		return false
	}

	_, allowed := nextStates[to]
	return allowed
}

func validateStateTransition(submission *Entities.SubmissionResult, to constants.SubmissionStatus) error {
	if !IsValidState(submission.Status) {
		return fmt.Errorf("invalid submission state: %q", submission.Status)
	}

	if !IsValidState(to) {
		return fmt.Errorf("invalid target submission state: %q", to)
	}

	if !canTransitionState(submission.Status, to) {
		return fmt.Errorf("invalid submission transition: %s -> %s", submission.Status, to)
	}

	return nil
}

func ApplyTransition(submission *Entities.SubmissionResult, to constants.SubmissionStatus) error {
	if err := validateStateTransition(submission, to); err != nil {
		return err
	}

	submission.Status = to
	return nil
}
