package models

import "gorm.io/gorm"

type Wishlist struct {
	gorm.Model
	UserID   uint    `gorm:"uniqueIndex:idx_user_course"`
	User     User    `gorm:"foreignKey:UserID"`
	CourseID uint    `gorm:"uniqueIndex:idx_user_course"`
	Course   *Course `gorm:"foreignKey:CourseID"`
}
