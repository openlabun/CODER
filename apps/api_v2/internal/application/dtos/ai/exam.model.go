package dtos

type GenerateExamInput struct {
	Topic string `json:"topic"`
	Difficulty string `json:"difficulty"`
}

type AIExamIdea struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeLimit   int    `json:"time_limit"` // in minutes
	TryLimit    int    `json:"try_limit"`
}

type GenerateExamOutput struct {
	Exam AIExamIdea `json:"exam"`
}