package states

import (
	"fmt"

	Entities "../entities/exam"
)

// State machine for challenge lifecycle
// States: draft -> published -> archived
//   - Can make direct transition to published
//   - Can return from archived to published, but not to draft

var challengeAllowedTransitions = map[Entities.ChallengeStatus]map[Entities.ChallengeStatus]struct{}{
	Entities.ChallengeStatusDraft: {
		Entities.ChallengeStatusPublished: {},
	},
	Entities.ChallengeStatusPublished: {
		Entities.ChallengeStatusArchived: {},
	},
	Entities.ChallengeStatusArchived: {
		Entities.ChallengeStatusPublished: {},
	},
}

func IsValidChallengeState(state Entities.ChallengeStatus) bool {
	switch state {
	case Entities.ChallengeStatusDraft:
		return true
	case Entities.ChallengeStatusPublished:
		return true
	case Entities.ChallengeStatusArchived:
		return true
	default:
		return false
	}
}

func CanTransitionChallengeState(from Entities.ChallengeStatus, to Entities.ChallengeStatus) bool {
	if !IsValidChallengeState(from) || !IsValidChallengeState(to) {
		return false
	}

	nextStates, ok := challengeAllowedTransitions[from]
	if !ok {
		return false
	}

	_, allowed := nextStates[to]
	return allowed
}

func ValidateChallengeStateTransition(from Entities.ChallengeStatus, to Entities.ChallengeStatus) error {
	if !IsValidChallengeState(from) {
		return fmt.Errorf("invalid challenge state: %q", from)
	}

	if !IsValidChallengeState(to) {
		return fmt.Errorf("invalid target challenge state: %q", to)
	}

	if !CanTransitionChallengeState(from, to) {
		return fmt.Errorf("invalid challenge transition: %s -> %s", from, to)
	}

	return nil
}
