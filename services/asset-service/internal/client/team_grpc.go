package client

import (
	"log"

	teampb "github.com/Thanhbinh1905/go-training-system/services/asset-service/pb/team"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TeamGRPCClient struct {
	Client teampb.TeamServiceClient
}

func NewTeamGRPCClient(addr string) *TeamGRPCClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user-service gRPC: %v", err)
	}

	client := teampb.NewTeamServiceClient(conn)
	return &TeamGRPCClient{Client: client}
}
