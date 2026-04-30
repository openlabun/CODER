package dtos

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	SubmissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

type GetExamScoringInput struct {
	ExamID 		string `json:"exam_id"`
	UserID  	string `json:"user_id"`
}

type GetExamScoresByUserInput struct {
	UserID 		string `json:"user_id"`
}

type GetUserScoresByExamInput struct {
	ExamID 		string 	`json:"exam_id"`
	CourseID 	*string `json:"course_id"`
}

type GetItemScoringInput struct {
	ExamScoreID  string `json:"exam_score_id"`
}

type ExamScoringDTO struct {
	ExamID     	string 				 `json:"exam_id"`
	BestScore  	int 			 	 `json:"best_score"`
	ExamScores 	[]Entities.ExamScore `json:"exam_scores"`
}

type ExamItemScoringDTO struct {
	ExamItem 	  Entities.ExamItem 				`json:"exam_item"`
	ExamItemScore Entities.ExamItemScore 			`json:"exam_item_score"`
	Submissions	  []SubmissionEntities.Submission 	`json:"submissions"`
}

type UserExamsScoresOutputDTO struct {
	UserID 		string 			 `json:"user_id"`
	PromScore 	float64 		 `json:"prom_score"`
	ExamScores  []ExamScoringDTO `json:"exam_scores"`
}

type UserExamScoresOutputDTO struct {
	UserID 		string 			 	`json:"user_id"`
	ExamScores  ExamScoringDTO 	 	`json:"exam_scores"`
}

type UsersExamScoresOutputDTO struct {
	ExamID 		string 			 		 	`json:"exam_id"`
	ExamScores  []UserExamScoresOutputDTO 	`json:"users"`
}

type ExamScoreDetailsOutputDTO struct {
	ExamScoreID 	string 					`json:"exam_score_id"`
	ExamID 			string 					`json:"exam_id"`
	ExamItems 		[]ExamItemScoringDTO 	`json:"exam_items"`
}