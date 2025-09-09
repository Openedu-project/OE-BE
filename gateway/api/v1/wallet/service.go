package wallets

import (
	"gateway/utils"
)

type WalletService struct {
	repo *WalletRepository
}

func NewWalletService(repo *WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

func (s *WalletService) CreateNearWallet(userId int) (*Wallet, error) {
	seedPhrase, secret, err := utils.GenerateSeedPhraseAndSecret()
	if err != nil {
		return nil, err
	}
	key := "12345678901234567890123456789012" // 32 bytes
	aes := utils.NewAES()
	// encSeed, _ := aes.Encrypt(key, seedPhrase)
	// encSecret, _ := aes.Encrypt(key, secret)

	// TODO: seedPhrase, secret save to BD
	account, err := utils.CreateImplicitAccount(seedPhrase, secret)

	if err != nil {
		return nil, err
	}

	encryptPrivateKey, _ := aes.Encrypt(key, account.PrivateKey)

	wallet := &Wallet{
		UserID:            uint(userId),
		AccountID:         account.AccountID,
		PublicKey:         account.PublicKey,
		EncryptPrivateKey: encryptPrivateKey,
	}

	if err := s.repo.Create(wallet); err != nil {
		return nil, err
	}

	return wallet, nil

}
