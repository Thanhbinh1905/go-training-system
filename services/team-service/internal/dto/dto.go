package dto

import "github.com/google/uuid"

type CreateTeamInput struct {
	TeamName  string      `json:"team_name" binding:"required"`
	Managers  []uuid.UUID `json:"managers" binding:"required,dive,required,uuid"`
	Members   []uuid.UUID `json:"members" binding:"dive,uuid"`
	CreatedBy uuid.UUID   `json:"-" binding:"uuid"`
}

type AddManagerInput struct {
	ManagerIDs []uuid.UUID `json:"manager_ids" binding:"required,dive,required,uuid"`
}

type AddMemberInput struct {
	MemberIDs []uuid.UUID `json:"member_ids" binding:"required,dive,required,uuid"`
}
