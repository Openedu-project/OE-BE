package enrollments

import (
	"errors"

	"gateway/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUserIDAndCourseID(userID uint, courseID uint) (*models.UserCourse, error) {
	var userCourse models.UserCourse
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&userCourse).Error; err != nil {
		return nil, err
	}
	return &userCourse, nil
}

func (r *Repository) Create(userCourse *models.UserCourse) error {
	var course models.Course
	if err := r.db.First(&course, userCourse.CourseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Cannot enroll in a non-existent course")
		}
		return err
	}
	if err := r.db.Create(userCourse).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindUserCoursesByUserID(userID uint) ([]models.UserCourse, error) {
	var userCourse []models.UserCourse
	if err := r.db.Preload("Course").Preload("Course.Lecturer").Where("user_id = ?", userID).Find(&userCourse).Error; err != nil {
		return nil, err
	}
	return userCourse, nil
}

type StatusCountResult struct {
	Status models.UserCourseStatus
	Count  int64
}

func (r *Repository) CountCoursesByStatus(userID uint) ([]StatusCountResult, error) {
	var results []StatusCountResult
	err := r.db.Model(&models.UserCourse{}).Select("status, count(*) as count").Where("user_id = ?", userID).Group("status").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
