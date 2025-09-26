package enrollments

import (
	"errors"

	"gateway/models"

	"gorm.io/gorm"
)

type EnrollRepository struct {
	db *gorm.DB
}

func NewEnrollRepository(db *gorm.DB) *EnrollRepository {
	return &EnrollRepository{db: db}
}

func (r *EnrollRepository) FindByUserIDAndCourseID(userID uint, courseID uint) (*models.UserCourse, error) {
	var userCourse models.UserCourse
	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&userCourse).Error; err != nil {
		return nil, err
	}
	return &userCourse, nil
}

func (r *EnrollRepository) Create(userCourse *models.UserCourse) error {
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

func (r *EnrollRepository) FindUserCoursesByUserID(userID uint) ([]models.UserCourse, error) {
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

func (r *EnrollRepository) CountCoursesByStatus(userID uint) ([]StatusCountResult, error) {
	var results []StatusCountResult
	err := r.db.Model(&models.UserCourse{}).Select("status, count(*) as count").Where("user_id = ?", userID).Group("status").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *EnrollRepository) FindUserCourseByUserIDAndStatus(userID uint, status models.UserCourseStatus, offset int, limit int) ([]models.UserCourse, error) {
	var userCourses []models.UserCourse
	err := r.db.Preload("Course").Preload("Course.Lecturer").Where("user_id = ?  AND status = ?", userID, status).Offset(offset).Limit(limit).Find(&userCourses).Error
	if err != nil {
		return nil, err
	}

	return userCourses, nil
}

func (r *EnrollRepository) Update(userCourse *models.UserCourse) error {
	return r.db.Save(userCourse).Error
}
