package newsfeed

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"gateway/models"
)

type BlogService struct {
	repo         *BlogRepository
	categoryRepo *CategoryRepository
}

func NewBlogService(repo *BlogRepository, categoryRepo *CategoryRepository) *BlogService {
	return &BlogService{repo: repo, categoryRepo: categoryRepo}
}

func (s *BlogService) CreateBlog(ctx context.Context, req CreateBlogsRequest, authorID uint) (*models.Blog, error) {
	blog := &models.Blog{
		Title:       req.Title,
		Slug:        generateSlug(req.Title),
		Description: req.Description,
		Content:     req.Content,
		Thumbnail:   req.Thumbnail,
		CategoryID:  req.CategoryID,
		AuthorID:    authorID,
		Language:    req.Language,
		Status:      "draft",
	}
	err := s.repo.CreateBlog(ctx, blog)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *BlogService) RequestPublish(ctx context.Context, blogID uint, userID uint) (*models.Blog, error) {
	blog, err := s.repo.GetBlogByID(ctx, blogID)
	if err != nil {
		return nil, errors.New("Blog not found")
	}

	// Creator chỉ được publish blog của chính họ
	if blog.AuthorID != userID {
		return nil, errors.New("Permission denied: cannot publish others' blog")
	}

	// Status logic
	if blog.Status == "draft" || blog.Status == "" {
		blog.Status = "pending"
		blog.UpdatedAt = time.Now()
	} else if blog.Status == "pending" {
		return nil, errors.New("Blog already pending approval")
	} else {
		return nil, errors.New("Invalid state for publish request")
	}

	if err := s.repo.UpdateBlog(ctx, blog); err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *BlogService) ApproveBlog(ctx context.Context, blogID uint, adminID uint, req ApproveBlogRequest) (*models.Blog, error) {
	blog, err := s.repo.GetBlogByID(ctx, blogID)
	if err != nil {
		return nil, errors.New("Blog not found")
	}

	if req.Approved {
		blog.Status = "published"
		now := time.Now()
		blog.ApprovedByID = &adminID
		blog.ApprovedAt = &now
		blog.RejectedReason = nil
	} else {
		blog.Status = "rejected"
		blog.RejectedReason = &req.Reason
	}

	if err := s.repo.UpdateBlog(ctx, blog); err != nil {
		return nil, err
	}

	return blog, nil
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")
	reSpace := regexp.MustCompile(`[\s\-]+`)
	slug = reSpace.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}

// Listing and detail
func (s *BlogService) ListBlogs(ctx context.Context, filters map[string]interface{}, search string, limit, offset int) ([]models.Blog, int64, error) {
	return s.repo.ListBlogs(ctx, filters, search, limit, offset)
}

func (s *BlogService) GetBlogDetail(ctx context.Context, id uint) (*models.Blog, error) {
	return s.repo.GetBlogDetail(ctx, id)
}

// my Blogs
func (s *BlogService) MyBlogs(ctx context.Context, authorID uint, limit, offset int) ([]models.Blog, int64, error) {
	return s.repo.ListBlogByAuthor(ctx, authorID, limit, offset)
}

func (s *BlogService) UpdateBlog(ctx context.Context, blogID uint, userID uint, req UpdateBlogRequest) (*models.Blog, error) {
	blog, err := s.repo.GetBlogByID(ctx, blogID)
	if err != nil {
		return nil, errors.New("Blog not found")
	}

	// Chỉ tác giả hoặc admin được sửa
	if blog.AuthorID != userID {
		return nil, errors.New("Permission denied: cannot edit others' blog")
	}

	// Update fields
	if req.Title != "" {
		blog.Title = req.Title
		blog.Slug = generateSlug(req.Title)
	}
	if req.Description != "" {
		blog.Description = req.Description
	}
	if req.Content != "" {
		blog.Content = req.Content
	}
	if req.Thumbnail != "" {
		blog.Thumbnail = req.Thumbnail
	}
	if req.CategoryID != 0 {
		blog.CategoryID = req.CategoryID
	}
	if req.Language != "" {
		blog.Language = req.Language
	}

	// Nếu blog đã published thì sửa xong phải revert về pending
	if blog.Status == "published" {
		blog.Status = "pending"
	}

	blog.UpdatedAt = time.Now()

	if err := s.repo.UpdateBlog(ctx, blog); err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *BlogService) DeleteBlog(ctx context.Context, blogID uint, authorID uint) error {
	blog, err := s.repo.GetBlogByID(ctx, blogID)
	if err != nil {
		return errors.New("Blog not found")
	}

	if blog.AuthorID != authorID {
		return errors.New("Permission denied: cannot delete others' blog")
	}

	// Allow delete only when not published
	if blog.Status == "published" {
		return errors.New("Cannot delete published blog")
	}
	if err := s.repo.DeleteBlog(ctx, blog); err != nil {
		return err
	}

	return nil
}

// Category management
func (s *BlogService) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*models.BlogCategory, error) {
	c := &models.BlogCategory{
		Name: req.Name,
		Slug: generateSlug(req.Name),
	}

	if err := s.categoryRepo.CreateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *BlogService) UpdateCategory(ctx context.Context, id uint, req UpdateCategoryRequest) (*models.BlogCategory, error) {
	c, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, errors.New("Category not found")
	}
	if req.Name != "" {
		c.Name = req.Name
		c.Slug = generateSlug(req.Name)
	}
	if req.Description != "" {
		c.Slug = generateSlug(req.Name)
	}
	if err := s.categoryRepo.UpdateCategory(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *BlogService) DeleteCategory(ctx context.Context, id uint) error {
	c, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return errors.New("Category not found")
	}
	if err := s.categoryRepo.DeleteCategory(ctx, c); err != nil {
		return err
	}
	return nil
}

func (s *BlogService) ListCategory(ctx context.Context) ([]models.BlogCategory, error) {
	return s.categoryRepo.ListCategory(ctx)
}
