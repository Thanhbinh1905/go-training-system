package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/authclient"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/repository"
	"github.com/google/uuid"
)

type TeamService interface {
	CreateTeam(ctx context.Context, token string, input *dto.CreateTeamInput) error
	AddManager(ctx context.Context, token string, teamID uuid.UUID, managerIDs []uuid.UUID) error
	RemoveManager(ctx context.Context, token string, teamID uuid.UUID, managerID uuid.UUID) error
	AddMember(ctx context.Context, token string, teamID uuid.UUID, managerIDs []uuid.UUID) error
	RemoveMember(ctx context.Context, token string, teamID uuid.UUID, managerID uuid.UUID) error
}

type teamService struct {
	repo       repository.TeamRepositorty
	authClient authclient.AuthServiceClient
}

func NewTeamService(repo repository.TeamRepositorty, authClient authclient.AuthServiceClient) TeamService {
	return &teamService{
		repo:       repo,
		authClient: authClient,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, token string, input *dto.CreateTeamInput) error {
	createdByID, err := s.authorizeManager(ctx, token)
	if err != nil {
		return err
	}

	newTeamId := uuid.New()

	team := &model.Team{
		ID:          newTeamId,
		TeamName:    input.TeamName,
		CreatedByID: createdByID,
	}

	if err := s.repo.Create(ctx, team); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	var joinedErr error

	input.Managers = append(input.Managers, createdByID)

	if err := s.repo.AddManager(ctx, createdByID, newTeamId, input.Managers); err != nil {
		joinedErr = errors.Join(joinedErr, fmt.Errorf("add managers failed: %w", err))
	}

	if err := s.repo.AddMember(ctx, createdByID, newTeamId, input.Members); err != nil {
		joinedErr = errors.Join(joinedErr, fmt.Errorf("add members failed: %w", err))
	}

	if joinedErr != nil {
		return joinedErr
	}

	return nil
}

func (s *teamService) AddManager(ctx context.Context, token string, teamID uuid.UUID, managerIDs []uuid.UUID) error {
	createdByID, err := s.authorizeManager(ctx, token)
	if err != nil {
		return err
	}

	if err := s.repo.AddManager(ctx, createdByID, teamID, managerIDs); err != nil {
		return fmt.Errorf("add managers failed: %w", err)
	}

	return nil
}

func (s *teamService) RemoveManager(ctx context.Context, token string, teamID uuid.UUID, managerID uuid.UUID) error {
	_, err := s.authorizeManager(ctx, token)
	if err != nil {
		return err
	}

	if err := s.repo.RemoveManager(ctx, teamID, managerID); err != nil {
		return fmt.Errorf("remove managers failed: %w", err)
	}

	return nil
}

func (s *teamService) AddMember(ctx context.Context, token string, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	createdByID, err := s.authorizeManager(ctx, token)
	if err != nil {
		return err
	}

	if err := s.repo.AddMember(ctx, createdByID, teamID, memberIDs); err != nil {
		return fmt.Errorf("add managers failed: %w", err)
	}

	return nil
}

func (s *teamService) RemoveMember(ctx context.Context, token string, teamID uuid.UUID, memberID uuid.UUID) error {
	_, err := s.authorizeManager(ctx, token)
	if err != nil {
		return err
	}

	if err := s.repo.RemoveMember(ctx, teamID, memberID); err != nil {
		return fmt.Errorf("remove managers failed: %w", err)
	}

	return nil
}

func (s *teamService) authorizeManager(ctx context.Context, token string) (uuid.UUID, error) {
	resp, err := s.authClient.VerifyToken(ctx, token)
	if err != nil {
		return uuid.Nil, fmt.Errorf("verify token failed: %w", err)
	}

	if !resp.VerifyToken.Valid || resp.VerifyToken.User.ID == uuid.Nil {
		return uuid.Nil, errors.New("unauthorized: invalid token")
	}

	if resp.VerifyToken.User.Role != "MANAGER" {
		return uuid.Nil, errors.New("forbidden: only MANAGER allowed")
	}

	return resp.VerifyToken.User.ID, nil
}
