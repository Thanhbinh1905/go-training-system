package kafka

import (
	"time"
)

// Team Event Types
const (
	TeamEventCreated      = "TEAM_CREATED"
	TeamEventMemberAdded  = "MEMBER_ADDED"
	TeamEventMemberRemoved = "MEMBER_REMOVED"
	TeamEventManagerAdded = "MANAGER_ADDED"
	TeamEventManagerRemoved = "MANAGER_REMOVED"
)

// Asset Event Types
const (
	AssetEventFolderCreated = "FOLDER_CREATED"
	AssetEventFolderUpdated = "FOLDER_UPDATED"
	AssetEventFolderDeleted = "FOLDER_DELETED"
	AssetEventNoteCreated   = "NOTE_CREATED"
	AssetEventNoteUpdated   = "NOTE_UPDATED"
	AssetEventNoteDeleted   = "NOTE_DELETED"
	AssetEventFolderShared  = "FOLDER_SHARED"
	AssetEventFolderUnshared = "FOLDER_UNSHARED"
	AssetEventNoteShared    = "NOTE_SHARED"
	AssetEventNoteUnshared  = "NOTE_UNSHARED"
)

// Kafka Topics
const (
	TeamActivityTopic  = "team.activity"
	AssetChangesTopic  = "asset.changes"
)

// BaseEvent represents common fields for all events
type BaseEvent struct {
	EventType  string    `json:"eventType"`
	Timestamp  time.Time `json:"timestamp"`
	PerformedBy string   `json:"performedBy"`
}

// TeamEvent represents team-related events
type TeamEvent struct {
	BaseEvent
	TeamID      string `json:"teamId"`
	TargetUserID string `json:"targetUserId,omitempty"`
}

// AssetEvent represents asset-related events
type AssetEvent struct {
	BaseEvent
	AssetType string `json:"assetType"` // "folder" or "note"
	AssetID   string `json:"assetId"`
	OwnerID   string `json:"ownerId"`
	ActionBy  string `json:"actionBy"`
}

// NewTeamEvent creates a new team event
func NewTeamEvent(eventType, teamID, performedBy, targetUserID string) *TeamEvent {
	return &TeamEvent{
		BaseEvent: BaseEvent{
			EventType:   eventType,
			Timestamp:   time.Now().UTC(),
			PerformedBy: performedBy,
		},
		TeamID:      teamID,
		TargetUserID: targetUserID,
	}
}

// NewAssetEvent creates a new asset event
func NewAssetEvent(eventType, assetType, assetID, ownerID, actionBy string) *AssetEvent {
	return &AssetEvent{
		BaseEvent: BaseEvent{
			EventType:   eventType,
			Timestamp:   time.Now().UTC(),
			PerformedBy: actionBy,
		},
		AssetType: assetType,
		AssetID:   assetID,
		OwnerID:   ownerID,
		ActionBy:  actionBy,
	}
}
