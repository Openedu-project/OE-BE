package models

import (
	"gorm.io/gorm"
)

type CourseSection struct {
	gorm.Model
	CourseID uint           `json:"course_id"`
	Name     string         `json:"name" gorm:"size:255;not null"`
	Status   string         `json:"status" gorm:"size:50;default:'draft'"`
	Lessons  []CourseLesson `json:"lessons" gorm:"foreignKey:CourseSectionID;references:ID"`
}

type CourseLesson struct {
	gorm.Model
	CourseSectionID uint   `json:"course_section_id"`
	Name            string `json:"name" gorm:"size:255;not null"`
	Status          string `json:"status" gorm:"size:50;default:'draft'"`
}
