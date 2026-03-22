package container

import (
	"fmt"

	ports "github.com/openlabun/CODER/apps/api_v2/internal/application/ports/user"
	user_repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/user"
	course_repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/course"
	exam_repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/exam"
	submission_repositories "github.com/openlabun/CODER/apps/api_v2/internal/domain/repositories/submission"
)


// ApplicationDependencies groups all contract-based dependencies required by
// application use cases.
type ApplicationDependencies struct {
	RegisterService ports.RegisterPort
	LoginService    ports.LoginPort
	UserService     ports.UserServicePort
	TokenService    ports.TokenServicePort
	PasswordHasher  ports.PasswordHasherPort

	UserRepository user_repositories.UserRepository
	CourseRepository course_repositories.CourseRepository
	ExamRepository exam_repositories.ExamRepository
	ChallengeRepository exam_repositories.ChallengeRepository
	TestCaseRepository exam_repositories.TestCaseRepository
	SubmissionRepository submission_repositories.SubmissionRepository
	SessionRepository submission_repositories.SessionRepository
	SubmissionResultRepository submission_repositories.SubmissionResultRepository
}

func NewApplicationDependencies(
	registerService ports.RegisterPort,
	loginService ports.LoginPort,
	userService ports.UserServicePort,
	tokenService ports.TokenServicePort,
	passwordHasher ports.PasswordHasherPort,
	userRepo user_repositories.UserRepository,
	courseRepo course_repositories.CourseRepository,
	examRepo exam_repositories.ExamRepository,
	challengeRepo exam_repositories.ChallengeRepository,
	testCaseRepo exam_repositories.TestCaseRepository,
	submissionRepo submission_repositories.SubmissionRepository,
	sessionRepo submission_repositories.SessionRepository,
	submissionResultRepo submission_repositories.SubmissionResultRepository,

) ApplicationDependencies {
	return ApplicationDependencies{
		RegisterService: registerService,
		LoginService:    loginService,
		UserService:     userService,
		TokenService:    tokenService,
		PasswordHasher:  passwordHasher,
		UserRepository:  userRepo,
		CourseRepository: courseRepo,
		ExamRepository: examRepo,
		ChallengeRepository: challengeRepo,
		TestCaseRepository: testCaseRepo,
		SubmissionRepository: submissionRepo,
		SessionRepository: sessionRepo,
		SubmissionResultRepository: submissionResultRepo,
	}
}

func (deps ApplicationDependencies) CheckDependencies() error {
	if deps.RegisterService == nil {
		return fmt.Errorf("RegisterService dependency is not provided")
	}

	if deps.LoginService == nil {
		return fmt.Errorf("LoginService dependency is not provided")
	}

	if deps.UserService == nil {
		return fmt.Errorf("UserService dependency is not provided")
	}

	if deps.TokenService == nil {
		return fmt.Errorf("TokenService dependency is not provided")
	}

	if deps.PasswordHasher == nil {
		return fmt.Errorf("PasswordHasher dependency is not provided")
	}

	if deps.UserRepository == nil {
		return fmt.Errorf("UserRepository dependency is not provided")
	}

	if deps.CourseRepository == nil {
		return fmt.Errorf("CourseRepository dependency is not provided")
	}

	if deps.ExamRepository == nil {
		return fmt.Errorf("ExamRepository dependency is not provided")
	}

	if deps.SubmissionRepository == nil {
		return fmt.Errorf("SubmissionRepository dependency is not provided")
	}

	if deps.ChallengeRepository == nil {
		return fmt.Errorf("ChallengeRepository dependency is not provided")
	}

	if deps.TestCaseRepository == nil {
		return fmt.Errorf("TestCaseRepository dependency is not provided")
	}

	if deps.SessionRepository == nil {
		return fmt.Errorf("SessionRepository dependency is not provided")
	}

	if deps.SubmissionResultRepository == nil {
		return fmt.Errorf("SubmissionResultRepository dependency is not provided")
	}

	return  nil
}