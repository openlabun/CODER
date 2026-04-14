package container

import (
	"fmt"

	course_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/course"
	course_crud_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/course/crud"
	user_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/user"

	exam_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam"

	ai_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/ai"

	challenge_crud_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam/challenge_crud"
	exam_crud_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam/exam_crud"
	exam_item_crud_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam/exam_item_crud"
	test_case_crud_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam/test_case_crud"
	submission_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/submission"
	session_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/submission/session"
)

// UserUseCases holds all user-related use cases available in the application.
type UserUseCases struct {
	Register     *user_usecases.RegisterUseCase
	Login        *user_usecases.LoginUseCase
	GetData      *user_usecases.GetDataUseCase
	RefreshToken *user_usecases.RefreshTokenUseCase
}

type CourseUseCases struct {
	CreateCourse       *course_crud_usecases.CreateCourseUseCase
	UpdateCourse       *course_crud_usecases.UpdateCourseUseCase
	DeleteCourse       *course_crud_usecases.DeleteCourseUseCase
	GetCourseDetails   *course_crud_usecases.GetCourseDetailsUseCase
	GetEnrolledCourses *course_crud_usecases.GetEnrolledCoursesUseCase
	GetOwnedCourses    *course_crud_usecases.GetOwnedCoursesUseCase

	GetCourseStudents       *course_usecases.GetCourseStudentsUseCase
	EnrollInCourse          *course_usecases.EnrollInCourseUseCase
	RemoveStudentFromCourse *course_usecases.RemoveStudentFromCourseUseCase
}

type ChallengeUseCases struct {
	CreateChallenge     *challenge_crud_usecases.CreateChallengeUseCase
	UpdateChallenge     *challenge_crud_usecases.UpdateChallengeUseCase
	PublishChallenge    *challenge_crud_usecases.PublishChallengeUseCase
	ArchiveChallenge    *challenge_crud_usecases.ArchiveChallengeUseCase
	DeleteChallenge     *challenge_crud_usecases.DeleteChallengeUseCase
	GetChallengeDetails *challenge_crud_usecases.GetChallengeDetailsUseCase
	GetChallengesByUser *challenge_crud_usecases.GetChallengesByUserUseCase
	GetPublicChallenges *exam_usecases.GetPublicChallengesUseCase
	ForkChallenge       *exam_usecases.ForkChallengeUseCase
}

type TestCaseUseCases struct {
	CreateTestCase          *test_case_crud_usecases.CreateTestCaseUseCase
	UpdateTestCase          *test_case_crud_usecases.UpdateTestCaseUseCase
	DeleteTestCase          *test_case_crud_usecases.DeleteTestCaseUseCase
	GetTestCasesByChallenge *test_case_crud_usecases.GetTestCasesByChallengeUseCase
}

type ExamItemUseCases struct {
	CreateExamItem *exam_item_crud_usecases.CreateExamItemUseCase
	UpdateExamItem *exam_item_crud_usecases.UpdateExamItemUseCase
	DeleteExamItem *exam_item_crud_usecases.DeleteExamItemUseCase
}

type SubmissionUseCases struct {
	CreateSubmission        *submission_usecases.CreateSubmissionUseCase
	CreateSubmissionWithoutScore *submission_usecases.CreateSubmissionWithoutScoreUseCase
	CreateCustomSubmission  *submission_usecases.CreateCustomSubmissionUseCase
	GetSubmissionStatus     *submission_usecases.GetSubmissionStatusUseCase
	GetChallengeSubmissions *submission_usecases.GetChallengeSubmissionsUseCase
	GetUserSubmissions      *submission_usecases.GetUserSubmissionsUseCase
	GetSessionSubmissions   *submission_usecases.GetSessionSubmissionsUseCase
	UpdateResult            *submission_usecases.UpdateResultUseCase
}

