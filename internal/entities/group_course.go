package entities

import "gorm.io/gorm"

type GroupCourse struct {
	Id        string
	CourseId  string
	GroupId   string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
