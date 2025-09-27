package certificates

import "time"

type CertificateDTO struct {
	CourseName string    `json:"course_name"`
	Code       string    `json:"code"`
	IssuedAt   time.Time `json:"issued_at"`
}
