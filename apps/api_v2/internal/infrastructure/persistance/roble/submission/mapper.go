package roble_infrastructure

import (
	"fmt"
	"strings"
	"time"

	submission_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	ExamEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	submission_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/submission"
)

const (
	submissionTableName       = "Submission"
	sessionTableName          = "Sessions"
	submissionResultTableName = "SubmissionResult"
)

func submissionToRecord(submission *Entities.Submission) map[string]any {
	record := map[string]any{
		"ID":          strings.TrimSpace(submission.ID),
		"Code":        submission.Code,
		"Language":    string(submission.Language),
		"Score":       submission.Score,
		"TimeMsTotal": submission.TimeMsTotal,
		"Scorable":    submission.Scorable,
		"CreatedAt":   submission.CreatedAt.UTC().Format(time.RFC3339),
		"UpdatedAt":   submission.UpdatedAt.UTC().Format(time.RFC3339),
		"ChallengeID": strings.TrimSpace(submission.ChallengeID),
		"SessionID":   strings.TrimSpace(submission.SessionID),
		"UserID":      strings.TrimSpace(submission.UserID),
	}

	if submission.ExamItemScoreID != nil {
		if examItemScoreID := strings.TrimSpace(*submission.ExamItemScoreID); examItemScoreID != "" {
			record["ExamItemScoreID"] = examItemScoreID
		}
	}

	return record
}

func submissionToUpdates(submission *Entities.Submission) map[string]any {
	updates := map[string]any{
		"Code":        submission.Code,
		"Language":    string(submission.Language),
		"Score":       submission.Score,
		"TimeMsTotal": submission.TimeMsTotal,
		"Scorable":    submission.Scorable,
		"ChallengeID": strings.TrimSpace(submission.ChallengeID),
		"SessionID":   strings.TrimSpace(submission.SessionID),
		"UserID":      strings.TrimSpace(submission.UserID),
	}

	if submission.ExamItemScoreID != nil {
		if examItemScoreID := strings.TrimSpace(*submission.ExamItemScoreID); examItemScoreID != "" {
			updates["ExamItemScoreID"] = examItemScoreID
		}
	} else {
		updates["ExamItemScoreID"] = nil
	}

	return updates
}

func recordToSubmission(record map[string]any) (*Entities.Submission, error) {
	createdAt, _ := asTime(record["CreatedAt"])
	updatedAt, _ := asTime(record["UpdatedAt"])

	var examItemScoreID *string
	if record["ExamItemScoreID"] != nil {
		if trimmedID := strings.TrimSpace(asString(record["ExamItemScoreID"])); trimmedID != "" {
			examItemScoreID = &trimmedID
		}
	}

	return submission_factory.ExistingSubmission(
		asString(record["ID"]),
		asString(record["Code"]),
		submission_constants.ProgrammingLanguage(asString(record["Language"])),
		asInt(record["Score"]),
		asInt(record["TimeMsTotal"]),
		asBool(record["Scorable"]),
		createdAt,
		updatedAt,
		asString(record["ChallengeID"]),
		asString(record["SessionID"]),
		asString(record["UserID"]),
		examItemScoreID,
	)
}

func sessionToRecord(session *Entities.Session) map[string]any {
	record := map[string]any{
		"ID":            strings.TrimSpace(session.ID),
		"StudentID":     strings.TrimSpace(session.StudentID),
		"ExamID":        strings.TrimSpace(session.ExamID),
		"Status":        string(session.Status),
		"Attempts":      session.Attempts,
		"TimeLeft":      session.TimeLeft,
		"StartedAt":     session.StartedAt.UTC().Format(time.RFC3339),
		"LastHeartbeat": session.LastHeartbeat.UTC().Format(time.RFC3339),
	}

	if session.FinishedAt != nil && !session.FinishedAt.IsZero() {
		record["FinishedAt"] = session.FinishedAt.UTC().Format(time.RFC3339)
	}

	return record
}

