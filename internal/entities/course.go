package entities

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	Id          string
	Title       string
	Description string
	AccessType  string
	PublishedAt *time.Time
	OwnerId     string
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
