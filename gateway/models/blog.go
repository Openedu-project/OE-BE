package models

import (
	"time"
)

type Blog struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	Title          string       `json:"title"`
	Slug           string       `gorm:"size:255;uniqueIndex;not null"`
	Description    string       `json:"description"`
	Content        string       `json:"content"`
	Thumbnail      string       `json:"thumbnail"`
	ImageDesc      string       `json:"image_desc"`
	CategoryID     uint         `json:"category_id"`
	Category       BlogCategory `gorm:"foreignKey:CategoryID" json:"category"`
	AuthorID       uint         `json:"author_id"`
	Author         User         `gorm:"foreignKey:AuthorID" json:"author"`
	Views          uint         `json:"views"`
	Likes          uint         `json:"likes"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Language       string       `json:"language"`
	Status         string       `gorm:"size:20;default:'draft'"` // draft, pending, published, rejected
	ApprovedByID   *uint
	ApprovedBy     *User
	ApprovedAt     *time.Time
	RejectedReason *string
}
