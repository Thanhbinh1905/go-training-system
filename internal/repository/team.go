package repository

import (
	"context"
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error)
	GetTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Team, error)
	AddMemberToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error
	AddManagerToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error
	RemoveMemberFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
	RemoveManagerFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) CreateTeam(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *teamRepository) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error) {
	var team model.Team
	err := r.db.WithContext(ctx).
		Preload("CreatedBy").
		Preload("Members").
		Preload("Managers").
		First(&team, "id = ?", teamID).Error

	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) GetTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Team, error) {
	var teams []model.Team
	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN team_members ON team_members.team_id = teams.id").
		Joins("LEFT JOIN team_managers ON team_managers.team_id = teams.id").
		Where("team_members.user_id = ? OR team_managers.user_id = ?", userID, userID).
		Preload("Members").
		Preload("Managers").
		Preload("CreatedBy").
		Group("teams.id").
		Find(&teams).Error

	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (r *teamRepository) AddMemberToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}
	if user.Role != model.UserRoleMember {
		return fmt.Errorf("user is not a member")
	}

	var count int64
	r.db.WithContext(ctx).Model(&model.TeamMember{}).
		Where("role = MEMBER AND team_id = ? AND user_id = ?", teamID, userID).
		Count(&count)

	if count == 0 {
		member := model.TeamMember{
			TeamID:    teamID,
			UserID:    userID,
			AddedByID: addedBy,
		}
		return r.db.WithContext(ctx).Create(&member).Error
	}
	return nil
}

func (r *teamRepository) AddManagerToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found")
	}
	if user.Role != model.UserRoleManager {
		return fmt.Errorf("user is not a manager")
	}

	var count int64
	r.db.WithContext(ctx).Model(&model.TeamManager{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count)

	if count == 0 {
		manager := model.TeamManager{
			TeamID:    teamID,
			UserID:    userID,
			AddedByID: addedBy,
		}
		return r.db.WithContext(ctx).Create(&manager).Error
	}
	return nil
}

func (r *teamRepository) RemoveMemberFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&model.TeamMember{}).Error
}

func (r *teamRepository) RemoveManagerFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&model.TeamManager{}).Error
}
