package roble_infrastructure

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	exam_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/exam"
)

// ExamItem mappers
func examItemToRecord(item *Entities.ExamItem) map[string]any {
	return map[string]any{
		"ID":          strings.TrimSpace(item.ID),
		"ChallengeID": strings.TrimSpace(item.ChallengeID),
		"ExamID":      strings.TrimSpace(item.ExamID),
		"Order":       item.Order,
		"Points":      item.Points,
	}
}

func examItemToUpdates(item *Entities.ExamItem) map[string]any {
	return map[string]any{
		"ChallengeID": strings.TrimSpace(item.ChallengeID),
		"ExamID":      strings.TrimSpace(item.ExamID),
		"Order":       item.Order,
		"Points":      item.Points,
	}
}

func recordToExamItem(record map[string]any) (*Entities.ExamItem, error) {
	return exam_factory.ExistingExamItem(
		asString(record["ID"]),
		asString(record["ChallengeID"]),
		asString(record["ExamID"]),
		asInt(record["Order"]),
		asInt(record["Points"]),
	)
}

func examToRecord(exam *Entities.Exam) map[string]any {
	record := map[string]any{
		"ID":                   strings.TrimSpace(exam.ID),
		"Title":                strings.TrimSpace(exam.Title),
		"Visibility":           string(exam.Visibility),
		"StartTime":            exam.StartTime.UTC().Format(time.RFC3339),
		"AllowLateSubmissions": exam.AllowLateSubmissions,
		"TimeLimit":            exam.TimeLimit,
		"TryLimit":             exam.TryLimit,
		"CreatedAt":            exam.CreatedAt.UTC().Format(time.RFC3339),
		"UpdatedAt":            exam.UpdatedAt.UTC().Format(time.RFC3339),
		"ProfessorID":          strings.TrimSpace(exam.ProfessorID),
	}

	if description := strings.TrimSpace(exam.Description); description != "" {
		record["Description"] = description
	}

	if courseID := exam.CourseID; courseID != nil && strings.TrimSpace(*courseID) != "" {
		record["CourseID"] = strings.TrimSpace(*courseID)
	}

	if exam.EndTime != nil && !exam.EndTime.IsZero() {
		record["EndTime"] = exam.EndTime.UTC().Format(time.RFC3339)
	}

	return record
}

func examToUpdates(exam *Entities.Exam) map[string]any {
	updates := map[string]any{
		"Title":                strings.TrimSpace(exam.Title),
		"Description":          strings.TrimSpace(exam.Description),
		"Visibility":           string(exam.Visibility),
		"StartTime":            exam.StartTime.UTC().Format(time.RFC3339),
		"AllowLateSubmissions": exam.AllowLateSubmissions,
		"TimeLimit":            exam.TimeLimit,
		"TryLimit":             exam.TryLimit,
		"ProfessorID":          strings.TrimSpace(exam.ProfessorID),
		"CourseID":             exam.CourseID,
	}

	if exam.EndTime != nil && !exam.EndTime.IsZero() {
		updates["EndTime"] = exam.EndTime.UTC().Format(time.RFC3339)
	} else {
		updates["EndTime"] = nil
	}

	return updates
}

func recordToExam(record map[string]any) (*Entities.Exam, error) {
	startTime, _ := asTime(record["StartTime"])
	if startTime.IsZero() {
		startTime = time.Now()
	}

	var endTime *time.Time
	if parsedEndTime, ok := asTime(record["EndTime"]); ok {
		endTime = &parsedEndTime
	}

	var courseID *string
	if rawCourseID, ok := record["CourseID"]; ok {
		if s := asString(rawCourseID); s != "" {
			courseID = &s
		}
	}

	createdAt, _ := asTime(record["CreatedAt"])
	updatedAt, _ := asTime(record["UpdatedAt"])

	return exam_factory.ExistingExam(
		asString(record["ID"]),
		asString(record["Title"]),
		asString(record["Description"]),
		Entities.Visibility(asString(record["Visibility"])),
		startTime,
		endTime,
		asBool(record["AllowLateSubmissions"]),
		asInt(record["TimeLimit"]),
		asInt(record["TryLimit"]),
		asString(record["ProfessorID"]),
		courseID,
		createdAt,
		updatedAt,
	)
}

