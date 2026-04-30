package services

import (
	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func CalculateBestScore(scores []*examEntities.ExamScore) int {
	bestScore := 0
	for _, score := range scores {
		if score.Score > bestScore {
			bestScore = score.Score
		}
	}
	return bestScore
}

func CalculatePromScore(scores []*examEntities.ExamScore) float64 {
	if len(scores) == 0 {
		return 0
	}

	totalScore := 0
	for _, score := range scores {
		totalScore += score.Score
	}

	return float64(totalScore) / float64(len(scores))
}

func AggregateScoresByExam(examScores []*examEntities.ExamScore) map[string][]*examEntities.ExamScore {
	aggregatedScores := make(map[string][]*examEntities.ExamScore)

	for _, score := range examScores {
		examID := score.ExamID
		aggregatedScores[examID] = append(aggregatedScores[examID], score)
	}

	return aggregatedScores
}

func AggregateScoresByUser(examScores []*examEntities.ExamScore) map[string][]*examEntities.ExamScore {
	aggregatedScores := make(map[string][]*examEntities.ExamScore)

	for _, score := range examScores {
		userID := score.StudentID
		aggregatedScores[userID] = append(aggregatedScores[userID], score)
	}

	return aggregatedScores
}

func FilterExamScoresByExams (aggregatedScores map[string][]*examEntities.ExamScore, exams []examEntities.Exam) map[string][]*examEntities.ExamScore {
	var filteredExamScores map[string][]*examEntities.ExamScore
	filteredExamScores = make(map[string][]*examEntities.ExamScore)

	for exam, scores := range aggregatedScores {
		for _, e := range exams {
			if e.ID == exam {
				filteredExamScores[exam] = scores
			}
		}
	}

	return filteredExamScores
}

