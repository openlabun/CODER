package dtos

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

type GetExamScoreInput struct {
	ExamID string `json:"exam_id"`
	UserID *string `json:"user_id,omitempty"`
}

type ExamItemScoreDTO struct {
	ExamItemID string `json:"exam_item_id"`
	Challenge  *Entities.Challenge `json:"challenge,omitempty"`
	Score      float64 `json:"score"`
}

type ExamScoreDTO struct {
	UserID string `json:"user_id"`
	Exam *Entities.Exam `json:"exam"`
	Score float64 `json:"score"`
	ExamItemScores []*ExamItemScoreDTO `json:"exam_item_scores"`
}