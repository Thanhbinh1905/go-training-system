package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/client"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/repository"
	userpb "github.com/Thanhbinh1905/go-training-system/services/team-service/pb/user"
	"github.com/Thanhbinh1905/go-training-system/shared/contextkey"
	"github.com/Thanhbinh1905/go-training-system/shared/kafka"
	"github.com/google/uuid"
)

type TeamService interface {
	CreateTeam(ctx context.Context, input *dto.CreateTeamInput) error
	AddManager(ctx context.Context, teamID uuid.UUID, managerIDs []uuid.UUID) error
	RemoveManager(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID) error
	AddMember(ctx context.Context, teamID uuid.UUID, managerIDs []uuid.UUID) error
	RemoveMember(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID) error
	GetManagersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*dto.TeamManagerResponse, error)
	GetMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*dto.TeamMemberResponse, error)
	GetUsersByTeamID(ctx context.Context, teamID uuid.UUID) (*dto.TeamUsersResponse, error)

	IsUserTeamManager(ctx context.Context, userID, teamID uuid.UUID) (bool, error)

	IsTeamExist(ctx context.Context, teamID uuid.UUID) (bool, error)
}

type teamService struct {
	repo       repository.TeamRepository
	userClient client.UserGRPCClient
	kafkaProducer kafka.Producer
}

func NewTeamService(repo repository.TeamRepository, userClient client.UserGRPCClient, kafkaProducer kafka.Producer) TeamService {
	return &teamService{
		repo:       repo,
		userClient: userClient,
		kafkaProducer: kafkaProducer,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, input *dto.CreateTeamInput) error {
	userID, err := contextkey.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id from context: %w", err)
	}

	newTeamId := uuid.New()

	team := &model.Team{
		ID:          newTeamId,
		TeamName:    input.TeamName,
		CreatedByID: userID,
	}

	if err := s.repo.Create(ctx, team); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	// Emit TEAM_CREATED event
	teamEvent := kafka.NewTeamEvent(kafka.TeamEventCreated, newTeamId.String(), userID.String(), "")
	if err := s.kafkaProducer.PublishTeamEvent(ctx, teamEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish team created event: %v\n", err)
	}

	var joinedErr error

	input.Managers = append(input.Managers, userID)

	if err := s.AddManager(ctx, newTeamId, input.Managers); err != nil {
		joinedErr = errors.Join(joinedErr, fmt.Errorf("add managers failed: %w", err))
	}

	if input.Members != nil {
		if err := s.AddMember(ctx, newTeamId, input.Members); err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("add members failed: %w", err))
		}
	}

	if joinedErr != nil {
		return joinedErr
	}

	return nil
}

func (s *teamService) AddManager(ctx context.Context, teamID uuid.UUID, managerIDs []uuid.UUID) error {
	userID, err := contextkey.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id from context: %w", err)
	}

	for _, managerID := range managerIDs {
		resp, err := s.userClient.Client.IsUserExist(ctx, &userpb.GetUserRequest{UserId: managerID.String()})
		if err != nil {
			return fmt.Errorf("check user existence failed: %w", err)
		}
		if !resp.IsExist {
			return fmt.Errorf("user with ID %s does not exist", managerID)
		}

		if err := s.repo.AddManager(ctx, userID, teamID, managerID); err != nil {
			return fmt.Errorf("add manager failed: %w", err)
		}

		// Emit MANAGER_ADDED event
		teamEvent := kafka.NewTeamEvent(kafka.TeamEventManagerAdded, teamID.String(), userID.String(), managerID.String())
		if err := s.kafkaProducer.PublishTeamEvent(ctx, teamEvent); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to publish manager added event: %v\n", err)
		}
	}

	return nil
}

func (s *teamService) RemoveManager(ctx context.Context, teamID uuid.UUID, managerID uuid.UUID) error {
	userID, err := contextkey.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id from context: %w", err)
	}

	if err := s.repo.RemoveManager(ctx, teamID, managerID); err != nil {
		return fmt.Errorf("remove managers failed: %w", err)
	}

	// Emit MANAGER_REMOVED event
	teamEvent := kafka.NewTeamEvent(kafka.TeamEventManagerRemoved, teamID.String(), userID.String(), managerID.String())
	if err := s.kafkaProducer.PublishTeamEvent(ctx, teamEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish manager removed event: %v\n", err)
	}

	return nil
}

func (s *teamService) AddMember(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	userID, err := contextkey.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id from context: %w", err)
	}

	for _, memberID := range memberIDs {
		resp, err := s.userClient.Client.IsUserExist(ctx, &userpb.GetUserRequest{UserId: memberID.String()})
		if err != nil {
			return fmt.Errorf("check user existence failed: %w", err)
		}
		if !resp.IsExist {
			return fmt.Errorf("user with ID %s does not exist", memberID)
		}

		if err := s.repo.AddMember(ctx, userID, teamID, memberID); err != nil {
			return fmt.Errorf("add member failed: %w", err)
		}

		// Emit MEMBER_ADDED event
		teamEvent := kafka.NewTeamEvent(kafka.TeamEventMemberAdded, teamID.String(), userID.String(), memberID.String())
		if err := s.kafkaProducer.PublishTeamEvent(ctx, teamEvent); err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to publish member added event: %v\n", err)
		}
	}

	return nil
}

func (s *teamService) RemoveMember(ctx context.Context, teamID uuid.UUID, memberID uuid.UUID) error {
	userID, err := contextkey.GetUserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id from context: %w", err)
	}

	if err := s.repo.RemoveMember(ctx, teamID, memberID); err != nil {
		return fmt.Errorf("remove managers failed: %w", err)
	}

	// Emit MEMBER_REMOVED event
	teamEvent := kafka.NewTeamEvent(kafka.TeamEventMemberRemoved, teamID.String(), userID.String(), memberID.String())
	if err := s.kafkaProducer.PublishTeamEvent(ctx, teamEvent); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish member removed event: %v\n", err)
	}

	return nil
}

func (s *teamService) GetManagersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*dto.TeamManagerResponse, error) {
	managers, err := s.repo.GetManagerIDsByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("get managers by team ID failed: %w", err)
	}

	var managerResponses []*dto.TeamManagerResponse
	for _, managerID := range managers {

		managerResponses = append(managerResponses, &dto.TeamManagerResponse{
			ManagerID: managerID,
		})
	}

	return managerResponses, nil
}

func (s *teamService) GetMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*dto.TeamMemberResponse, error) {
	members, err := s.repo.GetMemberIDsByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("get members by team ID failed: %w", err)
	}

	var memberResponses []*dto.TeamMemberResponse
	for _, memberID := range members {

		memberResponses = append(memberResponses, &dto.TeamMemberResponse{
			MemberID: memberID,
		})
	}

	return memberResponses, nil
}

func (s *teamService) GetUsersByTeamID(ctx context.Context, teamID uuid.UUID) (*dto.TeamUsersResponse, error) {
	managers, err := s.GetManagersByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("get managers by team ID failed: %w", err)
	}

	members, err := s.GetMembersByTeamID(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("get members by team ID failed: %w", err)
	}

	return &dto.TeamUsersResponse{
		Managers: managers,
		Members:  members,
	}, nil
}

func (s *teamService) IsUserTeamManager(ctx context.Context, userID, teamID uuid.UUID) (bool, error) {
	return s.repo.IsUserTeamManager(ctx, userID, teamID)
}

func (s *teamService) IsTeamExist(ctx context.Context, teamID uuid.UUID) (bool, error) {
	return s.repo.IsTeamExist(ctx, teamID)
}
