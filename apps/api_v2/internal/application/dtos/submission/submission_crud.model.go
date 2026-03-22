package dtos

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type CreateSubmissionInput struct {
	Code        string
	Language    string
	Score       int
	ChallengeID string
	SessionID   string
}

type GetSubmissionStatusInput struct {
	SubmissionID string
}

type GetUserSubmissionsInput struct {
	UserID      string
	Status      *string
	TestID      *string
	ChallengeID *string
}

type GetChallengeSubmissionsInput struct {
	ChallengeID string
	Status      *string
	TestID      *string
}

type SubmissionOutputDTO struct {
	Submission Entities.Submission
	Results    []Entities.SubmissionResult
}
