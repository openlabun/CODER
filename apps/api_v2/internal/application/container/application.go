package container

import (
	"fmt"

	course_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/course"
	course_crud_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/course/crud"
	user_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/user"

	"github.com/openlabun/CODER/apps/api_v2/internal/application/services"
	ai_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/ai"
	challenge_crud_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam/challenge_crud"
	exam_crud_usecases "github.com/openlabun/CODER/apps/api_v2/internal/application/usecases/exam/exam_crud"
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
	GetAllCourses      *course_crud_usecases.GetAllCoursesUseCase

	GetCourseStudents       *course_usecases.GetCourseStudentsUseCase
	EnrollInCourse          *course_usecases.EnrollInCourseUseCase
	RemoveStudentFromCourse *course_usecases.RemoveStudentFromCourseUseCase
	GetCourseChallenges     *course_usecases.GetCourseChallengesUseCase
	AssignChallengeToCourse *course_usecases.AssignChallengeToCourseUseCase
}

type ChallengeUseCases struct {
	CreateChallenge     *challenge_crud_usecases.CreateChallengeUseCase
	UpdateChallenge     *challenge_crud_usecases.UpdateChallengeUseCase
	PublishChallenge    *challenge_crud_usecases.PublishChallengeUseCase
	ArchiveChallenge    *challenge_crud_usecases.ArchiveChallengeUseCase
	DeleteChallenge     *challenge_crud_usecases.DeleteChallengeUseCase
	GetChallengeDetails *challenge_crud_usecases.GetChallengeDetailsUseCase
	GetChallengesByExam *challenge_crud_usecases.GetChallengesByExamUseCase
}

type TestCaseUseCases struct {
	CreateTestCase          *test_case_crud_usecases.CreateTestCaseUseCase
	UpdateTestCase          *test_case_crud_usecases.UpdateTestCaseUseCase
	DeleteTestCase          *test_case_crud_usecases.DeleteTestCaseUseCase
	GetTestCasesByChallenge *test_case_crud_usecases.GetTestCasesByChallengeUseCase
}

type SubmissionUseCases struct {
	CreateSubmission        *submission_usecases.CreateSubmissionUseCase
	GetSubmissionStatus     *submission_usecases.GetSubmissionStatusUseCase
	GetChallengeSubmissions *submission_usecases.GetChallengeSubmissionsUseCase
	GetUserSubmissions      *submission_usecases.GetUserSubmissionsUseCase
	UpdateResult            *submission_usecases.UpdateResultUseCase
}

type SessionUseCases struct {
	CreateSession    *session_usecases.CreateSessionUseCase
	GetSession       *session_usecases.GetSessionUseCase
	HeartBeatSession *session_usecases.HeartBeatSessionUseCase
	CloseSession     *session_usecases.CloseSessionUseCase
}

