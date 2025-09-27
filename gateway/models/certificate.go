package models

import (
	"time"

	"gorm.io/gorm"
)

type Certificate struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	User     User      `gorm:"foreignKey:UserID"`
	CourseID uint      `json:"course_id"`
	Course   *Course   `gorm:"foreignKey:CourseID"`
	Code     string    `gorm:"uniqueIndex;size:100"`
	IssuedAt time.Time `json:"issued_at"`
}
