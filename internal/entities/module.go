package entities

import "gorm.io/gorm"

type Module struct {
	Id          string
	Title       string
	Description string
	Position    int64
	CourseId    string
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
