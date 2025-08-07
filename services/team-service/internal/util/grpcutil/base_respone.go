package grpcutil

import pb "github.com/Thanhbinh1905/go-training-system/services/team-service/pb/team"

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
