package dto

import "github.com/google/uuid"

type CreateTeamInput struct {
	TeamName  string
	Managers  []uuid.UUID
	Members   []uuid.UUID
	CreatedBy uuid.UUID
}
