package wallets

import (
	"gateway/configs"
	"gateway/models"
)

// Export instance
var (
	WalletRepo *WalletRepository
	WalletSvc  *WalletService
)

func InitModule() {
	db := configs.DB
	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&models.Wallet{})
	}

	WalletRepo = NewWalletRepository(db)
	WalletSvc = NewWalletService(WalletRepo)
}
