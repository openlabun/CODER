package challenge_states

import (
	"fmt"

	constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

// State machine for challenge lifecycle
// States: draft -> published|private -> archived
//   - Can make direct transition to published
//   - Can return from archived to published, but not to draft

var challengeAllowedTransitions = map[constants.ChallengeStatus]map[constants.ChallengeStatus]struct{}{
	constants.ChallengeStatusDraft: {
		constants.ChallengeStatusPublished: {},
	},
	constants.ChallengeStatusPublished: {
		constants.ChallengeStatusArchived: {},
		constants.ChallengeStatusPrivate:  {},
	},
	constants.ChallengeStatusPrivate: {
		constants.ChallengeStatusPublished: {},
		constants.ChallengeStatusArchived: {},
	},
	constants.ChallengeStatusArchived: {
		constants.ChallengeStatusPublished: {},
	},
}

func IsValidState(state constants.ChallengeStatus) bool {
	switch state {
	case constants.ChallengeStatusDraft:
		return true
	case constants.ChallengeStatusPublished:
		return true
	case constants.ChallengeStatusArchived:
		return true
	case constants.ChallengeStatusPrivate:
		return true
	default:
		return false
	}
}

func canTransitionState(from constants.ChallengeStatus, to constants.ChallengeStatus) bool {
	if !IsValidState(from) || !IsValidState(to) {
		return false
	}

	nextStates, ok := challengeAllowedTransitions[from]
	if !ok {
		return false
	}

	_, allowed := nextStates[to]
	return allowed
}

func validateStateTransition(challenge *Entities.Challenge, to constants.ChallengeStatus) error {
	if !IsValidState(challenge.Status) {
		return fmt.Errorf("invalid challenge state: %q", challenge.Status)
	}

	if !IsValidState(to) {
		return fmt.Errorf("invalid target challenge state: %q", to)
	}

	if !canTransitionState(challenge.Status, to) {
		return fmt.Errorf("invalid challenge transition: %s -> %s", challenge.Status, to)
	}

	return nil
}

func ApplyTransition(challenge *Entities.Challenge, to constants.ChallengeStatus) error {
	if err := validateStateTransition(challenge, to); err != nil {
		return err
	}

	challenge.Status = to
	return nil
}
