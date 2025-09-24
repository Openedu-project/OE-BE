package enrollments

import "time"

type EnrollmentResponseDTO struct {
	ID       uint      `json:"id"`
	UserID   uint      `json:"user_id"`
	CourseID uint      `json:"course_id"`
	Status   string    `json:"status"`
	CreateAt time.Time `json:"create_at"`
}
