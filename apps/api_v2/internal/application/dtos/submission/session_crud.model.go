package dtos

type CreateSessionInput struct {
	UserID string		`json:"user_id"`
	ExamID string		`json:"exam_id"`
}

type HeartbeatSessionInput struct {
	SessionID string	`json:"session_id"`
}

type GetSessionInput struct {
	SessionID string 	`json:"session_id"`
}