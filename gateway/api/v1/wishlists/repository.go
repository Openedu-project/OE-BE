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
