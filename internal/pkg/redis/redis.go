package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"miniblog/config"
)

var Client *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

// Close 关闭 Redis 连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Publish 发布消息到频道
func Publish(ctx context.Context, channel string, message interface{}) error {
	return Client.Publish(ctx, channel, message).Err()
}

// Subscribe 订阅频道
func Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return Client.Subscribe(ctx, channels...)
}

