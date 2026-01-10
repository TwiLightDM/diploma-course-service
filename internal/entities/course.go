package entities

import "gorm.io/gorm"

type Course struct {
	Id          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	AccessType  string         `json:"access_type"`
	IsPublished bool           `json:"is_published"`
	OwnerId     string         `json:"owner_id"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
