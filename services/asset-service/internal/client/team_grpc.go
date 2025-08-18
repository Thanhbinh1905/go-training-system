package client

import (
	"log"

	teampb "github.com/Thanhbinh1905/go-training-system/services/asset-service/pb/team"
	"google.golang.org/grpc"
)

type TeamGRPCClient struct {
	Client teampb.TeamServiceClient
}

func NewTeamGRPCClient(addr string) *TeamGRPCClient {
	conn, err := grpc.NewClient(addr)
	if err != nil {
		log.Fatalf("Failed to connect to team-service gRPC: %v", err)
	}

	client := teampb.NewTeamServiceClient(conn)
	return &TeamGRPCClient{Client: client}
}
