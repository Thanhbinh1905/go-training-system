package repository

import (
	"context"

	"github.com/Thanhbinh1905/go-training-system/services/team-service/internal/model"
	"github.com/google/uuid"
)

type TeamRepository interface {
	Create(ctx context.Context, team *model.Team) error
	AddMember(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
	AddManager(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error
	RemoveManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error
	GetManagerIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error)
	GetMemberIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error)

	IsUserTeamManager(ctx context.Context, userID, teamID uuid.UUID) (bool, error)
	IsTeamExist(ctx context.Context, teamID uuid.UUID) (bool, error)
}

type teamRepository struct {
	dbRepo *TeamDbRepository
	cache  *TeamCache
}

func NewTeamRepository(dbRepo *TeamDbRepository, cache *TeamCache) TeamRepository {
	return &teamRepository{dbRepo: dbRepo, cache: cache}
}

func (r *teamRepository) AddMember(ctx context.Context, added_by, teamID, userID uuid.UUID) error {
	if err := r.dbRepo.AddMember(ctx, added_by, teamID, userID); err != nil {
		return err
	}
	_ = r.cache.AddMember(ctx, teamID, userID) // update cache
	return nil
}

func (r *teamRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	if err := r.dbRepo.RemoveMember(ctx, teamID, userID); err != nil {
		return err
	}
	_ = r.cache.RemoveMember(ctx, teamID, userID) // update cache
	return nil
}

func (r *teamRepository) GetMemberIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	// Try cache first
	ids, err := r.cache.GetMembers(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if len(ids) > 0 {
		return ids, nil
	}

	// Fallback DB
	ids, err = r.dbRepo.GetMemberIDsByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	// Rebuild cache
	_ = r.cache.SetMembers(ctx, teamID, ids)

	return ids, nil
}

func (r *teamRepository) Create(ctx context.Context, team *model.Team) error {
	return r.dbRepo.Create(ctx, team)
}

func (r *teamRepository) AddManager(ctx context.Context, added_by uuid.UUID, teamID uuid.UUID, userID uuid.UUID) error {
	return r.dbRepo.AddManager(ctx, added_by, teamID, userID)
}

func (r *teamRepository) RemoveManager(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return r.dbRepo.RemoveManager(ctx, teamID, userID)
}

func (r *teamRepository) GetManagerIDsByTeamID(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	return r.dbRepo.GetManagerIDsByTeamID(ctx, teamID)
}

func (r *teamRepository) IsUserTeamManager(ctx context.Context, userID, teamID uuid.UUID) (bool, error) {
	return r.dbRepo.IsUserTeamManager(ctx, userID, teamID)
}

func (r *teamRepository) IsTeamExist(ctx context.Context, teamID uuid.UUID) (bool, error) {
	return r.dbRepo.IsTeamExist(ctx, teamID)
}
