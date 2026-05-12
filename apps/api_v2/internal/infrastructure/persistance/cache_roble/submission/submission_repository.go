package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/submission"
)

type SubmissionRepository struct {
	robleRepository *repository.SubmissionRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewSubmissionRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *SubmissionRepository {
	return &SubmissionRepository{
		robleRepository: repository.NewSubmissionRepository(robleAdapter),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *SubmissionRepository) CreateSubmission(ctx context.Context, s *Entities.Submission) (*Entities.Submission, error) {
	result, err := r.robleRepository.CreateSubmission(ctx, s)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheSubmissionTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *SubmissionRepository) UpdateSubmission(ctx context.Context, s *Entities.Submission) (*Entities.Submission, error) {
	result, err := r.robleRepository.UpdateSubmission(ctx, s)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheSubmissionTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *SubmissionRepository) DeleteSubmission(ctx context.Context, submissionID string) error {
	if err := r.robleRepository.DeleteSubmission(ctx, submissionID); err != nil {
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
	s, err := r.robleRepository.GetSubmissionByID(ctx, submissionID)
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
	submissions, err := r.robleRepository.GetSubmissionsByExamItemScoreID(ctx, examItemScoreID)
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
	return r.robleRepository.GetLastSubmissionByExamItemScoreID(ctx, examItemScoreID)
}

func (r *SubmissionRepository) GetBestSubmissionByExamItemScoreID(ctx context.Context, examItemScoreID string) (*Entities.Submission, error) {
	return r.robleRepository.GetBestSubmissionByExamItemScoreID(ctx, examItemScoreID)
}

func (r *SubmissionRepository) GetSubmissionsBySessionID(ctx context.Context, sessionID string, status *string, testID *string, challengeID *string) ([]*Entities.Submission, error) {
	conditions := map[string]string{"session_id": strings.TrimSpace(sessionID)}
	if recs, err := r.cacheAdapter.FindBy(cache.CacheSubmissionTable, conditions); err == nil && recs != nil {
		return submissionsFromRecords(recs), nil
	}
	submissions, err := r.robleRepository.GetSubmissionsBySessionID(ctx, sessionID, status, testID, challengeID)
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
	return r.robleRepository.GetSubmissionsByUserID(ctx, userID, status, testID, challengeID)
}

func (r *SubmissionRepository) GetSubmissionsByChallengeID(ctx context.Context, challengeID string, status *string, testID *string) ([]*Entities.Submission, error) {
	return r.robleRepository.GetSubmissionsByChallengeID(ctx, challengeID, status, testID)
}