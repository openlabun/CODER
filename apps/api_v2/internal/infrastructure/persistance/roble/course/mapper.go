package roble_infrastructure

import (
	"strings"
	"time"

	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	course_factory "github.com/openlabun/CODER/apps/api_v2/internal/domain/factory/course"
	consts "github.com/openlabun/CODER/apps/api_v2/internal/domain/constants/course"
)

func courseToRecord(course *Entities.Course) map[string]any {
	record := map[string]any{
		"ID":             strings.TrimSpace(course.ID),
		"Name":           strings.TrimSpace(course.Name),
		"Visibility":     string(course.Visibility),
		"VisualIdentity": string(course.VisualIdentity),
		"Code":           strings.TrimSpace(course.Code),
		"EnrollmentCode": strings.TrimSpace(course.EnrollmentCode),
		"EnrollmentURL":  strings.TrimSpace(course.EnrollmentURL),
		"ProfessorID":    strings.TrimSpace(course.ProfessorID),
		"CreatedAt":      course.CreatedAt.UTC().Format(time.RFC3339),
		"UpdatedAt":      course.UpdatedAt.UTC().Format(time.RFC3339),
	}

	if enrollmentURL := strings.TrimSpace(course.EnrollmentURL); enrollmentURL != "" {
		record["EnrollmentURL"] = enrollmentURL
	}

	if course.Period != nil && course.Period.Year > 0 && strings.TrimSpace(string(course.Period.Semester)) != "" {
		record["PeriodYear"] = course.Period.Year
		record["PeriodSemester"] = string(course.Period.Semester)
	}

	return record
}

func courseToUpdates(course *Entities.Course) map[string]any {
	updates := map[string]any{
		"Name":           strings.TrimSpace(course.Name),
		"Description":    strings.TrimSpace(course.Description),
		"Visibility":     string(course.Visibility),
		"VisualIdentity": string(course.VisualIdentity),
		"Code":           strings.TrimSpace(course.Code),
		"EnrollmentCode": strings.TrimSpace(course.EnrollmentCode),
		"EnrollmentURL":  strings.TrimSpace(course.EnrollmentURL),
		"ProfessorID":    strings.TrimSpace(course.ProfessorID),
	}

	if course.Period != nil {
		updates["PeriodYear"] = course.Period.Year
		updates["PeriodSemester"] = string(course.Period.Semester)
	}

	return updates
}

func recordToCourse(record map[string]any) (*Entities.Course, error) {
	period := extractPeriod(record)
	createdAt, _ := asTime(record["CreatedAt"])
	updatedAt, _ := asTime(record["UpdatedAt"])

	return course_factory.ExistingCourse(
		asString(record["ID"]),
		asString(record["Name"]),
		asString(record["Description"]),
		consts.CourseVisibility(asString(record["Visibility"])),
		consts.CourseColour(asString(record["VisualIdentity"])),
		asString(record["Code"]),
		period,
		asString(record["EnrollmentCode"]),
		asString(record["EnrollmentURL"]),
		asString(record["ProfessorID"]),
		createdAt,
		updatedAt,
	)
}

func extractPeriod(record map[string]any) *Entities.Period {
	if rawPeriod, ok := record["Period"]; ok {
		if m, ok := rawPeriod.(map[string]any); ok {
			year := asInt(m["Year"])
			semester := consts.AcademicPeriod(asString(m["Semester"]))
			if year > 0 && semester != "" {
				return &Entities.Period{Year: year, Semester: semester}
			}
		}
	}

	year := asInt(record["PeriodYear"])
	semester := consts.AcademicPeriod(asString(record["PeriodSemester"]))
	if year > 0 && semester != "" {
		return &Entities.Period{Year: year, Semester: semester}
	}

	return nil
}