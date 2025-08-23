package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	connectTimeOut = 5 * time.Second
)

type RedisClient struct {
	Client *redis.Client
}

// Init khởi tạo Redis client
func Init(addr, password string, db int, log *zap.Logger) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // để trống nếu không có
		DB:       db,       // default = 0
	})

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeOut)
	defer cancel()

	// Test kết nối
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Error("Failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	log.Info("Redis connected successfully")
	return &RedisClient{Client: rdb}, nil
}

// Close đóng kết nối Redis
func (r *RedisClient) Close(log *zap.Logger) {
	if err := r.Client.Close(); err != nil {
		log.Warn("Failed to close Redis connection", zap.Error(err))
	} else {
		log.Info("Redis connection closed")
	}
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.Client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}
