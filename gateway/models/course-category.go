package models

import "gorm.io/gorm"

type CourseCategory struct {
	gorm.Model
	Name string `gorm:"size:255;not null" json:"name"`
}
