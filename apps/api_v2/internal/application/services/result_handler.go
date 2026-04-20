package services

import (
	"strings"

	exam_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/exam"
	state_machine "github.com/openlabun/CODER/apps/api_v2/internal/domain/states/submission"

	submission_constants "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/submission"
	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	submissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
)

func CheckSubmissionResult(submission *submissionEntities.SubmissionResult, testCase *examEntities.TestCase) (*submissionEntities.SubmissionResult, error) {
	if compareOutputValues(submission.ActualOutput.Value, testCase.ExpectedOutput.Value, testCase.ExpectedOutput.Type) {
		if submission.ActualOutput.Type == testCase.ExpectedOutput.Type {
			state_machine.ApplyTransition(submission, submission_constants.SubmissionStatusAccepted)
			return submission, nil
		}
	} 
	
	state_machine.ApplyTransition(submission, submission_constants.SubmissionStatusWrongAnswer)
	return submission, nil
}

func compareOutputValues(actual, expected string, valueType exam_constants.VariableFormat) bool {
	switch valueType {
	case exam_constants.VariableFormatArray:
		return normalizeArrayText(actual) == normalizeArrayText(expected)
	case exam_constants.VariableFormatBoolean:
		return normalizeBooleanText(actual) == normalizeBooleanText(expected)
	default:
		return strings.TrimSpace(actual) == strings.TrimSpace(expected)
	}
}

func normalizeArrayText(value string) string {
	// Treat equivalent array outputs with different spaces as equal, e.g. "[1,2,3]" and "[1, 2, 3]".
	return strings.Join(strings.Fields(strings.TrimSpace(value)), "")
}

func normalizeBooleanText(value string) string {
	v := strings.TrimSpace(strings.ToLower(value))
	switch v {
	case "1":
		return "true"
	case "0":
		return "false"
	default:
		return v
	}
}