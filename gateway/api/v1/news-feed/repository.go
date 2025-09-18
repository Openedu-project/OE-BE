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
