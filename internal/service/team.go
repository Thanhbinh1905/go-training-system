package service

import (
	"context"
	"errors"

	"github.com/Thanhbinh1905/go-training-system/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/internal/model"
	"github.com/Thanhbinh1905/go-training-system/internal/repository"
	"github.com/google/uuid"
)

type TeamService interface {
	CreateTeam(ctx context.Context, createdBy uuid.UUID, req *dto.CreateTeamRequest) error
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error)
	GetTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Team, error)
	AddMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error
	AddManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error
	RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
	RemoveManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
}

type teamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{repo: repo}
}

func (s *teamService) CreateTeam(ctx context.Context, createdBy uuid.UUID, req *dto.CreateTeamRequest) error {
	if req.TeamName == "" {
		return errors.New("team name is required")
	}

	teamID := uuid.New()
	team := &model.Team{
		ID:          teamID,
		CreatedByID: createdBy,
		TeamName:    req.TeamName,
	}

	err := s.repo.CreateTeam(ctx, team)
	if err != nil {
		return err
	}

	for _, m := range req.Managers {
		managerUUID, err := uuid.Parse(m.ManagerID)
		if err != nil {
			return err
		}
		if err := s.repo.AddManagerToTeam(ctx, teamID, managerUUID, createdBy); err != nil {
			return err
		}
	}

	for _, m := range req.Members {
		memberUUID, err := uuid.Parse(m.MemberID)
		if err != nil {
			return err
		}
		if err := s.repo.AddMemberToTeam(ctx, teamID, memberUUID, createdBy); err != nil {
			return err
		}
	}

	return nil
}

func (s *teamService) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*model.Team, error) {
	return s.repo.GetTeamByID(ctx, teamID)
}

func (s *teamService) GetTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Team, error) {
	return s.repo.GetTeamsByUserID(ctx, userID)
}
func (s *teamService) AddMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error {
	return s.repo.AddMemberToTeam(ctx, teamID, userID, addedBy)
}
func (s *teamService) AddManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, addedBy uuid.UUID) error {
	return s.repo.AddManagerToTeam(ctx, teamID, userID, addedBy)
}
func (s *teamService) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return s.repo.RemoveMemberFromTeam(ctx, teamID, userID)
}
func (s *teamService) RemoveManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return s.repo.RemoveManagerFromTeam(ctx, teamID, userID)
}
