package cache_syncronizer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"

	CourseEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	ExamEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/exam"
	SubmissionEntities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/submission"
	CourseRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	ExamRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	SubmissionRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
	UserRepositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	
	cache_publisher "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/cache-handler/publisher"
	cache_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
	rabbitmq_infrastructure "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/publisher/rabbitmq"
)

type CacheSyncronizerAdapter struct {
	// Placeholder for any dependencies or configuration needed for the cache syncronizer.
	// - Publisher (RabbitMQ) for consuming
	publisher 	*rabbitmq_infrastructure.RabbitMQAdapter
	queueName 	string

	// - Auth (JWT) for modifying Roble database with Admin account
	authAdapter ports.LoginPort

	// - Repositories for interacting with the Roble database
	repositories  map[string]any

	// - Cache Adapter for clearing the cache
	cacheAdapter *cache_infrastructure.CacheAdapter

	// Admin credentials for authenticating with Roble database
	adminUsername string
	adminPassword string
}

func NewCacheSyncronizerAdapter(
	publisher 					*rabbitmq_infrastructure.RabbitMQAdapter,
	authAdapter 				ports.LoginPort,
	cacheAdapter 				*cache_infrastructure.CacheAdapter,
	userRepository 				UserRepositories.UserRepository,
	courseRepository 			CourseRepositories.CourseRepository,
	examRepository 				ExamRepositories.ExamRepository,
	challengeRepository 		ExamRepositories.ChallengeRepository,
	examItemRepository 			ExamRepositories.ExamItemRepository,
	examItemScoreRepository 	ExamRepositories.ExamItemScoreRepository,
	examScoreRepository 		ExamRepositories.ExamScoreRepository,
	testCaseRepository 			ExamRepositories.TestCaseRepository,
	submissionRepository 		SubmissionRepositories.SubmissionRepository,
	submissionResultRepository 	SubmissionRepositories.SubmissionResultRepository,
	sessionRepository 			SubmissionRepositories.SessionRepository,
) *CacheSyncronizerAdapter {
	adminUsername := os.Getenv("INTERNAL_USER_EMAIL")
	if adminUsername == "" {
		panic("INTERNAL_USER_EMAIL environment variable is not set")
	}

	adminPassword := os.Getenv("INTERNAL_USER_PASSWORD")
	if adminPassword == "" {
		panic("INTERNAL_USER_PASSWORD environment variable is not set")
	}

	queueName := os.Getenv("CACHE_SYNC_QUEUE_NAME")
	if queueName == "" {
		panic("CACHE_SYNC_QUEUE_NAME environment variable is not set")
	}

	var repositories = map[string]any{
		cache_infrastructure.UserEntity: 			userRepository,
		cache_infrastructure.CourseEntity: 			courseRepository,
		cache_infrastructure.ExamEntity: 			examRepository,
		cache_infrastructure.ChallengeEntity: 		challengeRepository,
		cache_infrastructure.ExamItemEntity: 		examItemRepository,
		cache_infrastructure.ExamItemScoreEntity: 	examItemScoreRepository,
		cache_infrastructure.ExamScoreEntity: 		examScoreRepository,
		cache_infrastructure.TestCaseEntity: 		testCaseRepository,
		cache_infrastructure.SubmissionEntity: 		submissionRepository,
		cache_infrastructure.SubmissionResultEntity:submissionResultRepository,
		cache_infrastructure.SessionEntity: 		sessionRepository,
	}

	return &CacheSyncronizerAdapter{
		// Initialize any dependencies or configuration here.
		publisher: publisher,
		queueName: queueName,
		authAdapter: authAdapter,
		repositories: repositories,
		cacheAdapter: cacheAdapter,
		adminUsername: adminUsername,
		adminPassword: adminPassword,
	}
}

func (a *CacheSyncronizerAdapter) Start() {
	msgs, err := a.publisher.ConsumeFromQueue(a.queueName)
	if err != nil {
		fmt.Printf("[cache-syncronizer] failed to start consumer: %v\n", err)
		return
	}
	go func() {
		for delivery := range msgs {
			if err := a.ConsumeMessage(delivery.Body); err != nil {
				fmt.Printf("[cache-syncronizer] ConsumeMessage error: %v\n", err)
				_ = delivery.Nack(false, true)
				continue
			}
			_ = delivery.Ack(false)
		}
	}()
}

