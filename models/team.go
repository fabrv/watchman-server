package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TeamPayload struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=3,max=100"`
}

type Team struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null; unique" validate:"required,min=1,max=255"`
	Description string `json:"description" gorm:"not null" validate:"required,min=3,max=100"`
}

type TeamResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
