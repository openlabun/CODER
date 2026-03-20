package roble_infrastructure

import (
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
	infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/roble"
)

const userTableName = "UserModel"

type UserRepository struct {
	adapter *infrastructure.RobleDatabaseAdapter
}

func NewUserRepository(adapter *infrastructure.RobleDatabaseAdapter) *UserRepository {
	return &UserRepository{adapter: adapter}
}

func (r *UserRepository) SaveUser(user *Entities.User) (*Entities.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	record := UserToRecord(user)

	_, err := r.adapter.Insert(userTableName, []map[string]any{record})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(userID string) (*Entities.User, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("userID is required")
	}

	res, err := r.adapter.Read(userTableName, map[string]string{"ID": userID})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return RecordToUser(record, false), nil
}

func (r *UserRepository) GetUserByEmail(email string) (*Entities.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("email is required")
	}

	res, err := r.adapter.Read(userTableName, map[string]string{"Email": strings.ToLower(strings.TrimSpace(email))})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return RecordToUser(record, false), nil
}

func (r *UserRepository) GetUserByUsername(username string) (*Entities.User, error) {
	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("username is required")
	}

	res, err := r.adapter.Read(userTableName, map[string]string{"Username": strings.TrimSpace(username)})
	if err != nil {
		return nil, err
	}

	record, err := firstRecord(res)
	if err != nil {
		return nil, nil
	}

	return RecordToUser(record, false), nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	user, err := r.GetUserByEmail(email)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	user, err := r.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func (r *UserRepository) ExistsByID(userID string) (bool, error) {
	user, err := r.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func firstRecord(res map[string]any) (map[string]any, error) {
	if data, ok := res["data"]; ok {
		if arr := asRecordSlice(data); len(arr) > 0 {
			return arr[0], nil
		}
	}

	if records, ok := res["records"]; ok {
		if arr := asRecordSlice(records); len(arr) > 0 {
			return arr[0], nil
		}
	}

	if arr := asRecordSlice(res); len(arr) > 0 {
		return arr[0], nil
	}

	if isUserRecordShape(res) {
		return res, nil
	}

	return nil, fmt.Errorf("no records found")
}

func asRecordSlice(value any) []map[string]any {
	items, ok := value.([]any)
	if !ok {
		return nil
	}

	results := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if row, ok := item.(map[string]any); ok {
			results = append(results, row)
		}
	}

	return results
}

func isUserRecordShape(record map[string]any) bool {
	_, hasID := record["ID"]
	_, hasEmail := record["Email"]
	return hasID && hasEmail
}

func asString(v any) string {
	if v == nil {
		return ""
	}

	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}

	return strings.TrimSpace(fmt.Sprintf("%v", v))
}

func asTime(v any) (time.Time, bool) {
	s := asString(v)
	if s == "" {
		return time.Time{}, false
	}

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, true
	}

	return time.Time{}, false
}
