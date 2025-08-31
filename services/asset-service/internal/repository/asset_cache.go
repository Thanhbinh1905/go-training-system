package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Thanhbinh1905/go-training-system/services/asset-service/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AssetCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewAssetCache(rdb *redis.Client, ttl time.Duration) *AssetCache {
	return &AssetCache{rdb: rdb, ttl: ttl}
}

// -------------------- Folder --------------------

func (c *AssetCache) GetFolder(ctx context.Context, folderID uuid.UUID) (*model.Folder, error) {
	key := fmt.Sprintf("folder:%s", folderID.String())
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, err
	}

	var folder model.Folder
	if err := json.Unmarshal([]byte(data), &folder); err != nil {
		return nil, err
	}
	return &folder, nil
}

func (c *AssetCache) SetFolder(ctx context.Context, folder *model.Folder) error {
	key := fmt.Sprintf("folder:%s", folder.ID.String())
	data, err := json.Marshal(folder)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, c.ttl).Err()
}

func (c *AssetCache) InvalidateFolder(ctx context.Context, folderID uuid.UUID) error {
	key := fmt.Sprintf("folder:%s", folderID.String())
	return c.rdb.Del(ctx, key).Err()
}

func (c *AssetCache) GetFolders(ctx context.Context, folderIDs []uuid.UUID) ([]*model.Folder, error) {
	keys := make([]string, len(folderIDs))
	for i, id := range folderIDs {
		keys[i] = fmt.Sprintf("folder:%s", id.String())
	}

	data, err := c.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var results []*model.Folder
	for _, item := range data {
		if item == nil {
			results = append(results, nil) // cache miss
			continue
		}

		var folder model.Folder
		if err := json.Unmarshal([]byte(item.(string)), &folder); err != nil {
			return nil, err
		}
		results = append(results, &folder)
	}
	return results, nil
}

func (c *AssetCache) SetFolders(ctx context.Context, folders []*model.Folder) error {
	pipe := c.rdb.Pipeline()
	for _, folder := range folders {
		if folder == nil {
			continue
		}
		key := fmt.Sprintf("folder:%s", folder.ID.String())
		data, err := json.Marshal(folder)
		if err != nil {
			return err
		}
		pipe.Set(ctx, key, data, c.ttl)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// -------------------- Note --------------------

func (c *AssetCache) GetNote(ctx context.Context, noteID uuid.UUID) (*model.Note, error) {
	key := fmt.Sprintf("note:%s", noteID.String())
	data, err := c.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, err
	}

	var note model.Note
	if err := json.Unmarshal([]byte(data), &note); err != nil {
		return nil, err
	}
	return &note, nil
}

func (c *AssetCache) SetNote(ctx context.Context, note *model.Note) error {
	key := fmt.Sprintf("note:%s", note.ID.String())
	data, err := json.Marshal(note)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, c.ttl).Err()
}

func (c *AssetCache) InvalidateNote(ctx context.Context, noteID uuid.UUID) error {
	key := fmt.Sprintf("note:%s", noteID.String())
	return c.rdb.Del(ctx, key).Err()
}

func (c *AssetCache) GetNotes(ctx context.Context, noteIDs []uuid.UUID) ([]*model.Note, error) {
	keys := make([]string, len(noteIDs))
	for i, id := range noteIDs {
		keys[i] = fmt.Sprintf("note:%s", id.String())
	}

	data, err := c.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	var results []*model.Note
	for _, item := range data {
		if item == nil {
			results = append(results, nil) // cache miss
			continue
		}

		var note model.Note
		if err := json.Unmarshal([]byte(item.(string)), &note); err != nil {
			return nil, err
		}
		results = append(results, &note)
	}
	return results, nil
}

func (c *AssetCache) SetNotes(ctx context.Context, notes []*model.Note) error {
	pipe := c.rdb.Pipeline()
	for _, note := range notes {
		if note == nil {
			continue
		}
		key := fmt.Sprintf("note:%s", note.ID.String())
		data, err := json.Marshal(note)
		if err != nil {
			return err
		}
		pipe.Set(ctx, key, data, c.ttl)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// -------------------- ACL (Access Control) --------------------

// acl key format: asset:{assetId}:acl => hash(userId => accessType)
func (c *AssetCache) aclKey(assetID uuid.UUID) string {
	return fmt.Sprintf("asset:%s:acl", assetID.String())
}

// SetAssetACLForUser sets access for a single user on an asset
func (c *AssetCache) SetAssetACLForUser(ctx context.Context, assetID, userID uuid.UUID, access model.AccessLevel) error {
	key := c.aclKey(assetID)
	return c.rdb.HSet(ctx, key, userID.String(), string(access)).Err()
}

// SetAssetACL sets access for multiple users on an asset
func (c *AssetCache) SetAssetACLForUsers(ctx context.Context, assetID uuid.UUID, userIDs []uuid.UUID, access model.AccessLevel) error {
	key := c.aclKey(assetID)
	if len(userIDs) == 0 {
		return nil
	}
	fields := make(map[string]interface{}, len(userIDs))
	for _, uid := range userIDs {
		fields[uid.String()] = string(access)
	}
	return c.rdb.HSet(ctx, key, fields).Err()
}

// RemoveAssetACLForUser removes a user's access from an asset
func (c *AssetCache) RemoveAssetACLForUser(ctx context.Context, assetID, userID uuid.UUID) error {
	key := c.aclKey(assetID)
	return c.rdb.HDel(ctx, key, userID.String()).Err()
}

// GetAssetACLForUser retrieves a user's access for an asset
func (c *AssetCache) GetAssetACLForUser(ctx context.Context, assetID, userID uuid.UUID) (model.AccessLevel, error) {
	key := c.aclKey(assetID)
	val, err := c.rdb.HGet(ctx, key, userID.String()).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return model.AccessLevel(val), nil
}
