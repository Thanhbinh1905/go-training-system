package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	ttl = 30 * time.Minute
)

type TeamCache struct {
	rdb *redis.Client
}

func NewTeamCache(rdb *redis.Client) *TeamCache {
	return &TeamCache{rdb: rdb}
}

func (c *TeamCache) keyMembers(teamID uuid.UUID) string {
	return fmt.Sprintf("team:%s:members", teamID.String())
}

// Add member to cache
func (c *TeamCache) AddMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return c.rdb.SAdd(ctx, c.keyMembers(teamID), userID).Err()
}

// Remove member out of cache
func (c *TeamCache) RemoveMember(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	return c.rdb.SRem(ctx, c.keyMembers(teamID), userID.String()).Err()
}

// Get memberIDs from cached
func (c *TeamCache) GetMembers(ctx context.Context, teamID uuid.UUID) ([]uuid.UUID, error) {
	members, err := c.rdb.SMembers(ctx, c.keyMembers(teamID)).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(members))
	for _, m := range members {
		if uid, err := uuid.Parse(m); err == nil {
			ids = append(ids, uid)
		}
	}
	return ids, nil
}

// Build member cache from DB
func (c *TeamCache) SetMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	if len(memberIDs) == 0 {
		return nil
	}

	ids := make([]interface{}, len(memberIDs))
	for i, id := range memberIDs {
		ids[i] = id.String()
	}

	// Thêm TTL để tránh cache stale data
	pipe := c.rdb.TxPipeline()

	pipe.SAdd(ctx, c.keyMembers(teamID), ids...)
	pipe.Expire(ctx, c.keyMembers(teamID), ttl)

	_, err := pipe.Exec(ctx)
	return err
}

// Thêm method để check cache health
func (c *TeamCache) IsHealthy(ctx context.Context) bool {
	return c.rdb.Ping(ctx).Err() == nil
}

// Thêm method để clear cache khi cần
func (c *TeamCache) ClearTeamCache(ctx context.Context, teamID uuid.UUID) error {
	return c.rdb.Del(ctx, c.keyMembers(teamID)).Err()
}

// Thêm method để get cache stats
func (c *TeamCache) GetCacheStats(ctx context.Context, teamID uuid.UUID) (int64, error) {
	return c.rdb.SCard(ctx, c.keyMembers(teamID)).Result()
}
