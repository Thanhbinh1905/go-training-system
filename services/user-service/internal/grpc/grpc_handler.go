package grpc

import (
	"context"
	"time"

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
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}
	user, err := h.userService.User(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			UserId:    user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			Role:      convertRole(string(user.Role)),
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, err

}

func convertRole(roleStr string) pb.Role {
	switch roleStr {
	case "MANAGER":
		return pb.Role_MANAGER
	case "MEMBER":
		return pb.Role_MEMBER
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}
