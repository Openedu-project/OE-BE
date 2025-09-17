package launchpad

// Create Launchpad include optional VotingPlans
type CreateLaunchpadDTO struct {
	CourseID    uint                  `json:"course_id" binding:"required"`
	Title       string                `json:"title" binding:"required,max=255"`
	Description string                `json:"description"`
	FundingGoal float64               `json:"funding_goal" binding:"required,gt=0"`
	VotingPlans []VotingGoalCreateDTO `json:"voting_plan"`
}

type VotingGoalCreateDTO struct {
	Step     int `json:"step" binding:"required"`
	Sections int `json:"sections" binding:"required"`
	// Date string in "YYYY-MM-DD" (preferred) or "M/D/YYYY" (accepted)
	ScheduleAt string `json:"schedule_at" binding:"required"`
	Title      string `json:"title"`
}

type InvestDTO struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
