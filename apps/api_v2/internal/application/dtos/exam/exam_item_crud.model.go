package dtos

type CreateExamItemInput struct {
	ExamID      string `json:"exam_id"`
	ChallengeID string `json:"challenge_id"`
	Order       int    `json:"order"`
	Points      int    `json:"points"`
}

type UpdateExamItemInput struct {
	ID          string  `json:"id"`
	Order       *int    `json:"order"`
	Points      *int    `json:"points"`
}

type DeleteExamItemInput struct {
	ID string `json:"id"`
}