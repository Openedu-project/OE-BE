package sections

type CreateCourseSectionDTO struct {
	Name string `json:"name" binding:"required,max=255"`
}

type UpdateCourseSectionDTO struct {
	Name   string `json:"name" binding:"required,max=255"`
	Status string `json:"status" binding:"max=50"`
}
