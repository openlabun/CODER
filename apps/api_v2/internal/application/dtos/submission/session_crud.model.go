package dtos

type CreateSessionInput struct {
	UserID string		`json:"user_id"`
	ExamID string		`json:"exam_id"`
}

type HeartbeatSessionInput struct {
	SessionID string	`json:"session_id"`
}

type GetActiveSessionInput struct {
	UserID *string 	`json:"user_id"`
}

type CloseSessionInput struct {
	SessionID string 	`json:"session_id"`
}

type BlockSessionInput struct {
	SessionID string 	`json:"session_id"`
}