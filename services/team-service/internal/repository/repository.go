package repository

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepositorty interface {
	Create(ctx context.Context, team *model.Team) error
	AddMember(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
	AddManager(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error
	RemoveManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
}

type teamRepositorty struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepositorty {
	return &teamRepositorty{db}
}

func (r *teamRepositorty) Create(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}
func (r *teamRepositorty) AddMember(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error {
	member := model.TeamMember{
		TeamID:    teamID,
		UserID:    userID,
		AddedByID: added_by,
	}
	return r.db.WithContext(ctx).Create(&member).Error
}

func (r *teamRepositorty) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&model.TeamMember{}).Error
}

func (r *teamRepositorty) AddManager(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error {
	manager := model.TeamManager{
		TeamID:    teamID,
		UserID:    userID,
		AddedByID: added_by,
	}
	return r.db.WithContext(ctx).Create(&manager).Error
}

func (r *teamRepositorty) RemoveManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&model.TeamManager{}).Error
}
