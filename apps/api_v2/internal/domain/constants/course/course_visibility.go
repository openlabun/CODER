package constants

type CourseVisibility string

const (
	CourseVisibilityPublic  CourseVisibility = "public"
	CourseVisibilityPrivate CourseVisibility = "private"
	CourseVisibilityBlocked CourseVisibility = "blocked"
)