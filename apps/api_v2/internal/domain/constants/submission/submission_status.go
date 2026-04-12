package constants

type SubmissionStatus string

const (
	SubmissionStatusQueued      SubmissionStatus = "queued"
	SubmissionStatusRunning     SubmissionStatus = "running"
	SubmissionStatusTimeout     SubmissionStatus = "timeout"
	SubmissionStatusExecuted    SubmissionStatus = "executed"
	SubmissionStatusAccepted    SubmissionStatus = "accepted"
	SubmissionStatusWrongAnswer SubmissionStatus = "wrong_answer"
	SubmissionStatusError       SubmissionStatus = "error"
)
