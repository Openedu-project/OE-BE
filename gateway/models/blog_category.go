package models

type BlogCategory struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"`
	Slug string `gorm:"type:varchar(100);unique;not null" json:"slug"`
}
