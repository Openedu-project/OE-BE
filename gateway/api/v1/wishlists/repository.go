package wishlists

import (
	"errors"

	"gateway/models"

	"gorm.io/gorm"
)

type WishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (r *WishlistRepository) Create(wishlistItem *models.Wishlist) error {
	var course models.Course

	if err := r.db.First(&course, wishlistItem.CourseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}
	if err := r.db.Create(wishlistItem).Error; err != nil {
		return err
	}
	return nil
}

func (r *WishlistRepository) FindByUserID(userID uint, offset int, limit int) ([]models.Wishlist, error) {
	var wishlistItems []models.Wishlist
	err := r.db.Preload("Course").Preload("Course.Lecturer").Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&wishlistItems).Error
	if err != nil {
		return nil, err
	}
	return wishlistItems, nil
}

func (r *WishlistRepository) DeleteByUserIDAndCourseID(userID uint, courseID uint) (int64, error) {
	result := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).Delete(&models.Wishlist{})

	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *WishlistRepository) Update(wishlistItem *models.Wishlist) error {
	return r.db.Save(wishlistItem).Error
}

func (r *WishlistRepository) FindByUserIDAndCourseIDIncludingDeleted(userID uint, courseID uint) (*models.Wishlist, error) {
	var wishlistItem models.Wishlist

	err := r.db.Unscoped().Where("user_id = ? AND course_id = ?", userID, courseID).First(&wishlistItem).Error
	if err != nil {
		return nil, err
	}
	return &wishlistItem, nil
}

func (r *WishlistRepository) Restore(wishlistItem *models.Wishlist) error {
	return r.db.Unscoped().Model(wishlistItem).Update("deleted_at", nil).Error
}
