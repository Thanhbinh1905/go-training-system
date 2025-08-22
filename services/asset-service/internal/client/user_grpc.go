package client

import (
	"log"

	userpb "github.com/Thanhbinh1905/go-training-system/services/asset-service/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGRPCClient struct {
	Client userpb.UserServiceClient
}

func NewUserGRPCClient(addr string) *UserGRPCClient {
	// Dùng insecure.NewCredentials() nếu bạn chưa setup TLS
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user-service gRPC: %v", err)
	}

	client := userpb.NewUserServiceClient(conn)
	return &UserGRPCClient{Client: client}
}
