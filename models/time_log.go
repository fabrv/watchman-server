package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TimeLogPayload struct {
	UserID      uint      `json:"user_id" validate:"required"`
	ProjectID   uint      `json:"project_id" validate:"required"`
	TeamID      uint      `json:"team_id" validate:"required"`
	LogTypeID   uint      `json:"log_type_id" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	Description string    `json:"description" validate:"min=0,max=255"`
}

type TimeLog struct {
	gorm.Model
	UserID      uint      `json:"user_id" gorm:"TYPE: integer REFERENCES users(id)" validate:"required"`
	LogTypeID   uint      `json:"log_type_id" gorm:"TYPE: integer REFERENCES log_types(id)" validate:"required"`
	LogType     LogType   `json:"log_type"`
	ProjectID   uint      `json:"project_id" gorm:"TYPE: integer REFERENCES projects(id)" validate:"required"`
	Project     Project   `json:"project"`
	TeamID      uint      `json:"team_id" gorm:"TYPE: integer REFERENCES teams(id)" validate:"required"`
	Team        Team      `json:"team"`
	StartTime   time.Time `json:"start_time" gorm:"not null" validate:"required"`
	EndTime     time.Time `json:"end_time" gorm:"default: null"`
	Description string    `json:"description" gorm:"not null" validate:"required"`
}

type TimeLogResponse struct {
	ID          uint      `json:"id"`
	UserId      uint      `json:"user_id"`
	LogTypeId   uint      `json:"log_type_id"`
	ProjectId   uint      `json:"project_id"`
	TeamId      uint      `json:"team_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Description string    `json:"description"`
}
