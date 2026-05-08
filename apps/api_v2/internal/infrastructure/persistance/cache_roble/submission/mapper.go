package cache_roble_infrastructure

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

func resultsFromRecords(records []map[string]any) []*Entities.SubmissionResult {
	results := make([]*Entities.SubmissionResult, 0, len(records))
	for _, rec := range records {
		if res, err := cache.MapToEntity[Entities.SubmissionResult](rec); err == nil && res != nil {
			results = append(results, res)
		}
	}
	return results
}

func sessionsFromRecords(records []map[string]any) []*Entities.Session {
	sessions := make([]*Entities.Session, 0, len(records))
	for _, rec := range records {
		if s, err := cache.MapToEntity[Entities.Session](rec); err == nil && s != nil {
			sessions = append(sessions, s)
		}
	}
	return sessions
}

func submissionsFromRecords(records []map[string]any) []*Entities.Submission {
	submissions := make([]*Entities.Submission, 0, len(records))
	for _, rec := range records {
		cache.NormalizeBools(rec, "scorable")
		if s, err := cache.MapToEntity[Entities.Submission](rec); err == nil && s != nil {
			submissions = append(submissions, s)
		}
	}
	return submissions
}