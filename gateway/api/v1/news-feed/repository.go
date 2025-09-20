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

func (r *BlogRepository) CreateBlog(ctx context.Context, blog *models.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *BlogRepository) GetBlogByID(ctx context.Context, id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := r.db.WithContext(ctx).Preload("Category").Preload("Author").First(&blog, id).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) UpdateBlog(ctx context.Context, blog *models.Blog) error {
	return r.db.WithContext(ctx).Save(blog).Error
}

func (r *BlogRepository) DeleteBlog(ctx context.Context, blog *models.Blog) error {
	return r.db.WithContext(ctx).Delete(blog).Error
}

func (r *BlogRepository) ListBlogs(ctx context.Context, filters map[string]interface{}, search string, limit, offset int) ([]models.Blog, int64, error) {
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

func (r *BlogRepository) GetBlogDetail(ctx context.Context, id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := r.db.WithContext(ctx).Preload("Category").Preload("Author").Where("id = ? AND status = ?", id, "published").First(&blog).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) ListBlogByAuthor(ctx context.Context, authorID uint, limit, offset int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	query := r.db.WithContext(ctx).Preload("Category").Preload("Author").Where("author_id = ", authorID)

	if err := query.Model(&models.Blog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("createed_at DESC").Find(&blogs).Error; err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
}

// Category
type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, c *models.BlogCategory) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id uint) (*models.BlogCategory, error) {
	var c models.BlogCategory
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, c *models.BlogCategory) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, c *models.BlogCategory) error {
	return r.db.WithContext(ctx).Delete(c).Error
}

func (r *CategoryRepository) ListCategory(ctx context.Context) ([]models.BlogCategory, error) {
	var cats []models.BlogCategory
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}
