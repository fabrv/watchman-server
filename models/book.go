package models

import "github.com/jinzhu/gorm"

type Book struct {
	gorm.Model
	Title  string `json:"title" gorm:"not null; unique" validate:"required,min=1,max=255"`
	Author string `json:"author" gorm:"not null" validate:"required,min=3,max=100"`
	Rating int    `json:"rating" gorm:"not null" validate:"required,min=0,max=5"`
}