type SessionUseCases struct {
	CreateSession    *session_usecases.CreateSessionUseCase
	GetActiveSession *session_usecases.GetActiveSessionUseCase
	HeartBeatSession *session_usecases.HeartBeatSessionUseCase
	BlockSession     *session_usecases.BlockSessionUseCase
	CloseSession     *session_usecases.CloseSessionUseCase
}

type ExamUseCases struct {
	CreateExam              *exam_crud_usecases.CreateExamUseCase
	UpdateExam              *exam_crud_usecases.UpdateExamUseCase
	CloseExam               *exam_crud_usecases.CloseExamUseCase
	DeleteExam              *exam_crud_usecases.DeleteExamUseCase
	GetExamDetails          *exam_crud_usecases.GetExamDetailsUseCase
	GetExamsByCourse        *exam_crud_usecases.GetExamsByCourseUseCase
	GetExamItems            *exam_crud_usecases.GetExamItemsUseCase
	GetPublicExams          *exam_usecases.GetPublicExamsUseCase
	GetCodeDefaultTemplates *exam_usecases.GetCodeDefaultTemplates
}

type AIUseCases struct {
	GenerateFullChallenge *ai_usecases.GenerateFullChallengeUseCase
	GenerateExam          *ai_usecases.GenerateExamUseCase
}

type Application struct {
	Dependencies       ApplicationDependencies
	UserModule         UserUseCases
	CourseModule       CourseUseCases
	ExamModule         ExamUseCases
	ChallengeModule    ChallengeUseCases
	TestCaseModule     TestCaseUseCases
	ExamItemModule     ExamItemUseCases
	SessionModule      SessionUseCases
	SubmissionUseCases SubmissionUseCases
	AIModule           AIUseCases
}

