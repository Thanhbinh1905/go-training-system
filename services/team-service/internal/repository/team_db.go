package repository

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamDbRepository struct {
	db *gorm.DB
}

func NewTeamDbRepository(db *gorm.DB) *TeamDbRepository {
	return &TeamDbRepository{db}
}

func (r *TeamDbRepository) Create(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}
func (r *TeamDbRepository) AddMember(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error {
	member := model.TeamMember{
		TeamID:    teamID,
		UserID:    userID,
		AddedByID: added_by,
	}
	return r.db.WithContext(ctx).Create(&member).Error
}

func (r *TeamDbRepository) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&model.TeamMember{}).Error
}

func (r *TeamDbRepository) AddManager(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error {
	manager := model.TeamManager{
		TeamID:    teamID,
		UserID:    userID,
		AddedByID: added_by,
	}
	return r.db.WithContext(ctx).Create(&manager).Error
}

func (r *TeamDbRepository) RemoveManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("team_id = ? AND user_id = ?", teamID, userID).Delete(&model.TeamManager{}).Error
}

func (r *TeamDbRepository) GetManagersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*model.TeamManager, error) {
	var manager []*model.TeamManager
	err := r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&manager).Error
	if err != nil {
		return nil, err
	}
	return manager, nil
}

func (r *TeamDbRepository) GetManagerIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	var managerIDs []uuid.UUID
	err := r.db.WithContext(ctx).Model(&model.TeamManager{}).
		Where("team_id = ?", teamID).
		Pluck("user_id", &managerIDs).Error
	if err != nil {
		return nil, err
	}
	return managerIDs, nil
}

func (r *TeamDbRepository) GetMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*model.TeamMember, error) {
	var member []*model.TeamMember
	err := r.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&member).Error
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *TeamDbRepository) GetMemberIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	var memberIDs []uuid.UUID
	err := r.db.WithContext(ctx).Model(&model.TeamMember{}).
		Where("team_id = ?", teamID).
		Pluck("user_id", &memberIDs).Error
	if err != nil {
		return nil, err
	}
	return memberIDs, nil
}

func (r *TeamDbRepository) IsUserTeamManager(ctx context.Context, userID, teamID uuid.UUID) (bool, error) {
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

func (r *TeamDbRepository) IsTeamExist(ctx context.Context, teamID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Team{}).Where("id = ?", teamID).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
