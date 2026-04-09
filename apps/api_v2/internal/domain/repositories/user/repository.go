package user_repository

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
)

type UserRepository interface {
	// SaveUser persists a new user or updates an existing one.
	SaveUser(ctx context.Context, user *Entities.User) (*Entities.User, error)

	GetUserByID(ctx context.Context, userID string) (*Entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*Entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*Entities.User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
