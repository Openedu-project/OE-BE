package courses

import (
	"encoding/json"
	"errors"
	"fmt"

	"gateway/api/v1/enrollments"
	"gateway/guards"
	"gateway/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CourseService struct {
	repo       *CourseRepository
	enrollRepo *enrollments.EnrollRepository
}

func NewCourseService(repo *CourseRepository, enrollRepo *enrollments.EnrollRepository) *CourseService {
	return &CourseService{
		repo:       repo,
		enrollRepo: enrollRepo,
	}
}

func (s *CourseService) CreateCourse(dto CreateCourseDTO, userId uint) (*models.Course, error) {
	course := &models.Course{
		Name:             dto.Name,
		ShortDescription: dto.ShortDescription,
		LecturerId:       userId,
		Status:           "draft", // default
	}
	if err := s.repo.Create(course); err != nil {
		return nil, err
	}

	initSection := &models.CourseSection{
		CourseID: course.ID,
		Name:     "Create your first section",
		Status:   "draft",
	}
	if err := s.repo.db.Create(initSection).Error; err != nil {
		return nil, err
	}

	initLesson := &models.CourseLesson{
		CourseSectionID: initSection.ID,
		Name:            "Create your first lesson",
		Status:          "draft",
	}
	if err := s.repo.db.Create(initLesson).Error; err != nil {
		return nil, err
	}

	return course, nil
}

func (s *CourseService) UpdateCourseInfo(courseId uint, data map[string]interface{}, userId uint, role string) (*models.Course, error) {
	var course models.Course
	if err := s.repo.db.First(&course, courseId).Error; err != nil {
		return nil, err
	}

	if role != "admin" && role != "sysadmin" && course.LecturerId != userId {
		return nil, errors.New("forbidden: you are not allowed to update this course")
	}

	if course.IsCompleted {
		return nil, errors.New("course already completed, cannot update")
	}

	// update banner (JSON)
	if v, ok := data["banner"].(string); ok && v != "" {
		bytes, _ := json.Marshal(map[string]string{
			"path": v,
		})
		course.Banner = datatypes.JSON(bytes)
	}

	// update video_preview (JSON)
	if v, ok := data["video_preview"].(string); ok && v != "" {
		bytes, _ := json.Marshal(map[string]string{
			"path": v,
		})
		course.VideoPreview = datatypes.JSON(bytes)
	}
	if v, ok := data["category"].(string); ok && v != "" {
		course.CategoryId = v
	}
	if v, ok := data["level"].(string); ok && v != "" {
		course.Level = v
	}

	if err := s.repo.db.Save(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (s *CourseService) TogglePublishCourse(courseId uint, isPublish bool) (*models.Course, error) {
	var course models.Course
	if err := s.repo.db.First(&course, courseId).Error; err != nil {
		return nil, err
	}

	fmt.Println("isPublish", isPublish)
	if isPublish {
		course.Status = "publish"
	} else {
		course.Status = "draft"
	}

	if err := s.repo.db.Save(&course).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (s *CourseService) GetCourseByID(courseId uint, userId uint, role string) (*models.Course, error) {
	course, err := s.repo.FindByID(courseId)
	if err != nil {
		return nil, err
	}
	if role == string(guards.RoleAdmin) || role == string(guards.RoleSysAdmin) || course.LecturerId == userId {
		return course, nil
	}

	if role == string(guards.RoleLearner) {
		_, err := s.enrollRepo.FindByUserIDAndCourseID(userId, courseId)
		if err == nil {
			return course, nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("forbidden: you are not enrolled in this course")
		}
		return nil, err
	}
	return nil, errors.New("forbidden: you do not have permission to view this course")
}
