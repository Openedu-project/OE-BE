package certificates

import (
	"gateway/models"

	"gorm.io/gorm"
)

type CertificateRepository struct {
	db *gorm.DB
}

func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}

func (r *CertificateRepository) FindCertificatesByUserID(userID uint) ([]models.Certificate, error) {
	var certificates []models.Certificate
	err := r.db.Preload("Course").Where("user_id = ?", userID).Order("issued_at DESC").Find(&certificates).Error
	if err != nil {
		return nil, err
	}

	return certificates, nil
}
