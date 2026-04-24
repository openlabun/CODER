package dtos

type CreateExamItemInput struct {
	ExamID      string `json:"exam_id"`
	ChallengeID string `json:"challenge_id"`
	Order       int    `json:"order"`
	Points      int    `json:"points"`
	TryLimit    *int   `json:"try_limit"`	// Optional, -1 for unlimited
}

type UpdateExamItemInput struct {
	ID          string  `json:"id"`
	Order       *int    `json:"order"`
	Points      *int    `json:"points"`
	TryLimit    *int    `json:"try_limit"`	// Optional, -1 for unlimited
}

type DeleteExamItemInput struct {
	ID string 			`json:"id"`
}