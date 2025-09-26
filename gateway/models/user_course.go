package models

import (
	"time"

	"gorm.io/gorm"
)

type UserCourseStatus string

const (
	StatusInProgress UserCourseStatus = "in_progress"
	StatusCompleted  UserCourseStatus = "completed"
)

type UserCourse struct {
	gorm.Model
	UserID      uint             `json:"user_id"`
	User        User             `gorm:"foreignKey:UserID"`
	CourseID    uint             `json:"course_id"`
	Course      *Course          `gorm:"foreignKey:CourseID"`
	Status      UserCourseStatus `gorm:"default:'in_progress"`
	CompletedAt *time.Time
}
