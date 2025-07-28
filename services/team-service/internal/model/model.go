package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TeamName    string         `json:"team_name" gorm:"not null"`
	CreatedByID uuid.UUID      `json:"created_by_id" gorm:"type:uuid;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type TeamManager struct {
	TeamID    uuid.UUID `json:"team_id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key"`
	AddedAt   time.Time `json:"added_at" gorm:"default:CURRENT_TIMESTAMP"`
	AddedByID uuid.UUID `json:"added_by_id" gorm:"type:uuid"`
}

// TeamMember represents the many-to-many relationship between teams and members
type TeamMember struct {
	TeamID    uuid.UUID `json:"team_id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key"`
	AddedAt   time.Time `json:"added_at" gorm:"default:CURRENT_TIMESTAMP"`
	AddedByID uuid.UUID `json:"added_by_id" gorm:"type:uuid"`
}

func (TeamManager) TableName() string {
	return "team_managers"
}

func (TeamMember) TableName() string {
	return "team_members"
}
