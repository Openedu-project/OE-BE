package categories

import (
	"gateway/models"

	"gorm.io/gorm"
)

type CourseCategoryRepository struct {
	db *gorm.DB
}

func NewCourseCategoryRepository(db *gorm.DB) *CourseCategoryRepository {
	return &CourseCategoryRepository{
		db: db,
	}
}
func (r *CourseCategoryRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}
