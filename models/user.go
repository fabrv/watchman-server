package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserPayload struct {
	Name     string `json:"name" validate:"required,min=1,max=255"`
	Email    string `json:"email" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=3,max=100"`
	RoleID   uint   `json:"role_id" validate:"required"`
}

type User struct {
	gorm.Model
	Name     string    `json:"name" gorm:"not null;" validate:"required,min=1,max=255"`
	Email    string    `json:"email" gorm:"not null; unique" validate:"required,min=1,max=255"`
	Password string    `json:"password" gorm:"not null" validate:"required,min=3,max=100"`
	RoleID   uint      `json:"role_id" gorm:"TYPE:integer REFERENCES roles(id)"`
	Role     Role      `json:"role"`
	Projects []Project `json:"projects" gorm:"many2many:user_projects"`
	Teams    []Team    `json:"teams" gorm:"many2many:user_teams"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	RoleID    uint      `json:"role_id"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
