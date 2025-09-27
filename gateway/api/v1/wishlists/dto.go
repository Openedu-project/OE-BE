package wishlists

type AddToWishlistDTO struct {
	CourseID uint `json:"course_id" binding:"required"`
}
