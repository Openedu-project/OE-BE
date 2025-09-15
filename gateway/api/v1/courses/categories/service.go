package categories

import (
	"gateway/models"
)

type CourseCategoryService struct {
	repo *CourseCategoryRepository
}

func NewCourseCategoryService(repo *CourseCategoryRepository) *CourseCategoryService {
	return &CourseCategoryService{repo: repo}
}

func (s *CourseCategoryService) CreateCategory(dto CreateCourseCategoryDTO) (*models.CourseCategory, error) {
	category := models.CourseCategory{Name: dto.Name}
	if err := s.repo.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CourseCategoryService) GetCategoryByID(id uint) (*models.CourseCategory, error) {
	var category models.CourseCategory
	if err := s.repo.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CourseCategoryService) GetAllCategories() ([]models.CourseCategory, error) {
	var categories []models.CourseCategory
	if err := s.repo.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CourseCategoryService) UpdateCategory(id uint, dto UpdateCourseCategoryDTO) (*models.CourseCategory, error) {
	category, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	category.Name = dto.Name
	if err := s.repo.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CourseCategoryService) DeleteCategory(id uint) error {
	category, err := s.GetCategoryByID(id)
	if err != nil {
		return err
	}
	if err := s.repo.db.Delete(category).Error; err != nil {
		return err
	}
	return nil
}
