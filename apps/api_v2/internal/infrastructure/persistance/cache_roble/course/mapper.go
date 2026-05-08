package cache_roble_infrastructure

import (
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/course"
	cache "github.com/openlabun/CODER/apps/api_v2/internal/infrastructure/persistance/cache_roble"
)

func coursesFromRecords(records []map[string]any) []*Entities.Course {
	courses := make([]*Entities.Course, 0, len(records))
	for _, rec := range records {
		if c, err := cache.MapToEntity[Entities.Course](rec); err == nil && c != nil {
			courses = append(courses, c)
		}
	}
	return courses
}