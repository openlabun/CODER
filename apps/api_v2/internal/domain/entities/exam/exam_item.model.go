package challenge_entities

type ExamItem struct {
	ID 	  	 	string			`json:"id"`
	ChallengeID string			`json:"challenge_id"`
	ExamID      string			`json:"exam_id"`

	Order       int				`json:"order"`
	Points	  	int				`json:"points"`
}