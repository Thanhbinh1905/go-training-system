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
	GetManagerIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error)
	GetMemberIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error)

	IsUserTeamManager(ctx context.Context, userID, teamID uuid.UUID) (bool, error)
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

func (r *teamRepositorty) GetManagersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*model.TeamManager, error) {
	var manager []*model.TeamManager
	err := r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&manager).Error
	if err != nil {
		return nil, err
	}
	return manager, nil
}

func (r *teamRepositorty) GetManagerIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	var managerIDs []uuid.UUID
	err := r.db.WithContext(ctx).Model(&model.TeamManager{}).
		Where("team_id = ?", teamID).
		Pluck("user_id", &managerIDs).Error
	if err != nil {
		return nil, err
	}
	return managerIDs, nil
}

func (r *teamRepositorty) GetMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*model.TeamMember, error) {
	var member []*model.TeamMember
	err := r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&member).Error
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *teamRepositorty) GetMemberIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	var memberIDs []uuid.UUID
	err := r.db.WithContext(ctx).Model(&model.TeamMember{}).
		Where("team_id = ?", teamID).
		Pluck("user_id", &memberIDs).Error
	if err != nil {
		return nil, err
	}
	return memberIDs, nil
}

func (r *teamRepositorty) IsUserTeamManager(ctx context.Context, userID, teamID uuid.UUID) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.TeamManager{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
