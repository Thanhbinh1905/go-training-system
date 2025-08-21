package grpc

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/util/grpcutil"
	teampb "github.com/Thanhbinh1905/go-training-system/services/team-service/pb/team"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type teamgRPCHandler struct {
	teampb.UnimplementedTeamServiceServer
	teamService service.TeamService
}

func NewTeamgRPCHandler(teamService service.TeamService) *teamgRPCHandler {
	return &teamgRPCHandler{
		teamService: teamService,
	}
}

func (h *teamgRPCHandler) GetManagersByTeamID(ctx context.Context, req *teampb.GetUserIDsByTeamIDRequest) (*teampb.GetManagersByTeamIDResponse, error) {
	teamID, err := uuid.Parse(req.GetTeamId())
	if err != nil {
		return &teampb.GetManagersByTeamIDResponse{
			Base: grpcutil.BaseError("Invalid team ID"),
		}, status.Errorf(codes.InvalidArgument, "invalid team ID: %v", err)
	}

	managers, err := h.teamService.GetManagersByTeamID(ctx, teamID)
	if err != nil {
		return &teampb.GetManagersByTeamIDResponse{
			Base: grpcutil.BaseError("Failed to get managers"),
		}, status.Errorf(codes.Internal, "failed to get managers: %v", err)
	}

	managerIDs := make([]string, len(managers))
	for i, m := range managers {
		managerIDs[i] = m.ManagerID.String()
	}

	return &teampb.GetManagersByTeamIDResponse{
		Base:       grpcutil.BaseSuccess("Managers retrieved successfully"),
		ManagerIds: managerIDs,
	}, nil
}

func (h *teamgRPCHandler) GetMembersByTeamID(ctx context.Context, req *teampb.GetUserIDsByTeamIDRequest) (*teampb.GetMembersByTeamIDResponse, error) {
	teamID, err := uuid.Parse(req.GetTeamId())
	if err != nil {
		return &teampb.GetMembersByTeamIDResponse{
			Base: grpcutil.BaseError("Invalid team ID"),
		}, status.Errorf(codes.InvalidArgument, "invalid team ID: %v", err)
	}

	members, err := h.teamService.GetMembersByTeamID(ctx, teamID)
	if err != nil {
		return &teampb.GetMembersByTeamIDResponse{
			Base: grpcutil.BaseError("Failed to get Members"),
		}, status.Errorf(codes.Internal, "failed to get Members: %v", err)
	}

	MemberIDs := make([]string, len(members))
	for i, m := range members {
		MemberIDs[i] = m.MemberID.String()
	}

	return &teampb.GetMembersByTeamIDResponse{
		Base:      grpcutil.BaseSuccess("Members retrieved successfully"),
		MemberIds: MemberIDs,
	}, nil
}

func (h *teamgRPCHandler) GetUsersByTeamID(ctx context.Context, req *teampb.GetUserIDsByTeamIDRequest) (*teampb.GetUserIDsByTeamIDResponse, error) {
	teamID, err := uuid.Parse(req.GetTeamId())
	if err != nil {
		return &teampb.GetUserIDsByTeamIDResponse{
			Base: grpcutil.BaseError("Invalid team ID"),
		}, status.Errorf(codes.InvalidArgument, "invalid team ID: %v", err)
	}

	users, err := h.teamService.GetUsersByTeamID(ctx, teamID)
	if err != nil {
		return &teampb.GetUserIDsByTeamIDResponse{
			Base: grpcutil.BaseError("Failed to get users"),
		}, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	MemberIDs := make([]string, len(users.Members))
	for i, m := range users.Members {
		MemberIDs[i] = m.MemberID.String()
	}
	ManagerIDs := make([]string, len(users.Managers))
	for i, m := range users.Managers {
		ManagerIDs[i] = m.ManagerID.String()
	}

	return &teampb.GetUserIDsByTeamIDResponse{
		Base:       grpcutil.BaseSuccess("Users retrieved successfully"),
		ManagerIds: ManagerIDs,
		MemberIds:  MemberIDs,
	}, nil
}

func (h *teamgRPCHandler) IsUserTeamManager(ctx context.Context, req *teampb.IsUserTeamManagerRequest) (*teampb.IsUserTeamManagerResponse, error) {
	// Parse teamID
	teamID, err := uuid.Parse(req.GetTeamID())
	if err != nil {
		return &teampb.IsUserTeamManagerResponse{
			Base: grpcutil.BaseError("Invalid team ID"),
		}, status.Errorf(codes.InvalidArgument, "invalid team ID: %v", err)
	}

	// Parse userID
	userID, err := uuid.Parse(req.GetUserID())
	if err != nil {
		return &teampb.IsUserTeamManagerResponse{
			Base: grpcutil.BaseError("Invalid user ID"),
		}, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	// Chek is user manager
	isManager, err := h.teamService.IsUserTeamManager(ctx, userID, teamID)
	if err != nil {
		return &teampb.IsUserTeamManagerResponse{
			Base: grpcutil.BaseError("Internal server error"),
		}, status.Errorf(codes.Internal, "failed to check manager: %v", err)
	}

	if !isManager {
		return &teampb.IsUserTeamManagerResponse{
			Base:      grpcutil.BaseSuccess("User is not a manager"),
			IsManager: false,
		}, nil
	}

	return &teampb.IsUserTeamManagerResponse{
		Base:      grpcutil.BaseSuccess("User is a manager"),
		IsManager: true,
	}, nil
}
