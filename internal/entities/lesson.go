package entities

import "gorm.io/gorm"

type Lesson struct {
	Id          string
	Title       string
	Description string
	Content     string
	Position    int64
	ModuleId    string
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
