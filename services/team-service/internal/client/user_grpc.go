package client

import (
	"log"

	userpb "github.com/Thanhbinh1905/go-training-system/services/team-service/pb/user"
	"google.golang.org/grpc"
)

type UserGRPCClient struct {
	Client userpb.UserServiceClient
}

func NewUserGRPCClient(addr string) *UserGRPCClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user-service gRPC: %v", err)
	}

	client := userpb.NewUserServiceClient(conn)
	return &UserGRPCClient{Client: client}
}
