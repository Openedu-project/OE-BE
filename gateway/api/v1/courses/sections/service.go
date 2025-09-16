package sections

import (
	"gateway/models"
)

type CourseSectionService struct {
	repo *CourseSectionRepository
}

func NewCourseSectionService(repo *CourseSectionRepository) *CourseSectionService {
	return &CourseSectionService{repo: repo}
}

func (s *CourseSectionService) CreateCourseSection(dto CreateCourseSectionDTO, courseId uint) (*models.CourseSection, error) {
	courseSection := models.CourseSection{Name: dto.Name, CourseID: courseId}
	if err := s.repo.db.Create(&courseSection).Error; err != nil {
		return nil, err
	}
	return &courseSection, nil
}

func (s *CourseSectionService) GetCourseSectionByID(id uint) (*models.CourseSection, error) {
	var courseSection models.CourseSection
	if err := s.repo.db.First(&courseSection, id).Error; err != nil {
		return nil, err
	}
	return &courseSection, nil
}

func (s *CourseSectionService) UpdateCourseSection(id uint, dto UpdateCourseSectionDTO) (*models.CourseSection, error) {
	courseSection, err := s.GetCourseSectionByID(id)
	if err != nil {
		return nil, err
	}
	courseSection.Name = dto.Name
	if err := s.repo.db.Save(courseSection).Error; err != nil {
		return nil, err
	}
	return courseSection, nil
}

func (s *CourseSectionService) DeleteCourseSection(id uint) error {
	courseSection, err := s.GetCourseSectionByID(id)
	if err != nil {
		return err
	}
	if err := s.repo.db.Delete(courseSection).Error; err != nil {
		return err
	}
	return nil
}
