package models

import (
	"gorm.io/gorm"
)

type CourseSection struct {
	CourseID uint   `json:"course_id"`
	Name     string `json:"name" gorm:"size:255;not null"`
	Status   string `json:"status" gorm:"size:50;default:'draft'"`

	Lessons []CourseLesson `json:"lessons" gorm:"foreignKey:SectionID"`
	gorm.Model
}

type CourseLesson struct {
	SectionID uint   `json:"section_id"`
	Name      string `json:"name" gorm:"size:255;not null"`
	Status    string `json:"status" gorm:"size:50;default:'draft'"`
	QuizID    *uint  `json:"quiz_id"`
	gorm.Model
}
