package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/client"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/repository"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/pb"
	"github.com/google/uuid"
)

type TeamService interface {
	CreateTeam(ctx context.Context, added_by uuid.UUID, input *dto.CreateTeamInput) error
	AddManager(ctx context.Context, createdByID, teamID uuid.UUID, managerIDs []uuid.UUID) error
	RemoveManager(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID) error
	AddMember(ctx context.Context, createdByID, teamID uuid.UUID, managerIDs []uuid.UUID) error
	RemoveMember(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID) error
}

type teamService struct {
	repo       repository.TeamRepositorty
	userClient client.UserGRPCClient
}

func NewTeamService(repo repository.TeamRepositorty, userClient client.UserGRPCClient) TeamService {
	return &teamService{
		repo:       repo,
		userClient: userClient,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, createdByID uuid.UUID, input *dto.CreateTeamInput) error {
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

	if err := s.AddManager(ctx, createdByID, newTeamId, input.Managers); err != nil {
		joinedErr = errors.Join(joinedErr, fmt.Errorf("add managers failed: %w", err))
	}

	if input.Members != nil {
		if err := s.AddMember(ctx, createdByID, newTeamId, input.Members); err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("add members failed: %w", err))
		}
	}

	if joinedErr != nil {
		return joinedErr
	}

	return nil
}

func (s *teamService) AddManager(ctx context.Context, createdByID, teamID uuid.UUID, managerIDs []uuid.UUID) error {
	for _, managerID := range managerIDs {
		resp, err := s.userClient.Client.IsUserExist(ctx, &pb.GetUserRequest{UserId: managerID.String()})
		if err != nil {
			return fmt.Errorf("check user existence failed: %w", err)
		}
		if !resp.IsExist {
			return fmt.Errorf("user with ID %s does not exist", managerID)
		}

		if err := s.repo.AddMember(ctx, createdByID, teamID, managerID); err != nil {
			return fmt.Errorf("add manager failed: %w", err)
		}
	}

	return nil
}

func (s *teamService) RemoveManager(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID) error {
	if err := s.repo.RemoveManager(ctx, teamID, managerID); err != nil {
		return fmt.Errorf("remove managers failed: %w", err)
	}

	return nil
}

func (s *teamService) AddMember(ctx context.Context, createdByID, teamID uuid.UUID, memberIDs []uuid.UUID) error {

	for _, memberID := range memberIDs {
		resp, err := s.userClient.Client.IsUserExist(ctx, &pb.GetUserRequest{UserId: memberID.String()})
		if err != nil {
			return fmt.Errorf("check user existence failed: %w", err)
		}
		if !resp.IsExist {
			return fmt.Errorf("user with ID %s does not exist", memberID)
		}

		if err := s.repo.AddMember(ctx, createdByID, teamID, memberID); err != nil {
			return fmt.Errorf("add member failed: %w", err)
		}
	}

	return nil
}

func (s *teamService) RemoveMember(ctx context.Context, teamID uuid.UUID, memberID uuid.UUID) error {
	if err := s.repo.RemoveMember(ctx, teamID, memberID); err != nil {
		return fmt.Errorf("remove managers failed: %w", err)
	}

	return nil
}