func NewApplication(deps ApplicationDependencies) (*Application, error) {

	if err := deps.CheckDependencies(); err != nil {
		return nil, fmt.Errorf("application dependencies check failed: %w", err)
	}

	app := &Application{Dependencies: deps}

	app.CourseModule = CourseUseCases{
		CreateCourse:            course_crud_usecases.NewCreateCourseUseCase(deps.CourseRepository, deps.UserRepository),
		UpdateCourse:            course_crud_usecases.NewUpdateCourseUseCase(deps.CourseRepository, deps.UserRepository),
		DeleteCourse:            course_crud_usecases.NewDeleteCourseUseCase(deps.CourseRepository, deps.UserRepository, deps.ExamRepository, deps.ExamItemRepository, deps.ExamScoreRepository, deps.ExamItemScoreRepository),
		GetCourseDetails:        course_crud_usecases.NewGetCourseDetailsUseCase(deps.CourseRepository, deps.UserRepository),
		GetEnrolledCourses:      course_crud_usecases.NewGetEnrolledCoursesUseCase(deps.CourseRepository, deps.UserRepository),
		GetOwnedCourses:         course_crud_usecases.NewGetOwnedCoursesUseCase(deps.CourseRepository, deps.UserRepository),
		GetCourseStudents:       course_usecases.NewGetCourseStudentsUseCase(deps.CourseRepository, deps.UserRepository),
		EnrollInCourse:          course_usecases.NewEnrollInCourseUseCase(deps.CourseRepository, deps.UserRepository),
		RemoveStudentFromCourse: course_usecases.NewRemoveStudentFromCourseUseCase(deps.CourseRepository, deps.UserRepository),
	}

	app.ChallengeModule = ChallengeUseCases{
		CreateChallenge:     challenge_crud_usecases.NewCreateChallengeUseCase(deps.ChallengeRepository, deps.IOVariableRepository, deps.UserRepository),
		UpdateChallenge:     challenge_crud_usecases.NewUpdateChallengeUseCase(deps.ChallengeRepository, deps.IOVariableRepository, deps.UserRepository),
		PublishChallenge:    challenge_crud_usecases.NewPublishChallengeUseCase(deps.ChallengeRepository, deps.UserRepository),
		ArchiveChallenge:    challenge_crud_usecases.NewArchiveChallengeUseCase(deps.ChallengeRepository, deps.UserRepository),
		DeleteChallenge:     challenge_crud_usecases.NewDeleteChallengeUseCase(deps.ChallengeRepository, deps.UserRepository, deps.ExamRepository, deps.TestCaseRepository, deps.IOVariableRepository, deps.ExamItemRepository, deps.SubmissionRepository, deps.SubmissionResultRepository),
		GetChallengeDetails: challenge_crud_usecases.NewGetChallengeDetailsUseCase(deps.ChallengeRepository, deps.UserRepository, deps.ExamRepository),
		GetChallengesByUser: challenge_crud_usecases.NewGetChallengesByUserUseCase(deps.ChallengeRepository, deps.UserRepository, deps.ExamRepository),
		GetPublicChallenges: exam_usecases.NewGetPublicChallengesUseCase(deps.ChallengeRepository, deps.UserRepository),
		ForkChallenge:       exam_usecases.NewForkChallengeUseCase(deps.ChallengeRepository, deps.IOVariableRepository, deps.UserRepository, deps.TestCaseRepository),
	}

	app.TestCaseModule = TestCaseUseCases{
		CreateTestCase:          test_case_crud_usecases.NewCreateTestCaseUseCase(deps.UserRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.IOVariableRepository),
		UpdateTestCase:          test_case_crud_usecases.NewUpdateTestCaseUseCase(deps.UserRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.IOVariableRepository),
		DeleteTestCase:          test_case_crud_usecases.NewDeleteTestCaseUseCase(deps.UserRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.IOVariableRepository),
		GetTestCasesByChallenge: test_case_crud_usecases.NewGetTestCasesByChallengeUseCase(deps.UserRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.CourseRepository, deps.ExamItemRepository),
	}

	app.ExamItemModule = ExamItemUseCases{
		CreateExamItem: exam_item_crud_usecases.NewCreateExamItemUseCase(deps.UserRepository, deps.ExamRepository, deps.ExamItemRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.IOVariableRepository),
		UpdateExamItem: exam_item_crud_usecases.NewUpdateExamItemUseCase(deps.UserRepository, deps.ExamRepository, deps.ExamItemRepository, deps.ChallengeRepository),
		DeleteExamItem: exam_item_crud_usecases.NewDeleteExamItemUseCase(deps.UserRepository, deps.ExamRepository, deps.ExamItemRepository, deps.ChallengeRepository),
	}

	app.SessionModule = SessionUseCases{
		CreateSession:    session_usecases.NewCreateSessionUseCase(deps.UserRepository, deps.SessionRepository, deps.ExamRepository, deps.ExamScoreRepository, deps.ExamItemRepository, deps.ExamItemScoreRepository),
		GetActiveSession: session_usecases.NewGetActiveSessionUseCase(deps.SessionRepository, deps.UserRepository, deps.ExamRepository),
		HeartBeatSession: session_usecases.NewHeartBeatSessionUseCase(deps.UserRepository, deps.ExamRepository, deps.SessionRepository),
		BlockSession:     session_usecases.NewBlockSessionUseCase(deps.SessionRepository, deps.UserRepository),
		CloseSession:     session_usecases.NewCloseSessionUseCase(deps.SessionRepository, deps.ExamScoreRepository, deps.ExamItemRepository, deps.ExamItemScoreRepository, deps.SubmissionRepository, deps.UserRepository),
	}

	app.SubmissionUseCases = SubmissionUseCases{
		CreateSubmission:        submission_usecases.NewCreateSubmissionUseCase(deps.UserRepository, deps.SubmissionRepository, deps.SessionRepository, deps.ExamRepository, deps.ExamScoreRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.SubmissionResultRepository, deps.IOVariableRepository, deps.ExamItemRepository, deps.ExamItemScoreRepository, deps.PublisherPort),
		CreateSubmissionWithoutScore: submission_usecases.NewCreateSubmissionWithoutScoreUseCase(deps.UserRepository, deps.SubmissionRepository, deps.SessionRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.SubmissionResultRepository, deps.IOVariableRepository, deps.PublisherPort),
		CreateCustomSubmission:  submission_usecases.NewCreateCustomSubmissionUseCase(deps.UserRepository, deps.SubmissionRepository, deps.SessionRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.SubmissionResultRepository, deps.IOVariableRepository, deps.PublisherPort),
		GetSubmissionStatus:     submission_usecases.NewGetSubmissionStatusUseCase(deps.UserRepository, deps.SubmissionResultRepository, deps.SubmissionRepository),
		GetChallengeSubmissions: submission_usecases.NewGetChallengeSubmissionsUseCase(deps.UserRepository, deps.ChallengeRepository, deps.ExamRepository, deps.SubmissionRepository, deps.SubmissionResultRepository),
		GetSessionSubmissions:   submission_usecases.NewGetSessionSubmissionsUseCase(deps.UserRepository, deps.ChallengeRepository, deps.ExamRepository, deps.SubmissionRepository, deps.SubmissionResultRepository, deps.SessionRepository),
		GetUserSubmissions:      submission_usecases.NewGetUserSubmissionsUseCase(deps.UserRepository, deps.ChallengeRepository, deps.ExamRepository, deps.SubmissionRepository),
		UpdateResult:            submission_usecases.NewUpdateResultUseCase(deps.UserRepository, deps.SubmissionRepository, deps.SessionRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.IOVariableRepository, deps.SubmissionResultRepository, deps.LoginService, deps.PasswordHasher),
	}

	app.UserModule = UserUseCases{
		Register: user_usecases.NewRegisterUseCase(
			deps.RegisterService,
			deps.UserService,
			deps.PasswordHasher,
		),
		Login: user_usecases.NewLoginUseCase(
			deps.LoginService,
			deps.PasswordHasher,
		),
		GetData: user_usecases.NewGetDataUseCase(deps.LoginService),
		RefreshToken: user_usecases.NewRefreshTokenUseCase(
			deps.TokenService,
		),
	}

	app.ExamModule = ExamUseCases{
		CreateExam:              exam_crud_usecases.NewCreateExamUseCase(deps.UserRepository, deps.CourseRepository, deps.ExamRepository),
		UpdateExam:              exam_crud_usecases.NewUpdateExamUseCase(deps.UserRepository, deps.ExamRepository),
		CloseExam:               exam_crud_usecases.NewCloseExamUseCase(deps.UserRepository, deps.ExamRepository),
		DeleteExam:              exam_crud_usecases.NewDeleteExamUseCase(deps.UserRepository, deps.ExamRepository, deps.ExamItemRepository, deps.ExamScoreRepository, deps.ExamItemScoreRepository),
		GetExamDetails:          exam_crud_usecases.NewGetExamDetailsUseCase(deps.UserRepository, deps.ExamRepository, deps.CourseRepository),
		GetExamsByCourse:        exam_crud_usecases.NewGetExamsByCourseUseCase(deps.UserRepository, deps.ExamRepository, deps.CourseRepository),
		GetExamItems:            exam_crud_usecases.NewGetExamItemsUseCase(deps.ChallengeRepository, deps.UserRepository, deps.ExamRepository, deps.ExamItemRepository),
		GetPublicExams:          exam_usecases.NewGetPublicExamsUseCase(deps.UserRepository, deps.ExamRepository),
		GetCodeDefaultTemplates: exam_usecases.NewGetCodeDefaultTemplates(deps.UserRepository, deps.ExamRepository),
	}

	app.AIModule = AIUseCases{
		GenerateFullChallenge: ai_usecases.NewGenerateFullChallengeUseCase(deps.AIAdapter),
		GenerateExam:          ai_usecases.NewGenerateExamUseCase(deps.AIAdapter),
	}

	return app, nil
}
