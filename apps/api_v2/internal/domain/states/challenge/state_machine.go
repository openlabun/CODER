package challenge_states

import (
	"fmt"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
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

func IsValidState(state Entities.ChallengeStatus) bool {
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

func canTransitionState(from Entities.ChallengeStatus, to Entities.ChallengeStatus) bool {
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

func validateStateTransition(challenge *Entities.Challenge, to Entities.ChallengeStatus) error {
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

func ApplyTranstion(challenge *Entities.Challenge, to Entities.ChallengeStatus) error {
	if err := validateStateTransition(challenge, to); err != nil {
		return err
	}

	challenge.Status = to
	return nil
}
