package newsfeed

type CreateBlogsRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Thumbnail   string `json:"thumbnail"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Language    string `json:"language"`
}

type RequestPublishBlogResponse struct {
	Message string      `json:"message"`
	Blog    interface{} `json:"blog"`
}

type ApproveBlogRequest struct {
	Approved bool   `json:"approved" binding:"required"`
	Reason   string `json:"reason"`
}

type ApproveBlogResponse struct {
	Message string      `json:"message"`
	Blog    interface{} `json:"blog"`
}

type UpdateBlogRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Thumbnail   string `json:"thumbnail"`
	CategoryID  uint   `json:"category_id"`
	Language    string `json:"language"`
}

type DeleteBlogResponse struct {
	Message string `json:"message"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
