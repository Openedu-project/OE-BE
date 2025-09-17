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
