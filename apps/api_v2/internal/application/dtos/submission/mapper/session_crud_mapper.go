package mapper

import (
	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/submission"

	factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/submission"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	examEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
)

func MapCreateSessionInputToSessionRecord(input dtos.CreateSessionInput, exam *examEntities.Exam) (*Entities.Session, error) {
	session, err := factory.NewSession(
		input.UserID,
		exam,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}
