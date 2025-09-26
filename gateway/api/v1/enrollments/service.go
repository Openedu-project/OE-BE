package enrollments

import (
	"errors"
	"time"

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

func (s *Service) GetMyCourses(userId uint) (*MyCoureseResponseDTO, error) {
	userCourses, err := s.repo.FindUserCoursesByUserID(userId)
	if err != nil {
		return nil, err
	}

	response := MyCoureseResponseDTO{
		InProgressCourses: []CourseInfoDTO{},
		CompletedCourses:  []CourseInfoDTO{},
		NotStartedCourses: []CourseInfoDTO{},
	}

	for _, uc := range userCourses {
		if uc.Course == nil {
			continue
		}
		lecturerName := ""
		if uc.Course.Lecturer != nil {
			lecturerName = uc.Course.Lecturer.Name
		}
		courseInfo := CourseInfoDTO{
			ID:               uc.Course.ID,
			Name:             uc.Course.Name,
			ShortDescription: uc.Course.ShortDescription,
			Banner:           uc.Course.Banner,
			LecturerName:     lecturerName,
		}

		switch uc.Status {
		case models.StatusInProgress:
			response.InProgressCourses = append(response.InProgressCourses, courseInfo)
		case models.StatusCompleted:
			response.CompletedCourses = append(response.CompletedCourses, courseInfo)
		}
	}

	return &response, nil
}

func (s *Service) GetDashboardSummary(userId uint) (*DashboardSummaryDTO, error) {
	counts, err := s.repo.CountCoursesByStatus(userId)
	if err != nil {
		return nil, err
	}

	summary := &DashboardSummaryDTO{
		InProgressCount: 0,
		CompletedCount:  0,
		NotStartedCount: 0,
	}

	for _, result := range counts {
		switch result.Status {
		case models.StatusInProgress:
			summary.InProgressCount = result.Count
		case models.StatusCompleted:
			summary.CompletedCount = result.Count
		}
	}

	return summary, nil
}

func (s *Service) GetMyCoursesByStatus(userId uint, status models.UserCourseStatus, page int, pageSize int) ([]CourseInfoDTO, error) {
	offset := (page - 1) * pageSize
	limit := pageSize

	userCourses, err := s.repo.FindUserCourseByUserIDAndStatus(userId, status, offset, limit)
	if err != nil {
		return nil, err
	}

	var coursesDTO []CourseInfoDTO
	for _, uc := range userCourses {
		if uc.Course == nil {
			continue
		}
		lecturerName := ""

		if uc.Course.Lecturer != nil {
			lecturerName = uc.Course.Lecturer.Name
		}
		courseInfo := CourseInfoDTO{
			ID:               uc.Course.ID,
			Name:             uc.Course.Name,
			ShortDescription: uc.Course.ShortDescription,
			Banner:           uc.Course.Banner,
			LecturerName:     lecturerName,
		}
		coursesDTO = append(coursesDTO, courseInfo)
	}

	if coursesDTO == nil {
		coursesDTO = []CourseInfoDTO{}
	}

	return coursesDTO, nil
}

func (s *Service) CompletedCourse(userId uint, courseId uint) (*models.UserCourse, error) {
	userCourse, err := s.repo.FindByUserIDAndCourseID(userId, courseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User is not enrolled in this course")
		}
		return nil, err
	}

	if userCourse.Status == models.StatusCompleted {
		return nil, errors.New("Course is already completed")
	}

	now := time.Now()
	userCourse.Status = models.StatusCompleted
	userCourse.CompletedAt = &now

	if err := s.repo.Update(userCourse); err != nil {
		return nil, err
	}

	return userCourse, nil
}
