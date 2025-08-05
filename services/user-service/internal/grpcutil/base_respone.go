package grpcutil

import "github.com/Thanhbinh1905/go-training-system/services/user-service/pb"

func BaseSuccess(message string) *pb.BaseResponse {
	return &pb.BaseResponse{
		Success: true,
		Message: message,
	}
}

func BaseError(message string) *pb.BaseResponse {
	return &pb.BaseResponse{
		Success: false,
		Message: message,
	}
}

func ConvertRole(roleStr string) pb.Role {
	switch roleStr {
	case "MANAGER":
		return pb.Role_MANAGER
	case "MEMBER":
		return pb.Role_MEMBER
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}
