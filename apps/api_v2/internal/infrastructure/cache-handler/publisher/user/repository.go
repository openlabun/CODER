package cache_publisher

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	UserRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	PublisherAdapter "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
)

type UserRepository struct {
	publisherAdapter *PublisherAdapter.CachePublisherAdapter
	repository       UserRepositories.UserRepository
}

func NewUserRepository(publisherAdapter *PublisherAdapter.CachePublisherAdapter, repository UserRepositories.UserRepository) *UserRepository {
	return &UserRepository{publisherAdapter: publisherAdapter, repository: repository}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *Entities.User) (*Entities.User, error) {
	message, err := PublisherAdapter.MapUserToSyncMessage(user, cache_infrastructure.InsertOperation)
	if err != nil {
		return nil, err
	}
	if err := r.publisherAdapter.Save(message); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*Entities.User, error) {
	return r.repository.GetUserByID(ctx, userID)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*Entities.User, error) {
	return r.repository.GetUserByEmail(ctx, email)
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*Entities.User, error) {
	return r.repository.GetUserByUsername(ctx, username)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.repository.ExistsByEmail(ctx, email)
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.repository.ExistsByUsername(ctx, username)
}
