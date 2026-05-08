package cache_roble_infrastructure

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

func testCasesFromRecords(records []map[string]any) []*Entities.TestCase {
	testCases := make([]*Entities.TestCase, 0, len(records))
	for _, rec := range records {
		cache.NormalizeBools(rec, "is_sample", "custom")
		if t, err := cache.MapToEntity[Entities.TestCase](rec); err == nil && t != nil {
			testCases = append(testCases, t)
		}
	}
	return testCases
}

func challengesFromRecords(records []map[string]any) []*Entities.Challenge {
	challenges := make([]*Entities.Challenge, 0, len(records))
	for _, rec := range records {
		if c, err := cache.MapToEntity[Entities.Challenge](rec); err == nil && c != nil {
			challenges = append(challenges, c)
		}
	}
	return challenges
}

func examItemsFromRecords(records []map[string]any) []*Entities.ExamItem {
	items := make([]*Entities.ExamItem, 0, len(records))
	for _, rec := range records {
		if item, err := cache.MapToEntity[Entities.ExamItem](rec); err == nil && item != nil {
			items = append(items, item)
		}
	}
	return items
}

func examItemScoresFromRecords(records []map[string]any) []*Entities.ExamItemScore {
	scores := make([]*Entities.ExamItemScore, 0, len(records))
	for _, rec := range records {
		if s, err := cache.MapToEntity[Entities.ExamItemScore](rec); err == nil && s != nil {
			scores = append(scores, s)
		}
	}
	return scores
}

func examsFromRecords(records []map[string]any) []*Entities.Exam {
	exams := make([]*Entities.Exam, 0, len(records))
	for _, rec := range records {
		cache.NormalizeBools(rec, "allow_late_submissions")
		if exam, err := cache.MapToEntity[Entities.Exam](rec); err == nil && exam != nil {
			exams = append(exams, exam)
		}
	}
	return exams
}

func examScoresFromRecords(records []map[string]any) []*Entities.ExamScore {
	scores := make([]*Entities.ExamScore, 0, len(records))
	for _, rec := range records {
		if s, err := cache.MapToEntity[Entities.ExamScore](rec); err == nil && s != nil {
			scores = append(scores, s)
		}
	}
	return scores
}