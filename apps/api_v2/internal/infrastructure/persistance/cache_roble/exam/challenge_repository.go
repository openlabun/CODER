package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/exam"
)

const cacheChallengeTable = "cache_challenge"

type ChallengeRepository struct {
	robleRepository *repository.ChallengeRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewChallengeRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *ChallengeRepository {
	ioVariableRepository := repository.NewIOVariableRepository(robleAdapter)
	return &ChallengeRepository{
		robleRepository: repository.NewChallengeRepository(robleAdapter, ioVariableRepository),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *ChallengeRepository) CreateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error) {
	result, err := r.robleRepository.CreateChallenge(ctx, challenge)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cacheChallengeTable, "insert", rec)
	}
	return result, nil
}

func (r *ChallengeRepository) UpdateChallenge(ctx context.Context, challenge *Entities.Challenge) (*Entities.Challenge, error) {
	result, err := r.robleRepository.UpdateChallenge(ctx, challenge)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cacheChallengeTable, "update", rec)
	}
	return result, nil
}

func (r *ChallengeRepository) DeleteChallenge(ctx context.Context, challengeID string) error {
	if err := r.robleRepository.DeleteChallenge(ctx, challengeID); err != nil {
		return err
	}
	r.cacheAdapter.DeleteByID(cacheChallengeTable, strings.TrimSpace(challengeID))
	return nil
}

func (r *ChallengeRepository) GetChallengeByID(ctx context.Context, challengeID string) (*Entities.Challenge, error) {
	id := strings.TrimSpace(challengeID)
	if rec, err := r.cacheAdapter.FindByID(cacheChallengeTable, id); err == nil && rec != nil {
		if c, e := cache.MapToEntity[Entities.Challenge](rec); e == nil {
			return c, nil
		}
	}
	c, err := r.robleRepository.GetChallengeByID(ctx, challengeID)
	if err != nil || c == nil {
		return c, err
	}
	if rec, e := cache.EntityToMap(c); e == nil {
		r.cacheAdapter.Save(cacheChallengeTable, "insert", rec)
	}
	return c, nil
}

func (r *ChallengeRepository) GetChallenges(ctx context.Context, status, tag, difficulty *string) ([]*Entities.Challenge, error) {
	challenges, err := r.robleRepository.GetChallenges(ctx, status, tag, difficulty)
	if err != nil {
		return nil, err
	}
	for _, c := range challenges {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cacheChallengeTable, "insert", rec)
		}
	}
	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByUserID(ctx context.Context, userID string, examID *string) ([]*Entities.Challenge, error) {
	id := strings.TrimSpace(userID)
	if examID == nil {
		if recs, err := r.cacheAdapter.FindBy(cacheChallengeTable, map[string]string{"user_id": id}); err == nil && recs != nil {
			return challengesFromRecords(recs), nil
		}
	}
	challenges, err := r.robleRepository.GetChallengesByUserID(ctx, userID, examID)
	if err != nil {
		return nil, err
	}
	for _, c := range challenges {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cacheChallengeTable, "insert", rec)
		}
	}
	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByExamID(ctx context.Context, examID string) ([]*Entities.Challenge, error) {
	challenges, err := r.robleRepository.GetChallengesByExamID(ctx, examID)
	if err != nil {
		return nil, err
	}
	for _, c := range challenges {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cacheChallengeTable, "insert", rec)
		}
	}
	return challenges, nil
}

func (r *ChallengeRepository) GetChallengesByTag(ctx context.Context, tag string) ([]*Entities.Challenge, error) {
	challenges, err := r.robleRepository.GetChallengesByTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	for _, c := range challenges {
		if rec, e := cache.EntityToMap(c); e == nil {
			r.cacheAdapter.Save(cacheChallengeTable, "insert", rec)
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
