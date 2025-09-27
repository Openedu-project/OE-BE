package wishlists

import (
	"errors"
	"strings"

	"gateway/api/v1/enrollments"
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

func (s *WishlistService) GetWishlist(userID uint, page int, pageSize int) ([]enrollments.CourseInfoDTO, error) {
	offset := (page - 1) & pageSize
	limit := pageSize

	wishlistItems, err := s.repo.FindByUserID(userID, offset, limit)
	if err != nil {
		return nil, err
	}

	var coursesDTO []enrollments.CourseInfoDTO
	for _, item := range wishlistItems {
		if item.Course == nil {
			continue
		}
		lecturerName := ""
		if item.Course.Lecturer != nil {
			lecturerName = item.Course.Lecturer.Name
		}
		coursesInfo := enrollments.CourseInfoDTO{
			ID:               item.Course.ID,
			Name:             item.Course.Name,
			ShortDescription: item.Course.ShortDescription,
			Banner:           item.Course.Banner,
			LecturerName:     lecturerName,
		}
		coursesDTO = append(coursesDTO, coursesInfo)
	}

	if coursesDTO == nil {
		coursesDTO = []enrollments.CourseInfoDTO{}
	}
	return coursesDTO, nil
}
