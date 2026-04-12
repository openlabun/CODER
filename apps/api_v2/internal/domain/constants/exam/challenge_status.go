package constants

type ChallengeStatus string

const (
	ChallengeStatusDraft     ChallengeStatus = "draft"
	ChallengeStatusPublished ChallengeStatus = "published"
	ChallengeStatusPrivate   ChallengeStatus = "private"
	ChallengeStatusArchived  ChallengeStatus = "archived"
)