func challengeToRecord(challenge *Entities.Challenge) map[string]any {
	return map[string]any{
		"ID":                strings.TrimSpace(challenge.ID),
		"Title":             strings.TrimSpace(challenge.Title),
		"Description":       strings.TrimSpace(challenge.Description),
		"Tags":              listFieldValue(challenge.Tags),
		"Status":            string(challenge.Status),
		"Difficulty":        string(challenge.Difficulty),
		"WorkerTimeLimit":   challenge.WorkerTimeLimit,
		"WorkerMemoryLimit": challenge.WorkerMemoryLimit,
		"InputVariables":             listFieldValue(ioVariableIDs(challenge.InputVariables)),
		"OutputVariable":            strings.TrimSpace(challenge.OutputVariable.ID),
		"Constraints":       strings.TrimSpace(challenge.Constraints),
		"CreatedAt":         challenge.CreatedAt.UTC().Format(time.RFC3339),
		"UpdatedAt":         challenge.UpdatedAt.UTC().Format(time.RFC3339),
		"UserID":            strings.TrimSpace(challenge.UserID),
	}
}

func challengeToUpdates(challenge *Entities.Challenge) map[string]any {
	return map[string]any{
		"Title":             strings.TrimSpace(challenge.Title),
		"Description":       strings.TrimSpace(challenge.Description),
		"Tags":              listFieldValue(challenge.Tags),
		"Status":            string(challenge.Status),
		"Difficulty":        string(challenge.Difficulty),
		"WorkerTimeLimit":   challenge.WorkerTimeLimit,
		"WorkerMemoryLimit": challenge.WorkerMemoryLimit,
		"InputVariables":             listFieldValue(ioVariableIDs(challenge.InputVariables)),
		"OutputVariable":            strings.TrimSpace(challenge.OutputVariable.ID),
		"Constraints":       strings.TrimSpace(challenge.Constraints),
		"UserID":            strings.TrimSpace(challenge.UserID),
	}
}

func recordToChallenge(record map[string]any, inputVariables []Entities.IOVariable, outputVariable *Entities.IOVariable) (*Entities.Challenge, error) {
	createdAt, _ := asTime(record["CreatedAt"])
	updatedAt, _ := asTime(record["UpdatedAt"])

	tags := asStringList(record["Tags"])
	status := Entities.ChallengeStatus(asString(record["Status"]))
	if status == "" {
		status = Entities.ChallengeStatusDraft
	}

	difficulty := Entities.ChallengeDifficulty(asString(record["Difficulty"]))
	if difficulty == "" {
		difficulty = Entities.ChallengeDifficultyEasy
	}

	output := Entities.IOVariable{}
	if outputVariable != nil {
		output = *outputVariable
	}

	return exam_factory.ExistingChallenge(
		asString(record["ID"]),
		asString(record["Title"]),
		asString(record["Description"]),
		tags,
		status,
		difficulty,
		asInt(record["WorkerTimeLimit"]),
		asInt(record["WorkerMemoryLimit"]),
		inputVariables,
		output,
		asString(record["Constraints"]),
		asString(record["UserID"]),
		createdAt,
		updatedAt,
	)
}

func testCaseToRecord(testCase *Entities.TestCase) map[string]any {
	return map[string]any{
		"ID":          strings.TrimSpace(testCase.ID),
		"Name":        strings.TrimSpace(testCase.Name),
		"Input":       listFieldValue(ioVariableIDs(testCase.Input)),
		"ExpectedOutput":      strings.TrimSpace(testCase.ExpectedOutput.ID),
		"IsSample":    testCase.IsSample,
		"Points":      testCase.Points,
		"CreatedAt":   testCase.CreatedAt.UTC().Format(time.RFC3339),
		"ChallengeID": strings.TrimSpace(testCase.ChallengeID),
	}
}

func testCaseToUpdates(testCase *Entities.TestCase) map[string]any {
	return map[string]any{
		"Name":        strings.TrimSpace(testCase.Name),
		"Input":       listFieldValue(ioVariableIDs(testCase.Input)),
		"ExpectedOutput":      strings.TrimSpace(testCase.ExpectedOutput.ID),
		"IsSample":    testCase.IsSample,
		"Points":      testCase.Points,
		"ChallengeID": strings.TrimSpace(testCase.ChallengeID),
	}
}

