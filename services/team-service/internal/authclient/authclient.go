package authclient

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/machinebox/graphql"
)

const verifyRequest string = `
		query VerifyToken($input: String!) {
			verifyToken(input: { token: $input }) {
				valid
				user {
					id
					role
				}
			}
		}
	`

type VerifyTokenResponse struct {
	VerifyToken struct {
		Valid bool `json:"valid"`
		User  struct {
			ID   uuid.UUID `json:"id"`
			Role string    `json:"role"`
		} `json:"user"`
	} `json:"verifyToken"`
}

type AuthServiceClient struct {
	client *graphql.Client
}

func NewAuthServiceClient(authServiceURL string) *AuthServiceClient {
	return &AuthServiceClient{
		client: graphql.NewClient(authServiceURL),
	}
}

func (a *AuthServiceClient) VerifyToken(ctx context.Context, token string) (*VerifyTokenResponse, error) {
	req := graphql.NewRequest(verifyRequest)

	req.Var("input", token)

	var resp VerifyTokenResponse
	if err := a.client.Run(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	return &resp, nil
}
