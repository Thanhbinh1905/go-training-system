package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"fmt"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/dto"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/graph/helper"
	gqlmodel "github.com/Thanhbinh1905/go-training-system/services/user-service/internal/graph/model"
	"github.com/Thanhbinh1905/go-training-system/services/user-service/internal/model"
	"github.com/Thanhbinh1905/go-training-system/shared/apperror"
	"github.com/google/uuid"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input gqlmodel.CreateUserInput) (*gqlmodel.UserMutationResponse, error) {
	arg := &dto.CreateUserInput{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		Role:     model.UserRole(input.Role),
	}
	user, err := r.Service.Register(ctx, arg)
	if err != nil {
		if err == apperror.ErrEmailTaken {
			msg := err.Error()
			return helper.NewUserMutationError("400", &msg, nil), nil
		}
		msg := err.Error()
		return helper.NewUserMutationError("500", &msg, nil), nil
	}

	if user == nil {
		msg := "Failed to create user"
		return helper.NewUserMutationError("500", &msg, nil), nil
	}
	gqlUser := &gqlmodel.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     input.Role,
	}

	return helper.NewUserMutationSuccess(gqlUser), nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input gqlmodel.UserInput) (*gqlmodel.AuthMutationResponse, error) {
	args := &dto.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	}
	authRes, err := r.Service.Login(ctx, args)
	if err != nil {
		if err == apperror.ErrInvalidLogin {
			msg := err.Error()
			return helper.AuthMutationError("401", msg, nil), nil
		}
		msg := "Internal server error"
		return helper.AuthMutationError("500", msg, nil), nil
	}

	return helper.AuthMutationSuccess(authRes), nil
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented: Logout - logout"))
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, pagination *gqlmodel.UserPaginationInput) (*gqlmodel.PaginatedUsers, error) {
	args := &dto.UserPaginationInput{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
		Role: func() *model.UserRole {
			if pagination.Role == nil {
				return nil
			}
			role := model.UserRole(*pagination.Role)
			return &role
		}(),
	}

	paginatedUsers, err := r.Service.FetchUsers(ctx, args)
	if err != nil {
		return nil, err
	}
	result := make([]*gqlmodel.User, 0, len(paginatedUsers.Users))
	for _, user := range paginatedUsers.Users {
		createdAt := user.CreatedAt.Format(time.RFC3339)
		result = append(result, &gqlmodel.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      gqlmodel.UserType(user.Role),
			CreatedAt: &createdAt,
		})
	}
	return &gqlmodel.PaginatedUsers{
		Users:  result,
		Total:  paginatedUsers.Total,
		Limit:  paginatedUsers.Limit,
		Offset: paginatedUsers.Offset,
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id uuid.UUID) (*gqlmodel.User, error) {
	user, err := r.Service.User(ctx, id)
	if err != nil {
		return nil, err
	}

	createdAt := user.CreatedAt.Format(time.RFC3339)
	return &gqlmodel.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      gqlmodel.UserType(user.Role),
		CreatedAt: &createdAt,
	}, nil
}

// VerifyToken is the resolver for the verifyToken field.
func (r *queryResolver) VerifyToken(ctx context.Context, input gqlmodel.TokenInput) (*gqlmodel.TokenValidationResponse, error) {
	validated, err := r.Service.ValidateToken(&dto.TokenVerifyInput{
		Token: input.Token,
	})
	if err != nil || validated == nil || validated.User == nil {
		return &gqlmodel.TokenValidationResponse{
			Valid: false,
			User:  nil,
		}, nil
	}

	return &gqlmodel.TokenValidationResponse{
		Valid: true,
		User: &gqlmodel.UserClaims{
			ID:   validated.User.ID,
			Role: gqlmodel.UserType(validated.User.Role),
		},
	}, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, id uuid.UUID, input gqlmodel.UpdateUserInput) (*gqlmodel.UserMutationResponse, error) {
	panic(fmt.Errorf("not implemented: UpdateUser - updateUser"))
}

// Teams is the resolver for the teams field.
func (r *queryResolver) Teams(ctx context.Context) ([]*gqlmodel.Team, error) {
	panic(fmt.Errorf("not implemented: Teams - teams"))
}

// Team is the resolver for the team field.
func (r *queryResolver) Team(ctx context.Context, teamID uuid.UUID) (*gqlmodel.Team, error) {
	panic(fmt.Errorf("not implemented: Team - team"))
}

// MyTeams is the resolver for the myTeams field.
func (r *queryResolver) MyTeams(ctx context.Context) ([]*gqlmodel.Team, error) {
	panic(fmt.Errorf("not implemented: MyTeams - myTeams"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
