package sections

import (
	"gateway/models"

	"gorm.io/gorm"
)

type CourseSectionRepository struct {
	db *gorm.DB
}

func NewCourseSectionRepository(db *gorm.DB) *CourseSectionRepository {
	return &CourseSectionRepository{
		db: db,
	}
}
func (r *CourseSectionRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}
