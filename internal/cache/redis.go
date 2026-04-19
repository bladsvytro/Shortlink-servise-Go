package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"url-shortener/internal/config"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// Cache represents a Redis cache client
type Cache struct {
	client *redis.Client
	logger *zap.Logger
	enabled bool
}

// New creates a new Redis cache client
func New(cfg config.RedisConfig, logger *zap.Logger) (*Cache, error) {
	if !cfg.Enabled {
		return &Cache{
			enabled: false,
			logger:  logger,
		}, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Redis cache connected",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
	)

	return &Cache{
		client:  client,
		logger:  logger,
		enabled: true,
	}, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	if !c.enabled {
		return "", redis.Nil
	}

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", redis.Nil
	}
	if err != nil {
		c.logger.Error("Redis get error", zap.Error(err), zap.String("key", key))
		return "", err
	}
	return val, nil
}

// Set sets a value in cache with expiration
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if !c.enabled {
		return nil
	}

	var val string
	switch v := value.(type) {
	case string:
		val = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		val = string(jsonBytes)
	}

	err := c.client.Set(ctx, key, val, expiration).Err()
	if err != nil {
		c.logger.Error("Redis set error", zap.Error(err), zap.String("key", key))
	}
	return err
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	if !c.enabled {
		return nil
	}

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.Error("Redis delete error", zap.Error(err), zap.String("key", key))
	}
	return err
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	if !c.enabled {
		return nil
	}
	return c.client.Close()
}

// IsEnabled returns whether caching is enabled
func (c *Cache) IsEnabled() bool {
	return c.enabled
}