package wallets

import (
	"gateway/api/v1/users"

	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model
	UserID            uint       `gorm:"uniqueIndex"`
	AccountID         string     `gorm:"size:128"`          // Implicit account ID (hex of public key)
	PublicKey         string     `gorm:"size:128"`          // ed25519:base64
	EncryptPrivateKey string     `gorm:"size:255" json:"-"` // ed25519:base64; encrypt in production!
	EncryptSeedPhase  string     `gorm:"size:255" json:"-"`
	EncryptSecret     string     `gorm:"size:255" json:"-"`
	User              users.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
}