type ExamUseCases struct {
	CreateExam       *exam_crud_usecases.CreateExamUseCase
	UpdateExam       *exam_crud_usecases.UpdateExamUseCase
	DeleteExam       *exam_crud_usecases.DeleteExamUseCase
	GetExamDetails   *exam_crud_usecases.GetExamDetailsUseCase
	GetExamsByCourse *exam_crud_usecases.GetExamsByCourseUseCase
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
		DeleteCourse:            course_crud_usecases.NewDeleteCourseUseCase(deps.CourseRepository, deps.UserRepository),
		GetCourseDetails:        course_crud_usecases.NewGetCourseDetailsUseCase(deps.CourseRepository, deps.UserRepository),
		GetEnrolledCourses:      course_crud_usecases.NewGetEnrolledCoursesUseCase(deps.CourseRepository, deps.UserRepository),
		GetOwnedCourses:         course_crud_usecases.NewGetOwnedCoursesUseCase(deps.CourseRepository, deps.UserRepository),
		GetAllCourses:           course_crud_usecases.NewGetAllCoursesUseCase(deps.CourseRepository, deps.UserRepository),
		GetCourseStudents:       course_usecases.NewGetCourseStudentsUseCase(deps.CourseRepository, deps.UserRepository),
		EnrollInCourse:          course_usecases.NewEnrollInCourseUseCase(deps.CourseRepository, deps.UserRepository),
		RemoveStudentFromCourse: course_usecases.NewRemoveStudentFromCourseUseCase(deps.CourseRepository, deps.UserRepository),
		GetCourseChallenges:     course_usecases.NewGetCourseChallengesUseCase(deps.ChallengeRepository, deps.UserRepository),
		AssignChallengeToCourse: course_usecases.NewAssignChallengeToCourseUseCase(deps.ChallengeRepository, deps.UserRepository),
	}

	app.ChallengeModule = ChallengeUseCases{
		CreateChallenge:     challenge_crud_usecases.NewCreateChallengeUseCase(deps.ChallengeRepository, deps.UserRepository),
		UpdateChallenge:     challenge_crud_usecases.NewUpdateChallengeUseCase(deps.ChallengeRepository, deps.UserRepository),
		PublishChallenge:    challenge_crud_usecases.NewPublishChallengeUseCase(deps.ChallengeRepository, deps.UserRepository),
		ArchiveChallenge:    challenge_crud_usecases.NewArchiveChallengeUseCase(deps.ChallengeRepository, deps.UserRepository),
		DeleteChallenge:     challenge_crud_usecases.NewDeleteChallengeUseCase(deps.ChallengeRepository, deps.UserRepository, deps.ExamRepository),
		GetChallengeDetails: challenge_crud_usecases.NewGetChallengeDetailsUseCase(deps.ChallengeRepository, deps.UserRepository, deps.ExamRepository),
		GetChallengesByExam: challenge_crud_usecases.NewGetChallengesByExamUseCase(deps.ChallengeRepository, deps.UserRepository, deps.ExamRepository),
	}

	app.TestCaseModule = TestCaseUseCases{
		CreateTestCase:          test_case_crud_usecases.NewCreateTestCaseUseCase(deps.UserRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository),
		UpdateTestCase:          test_case_crud_usecases.NewUpdateTestCaseUseCase(deps.UserRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository),
		DeleteTestCase:          test_case_crud_usecases.NewDeleteTestCaseUseCase(deps.UserRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository),
		GetTestCasesByChallenge: test_case_crud_usecases.NewGetTestCasesByChallengeUseCase(deps.UserRepository, deps.ExamRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.CourseRepository),
	}

	app.SessionModule = SessionUseCases{
		CreateSession:    session_usecases.NewCreateSessionUseCase(deps.UserRepository, deps.SessionRepository, deps.ExamRepository),
		GetSession:       session_usecases.NewGetSessionUseCase(deps.SessionRepository, deps.UserRepository),
		HeartBeatSession: session_usecases.NewHeartBeatSessionUseCase(deps.UserRepository, deps.SessionRepository),
		CloseSession:     session_usecases.NewCloseSessionUseCase(deps.UserRepository, deps.SessionRepository, deps.ExamRepository),
	}

	app.SubmissionUseCases = SubmissionUseCases{
		CreateSubmission:        submission_usecases.NewCreateSubmissionUseCase(deps.UserRepository, deps.SubmissionRepository, deps.SessionRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.SubmissionResultRepository, deps.PublisherPort),
		GetSubmissionStatus:     submission_usecases.NewGetSubmissionStatusUseCase(deps.UserRepository, deps.SubmissionResultRepository, deps.SubmissionRepository),
		GetChallengeSubmissions: submission_usecases.NewGetChallengeSubmissionsUseCase(deps.UserRepository, deps.ChallengeRepository, deps.ExamRepository, deps.SubmissionRepository, deps.SubmissionResultRepository),
		GetUserSubmissions:      submission_usecases.NewGetUserSubmissionsUseCase(deps.UserRepository, deps.ChallengeRepository, deps.ExamRepository, deps.SubmissionRepository),
		UpdateResult:            submission_usecases.NewUpdateResultUseCase(deps.UserRepository, deps.SubmissionRepository, deps.SessionRepository, deps.ChallengeRepository, deps.TestCaseRepository, deps.SubmissionResultRepository, deps.LoginService, deps.PasswordHasher),
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
		CreateExam:       exam_crud_usecases.NewCreateExamUseCase(deps.UserRepository, deps.ExamRepository),
		UpdateExam:       exam_crud_usecases.NewUpdateExamUseCase(deps.UserRepository, deps.ExamRepository),
		DeleteExam:       exam_crud_usecases.NewDeleteExamUseCase(deps.UserRepository, deps.ExamRepository),
		GetExamDetails:   exam_crud_usecases.NewGetExamDetailsUseCase(deps.UserRepository, deps.ExamRepository, deps.CourseRepository),
		GetExamsByCourse: exam_crud_usecases.NewGetExamsByCourseUseCase(deps.UserRepository, deps.ExamRepository, deps.CourseRepository),
	}

	geminiService := services.NewGeminiService()
	app.AIModule = AIUseCases{
		GenerateFullChallenge: ai_usecases.NewGenerateFullChallengeUseCase(geminiService),
		GenerateExam:          ai_usecases.NewGenerateExamUseCase(geminiService),
	}

	return app, nil
}
