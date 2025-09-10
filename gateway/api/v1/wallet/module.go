package wallets

import (
	"gateway/configs"
)

// Export instance
var (
	WalletRepo *WalletRepository
	WalletSvc  *WalletService
)

func InitModule() {
	db := configs.DB
	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&Wallet{})
	}

	WalletRepo := NewWalletRepository(db)
	NewWalletService(WalletRepo)
}
