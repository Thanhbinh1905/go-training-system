package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleManager UserRole = "MANAGER"
	UserRoleMember  UserRole = "MEMBER"
)

type User struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	Role         UserRole       `json:"role" gorm:"type:varchar(20);not null;check:role IN ('MANAGER', 'MEMBER')"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	OwnedTeams   []Team        `json:"owned_teams,omitempty" gorm:"foreignKey:CreatedByID"`
	ManagedTeams []Team        `json:"managed_teams,omitempty" gorm:"many2many:team_managers;"`
	MemberTeams  []Team        `json:"member_teams,omitempty" gorm:"many2many:team_members;"`
	OwnedFolders []Folder      `json:"owned_folders,omitempty" gorm:"foreignKey:OwnerID"`
	OwnedNotes   []Note        `json:"owned_notes,omitempty" gorm:"foreignKey:OwnerID"`
	FolderShares []FolderShare `json:"folder_shares,omitempty" gorm:"foreignKey:UserID"`
	NoteShares   []NoteShare   `json:"note_shares,omitempty" gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
