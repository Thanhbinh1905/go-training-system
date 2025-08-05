package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/grpcutil"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/service"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRPCHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewUserRPCHandler(userSvc service.UserService) *UserRPCHandler {
	return &UserRPCHandler{
		userService: userSvc,
	}
}

func (h *UserRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return &pb.GetUserResponse{
			Base: grpcutil.BaseError("Invalid user ID"),
		}, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	user, err := h.userService.User(ctx, userID)
	if err != nil {
		return &pb.GetUserResponse{
			Base: grpcutil.BaseError("Failed to get user"),
		}, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if user == nil {
		return &pb.GetUserResponse{
			Base: grpcutil.BaseError("User not found"),
		}, nil
	}

	return &pb.GetUserResponse{
		Base: grpcutil.BaseSuccess("User retrieved successfully"),
		User: &pb.User{
			UserId:    user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			Role:      grpcutil.ConvertRole(string(user.Role)),
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *UserRPCHandler) IsUserExist(ctx context.Context, req *pb.GetUserRequest) (*pb.IsUserExistResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return &pb.IsUserExistResponse{
			Base: grpcutil.BaseError("Invalid user ID"),
		}, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	exist, err := h.userService.CheckUserExist(ctx, userID)
	if err != nil {
		return &pb.IsUserExistResponse{
			Base: grpcutil.BaseError("Failed to check user existence"),
		}, status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}

	return &pb.IsUserExistResponse{
		IsExist: exist,
		Base:    grpcutil.BaseSuccess("User existence checked successfully"),
	}, nil
}

func (h *UserRPCHandler) VerifyAccessToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if req.GetAccessToken() == "" {
		return &pb.VerifyTokenResponse{
			Base:    grpcutil.BaseError("Access token is required"),
			IsValid: false,
		}, status.Errorf(codes.InvalidArgument, "access token is required")
	}

	valid, err := h.userService.ValidateToken(&dto.TokenVerifyInput{Token: req.GetAccessToken()})
	if err != nil {
		return &pb.VerifyTokenResponse{
			Base:    grpcutil.BaseError("Failed to verify access token"),
			IsValid: false,
		}, status.Errorf(codes.Internal, "failed to verify access token: %v", err)
	}

	return &pb.VerifyTokenResponse{
		Base:    grpcutil.BaseSuccess("Access token verified successfully"),
		IsValid: valid.IsValid,
		UserInfo: &pb.UserInfo{
			UserId: valid.User.ID.String(),
			Role:   grpcutil.ConvertRole(valid.User.Role),
		},
	}, nil
}

func (h *UserRPCHandler) GetUsersByIds(ctx context.Context, req *pb.GetUsersByIdsRequest) (*pb.GetUsersByIdsResponse, error) {
	var users []*pb.GetUserResponse
	var failedCount int

	for _, id := range req.GetUserIds() {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			failedCount++
			continue
		}

		user, err := h.userService.User(ctx, parsedID)
		if err != nil || user == nil {
			failedCount++
			continue
		}

		users = append(users, &pb.GetUserResponse{
			User: &pb.User{
				UserId:    user.ID.String(),
				Username:  user.Username,
				Email:     user.Email,
				Role:      grpcutil.ConvertRole(string(user.Role)),
				CreatedAt: user.CreatedAt.Format(time.RFC3339),
			},
		})
	}

	if len(users) == 0 {
		return &pb.GetUsersByIdsResponse{
			Base:  grpcutil.BaseError("No users found or all IDs invalid"),
			Users: []*pb.GetUserResponse{},
		}, nil
	}

	message := "Users retrieved successfully"
	if failedCount > 0 {
		message = fmt.Sprintf("%d users retrieved, %d failed", len(users), failedCount)
	}

	return &pb.GetUsersByIdsResponse{
		Base:  grpcutil.BaseSuccess(message),
		Users: users,
	}, nil
}
