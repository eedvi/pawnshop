package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrCacheSet  = errors.New("failed to set cache")
)

// Cache provides caching functionality using Redis
type Cache struct {
	client *redis.Client
	prefix string
}

// Config holds Redis cache configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	Prefix   string
}

// New creates a new Redis cache instance
func New(cfg Config) (*Cache, error) {
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

	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "pawnshop"
	}

	return &Cache{
		client: client,
		prefix: prefix,
	}, nil
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// key returns the full cache key with prefix
func (c *Cache) key(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// Set stores a value in the cache with the specified TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.client.Set(ctx, c.key(key), data, ttl).Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrCacheSet, err)
	}

	return nil
}

// Get retrieves a value from the cache and unmarshals it into the provided pointer
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, c.key(key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a value from the cache
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	fullKeys := make([]string, len(keys))
	for i, k := range keys {
		fullKeys[i] = c.key(k)
	}

	if err := c.client.Del(ctx, fullKeys...).Err(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	return nil
}

// DeleteByPattern removes all keys matching the pattern
func (c *Cache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, c.key(pattern), 0).Iterator()
	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan keys: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
	}

	return nil
}

// Exists checks if a key exists in the cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, c.key(key)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return n > 0, nil
}

// SetNX sets a value only if the key does not exist (for distributed locks)
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.SetNX(ctx, c.key(key), data, ttl).Result()
}

// Increment increments an integer value in the cache
func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, c.key(key)).Result()
}

// Decrement decrements an integer value in the cache
func (c *Cache) Decrement(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, c.key(key)).Result()
}

// Expire sets or updates the TTL for a key
func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, c.key(key), ttl).Err()
}

// TTL returns the remaining time to live for a key
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, c.key(key)).Result()
}

// GetOrSet retrieves a value from cache, or sets it using the provided function if not found
func (c *Cache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fn func() (interface{}, error)) error {
	// Try to get from cache first
	if err := c.Get(ctx, key, dest); err == nil {
		return nil
	}

	// Cache miss, get value from function
	value, err := fn()
	if err != nil {
		return err
	}

	// Store in cache
	if err := c.Set(ctx, key, value, ttl); err != nil {
		// Log but don't fail - we have the value
		_ = err
	}

	// Marshal and unmarshal to populate dest
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Client returns the underlying Redis client for advanced operations
func (c *Cache) Client() *redis.Client {
	return c.client
}
