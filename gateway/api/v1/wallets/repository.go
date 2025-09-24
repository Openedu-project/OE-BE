package wallets

import (
	"gateway/models"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) FindAll() ([]models.Wallet, error) {
	var users []models.Wallet
	err := r.db.Find(&users).Error
	return users, err
}

func (r *WalletRepository) FindByID(id uint) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.First(&wallet, id).Error
	return &wallet, err
}

func (r *WalletRepository) FindByEmail(email string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("email = ?", email).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) FindByUserId(userId uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userId).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}
