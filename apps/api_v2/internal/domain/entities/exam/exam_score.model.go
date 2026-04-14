package challenge_entities

type ExamScore struct {
	ID 			string			`json:"id"`

	// Relationships
	ExamID      string			`json:"exam_id"`
	SessionID 	string			`json:"session_id"`

	// Exam Execution Details
	Score 		int				`json:"score"`

	// Metadata
	CreatedAt 	string			`json:"created_at"`
	UpdatedAt 	string			`json:"updated_at"`
	StudentID 	string			`json:"student_id"`
}