package client

import (
	"log"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/pb"
	"google.golang.org/grpc"
)

type UserGRPCClient struct {
	Client pb.UserServiceClient
}

func NewUserGRPCClient(addr string) *UserGRPCClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user-service gRPC: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return &UserGRPCClient{Client: client}
}
