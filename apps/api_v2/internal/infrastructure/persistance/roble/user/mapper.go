package roble_infrastructure

import (
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	UserFactory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/user"
)

func UserToRecord(user *Entities.User) map[string]any {
	return map[string]any{
		"ID":             strings.TrimSpace(user.ID),
		"Username":       strings.TrimSpace(user.Username),
		"Email":          strings.ToLower(strings.TrimSpace(user.Email)),
		"Role":           string(user.Role),
		"CreatedAt":      user.CreatedAt.UTC().Format(time.RFC3339),
		"LastConnection": user.LastConnection.UTC().Format(time.RFC3339),
	}
}

func RecordToUser(record map[string]any, connection bool) *Entities.User {
	time_created := time.Time{}
	if t, ok := asTime(record["CreatedAt"]); ok {
		time_created = t
	}

	time_last_connection := time.Time{}
	if t, ok := asTime(record["LastConnection"]); ok {
		time_last_connection = t
	}

	user, err := UserFactory.ExistingUser(
		asString(record["ID"]),
		asString(record["Username"]),
		asString(record["Email"]),
		Entities.UserRole(asString(record["Role"])),
		time_created,
		time_last_connection,
		connection,
	)

	if err != nil {
		return nil
	}

	return user
}
