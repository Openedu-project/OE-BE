package newsfeed

import (
	"context"

	"gateway/models"

	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(ctx context.Context, blog *models.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *BlogRepository) GetByID(ctx context.Context, id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := r.db.WithContext(ctx).Preload("Category").Preload("Author").First(&blog, id).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) Update(ctx context.Context, blog *models.Blog) error {
	return r.db.WithContext(ctx).Save(blog).Error
}

func (r *BlogRepository) List(ctx context.Context, filters map[string]interface{}, search string, limit, offset int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	query := r.db.WithContext(ctx).Preload("Category").Preload("Author").Where("status = ?", "published")

	// Filter by category
	if catID, ok := filters["category_id"]; ok {
		query = query.Where("category_id = ?", catID)
	}

	// Filter by author
	if authorID, ok := filters["author_id"]; ok {
		query = query.Where("author_id = ?", authorID)
	}

	// Search in title or description
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}

	// Count total
	if err := query.Model(&models.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

func (r *BlogRepository) GetDetail(ctx context.Context, id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := r.db.WithContext(ctx).Preload("Category").Preload("Author").Where("id = ? AND status = ?", id, "published").First(&blog).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}
