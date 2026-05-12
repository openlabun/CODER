package cache_syncronizer

import (
	"encoding/json"
	"fmt"
	"context"
	"strings"

	cache_publisher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"

	services "github.com/openlabun/CODER/apps/api_v2/internal/application/services"

	CourseEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	ExamEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	SubmissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	UserEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
)

func MapMessageToTask(message cache_publisher.SyncMessage) (SyncTask, error) {
	data, err := mapToModel(message.Model, message.Data)
	if err != nil {
		return SyncTask{}, err
	}

	return SyncTask{
		Model: message.Model,
		Request: message.Request,
		Data:  data,
	}, nil
}

func mapToModel(model string, data json.RawMessage) (any, error) {
	switch model {
	case cache_infrastructure.UserEntity:
		var e UserEntities.User
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.CourseEntity:
		var e CourseEntities.Course
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.ExamEntity:
		var e ExamEntities.Exam
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.ChallengeEntity:
		var e ExamEntities.Challenge
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.ExamItemEntity:
		var e ExamEntities.ExamItem
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.ExamItemScoreEntity:
		var e ExamEntities.ExamItemScore
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.ExamScoreEntity:
		var e ExamEntities.ExamScore
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.TestCaseEntity:
		var e ExamEntities.TestCase
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.IOVariableEntity:
		var e ExamEntities.IOVariable
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.SessionEntity:
		var e SubmissionEntities.Session
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.SubmissionEntity:
		var e SubmissionEntities.Submission
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case cache_infrastructure.SubmissionResultEntity:
		var e SubmissionEntities.SubmissionResult
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	default:
		return nil, fmt.Errorf("unknown model: %s", model)
	}
}

func BuildRequestContext(userEmail, token string) context.Context {
	ctx := context.Background()

	clean_token := strings.TrimSpace(token)
	if strings.HasPrefix(strings.ToLower(clean_token), "bearer ") {
		if clean_token != "" {
			ctx = services.WithAccessToken(ctx, clean_token)
		}
	}

	email := strings.TrimSpace(userEmail)
	if email != "" {
		ctx = services.WithUserEmail(ctx, email)
	}

	return ctx
}

