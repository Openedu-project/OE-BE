package wishlists

import (
	"errors"
	"strings"

	"gateway/models"
)

type WishlistService struct {
	repo *WishlistRepository
}

func NewWishlistService(repo *WishlistRepository) *WishlistService {
	return &WishlistService{repo: repo}
}

func (s *WishlistService) AddToWishlist(userID uint, courseID uint) (*models.Wishlist, error) {
	wishlistItem := &models.Wishlist{
		UserID:   userID,
		CourseID: courseID,
	}

	err := s.repo.Create(wishlistItem)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique contraint") {
			return nil, errors.New("course already in wishlist")
		}
		return nil, err
	}
	return wishlistItem, nil
}
