package models

import (
	"time"

	"gorm.io/gorm"
)

type VotingPlan struct {
	gorm.Model
	LaunchpadID uint      `json:"launchpad_id"`
	Step        int       `json:"step"`     // order: 1,2,3...
	Sections    int       `json:"sections"` // Số phần trong bước nhảy
	ScheduleAt  time.Time `json:"schedule_at"`
	Title       string    `gorm:"size:255" json:"title"`
}
