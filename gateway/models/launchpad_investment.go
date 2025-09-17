package models

import (
	"time"

	"gorm.io/gorm"
)

type LaunchpadInvestment struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	LaunchpadID uint           `json:"launchpad_id"`
	Launchpad   Launchpad      `gorm:"foreignKey:LaunchpadID" json:"launchpad"`
	Amount      float64        `json:"amount"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
