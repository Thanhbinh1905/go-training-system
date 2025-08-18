package dto

import "github.com/google/uuid"

type CreateTeamInput struct {
	TeamName string      `json:"team_name" binding:"required"`
	Managers []uuid.UUID `json:"managers" binding:"required,dive,required,uuid"`
	Members  []uuid.UUID `json:"members" binding:"dive,uuid"`
}

type AddManagerInput struct {
	ManagerIDs []uuid.UUID `json:"manager_ids" binding:"required,dive,required,uuid"`
}

type AddMemberInput struct {
	MemberIDs []uuid.UUID `json:"member_ids" binding:"required,dive,required,uuid"`
}

type TeamManagerResponse struct {
	ManagerID uuid.UUID `json:"manager_id"`
}

type TeamMemberResponse struct {
	MemberID uuid.UUID `json:"member_id"`
}

type TeamUsersResponse struct {
	Managers []*TeamManagerResponse `json:"managers"`
	Members  []*TeamMemberResponse  `json:"members"`
}

type TeamResponse struct {
	ID          uuid.UUID             `json:"id"`
	TeamName    string                `json:"team_name"`
	CreatedByID uuid.UUID             `json:"created_by_id"`
	Managers    []TeamManagerResponse `json:"managers"`
	Members     []TeamMemberResponse  `json:"members"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
	DeletedAt   *string               `json:"deleted_at,omitempty"`
}
