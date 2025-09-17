package newsfeed

import (
	"context"
	"regexp"
	"strings"

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

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")
	reSpace := regexp.MustCompile(`[\s\-]+`)
	slug = reSpace.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}
