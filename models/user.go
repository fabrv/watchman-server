package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string    `json:"name" gorm:"not null; unique" validate:"required,min=1,max=255"`
	Email    string    `json:"email" gorm:"not null; unique" validate:"required,min=1,max=255"`
	Password string    `json:"password" gorm:"not null" validate:"required,min=3,max=100"`
	RoleID   uint      `json:"role_id" gorm:"not null"`
	Role     Role      `json:"role" gorm:"foreignkey:RoleID"`
	Projects []Project `json:"projects" gorm:"many2many:user_projects"`
	Teams    []Team    `json:"teams" gorm:"many2many:user_teams"`
}
