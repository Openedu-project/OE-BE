package categories

type CreateCourseCategoryDTO struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCourseCategoryDTO struct {
	Name string `json:"name" binding:"required"`
}
