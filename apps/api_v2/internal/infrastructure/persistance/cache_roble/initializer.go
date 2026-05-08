package cache_roble_infrastructure

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	challenge_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	course_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	submission_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	user_entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
)

type CacheTableConfig struct {
	Entity      any
	TTL         time.Duration
	CompositePK []string // set when entity has no single "id" field (e.g. join tables)
}

type SchemaInitializer struct {
	client  *MariaDBClient
	configs []CacheTableConfig
}

func NewSchemaInitializer(client *MariaDBClient, configs []CacheTableConfig) *SchemaInitializer {
	return &SchemaInitializer{client: client, configs: configs}
}

func (s *SchemaInitializer) Initialize() error {
	for _, cfg := range s.configs {
		t := reflect.TypeOf(cfg.Entity)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		stmt, err := buildCreateTableSQL(t, cfg.CompositePK)
		if err != nil {
			return fmt.Errorf("build schema for %s: %w", t.Name(), err)
		}
		if _, err := s.client.DB().Exec(stmt); err != nil {
			return fmt.Errorf("create table cache_%s: %w", toSnakeCase(t.Name()), err)
		}
	}
	return nil
}

func DefaultEntityConfigs() []CacheTableConfig {
	return []CacheTableConfig{
		{Entity: challenge_entities.Exam{}, TTL: 5 * time.Minute},
		{Entity: challenge_entities.ExamItem{}, TTL: 5 * time.Minute},
		{Entity: challenge_entities.ExamScore{}, TTL: 30 * time.Second},
		{Entity: challenge_entities.ExamItemScore{}, TTL: 30 * time.Second},
		{Entity: challenge_entities.Challenge{}, TTL: 5 * time.Minute},
		{Entity: challenge_entities.TestCase{}, TTL: 5 * time.Minute},
		{Entity: user_entities.User{}, TTL: 60 * time.Second},
		{Entity: course_entities.Course{}, TTL: 5 * time.Minute},
		{Entity: course_entities.CourseStudent{}, TTL: 30 * time.Second, CompositePK: []string{"course_id", "student_id"}},
		{Entity: submission_entities.Submission{}, TTL: 30 * time.Second},
		{Entity: submission_entities.Session{}, TTL: 10 * time.Second},
		{Entity: submission_entities.SubmissionResult{}, TTL: 30 * time.Second},
	}
}

var timeType = reflect.TypeOf(time.Time{})

func buildCreateTableSQL(t reflect.Type, compositePK []string) (string, error) {
	tableName := "cache_" + toSnakeCase(t.Name())
	useInlinePK := len(compositePK) == 0 && hasFieldWithJSONTag(t, "id")

	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		colName := columnName(field)
		if colName == "-" || colName == "" {
			continue
		}

		isPK := (colName == "id" && useInlinePK) || sliceContains(compositePK, colName)
		var colType string
		if isPK {
			colType = "VARCHAR(255)"
		} else {
			var err error
			colType, err = goTypeToSQL(field.Type)
			if err != nil {
				return "", fmt.Errorf("field %s: %w", field.Name, err)
			}
		}

		col := fmt.Sprintf("`%s` %s", colName, colType)
		if colName == "id" && useInlinePK {
			col += " PRIMARY KEY"
		} else if isPK {
			col += " NOT NULL"
		} else if !isNullable(field.Type) {
			col += " NOT NULL"
		}
		columns = append(columns, col)
	}

	columns = append(columns,
		"`cached_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP",
		"`synced` BOOLEAN NOT NULL DEFAULT FALSE",
		"`operation` ENUM('insert','update') NOT NULL DEFAULT 'insert'",
	)

	if len(compositePK) > 0 {
		quoted := make([]string, len(compositePK))
		for i, k := range compositePK {
			quoted[i] = "`" + k + "`"
		}
		columns = append(columns, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(quoted, ", ")))
	}

	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s` (\n  %s\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
		tableName,
		strings.Join(columns, ",\n  "),
	), nil
}

func goTypeToSQL(t reflect.Type) (string, error) {
	if t.Kind() == reflect.Ptr {
		return goTypeToSQL(t.Elem())
	}
	if t == timeType {
		return "DATETIME", nil
	}
	switch t.Kind() {
	case reflect.String:
		return "TEXT", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "INT", nil
	case reflect.Int64:
		return "BIGINT", nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "INT UNSIGNED", nil
	case reflect.Uint64:
		return "BIGINT UNSIGNED", nil
	case reflect.Bool:
		return "BOOLEAN", nil
	case reflect.Float32, reflect.Float64:
		return "DOUBLE", nil
	case reflect.Slice, reflect.Struct, reflect.Map:
		return "JSON", nil
	default:
		return "", fmt.Errorf("unsupported type kind %s", t.Kind())
	}
}

func isNullable(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice || t.Kind() == reflect.Map
}

func columnName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		return strings.ToLower(field.Name)
	}
	name, _, _ := strings.Cut(tag, ",")
	return name
}

func hasFieldWithJSONTag(t reflect.Type, tag string) bool {
	for i := 0; i < t.NumField(); i++ {
		if columnName(t.Field(i)) == tag {
			return true
		}
	}
	return false
}

func sliceContains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
