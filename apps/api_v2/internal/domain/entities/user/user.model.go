package user_entities

import "time"

type UserRole string

const (
	UserRoleStudent   UserRole = "student"
	UserRoleProfessor UserRole = "professor"
	UserRoleAdmin     UserRole = "admin"
)

type User struct {
	ID           string			`json:"id"`
	Username     string			`json:"username"`
	Email        string			`json:"email"`
	Role         UserRole		`json:"role"`

	//Metadata
	CreatedAt      time.Time	`json:"created_at"`
	LastConnection time.Time	`json:"last_connection"`
}
