package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	Repository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

type SubmissionRepository struct {
	repository   Repository.SubmissionRepository
	cacheAdapter *cache.CacheAdapter
}

func NewSubmissionRepository(repository Repository.SubmissionRepository, cacheAdapter *cache.CacheAdapter) *SubmissionRepository {
	return &SubmissionRepository{
		repository:   repository,
		cacheAdapter: cacheAdapter,
	}
}

func (r *SubmissionRepository) CreateSubmission(ctx context.Context, s *Entities.Submission) (*Entities.Submission, error) {
	result, err := r.repository.CreateSubmission(ctx, s)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheSubmissionTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *SubmissionRepository) UpdateSubmission(ctx context.Context, s *Entities.Submission) (*Entities.Submission, error) {
	result, err := r.repository.UpdateSubmission(ctx, s)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheSubmissionTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *SubmissionRepository) DeleteSubmission(ctx context.Context, submissionID string) error {
	if err := r.repository.DeleteSubmission(ctx, submissionID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheSubmissionTable, strings.TrimSpace(submissionID))
	return nil
}

func (r *SubmissionRepository) GetSubmissionByID(ctx context.Context, submissionID string) (*Entities.Submission, error) {
	id := strings.TrimSpace(submissionID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheSubmissionTable, id); err == nil && rec != nil {
		cache.NormalizeBools(rec, "scorable")
		if s, e := cache.MapToEntity[Entities.Submission](rec); e == nil {
			return s, nil
		}
	}
	s, err := r.repository.GetSubmissionByID(ctx, submissionID)
	if err != nil || s == nil {
		return s, err
	}
	if rec, e := cache.EntityToMap(s); e == nil {
		r.cacheAdapter.Save(cache.CacheSubmissionTable, cache.InsertOperation, rec)
	}
	return s, nil
}

func (r *SubmissionRepository) GetSubmissionsByExamItemScoreID(ctx context.Context, examItemScoreID string) ([]*Entities.Submission, error) {
	id := strings.TrimSpace(examItemScoreID)
	if recs, err := r.cacheAdapter.FindBy(cache.CacheSubmissionTable, map[string]string{"exam_item_score_id": id}); err == nil && recs != nil {
		return submissionsFromRecords(recs), nil
	}
	submissions, err := r.repository.GetSubmissionsByExamItemScoreID(ctx, examItemScoreID)
	if err != nil {
		return nil, err
	}
	for _, s := range submissions {
		if rec, e := cache.EntityToMap(s); e == nil {
			r.cacheAdapter.Save(cache.CacheSubmissionTable,  cache.InsertOperation, rec)
		}
	}
	return submissions, nil
}

func (r *SubmissionRepository) GetLastSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error) {
	return r.repository.GetLastSubmissionByExamItemScoreID(ctx, examItemScoreID)
}

func (r *SubmissionRepository) GetBestSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error) {
	return r.repository.GetBestSubmissionByExamItemScoreID(ctx, examItemScoreID)
}

func (r *SubmissionRepository) GetSubmissionsBySessionID(ctx context.Context, sessionID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error) {
	conditions := map[string]string{"session_id": strings.TrimSpace(sessionID)}
	if recs, err := r.cacheAdapter.FindBy(cache.CacheSubmissionTable, conditions); err == nil && recs != nil {
		return submissionsFromRecords(recs), nil
	}
	submissions, err := r.repository.GetSubmissionsBySessionID(ctx, sessionID, status, testID, challengeID)
	if err != nil {
		return nil, err
	}
	for _, s := range submissions {
		if rec, e := cache.EntityToMap(s); e == nil {
			r.cacheAdapter.Save(cache.CacheSubmissionTable,  cache.InsertOperation, rec)
		}
	}
	return submissions, nil
}

func (r *SubmissionRepository) GetSubmissionsByUserID(ctx context.Context, userID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error) {
	return r.repository.GetSubmissionsByUserID(ctx, userID, status, testID, challengeID)
}

func (r *SubmissionRepository) GetSubmissionsByChallengeID(ctx context.Context, challengeID string, status *string, testID *string) ([]*Entities.Submission, error) {
	return r.repository.GetSubmissionsByChallengeID(ctx, challengeID, status, testID)
}