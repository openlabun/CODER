package cache_roble_infrastructure

// Domain entities
const UserEntity = 				"User"
const CourseEntity = 			"Course"

const ExamEntity = 				"Exam"
const ChallengeEntity = 		"Challenge"
const ExamItemEntity = 			"ExamItem"
const ExamItemScoreEntity = 	"ExamItemScore"
const ExamScoreEntity = 		"ExamScore"
const IOVariableEntity = 		"IOVariable"
const TestCaseEntity = 			"TestCase"

const SubmissionResultEntity = 	"SubmissionResult"
const SessionEntity = 			"Session"
const SubmissionEntity = 		"Submission"

// Table names for cache storage
const CacheUserTable = 				"cache_user"
const CacheCourseTable = 			"cache_course"

const CacheChallengeTable = 		"cache_challenge"
const CacheExamItemTable = 			"cache_exam_item"
const CacheExamItemScoreTable = 	"cache_exam_item_score"
const CacheExamTable = 				"cache_exam"
const CacheExamScoreTable = 		"cache_exam_score"
const CacheIOVariableTable = 		"cache_io_variable"
const CacheTestCaseTable = 			"cache_test_case"

const CacheSubmissionResultTable = 	"cache_submission_result"
const CacheSessionTable = 			"cache_session"
const CacheSubmissionTable = 		"cache_submission"

// Operations
const InsertOperation = "insert"
const UpdateOperation = "update"

// Dictionary Domain entities -> Cache table names
var EntityToCacheTable = map[string]string{
	CourseEntity: CacheCourseTable,

	ExamEntity: CacheExamTable,
	ChallengeEntity: CacheChallengeTable,
	ExamItemEntity: CacheExamItemTable,
	ExamItemScoreEntity: CacheExamItemScoreTable,
	ExamScoreEntity: CacheExamScoreTable,
	IOVariableEntity: CacheIOVariableTable,
	TestCaseEntity: CacheTestCaseTable,

	SubmissionResultEntity: CacheSubmissionResultTable,
	SessionEntity: CacheSessionTable,
	SubmissionEntity: CacheSubmissionTable,
}