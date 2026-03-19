package states

import (
	"fmt"

	Entities "../entities/submission"
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

func IsValidSubmissionState(state Entities.SubmissionStatus) bool {
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

func CanTransitionSubmissionState(from Entities.SubmissionStatus, to Entities.SubmissionStatus) bool {
	if !IsValidSubmissionState(from) || !IsValidSubmissionState(to) {
		return false
	}

	nextStates, ok := submissionAllowedTransitions[from]
	if !ok {
		return false
	}

	_, allowed := nextStates[to]
	return allowed
}

func ValidateSubmissionStateTransition(from Entities.SubmissionStatus, to Entities.SubmissionStatus) error {
	if !IsValidSubmissionState(from) {
		return fmt.Errorf("invalid submission state: %q", from)
	}

	if !IsValidSubmissionState(to) {
		return fmt.Errorf("invalid target submission state: %q", to)
	}

	if !CanTransitionSubmissionState(from, to) {
		return fmt.Errorf("invalid submission transition: %s -> %s", from, to)
	}

	return nil
}

type SubmissionStateData struct {
	HasExpectedOutput bool
	HasActualOutput   bool
	HasErrorMessage   bool
}

func ValidateSubmissionStateData(state Entities.SubmissionStatus, data SubmissionStateData) error {
	if !IsValidSubmissionState(state) {
		return fmt.Errorf("invalid submission state: %q", state)
	}

	if data.HasExpectedOutput != data.HasActualOutput {
		return fmt.Errorf("expected_output and actual_output must be both present or both absent")
	}

	switch state {
	case Entities.SubmissionStatusQueued, Entities.SubmissionStatusRunning:
		if data.HasExpectedOutput || data.HasActualOutput || data.HasErrorMessage {
			return fmt.Errorf("state %s cannot include outputs or error message", state)
		}
	case Entities.SubmissionStatusAccepted, Entities.SubmissionStatusWrongAnswer:
		if data.HasErrorMessage {
			return fmt.Errorf("state %s cannot include error message", state)
		}
	case Entities.SubmissionStatusError:
		if data.HasExpectedOutput || data.HasActualOutput {
			return fmt.Errorf("state %s cannot include expected/actual outputs", state)
		}
	}

	return nil
}
