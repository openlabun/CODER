package exam_factory

import (
	"strings"
	"time"

	Entities "../../entities/exam"
	Validations "../../validations/exam"
)


func NewTestCase(
	id, name string,
	input []Entities.IOVariable,
	expectedOutput Entities.IOVariable,
	isSample bool,
	points int,
	challengeID string,
) (*Entities.TestCase, error) {
	testCase := &Entities.TestCase{
		ID:             strings.TrimSpace(id),
		Name:           strings.TrimSpace(name),
		Input:          input,
		ExpectedOutput: expectedOutput,
		IsSample:       isSample,
		Points:         points,
		CreatedAt:      time.Now(),
		ChallengeID:    strings.TrimSpace(challengeID),
	}

	if err := Validations.ValidateTestCase(testCase); err != nil {
		return nil, err
	}

	return testCase, nil
}

func ExistingTestCase(
	id, name string,
	input []Entities.IOVariable,
	expectedOutput Entities.IOVariable,
	isSample bool,
	points int,
	challengeID string,
	createdAt time.Time,
) (*Entities.TestCase, error) {
	testCase := &Entities.TestCase{
		ID:             strings.TrimSpace(id),
		Name:           strings.TrimSpace(name),
		Input:          input,
		ExpectedOutput: expectedOutput,
		IsSample:       isSample,
		Points:         points,
		CreatedAt:      createdAt,
		ChallengeID:    strings.TrimSpace(challengeID),
	}

	if err := Validations.ValidateTestCase(testCase); err != nil {
		return nil, err
	}

	return testCase, nil
}
