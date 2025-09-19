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
	repo *BlogRepository
}

func NewBlogService(repo *BlogRepository) *BlogService {
	return &BlogService{repo: repo}
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
	err := s.repo.Create(ctx, blog)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *BlogService) RequestPublish(ctx context.Context, blogID uint, userID uint) (*models.Blog, error) {
	blog, err := s.repo.GetByID(ctx, blogID)
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

	if err := s.repo.Update(ctx, blog); err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *BlogService) ApproveBlog(ctx context.Context, blogID uint, adminID uint, req ApproveBlogRequest) (*models.Blog, error) {
	blog, err := s.repo.GetByID(ctx, blogID)
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

	if err := s.repo.Update(ctx, blog); err != nil {
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
