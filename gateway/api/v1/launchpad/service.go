package launchpad

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gateway/api/v1/wallets"
	"gateway/configs"
	"gateway/models"
	"gateway/utils"

	"github.com/aurora-is-near/near-api-go"
)

type LaunchpadService struct {
	repo       *LaunchpadRepository
	walletRepo *wallets.WalletRepository
}

func NewLaunchpadService(r *LaunchpadRepository, walletRepo *wallets.WalletRepository) *LaunchpadService {
	return &LaunchpadService{repo: r, walletRepo: walletRepo}
}

func (s *LaunchpadService) CreateLaunchpad(dto CreateLaunchpadDTO, userId uint) (*models.Launchpad, error) {
	ok, err := s.repo.CourseExists(dto.CourseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("Course not found")
	}

	lp := &models.Launchpad{
		CourseID:    dto.CourseID,
		Title:       dto.Title,
		Description: dto.Description,
		FundingGoal: dto.FundingGoal,
		Funded:      0,
		Backers:     0,
		Approved:    false,
		Status:      models.LaunchpadUpcoming,
	}

	if err := s.repo.Create(lp); err != nil {
		return nil, err
	}

	wallet, err := s.walletRepo.FindByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("wallet not found for user %d: %w", userId, err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet is nil for user %d", userId)
	}

	key := configs.Env.AESSecret
	aes := utils.NewAES()
	decryptPrivateKey, err := aes.Decrypt(key, wallet.EncryptPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	pk, err := utils.Ed25519PrivateKeyFromString("ed25519:" + decryptPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key format: %w", err)
	}

	nodeURL := "https://rpc.testnet.near.org" // better to use config
	connection := near.NewConnection(nodeURL)

	account := near.LoadAccountWithPrivateKey(connection, wallet.AccountID, pk)
	if account == nil {
		return nil, fmt.Errorf("failed to load NEAR account object for %s", account.AccountID())
	}

	args := map[string]any{
		"token_id":            "ft_1.nvvan.testnet",
		"campaign_id":         RandomCampaignID(),
		"target_funding":      "3000000000", // 3000 USDT (24 decimals)
		"min_multiple_pledge": 15000000,     // 15 USDT (24 decimals)
	}

	// Encode sang JSON
	argsBytes, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal args: %w", err)
	}

	gas := uint64(300_000_000_000_000) //
	deposit, _ := new(big.Int).SetString("1000000000000000000000000", 10)

	res, err := account.FunctionCall(
		"launchpad_2.nvvan.testnet",
		"init_pool",
		argsBytes,
		gas,
		*deposit,
	)
	if err != nil {
		return nil, fmt.Errorf("rpc error: %w", err)
	}

	// ép kiểu về map
	status, ok := res["status"].(map[string]interface{})
	if ok {
		if failure, exists := status["Failure"]; exists {
			return nil, fmt.Errorf("contract execution failed: %+v", failure)
		}
		if success, exists := status["SuccessValue"]; exists {
			// thường là base64 string
			if s, ok := success.(string); ok && s != "" {
				decoded, derr := base64.StdEncoding.DecodeString(s)
				if derr == nil {
					fmt.Println("Contract returned:", string(decoded))
				}
			}
		}
	}

	// ... voting plans logic unchanged ...

	return lp, nil
}

func (s *LaunchpadService) GetLaunchpadByID(id uint) (*models.Launchpad, error) {
	return s.repo.FindByID(id)
}

func (s *LaunchpadService) GetAllLaunchpadHome() (map[string][]models.Launchpad, error) {
	launchpads, err := s.repo.FindAllLaunchpadHome()
	if err != nil {
		return nil, err
	}

	groups := map[string][]models.Launchpad{
		"featuring": {},
		"upcoming":  {},
		"success":   {},
	}

	for _, lp := range launchpads {
		groups[string(lp.Status)] = append(groups[string(lp.Status)], lp)
	}
	return groups, nil
}

// func (s *LaunchpadService) GetLaunchpads() ([]models.Launchpad, error) {
// 	return s.repo.FindAll(true)
// }

func (s *LaunchpadService) ApproveLaunchpad(id uint) (*models.Launchpad, error) {
	// Find the launchpad first
	launchpad, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// update the approved status
	launchpad.Approved = true

	// save the changes
	if err := s.repo.Update(launchpad); err != nil {
		return nil, err
	}

	return launchpad, nil
}

// Goal will be input by investor
func calculateStatus(goal, funded float64) models.LaunchpadStatus {
	if funded <= 0 || goal <= 0 {
		return models.LaunchpadUpcoming
	}
	percent := (funded / goal) * 100
	if percent >= 80 {
		return models.LaunchpadSuccess
	}
	return models.LaunchpadFeaturing
}

func parseDateFlexible(s string) (time.Time, error) {
	// try ISO YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}

	// try M/D/YYYY like "3/6/2025"
	if t, err := time.Parse("2/1/2006", s); err == nil {
		return t, nil
	}

	// try RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("invalid date format, use YYYY-MM-DD")
}

func (s *LaunchpadService) InvestInLaunchpad(userID uint, launchpadID uint, amount float64) (*models.Launchpad, error) {
	// Get launchpad by ID
	launchpad, err := s.repo.FindByID(launchpadID)
	if err != nil {
		return nil, err
	}

	// check approved launchpad
	if !launchpad.Approved {
		return nil, errors.New("Cannot invest in a launchpad that is not approved")
	}

	// Phần trừ tiền trong wallet của user
	// Khi WalletService sẵn sàng, thì inject vào đây
	// _, err := s.WalletService.Debit(userID, amount)
	// if err != nil {
	// 	return nil, err
	// }

	// Check if user has already invested
	existingInvestment, err := s.repo.FindInvestment(userID, launchpadID)
	if err != nil {
		return nil, err
	}

	if existingInvestment == nil {
		// New investor
		launchpad.Backers += 1
		// Create new investment record
		newInvestment := &models.LaunchpadInvestment{
			UserID:      userID,
			LaunchpadID: launchpadID,
			Amount:      amount,
		}
		if err := s.repo.CreateInvestment(newInvestment); err != nil {
			return nil, err
		}
	} else {
		// Existing investor, just update amount
		existingInvestment.Amount += amount
		if err := s.repo.UpdateInvestment(existingInvestment); err != nil {
			return nil, err
		}
	}

	// Update total funded amount for the launchpad
	launchpad.Funded += amount

	// Update launchpad status based on new funds
	launchpad.Status = calculateStatus(launchpad.FundingGoal, launchpad.Funded)

	// Save changes
	if err := s.repo.Update(launchpad); err != nil {
		return nil, err
	}

	return launchpad, nil
}

func RandomCampaignID() string {
	b := make([]byte, 12) // độ dài 12 byte ~ 16 ký tự base64
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
