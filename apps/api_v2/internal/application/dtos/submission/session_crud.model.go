package dtos

type CreateSessionInput struct {
	UserID string
	ExamID string
}

type HeartbeatSessionInput struct {
	SessionID string
}

type GetSessionInput struct {
	SessionID string
}

type CloseSessionInput struct {
	SessionID string
}