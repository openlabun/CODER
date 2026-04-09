package dtos

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

type CreateExamInput struct {
	CourseID    		*string `json:"course_id"`
	Title           	string 	`json:"title"`
	Description     	string 	`json:"description"`
	Visibility      	string 	`json:"visibility"`
	StartTime       	string 	`json:"start_time"`
	EndTime         	*string `json:"end_time"`
	AllowLateSubmissions bool 	`json:"allow_late_submissions"`
	TimeLimit 			int 	`json:"time_limit"`
	TryLimit  			int 	`json:"try_limit"`
	ProfessorID 		string 	`json:"professor_id"`
}

type UpdateExamInput struct {
	ExamID      		string 		`json:"exam_id"`
	Title           	*string 	`json:"title"`
	Description     	*string 	`json:"description"`
	Visibility      	*string 	`json:"visibility"`
	StartTime       	*string 	`json:"start_time"`
	EndTime         	*string 	`json:"end_time"`
	AllowLateSubmissions *bool 		`json:"allow_late_submissions"`
	TimeLimit 			*int 		`json:"time_limit"`
	TryLimit  			*int 		`json:"try_limit"`
}

type DeleteExamInput struct {
	ExamID string `json:"exam_id"`
}

type ChangeExamVisibilityInput struct {
	ExamID     string `json:"exam_id"`
	Visibility string `json:"visibility"`
}

type CloseExamInput struct {
	ExamID string `json:"exam_id"`
}

type GetExamDetailsInput struct {
	ExamID string `json:"exam_id"`
}

type GetExamsByCourseInput struct {
	CourseID string `json:"course_id"`
}

type GetExamItemsInput struct {
	ExamID string `json:"exam_id"`
}

type ExamItemDTO struct {
	ID          string `json:"id"`
	Order 	 	int    `json:"order"`
	Points 	 	int    `json:"points"`
	ExamID      string `json:"exam_id"`
	Challenge   *Entities.Challenge `json:"challenge,omitempty"`
}