func recordToTestCase(record map[string]any, inputVariables []Entities.IOVariable, outputVariable *Entities.IOVariable) (*Entities.TestCase, error) {
	createdAt, _ := asTime(record["CreatedAt"])

	output := Entities.IOVariable{}
	if outputVariable != nil {
		output = *outputVariable
	}

	return exam_factory.ExistingTestCase(
		asString(record["ID"]),
		asString(record["Name"]),
		inputVariables,
		output,
		asBool(record["IsSample"]),
		asInt(record["Points"]),
		asString(record["ChallengeID"]),
		createdAt,
	)
}

func ioVariableToRecord(variable Entities.IOVariable) map[string]any {
	return map[string]any{
		"ID":    strings.TrimSpace(variable.ID),
		"Name":  strings.TrimSpace(variable.Name),
		"Type":  string(variable.Type),
		"Value": strings.TrimSpace(variable.Value),
	}
}

func ioVariableToUpdates(variable Entities.IOVariable) map[string]any {
	return map[string]any{
		"Name":  strings.TrimSpace(variable.Name),
		"Type":  string(variable.Type),
		"Value": strings.TrimSpace(variable.Value),
	}
}

func recordToIOVariable(record map[string]any) (*Entities.IOVariable, error) {
	return exam_factory.ExistingIOVariable(
		asString(record["ID"]),
		asString(record["Name"]),
		Entities.VariableFormat(asString(record["Type"])),
		asString(record["Value"]),
	)
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

	if isGenericRecordShape(res) {
		return res, nil
	}

	return nil, fmt.Errorf("no records found")
}

func extractRecords(res map[string]any) []map[string]any {
	if data, ok := res["data"]; ok {
		if arr := asRecordSlice(data); len(arr) > 0 {
			return arr
		}
	}

	if records, ok := res["records"]; ok {
		if arr := asRecordSlice(records); len(arr) > 0 {
			return arr
		}
	}

	if arr := asRecordSlice(res); len(arr) > 0 {
		return arr
	}

	if isGenericRecordShape(res) {
		return []map[string]any{res}
	}

	return nil
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

func isGenericRecordShape(record map[string]any) bool {
	_, hasID := record["ID"]
	return hasID
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

func asInt(v any) int {
	switch value := v.(type) {
	case int:
		return value
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float32:
		return int(value)
	case float64:
		return int(value)
	case string:
		var parsed int
		_, err := fmt.Sscanf(strings.TrimSpace(value), "%d", &parsed)
		if err == nil {
			return parsed
		}
	}

	return 0
}

func asBool(v any) bool {
	switch value := v.(type) {
	case bool:
		return value
	case string:
		normalized := strings.ToLower(strings.TrimSpace(value))
		return normalized == "true" || normalized == "1" || normalized == "yes"
	case int:
		return value != 0
	case int32:
		return value != 0
	case int64:
		return value != 0
	case float64:
		return value != 0
	}

	return false
}

func asTime(v any) (time.Time, bool) {
	s := asString(v)
	if s == "" {
		return time.Time{}, false
	}

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}
	if t, err := time.Parse("2006-01-02T15:04", s); err == nil {
		return t, true
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, true
	}

	return time.Time{}, false
}

func asStringList(v any) []string {
	if v == nil {
		return nil
	}

	switch value := v.(type) {
	case map[string]any:
		if wrapped, ok := value["values"]; ok {
			return asStringList(wrapped)
		}
		return nil
	case []string:
		out := make([]string, 0, len(value))
		for _, s := range value {
			if trimmed := strings.TrimSpace(s); trimmed != "" {
				out = append(out, trimmed)
			}
		}
		return out
	case []any:
		out := make([]string, 0, len(value))
		for _, item := range value {
			if s := asString(item); s != "" {
				out = append(out, s)
			}
		}
		return out
	case string:
		raw := strings.TrimSpace(value)
		if raw == "" {
			return nil
		}

		if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
			var parsed []string
			if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
				return asStringList(parsed)
			}

			var parsedAny []any
			if err := json.Unmarshal([]byte(raw), &parsedAny); err == nil {
				return asStringList(parsedAny)
			}
		}

		parts := strings.Split(raw, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				out = append(out, trimmed)
			}
		}
		return out
	}

	return nil
}

func listFieldValue(values []string) map[string]any {
	return map[string]any{"values": values}
}
