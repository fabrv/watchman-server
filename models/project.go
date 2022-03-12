package models

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null; unique" validate:"required,min=1,max=255"`
	Description string `json:"description" gorm:"not null" validate:"required,min=3,max=100"`
}