func (a *CacheSyncronizerAdapter) ConsumeMessage(body []byte) error {
	// [STEP 1] Consume messages from RabbitMQ and unmarshal into SyncMessage struct
	var message cache_publisher.SyncMessage
	if err := json.Unmarshal(body, &message); err != nil {
		return fmt.Errorf("unmarshal sync message: %w", err)
	}

	// [STEP 2] Process message and map into domain objects
	task, err := MapMessageToTask(message)
	if err != nil {
		return fmt.Errorf("map message to task: %w", err)
	}

	// [STEP 3] Syncronize data inside Roble infrastructure
	result_id, err := a.syncronizeData(task)
	if err != nil {
		return fmt.Errorf("syncronize data: %w", err)
	}
	if result_id == nil {
		return fmt.Errorf("syncronize data: no result ID returned")
	}

	// [STEP 4] Get cache table name from model
	cacheTableName := cache_infrastructure.EntityToCacheTable[task.Model]
	if cacheTableName == "" {
		return fmt.Errorf("unknown model for cache table: %s", task.Model)
	}

	// [STEP 4] Clear the cache for the affected model and ID
	if err := a.cacheAdapter.MarkSynced(cacheTableName, *result_id); err != nil {
		return fmt.Errorf("[cache-syncronizer] Failed to mark cache as synced: %v", err)
	}

	return nil
}

func (a *CacheSyncronizerAdapter) syncronizeData (task SyncTask) (*string, error) {
	// [STEP 1] Authenticate with JWT to modify Roble database with Admin account and create context
	access, err := a.authAdapter.LoginUser(a.adminUsername, a.adminPassword)
	if err != nil {
		// Handle authentication error
		return nil, err
	}
	ctx := BuildRequestContext(a.adminUsername, access.Token.AccessToken)
	if ctx == nil {
		return nil, fmt.Errorf("[cache-syncronizer] Failed to get token with Admin user")
	}
	

	// [STEP 2] Get the appropriate repository based on the model
	repository := a.repositories[task.Model]
	if repository == nil {
		return nil, fmt.Errorf("[cache-syncronizer] Unknown model: %s", task.Model)
	}


	// [STEP 3] Execute update or create operation based on the request type
	result_id, err := a.executeTask(ctx, task, repository)
	if err != nil {
		return nil, fmt.Errorf("[cache-syncronizer] Failed to execute task: %v", err)
	}

	return &result_id, nil
}

func (a *CacheSyncronizerAdapter) executeTask(ctx context.Context, task SyncTask, repository any) (string, error) {
	var err error
	result_id := ""
	switch task.Model {
		case cache_infrastructure.CourseEntity:
			repo := repository.(CourseRepositories.CourseRepository)
			data := task.Data.(CourseEntities.Course)
			var result *CourseEntities.Course
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateCourse(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateCourse(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.ExamEntity:
			repo := repository.(ExamRepositories.ExamRepository)
			data := task.Data.(ExamEntities.Exam)
			var result *ExamEntities.Exam
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateExam(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateExam(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.ChallengeEntity:
			repo := repository.(ExamRepositories.ChallengeRepository)
			data := task.Data.(ExamEntities.Challenge)
			var result *ExamEntities.Challenge
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateChallenge(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateChallenge(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.ExamItemEntity:
			repo := repository.(ExamRepositories.ExamItemRepository)
			data := task.Data.(ExamEntities.ExamItem)
			var result *ExamEntities.ExamItem
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateExamItem(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateExamItem(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.ExamItemScoreEntity:
			repo := repository.(ExamRepositories.ExamItemScoreRepository)
			data := task.Data.(ExamEntities.ExamItemScore)
			var result *ExamEntities.ExamItemScore
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateExamItemScore(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateExamItemScore(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.ExamScoreEntity:
			repo := repository.(ExamRepositories.ExamScoreRepository)
			data := task.Data.(ExamEntities.ExamScore)
			var result *ExamEntities.ExamScore
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateExamScore(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateExamScore(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.TestCaseEntity:
			repo := repository.(ExamRepositories.TestCaseRepository)
			data := task.Data.(ExamEntities.TestCase)
			var result *ExamEntities.TestCase
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateTestCase(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateTestCase(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.IOVariableEntity:
			repo := repository.(ExamRepositories.IOVariableRepository)
			data := task.Data.(ExamEntities.IOVariable)
			var result *ExamEntities.IOVariable
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateIOVariable(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateIOVariable(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.SessionEntity:
			repo := repository.(SubmissionRepositories.SessionRepository)
			data := task.Data.(SubmissionEntities.Session)
			var result *SubmissionEntities.Session
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateSession(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateSession(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.SubmissionEntity:
			repo := repository.(SubmissionRepositories.SubmissionRepository)
			data := task.Data.(SubmissionEntities.Submission)
			var result *SubmissionEntities.Submission
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateSubmission(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateSubmission(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}

		case cache_infrastructure.SubmissionResultEntity:
			repo := repository.(SubmissionRepositories.SubmissionResultRepository)
			data := task.Data.(SubmissionEntities.SubmissionResult)
			var result *SubmissionEntities.SubmissionResult
			switch task.Request {
			case cache_infrastructure.UpdateOperation:
				result, err = repo.UpdateResult(ctx, &data)
			case cache_infrastructure.InsertOperation:
				result, err = repo.CreateResult(ctx, &data)
			}
			if err == nil {
				result_id = result.ID
			}
	}

	return result_id, err
}