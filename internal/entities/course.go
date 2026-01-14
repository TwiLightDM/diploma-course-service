package entities

import "gorm.io/gorm"

type Course struct {
	Id          string
	Title       string
	Description string
	AccessType  string
	IsPublished bool
	OwnerId     string
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
