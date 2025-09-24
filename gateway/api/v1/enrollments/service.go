package enrollments

import (
	"errors"

	"gateway/models"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEnrollment(userId uint, courseId uint) (*models.UserCourse, error) {
	// Check đã có enroll chưa
	_, err := s.repo.FindByUserIDAndCourseID(userId, courseId)
	if err == nil {
		return nil, errors.New("user already enrolled in this course")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	newUserCourse := &models.UserCourse{
		UserID:   userId,
		CourseID: courseId,
		Status:   models.StatusInProgress,
	}
	if err := s.repo.Create(newUserCourse); err != nil {
		return nil, err
	}

	return newUserCourse, nil
}
