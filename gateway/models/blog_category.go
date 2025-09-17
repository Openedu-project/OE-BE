package models

type BlogCategory struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
