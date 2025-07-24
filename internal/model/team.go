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

	// Relationships
	CreatedBy User   `json:"created_by" gorm:"foreignKey:CreatedByID"`
	Managers  []User `json:"managers" gorm:"many2many:team_managers;"`
	Members   []User `json:"members" gorm:"many2many:team_members;"`
}

// TeamManager represents the many-to-many relationship between teams and managers
type TeamManager struct {
	TeamID    uuid.UUID `json:"team_id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key"`
	AddedAt   time.Time `json:"added_at" gorm:"default:CURRENT_TIMESTAMP"`
	AddedByID uuid.UUID `json:"added_by_id" gorm:"type:uuid"`

	Team    Team `json:"team" gorm:"foreignKey:TeamID"`
	User    User `json:"user" gorm:"foreignKey:UserID"`
	AddedBy User `json:"added_by" gorm:"foreignKey:AddedByID"`
}

// TeamMember represents the many-to-many relationship between teams and members
type TeamMember struct {
	TeamID    uuid.UUID `json:"team_id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;primary_key"`
	AddedAt   time.Time `json:"added_at" gorm:"default:CURRENT_TIMESTAMP"`
	AddedByID uuid.UUID `json:"added_by_id" gorm:"type:uuid"`

	Team    Team `json:"team" gorm:"foreignKey:TeamID"`
	User    User `json:"user" gorm:"foreignKey:UserID"`
	AddedBy User `json:"added_by" gorm:"foreignKey:AddedByID"`
}

func (TeamManager) TableName() string {
	return "team_managers"
}

func (TeamMember) TableName() string {
	return "team_members"
}
