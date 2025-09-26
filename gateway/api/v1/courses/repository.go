package courses

import (
	"gateway/models"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (r *CourseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *CourseRepository) FindByID(courseId uint) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Lecturer").First(&course, courseId).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}
