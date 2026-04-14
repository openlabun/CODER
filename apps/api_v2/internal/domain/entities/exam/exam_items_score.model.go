package challenge_entities

type ExamItemScore struct {
	ID 				string			`json:"id"`

	// Relationships
	ExamItemID 		string			`json:"exam_item_id"`
	ExamScoreID 	string			`json:"exam_score_id"`

	// Scoring Details
	Score 			int				`json:"score"`
	Tries 			int				`json:"tries"`

	// Metadata
	CreatedAt 		string			`json:"created_at"`
	UpdatedAt 		string			`json:"updated_at"`
}
