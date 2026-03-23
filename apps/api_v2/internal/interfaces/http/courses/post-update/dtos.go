package postupdate

import courseDtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/course"

type RequestDTO struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Visibility     *string `json:"visibility"`
	VisualIdentity *string `json:"visual_identity"`
	Code           *string `json:"code"`
	Year           *int    `json:"year"`
	Semester       *string `json:"semester"`
	EnrollmentCode *string `json:"enrollment_code"`
	EnrollmentURL  *string `json:"enrollment_url"`
	TeacherID      *string `json:"teacher_id"`
}

type PathDTO struct { ID string }

func ToInput(path PathDTO, body RequestDTO) courseDtos.UpdateCourseInput {
	return courseDtos.UpdateCourseInput{
		ID: path.ID,
		Name: body.Name,
		Description: body.Description,
		Visibility: body.Visibility,
		VisualIdentity: body.VisualIdentity,
		Code: body.Code,
		Year: body.Year,
		Semester: body.Semester,
		EnrollmentCode: body.EnrollmentCode,
		EnrollmentURL: body.EnrollmentURL,
		TeacherID: body.TeacherID,
	}
}
