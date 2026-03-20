package submission_states

import (
	"fmt"

	Entities "../../entities/submission"
)

// State machine for submission processing
// States: queued -> running -> accepted|wrong_answer|error
//   - queued: initial state when a submission is made
//   - running: when the submission is being evaluated
//   - accepted: when the submission passes all test cases
//   - wrong_answer: when the submission fails one or more test cases
//   - error: when there is an internal error during processing (e.g., runtime error, compilation error)

// If a submission is in the accepted/wrong_answer state: ExpectedOutput and ActualOutput may be populated for feedback purposes
// If a submission is in the error state: ErrorMessage may be populated with details about the error

var submissionAllowedTransitions = map[Entities.SubmissionStatus]map[Entities.SubmissionStatus]struct{}{
	Entities.SubmissionStatusQueued: {
		Entities.SubmissionStatusRunning: {},
	},
	Entities.SubmissionStatusRunning: {
		Entities.SubmissionStatusAccepted:    {},
		Entities.SubmissionStatusWrongAnswer: {},
		Entities.SubmissionStatusError:       {},
	},
}

func IsValidState(state Entities.SubmissionStatus) bool {
	switch state {
	case Entities.SubmissionStatusQueued:
		return true
	case Entities.SubmissionStatusRunning:
		return true
	case Entities.SubmissionStatusAccepted:
		return true
	case Entities.SubmissionStatusWrongAnswer:
		return true
	case Entities.SubmissionStatusError:
		return true
	default:
		return false
	}
}

func canTransitionState(from Entities.SubmissionStatus, to Entities.SubmissionStatus) bool {
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

func validateStateTransition(submission *Entities.SubmissionResult, to Entities.SubmissionStatus) error {
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

func ApplyTransition(submission *Entities.SubmissionResult, to Entities.SubmissionStatus) error {
	if err := validateStateTransition(submission, to); err != nil {
		return err
	}

	submission.Status = to
	return nil
}
