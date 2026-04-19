package entities

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	Id              string
	Title           string
	Description     string
	AccessType      string
	PublishedAt     *time.Time
	OwnerId         string
	DeletedAt       gorm.DeletedAt `gorm:"index"`
	AmountOfModules int            `gorm:"->"`
	AmountOfLessons int            `gorm:"->"`
}
