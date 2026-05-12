package cache_roble_infrastructure

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

type CacheAdapter struct {
	db *sql.DB
	initializer *SchemaInitializer
}

func NewCacheAdapter(client *MariaDBClient) *CacheAdapter {
	initializer := NewSchemaInitializer(client, DefaultEntityConfigs())
	if err := initializer.Initialize(); err != nil {
		panic(fmt.Sprintf("[cache-adapter] failed create cache schema initializer: %v", err))
	}

	if err := initializer.Initialize(); err != nil {
		panic(fmt.Sprintf("[cache-adapter] failed to initialize cache schema: %v", err))
	}

	return &CacheAdapter{
		db: client.DB(), 
		initializer: initializer,
	}
}

func (a *CacheAdapter) Save(tableName, operation string, record map[string]any) error {
	if len(record) == 0 {
		return fmt.Errorf("empty record")
	}

	keys := sortedKeys(record)
	cols := make([]string, 0, len(keys)+3)
	placeholders := make([]string, 0, len(keys)+3)
	args := make([]any, 0, len(keys)+3)
	updates := make([]string, 0, len(keys)+2)

	for _, k := range keys {
		cols = append(cols, "`"+k+"`")
		placeholders = append(placeholders, "?")
		args = append(args, record[k])
		if k != "id" {
			updates = append(updates, "`"+k+"` = VALUES(`"+k+"`)")
		}
	}

	cols = append(cols, "`cached_at`", "`synced`", "`operation`")
	placeholders = append(placeholders, "?", "?", "?")
	args = append(args, time.Now(), false, operation)
	updates = append(updates,
		"`cached_at` = VALUES(`cached_at`)",
		"`synced` = VALUES(`synced`)",
		"`operation` = VALUES(`operation`)",
	)

	query := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		tableName,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(updates, ", "),
	)

	_, err := a.db.Exec(query, args...)
	return err
}

func (a *CacheAdapter) MarkSynced(tableName, id string) error {
	_, err := a.db.Exec(
		"UPDATE `"+tableName+"` SET `synced` = TRUE WHERE `id` = ?",
		id,
	)
	return err
}

func (a *CacheAdapter) FindByID(tableName, id string) (map[string]any, error) {
	rows, err := a.db.Query("SELECT * FROM `"+tableName+"` WHERE `id` = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records, err := scanRows(rows)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	go a.db.Exec("UPDATE `"+tableName+"` SET `cached_at` = NOW() WHERE `id` = ?", id)

	return records[0], nil
}

func (a *CacheAdapter) FindBy(tableName string, conditions map[string]string) ([]map[string]any, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("conditions required")
	}

	keys := sortedStringKeys(conditions)
	wheres := make([]string, len(keys))
	args := make([]any, len(keys))
	for i, k := range keys {
		wheres[i] = "`" + k + "` = ?"
		args[i] = conditions[k]
	}

	rows, err := a.db.Query(
		"SELECT * FROM `"+tableName+"` WHERE "+strings.Join(wheres, " AND "),
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records, err := scanRows(rows)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}

	go func() {
		upd := "UPDATE `" + tableName + "` SET `cached_at` = NOW() WHERE " + strings.Join(wheres, " AND ")
		a.db.Exec(upd, args...)
	}()

	return records, nil
}

func (a *CacheAdapter) DeleteByID(tableName, id string) error {
	_, err := a.db.Exec("DELETE FROM `"+tableName+"` WHERE `id` = ?", id)
	return err
}

func (a *CacheAdapter) DeleteBy(tableName string, conditions map[string]string) error {
	if len(conditions) == 0 {
		return fmt.Errorf("conditions required")
	}

	keys := sortedStringKeys(conditions)
	wheres := make([]string, len(keys))
	args := make([]any, len(keys))
	for i, k := range keys {
		wheres[i] = "`" + k + "` = ?"
		args[i] = conditions[k]
	}

	_, err := a.db.Exec(
		"DELETE FROM `"+tableName+"` WHERE "+strings.Join(wheres, " AND "),
		args...,
	)
	return err
}

// EntityToMap converts a domain entity to a map for cache storage.
// Preserves native Go types (time.Time, bool, int) so the SQL driver handles them correctly.
// Slices, maps, and non-time structs are JSON-encoded to strings.
func EntityToMap(entity any) (map[string]any, error) {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("nil entity")
		}
		v = v.Elem()
	}
	t := v.Type()

	result := make(map[string]any, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		col := columnName(field)
		if col == "-" || col == "" {
			continue
		}
		result[col] = marshalFieldValue(v.Field(i))
	}
	return result, nil
}

// MapToEntity deserializes a cache record (returned by scanRows) into a domain entity via JSON.
func MapToEntity[T any](record map[string]any) (*T, error) {
	b, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	var entity T
	return &entity, json.Unmarshal(b, &entity)
}

// NormalizeBools converts int64 0/1 values to bool for the given columns.
// Required because MariaDB stores BOOLEAN as TINYINT(1), scanned by the driver as int64.
func NormalizeBools(record map[string]any, cols ...string) {
	for _, col := range cols {
		if v, ok := record[col]; ok {
			switch val := v.(type) {
			case int64:
				record[col] = val != 0
			case []byte:
				record[col] = string(val) == "1"
			}
		}
	}
}

func marshalFieldValue(v reflect.Value) any {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Struct:
		if v.Type() == timeType {
			return v.Interface()
		}
		b, _ := json.Marshal(v.Interface())
		return string(b)
	case reflect.Slice, reflect.Map:
		if v.IsNil() {
			return nil
		}
		b, _ := json.Marshal(v.Interface())
		return string(b)
	}
	return v.Interface()
}

func scanRows(rows *sql.Rows) ([]map[string]any, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]any
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(cols))
		for i, col := range cols {
			switch vv := vals[i].(type) {
			case []byte:
				if len(vv) > 0 && (vv[0] == '[' || vv[0] == '{') {
					var parsed any
					if jsonErr := json.Unmarshal(vv, &parsed); jsonErr == nil {
						row[col] = parsed
						continue
					}
				}
				row[col] = string(vv)
			default:
				row[col] = vv
			}
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
