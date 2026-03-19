package user_entities

import "time"

type UserRole string

const (
	UserRoleStudent   UserRole = "student"
	UserRoleProfessor UserRole = "professor"
	UserRoleAdmin     UserRole = "admin"
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	Role         UserRole

	//Metadata
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
