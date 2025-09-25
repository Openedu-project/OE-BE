package enrollments

import (
	"time"

	"gorm.io/datatypes"
)

type EnrollmentResponseDTO struct {
	ID       uint      `json:"id"`
	UserID   uint      `json:"user_id"`
	CourseID uint      `json:"course_id"`
	Status   string    `json:"status"`
	CreateAt time.Time `json:"create_at"`
}

type CourseInfoDTO struct {
	ID               uint           `json:"id"`
	Name             string         `json:"name"`
	ShortDescription string         `json:"short_description"`
	Banner           datatypes.JSON `json:"banner"`
	LecturerName     string         `json:"lecturer_name"`
}

type MyCoureseResponseDTO struct {
	InProgressCourses []CourseInfoDTO `json:"in_progress_courses"`
	CompletedCourses  []CourseInfoDTO `json:"completed_courses"`
	NotStartedCourses []CourseInfoDTO `json:"not_started_courses"`
}
