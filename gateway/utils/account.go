package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/cosmos/go-bip39"
)

type OpenEduAccount struct {
	AccountID  string
	PublicKey  string
	PrivateKey string
}

func CreateImplicitAccount(seedPhrase, secretString string) (*OpenEduAccount, error) {
	// Generate ED25519 key pair from seed phrase and secret string
	publicKey, privateKey, err := GenerateKeyPair(seedPhrase, secretString)
	if err != nil {
		return nil, fmt.Errorf("error generating public key: %v", err)
	}

	// Convert the public key to account ID
	accountID := PublicKeyToAccountID(publicKey)

	// Encode keys in base58 for better readability
	base58PublicKey := base58.Encode(publicKey)
	base58PrivateKey := base58.Encode(privateKey)

	return &OpenEduAccount{
		AccountID:  accountID,
		PublicKey:  base58PublicKey,
		PrivateKey: base58PrivateKey,
	}, nil
}

// GenerateKeyPair generates an ED25519 key pair from a seed phrase and secret string
func GenerateKeyPair(seedPhrase, secretString string) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	// ValidateSingleTransfer and normalize the seed phrase
	if !bip39.IsMnemonicValid(seedPhrase) {
		return nil, nil, fmt.Errorf("invalid seed phrase")
	}

	// Combine seed phrase and secret string
	combined := seedPhrase + secretString

	// Generate seed from the combined string
	seed := sha256.Sum256([]byte(combined))

	// Generate ED25519 key pair from the seed
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	return publicKey, privateKey, nil
}

// PublicKeyToAccountID converts a public key to its corresponding account ID
func PublicKeyToAccountID(publicKey ed25519.PublicKey) string {
	return hex.EncodeToString(publicKey)
}

// GenerateSeedPhraseAndSecret tạo seed phrase (BIP-39) và secret ngẫu nhiên
func GenerateSeedPhraseAndSecret() (string, string, error) {
	// 1. Tạo entropy 128-bit => seed phrase 12 từ
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	// 2. Sinh mnemonic (seed phrase)
	seedPhrase, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// 3. Sinh secret (32 bytes random)
	secretBytes := make([]byte, 32)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate secret: %w", err)
	}
	secret := base64.StdEncoding.EncodeToString(secretBytes)

	return seedPhrase, secret, nil
}
