package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccessLevel enum type
type AccessLevel string

const (
	AccessLevelRead  AccessLevel = "read"
	AccessLevelWrite AccessLevel = "write"
)

// Folder represents a folder that contains notes
type Folder struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	OwnerID     uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Owner  User          `json:"owner" gorm:"foreignKey:OwnerID"`
	Notes  []Note        `json:"notes,omitempty" gorm:"foreignKey:FolderID"`
	Shares []FolderShare `json:"shares,omitempty" gorm:"foreignKey:FolderID"`
}

// Note represents a note within a folder
type Note struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title     string         `json:"title" gorm:"not null"`
	Body      string         `json:"body"`
	FolderID  uuid.UUID      `json:"folder_id" gorm:"type:uuid;not null"`
	OwnerID   uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Folder Folder      `json:"folder" gorm:"foreignKey:FolderID"`
	Owner  User        `json:"owner" gorm:"foreignKey:OwnerID"`
	Shares []NoteShare `json:"shares,omitempty" gorm:"foreignKey:NoteID"`
}

// FolderShare represents sharing permissions for folders
type FolderShare struct {
	ID         uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FolderID   uuid.UUID   `json:"folder_id" gorm:"type:uuid;not null"`
	UserID     uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	Access     AccessLevel `json:"access" gorm:"type:varchar(10);not null;check:access IN ('read', 'write')"`
	SharedAt   time.Time   `json:"shared_at" gorm:"default:CURRENT_TIMESTAMP"`
	SharedByID uuid.UUID   `json:"shared_by_id" gorm:"type:uuid;not null"`

	// Relationships
	Folder   Folder `json:"folder" gorm:"foreignKey:FolderID"`
	User     User   `json:"user" gorm:"foreignKey:UserID"`
	SharedBy User   `json:"shared_by" gorm:"foreignKey:SharedByID"`

	// Unique constraint to prevent duplicate shares
	// UniqueIndex is defined in the migration or through gorm tags
}

// NoteShare represents sharing permissions for individual notes
type NoteShare struct {
	ID         uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	NoteID     uuid.UUID   `json:"note_id" gorm:"type:uuid;not null"`
	UserID     uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	Access     AccessLevel `json:"access" gorm:"type:varchar(10);not null;check:access IN ('read', 'write')"`
	SharedAt   time.Time   `json:"shared_at" gorm:"default:CURRENT_TIMESTAMP"`
	SharedByID uuid.UUID   `json:"shared_by_id" gorm:"type:uuid;not null"`

	// Relationships
	Note     Note `json:"note" gorm:"foreignKey:NoteID"`
	User     User `json:"user" gorm:"foreignKey:UserID"`
	SharedBy User `json:"shared_by" gorm:"foreignKey:SharedByID"`
}

func (FolderShare) TableName() string {
	return "folder_shares"
}

func (NoteShare) TableName() string {
	return "note_shares"
}

func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

func (n *Note) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

func (fs *FolderShare) BeforeCreate(tx *gorm.DB) error {
	if fs.ID == uuid.Nil {
		fs.ID = uuid.New()
	}
	return nil
}

func (ns *NoteShare) BeforeCreate(tx *gorm.DB) error {
	if ns.ID == uuid.Nil {
		ns.ID = uuid.New()
	}
	return nil
}
