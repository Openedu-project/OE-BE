package models

import (
	"time"

	"gorm.io/gorm"
)

type LaunchpadStatus string

const (
	LaunchpadUpcoming  LaunchpadStatus = "upcoming"
	LaunchpadFeaturing LaunchpadStatus = "featuring"
	LaunchpadSuccess   LaunchpadStatus = "success"
)

type Launchpad struct {
	gorm.Model
	CourseID     uint            `json:"course_id"`
	Course       *Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Title        string          `gorm:"size:255;not null" json:"title"`
	Description  string          `gorm:"text" json:"description"`
	FundingGoal  float64         `json:"funding_goal"`
	Funded       float64         `json:"funded"`
	Backers      int             `json:"backers"`
	Status       LaunchpadStatus `gorm:"type:varchar(20)" json:"status"`
	Approved     bool            `gorm:"default:false" json:"approved"`
	VotingPlans  []VotingPlan    `gorm:"foreignKey:LaunchpadID" json:"voting_plan,omitempty"`
	NextVotingAt *time.Time      `json:"next_voting_at,omitempty"`
}
