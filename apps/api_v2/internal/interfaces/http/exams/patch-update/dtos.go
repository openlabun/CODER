package patchupdate

import examDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/exam"

type RequestDTO struct {
Title *string `json:"title"`
Description *string `json:"description"`
Visibility *string `json:"visibility"`
StartTime *string `json:"start_time"`
EndTime *string `json:"end_time"`
AllowLateSubmissions *bool `json:"allow_late_submissions"`
TimeLimit *int `json:"time_limit"`
TryLimit *int `json:"try_limit"`
}

type PathDTO struct { ID string }

func ToInput(path PathDTO, body RequestDTO) examDtos.UpdateExamInput {
return examDtos.UpdateExamInput{
ExamID: path.ID,
Title: body.Title,
Description: body.Description,
Visibility: body.Visibility,
StartTime: body.StartTime,
EndTime: body.EndTime,
AllowLateSubmissions: body.AllowLateSubmissions,
TimeLimit: body.TimeLimit,
TryLimit: body.TryLimit,
}
}
