package wishlists

import (
	"errors"

	"gateway/api/v1/enrollments"
	"gateway/models"

	"gorm.io/gorm"
)

type WishlistService struct {
	repo *WishlistRepository
}

func NewWishlistService(repo *WishlistRepository) *WishlistService {
	return &WishlistService{repo: repo}
}

func (s *WishlistService) AddToWishlist(userID uint, courseID uint) (*models.Wishlist, error) {
	existingItem, err := s.repo.FindByUserIDAndCourseIDIncludingDeleted(userID, courseID)
	if err == nil {
		if existingItem.DeletedAt.Valid {
			if err := s.repo.Restore(existingItem); err != nil {
				return nil, err
			}
			existingItem.DeletedAt.Valid = false
			return existingItem, nil
		}
		return nil, errors.New("course already in wishlist")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	wishlistItem := &models.Wishlist{
		UserID:   userID,
		CourseID: courseID,
	}

	if err := s.repo.Create(wishlistItem); err != nil {
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

func (s *WishlistService) RemoveFromWishlist(userID uint, courseID uint) error {
	rowsAffected, err := s.repo.DeleteByUserIDAndCourseID(userID, courseID)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("item not found in wishlist")
	}
	return nil
}
