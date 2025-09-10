package wallets

import (
	"gateway/api/v1/users"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	UserID            uint   `gorm:"uniqueIndex"` // one-to-one: mỗi user chỉ có 1 wallet
	AccountID         string `gorm:"size:128"`    // Implicit account ID (hex of public key)
	PublicKey         string `gorm:"size:128"`    // ed25519:base64
	EncryptPrivateKey string `gorm:"size:128"`    // ed25519:base64; encrypt in production!

	User users.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
}