func sessionToUpdates(session *Entities.Session) map[string]any {
	updates := map[string]any{
		"StudentID":     strings.TrimSpace(session.StudentID),
		"ExamID":        strings.TrimSpace(session.ExamID),
		"Status":        string(session.Status),
		"Attempts":      session.Attempts,
		"TimeLeft":      session.TimeLeft,
		"StartedAt":     session.StartedAt.UTC().Format(time.RFC3339),
		"LastHeartbeat": session.LastHeartbeat.UTC().Format(time.RFC3339),
	}

	if session.FinishedAt != nil && !session.FinishedAt.IsZero() {
		updates["FinishedAt"] = session.FinishedAt.UTC().Format(time.RFC3339)
	} else {
		updates["FinishedAt"] = nil
	}

	return updates
}

func recordToSession(record map[string]any) (*Entities.Session, error) {
	startedAt, _ := asTime(record["StartedAt"])
	lastHeartbeat, _ := asTime(record["LastHeartbeat"])

	var finishedAt *time.Time
	if parsed, ok := asTime(record["FinishedAt"]); ok {
		finishedAt = &parsed
	}

	return submission_factory.ExistingSession(
		asString(record["ID"]),
		asString(record["StudentID"]),
		asString(record["ExamID"]),
		submission_constants.SessionStatus(asString(record["Status"])),
		asInt(record["Attempts"]),
		asInt(record["TimeLeft"]),
		startedAt,
		finishedAt,
		lastHeartbeat,
	)
}

func resultToRecord(result *Entities.SubmissionResult) map[string]any {
	record := map[string]any{
		"ID":           strings.TrimSpace(result.ID),
		"SubmissionID": strings.TrimSpace(result.SubmissionID),
		"TestCaseID":   strings.TrimSpace(result.TestCaseID),
		"Status":       string(result.Status),
	}

	if result.ActualOutput != nil {
		if actualOutputID := strings.TrimSpace(result.ActualOutput.ID); actualOutputID != "" {
			record["ActualOutput"] = actualOutputID
		}
	}

	if result.ErrorMessage != nil && strings.TrimSpace(*result.ErrorMessage) != "" {
		record["ErrorMessage"] = strings.TrimSpace(*result.ErrorMessage)
	}

	return record
}

func resultToUpdates(result *Entities.SubmissionResult) map[string]any {
	updates := map[string]any{
		"SubmissionID": strings.TrimSpace(result.SubmissionID),
		"TestCaseID":   strings.TrimSpace(result.TestCaseID),
		"Status":       string(result.Status),
	}

	if result.ActualOutput != nil {
		updates["ActualOutput"] = strings.TrimSpace(result.ActualOutput.ID)
	} else {
		updates["ActualOutput"] = nil
	}

	if result.ErrorMessage != nil && strings.TrimSpace(*result.ErrorMessage) != "" {
		updates["ErrorMessage"] = strings.TrimSpace(*result.ErrorMessage)
	} else {
		updates["ErrorMessage"] = nil
	}

	return updates
}

func recordToResult(record map[string]any, actualOutput *ExamEntities.IOVariable) (*Entities.SubmissionResult, error) {
	var errMsg *string
	if msg := strings.TrimSpace(asString(record["ErrorMessage"])); msg != "" {
		errMsg = &msg
	}

	return submission_factory.ExistingSubmissionResult(
		asString(record["ID"]),
		asString(record["SubmissionID"]),
		asString(record["TestCaseID"]),
		submission_constants.SubmissionStatus(asString(record["Status"])),
		actualOutput,
		errMsg,
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
		s := strings.ToLower(strings.TrimSpace(value))
		return s == "true" || s == "1" || s == "yes"
	case int:
		return value != 0
	case int32:
		return value != 0
	case int64:
		return value != 0
	case float32:
		return value != 0
	case float64:
		return value != 0
	default:
		return false
	}
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
