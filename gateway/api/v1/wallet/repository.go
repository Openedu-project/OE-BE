package wallets

import "gorm.io/gorm"

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) FindAll() ([]Wallet, error) {
	var users []Wallet
	err := r.db.Find(&users).Error
	return users, err
}

func (r *WalletRepository) FindByID(id uint) (*Wallet, error) {
	var wallet Wallet
	err := r.db.First(&wallet, id).Error
	return &wallet, err
}

func (r *WalletRepository) FindByEmail(email string) (*Wallet, error) {
	var wallet Wallet
	if err := r.db.Where("email = ?", email).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}
