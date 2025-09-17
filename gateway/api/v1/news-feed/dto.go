package newsfeed

type CreateBlogsRequest struct {
	Title       string `json:"title" bidning:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Thumbnail   string `json:"thumbnail"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Language    string `json:"language"`
}

type UpdateBlogRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Content     *string `json:"content"`
	Thumbnail   *string `json:"thumbnail"`
	CategoryID  *uint   `json:"category_id"`
	Language    *string `json:"language"`
	Status      *string `json:"status"`
}

type BlogResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Thumbnail   string `json:"thumbnail"`
	Category    string `json:"category"`
	CategoryID  uint   `json:"category_id"`
	Author      string `json:"author"`
	AuthorID    uint   `json:"author_id"`
	Language    string `json:"language"`
	Views       uint   `json:"views"`
	Likes       uint   `json:"likes"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at,omitempty"`
	Status      string `json:"status"`
}

type BlogCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PublishBlogRequest struct {
	BlogID uint   `json:"blog_id" binding:"required"`
	Note   string `json:"note"` // Ghi chú gửi admin (nếu cần)
}

type ApproveBlogRequest struct {
	BlogID  uint   `json:"blog_id" binding:"required"`
	Approve bool   `json:"approve"`
	Reason  string `json:"reason"` // Lý do nếu reject
}
