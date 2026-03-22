package user_ports

import (
	"context"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
)

// UserServicePort abstracts user persistence/query operations needed by user use cases.
type UserServicePort interface {
	GetUserByID(ctx context.Context, userID string) (*Entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*Entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*Entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

// PasswordHasherPort abstracts password hashing/verification.
type PasswordHasherPort interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) (bool, error)
}

// TokenServicePort abstracts token generation/refresh logic.
type TokenServicePort interface {
	RefreshUserToken(refreshToken string) (*dtos.Token, error)
}
