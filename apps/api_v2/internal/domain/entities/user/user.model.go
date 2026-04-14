package user_entities

import (
	"time"

	user_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/user"
)

type User struct {
	ID           string						`json:"id"`
	Username     string						`json:"username"`
	Email        string						`json:"email"`
	Role         user_constants.UserRole	`json:"role"`

	//Metadata
	CreatedAt      time.Time				`json:"created_at"`
	LastConnection time.Time				`json:"last_connection"`
}
