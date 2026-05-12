package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Repository "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

type ChallengeRepository struct {
	repository   Repository.ChallengeRepository
	cacheAdapter *cache.CacheAdapter
}

func NewChallengeRepository(repository Repository.ChallengeRepository, cacheAdapter *cache.CacheAdapter) *ChallengeRepository {
	return &ChallengeRepository{
		repository:   repository,
		cacheAdapter: cacheAdapter,
	}
}

func (r *ChallengeRepository) CreateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error) {
	result, err := r.repository.CreateChallenge(ctx, challenge)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheChallengeTable, cache.InsertOperation, rec)
	}
	return result, nil
}

func (r *ChallengeRepository) UpdateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error) {
	result, err := r.repository.UpdateChallenge(ctx, challenge)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cache.CacheChallengeTable, cache.UpdateOperation, rec)
	}
	return result, nil
}

func (r *ChallengeRepository) DeleteChallenge(ctx context.Context, challengeID string) error {
	if err := r.repository.DeleteChallenge(ctx, challengeID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cache.CacheChallengeTable, strings.TrimSpace(challengeID))
	return nil
}

func (r *ChallengeRepository) GetChallengeByID(ctx context.Context, challengeID string) (*Entities.Challenge, error) {
	id := strings.TrimSpace(challengeID)
	if rec, err := r.cacheAdapter.FindByID(cache.CacheChallengeTable, id); err == nil && rec != nil {
		if c, e := cache.MapToEntity[Entities.Challenge](rec); e == nil {
			return c, nil
		}
	}
	c, err := r.repository.GetChallengeByID(ctx, challengeID)
	if err != nil || c == nil {
		return c, err
	}
	if rec, e := cache.EntityToMap(c); e == nil {
		r.cacheAdapter.Save(cache.CacheChallengeTable, cache.InsertOperation, rec)
	}
	return c, nil
}

func (r *ChallengeRepository) GetChallenges(ctx context.Context, status, tag, difficulty *string) ([]*Entities.Challenge, error) {
	challenges, err := r.repository.GetChallenges(ctx, status, tag, difficulty)
	if err != nil {
		return nil, err
	}
	for _, c := range challenges {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cache.CacheChallengeTable, cache.InsertOperation, rec)
		}
	}
	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByUserID(ctx context.Context, userID string, examID *string) ([]*Entities.Challenge, error) {
	id := strings.TrimSpace(userID)
	if examID == nil {
		if recs, err := r.cacheAdapter.FindBy(cache.CacheChallengeTable, map[string]string{"user_id": id}); err == nil && recs != nil {
			return challengesFromRecords(recs), nil
		}
	}
	challenges, err := r.repository.GetChallengesByUserID(ctx, userID, examID)
	if err != nil {
		return nil, err
	}
	for _, c := range challenges {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cache.CacheChallengeTable, cache.InsertOperation, rec)
		}
	}
	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByExamID(ctx context.Context, examID string) ([]*Entities.Challenge, error) {
	challenges, err := r.repository.GetChallengesByExamID(ctx, examID)
	if err != nil {
		return nil, err
	}
	for _, c := range challenges {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cache.CacheChallengeTable, cache.InsertOperation, rec)
		}
	}
	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByTag(ctx context.Context, tag string) ([]*Entities.Challenge, error) {
	challenges, err := r.repository.GetChallengesByTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	for _, c := range challenges {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cache.CacheChallengeTable, cache.InsertOperation, rec)
		}
	}
	return challenges, nil
}

func (r *ChallengeRepository) GetInputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error) {
	c, err := r.GetChallengeByID(ctx, challengeID)
	if err != nil || c == nil {
		return []*Entities.IOVariable{}, err
	}
	vars := make([]*Entities.IOVariable, 0, len(c.InputVariables))
	for i := range c.InputVariables {
		v := c.InputVariables[i]
		vars = append(vars, &v)
	}
	return vars, nil
}

func (r *ChallengeRepository) GetOutputVariablesByChallengeID(ctx context.Context, challengeID string) ([]*Entities.IOVariable, error) {
	c, err := r.GetChallengeByID(ctx, challengeID)
	if err != nil || c == nil || strings.TrimSpace(c.OutputVariable.ID) == "" {
		return []*Entities.IOVariable{}, err
	}
	v := c.OutputVariable
	return []*Entities.IOVariable{&v}, nil
}
