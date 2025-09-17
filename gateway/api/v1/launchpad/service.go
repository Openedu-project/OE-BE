package launchpad

import (
	"errors"
	"time"

	"gateway/models"
)

type LaunchpadService struct {
	repo *LaunchpadRepository
}

func NewLaunchpadService(r *LaunchpadRepository) *LaunchpadService {
	return &LaunchpadService{repo: r}
}

func (s *LaunchpadService) CreateLaunchpad(dto CreateLaunchpadDTO) (*models.Launchpad, error) {
	// Check course exist
	ok, err := s.repo.CourseExists(dto.CourseID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("Course not found")
	}

	// Create Launchpad object
	lp := &models.Launchpad{
		CourseID:    dto.CourseID,
		Title:       dto.Title,
		Description: dto.Description,
		FundingGoal: dto.FundingGoal,
		Funded:      0,
		Backers:     0,
		Approved:    false, // admin will approve
		Status:      models.LaunchpadUpcoming,
	}

	if err := s.repo.Create(lp); err != nil {
		return nil, err
	}

	// create voting plans if provided
	if len(dto.VotingPlans) > 0 {
		var plans []models.VotingPlan
		for _, p := range dto.VotingPlans {
			t, parseErr := parseDateFlexible(p.ScheduleAt)
			if parseErr != nil {
				return nil, parseErr
			}
			plans = append(plans, models.VotingPlan{
				LaunchpadID: lp.ID,
				Step:        p.Step,
				Sections:    p.Sections,
				ScheduleAt:  t,
				Title:       p.Title,
			})
		}
		if err := s.repo.CreateVotingPlans(lp.ID, plans); err != nil {
			return nil, err
		}
		// optionally set next voting date
		if len(plans) > 0 {
			t := plans[0].ScheduleAt
			_ = s.repo.UpdateNextVotingAt(lp.ID, &t)
		}
	}

	return s.repo.FindByID(lp.ID)
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
