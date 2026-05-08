package cache_roble_infrastructure

import (
	"context"
	"strings"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	roble_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
	repository "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble/user"
)

const cacheUserTable = "cache_user"

type UserRepository struct {
	robleRepository *repository.UserRepository
	cacheAdapter    *cache.CacheAdapter
}

func NewUserRepository(robleAdapter *roble_infrastructure.RobleDatabaseAdapter, cacheAdapter *cache.CacheAdapter) *UserRepository {
	return &UserRepository{
		robleRepository: repository.NewUserRepository(robleAdapter),
		cacheAdapter:    cacheAdapter,
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *Entities.User) (*Entities.User, error) {
	result, err := r.robleRepository.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	if rec, e := cache.EntityToMap(result); e == nil {
		r.cacheAdapter.Save(cacheUserTable, "insert", rec)
	}
	return result, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*Entities.User, error) {
	id := strings.TrimSpace(userID)
	if rec, err := r.cacheAdapter.FindByID(cacheUserTable, id); err == nil && rec != nil {
		if u, e := cache.MapToEntity[Entities.User](rec); e == nil {
			return u, nil
		}
	}
	u, err := r.robleRepository.GetUserByID(ctx, userID)
	if err != nil || u == nil {
		return u, err
	}
	if rec, e := cache.EntityToMap(u); e == nil {
		r.cacheAdapter.Save(cacheUserTable, "insert", rec)
	}
	return u, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*Entities.User, error) {
	e := strings.TrimSpace(email)
	if recs, err := r.cacheAdapter.FindBy(cacheUserTable, map[string]string{"email": e}); err == nil && len(recs) > 0 {
		if u, err := cache.MapToEntity[Entities.User](recs[0]); err == nil {
			return u, nil
		}
	}
	u, err := r.robleRepository.GetUserByEmail(ctx, email)
	if err != nil || u == nil {
		return u, err
	}
	if rec, mapErr := cache.EntityToMap(u); mapErr == nil {
		r.cacheAdapter.Save(cacheUserTable, "insert", rec)
	}
	return u, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*Entities.User, error) {
	name := strings.TrimSpace(username)
	if recs, err := r.cacheAdapter.FindBy(cacheUserTable, map[string]string{"username": name}); err == nil && len(recs) > 0 {
		if u, err := cache.MapToEntity[Entities.User](recs[0]); err == nil {
			return u, nil
		}
	}
	u, err := r.robleRepository.GetUserByUsername(ctx, username)
	if err != nil || u == nil {
		return u, err
	}
	if rec, mapErr := cache.EntityToMap(u); mapErr == nil {
		r.cacheAdapter.Save(cacheUserTable, "insert", rec)
	}
	return u, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	u, err := r.GetUserByEmail(ctx, email)
	return u != nil, err
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	u, err := r.GetUserByUsername(ctx, username)
	return u != nil, err
}

func (r *UserRepository) ExistsByID(ctx context.Context, userID string) (bool, error) {
	u, err := r.GetUserByID(ctx, userID)
	return u != nil, err
}
