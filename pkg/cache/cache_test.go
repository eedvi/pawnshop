package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func setupTestCache(t *testing.T) (*Cache, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cache := &Cache{
		client: client,
		prefix: "test",
	}

	return cache, mr
}

func TestCache_SetAndGet(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	data := testStruct{Name: "test", Value: 42}

	// Set
	err := cache.Set(ctx, "mykey", data, 5*time.Minute)
	assert.NoError(t, err)

	// Get
	var result testStruct
	err = cache.Get(ctx, "mykey", &result)
	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 42, result.Value)
}

func TestCache_Get_CacheMiss(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	var result testStruct
	err := cache.Get(ctx, "nonexistent", &result)
	assert.ErrorIs(t, err, ErrCacheMiss)
}

func TestCache_Delete(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Set a value
	err := cache.Set(ctx, "deletekey", "value", 5*time.Minute)
	assert.NoError(t, err)

	// Verify it exists
	exists, err := cache.Exists(ctx, "deletekey")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Delete
	err = cache.Delete(ctx, "deletekey")
	assert.NoError(t, err)

	// Verify it's gone
	exists, err = cache.Exists(ctx, "deletekey")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCache_DeleteMultiple(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Set multiple values
	cache.Set(ctx, "key1", "value1", 5*time.Minute)
	cache.Set(ctx, "key2", "value2", 5*time.Minute)
	cache.Set(ctx, "key3", "value3", 5*time.Minute)

	// Delete multiple
	err := cache.Delete(ctx, "key1", "key2")
	assert.NoError(t, err)

	// Verify key1 and key2 are gone, key3 remains
	exists, _ := cache.Exists(ctx, "key1")
	assert.False(t, exists)
	exists, _ = cache.Exists(ctx, "key2")
	assert.False(t, exists)
	exists, _ = cache.Exists(ctx, "key3")
	assert.True(t, exists)
}

func TestCache_Exists(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Key doesn't exist
	exists, err := cache.Exists(ctx, "nokey")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Set key
	cache.Set(ctx, "existkey", "value", 5*time.Minute)

	// Key exists
	exists, err = cache.Exists(ctx, "existkey")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCache_SetNX(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// First SetNX should succeed
	ok, err := cache.SetNX(ctx, "nxkey", "first", 5*time.Minute)
	assert.NoError(t, err)
	assert.True(t, ok)

	// Second SetNX should fail (key exists)
	ok, err = cache.SetNX(ctx, "nxkey", "second", 5*time.Minute)
	assert.NoError(t, err)
	assert.False(t, ok)

	// Verify original value
	var result string
	err = cache.Get(ctx, "nxkey", &result)
	assert.NoError(t, err)
	assert.Equal(t, "first", result)
}

func TestCache_Increment(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Start at 0, increment to 1
	val, err := cache.Increment(ctx, "counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)

	// Increment again
	val, err = cache.Increment(ctx, "counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), val)
}

func TestCache_Decrement(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Start at 0, decrement to -1
	val, err := cache.Decrement(ctx, "counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), val)

	// Decrement again
	val, err = cache.Decrement(ctx, "counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(-2), val)
}

func TestCache_Key(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	// Test key prefix
	key := cache.key("mykey")
	assert.Equal(t, "test:mykey", key)
}

func TestCache_GetOrSet_CacheHit(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Pre-populate cache
	cache.Set(ctx, "getorset", testStruct{Name: "cached", Value: 100}, 5*time.Minute)

	// GetOrSet should return cached value
	var result testStruct
	fnCalled := false
	err := cache.GetOrSet(ctx, "getorset", &result, 5*time.Minute, func() (interface{}, error) {
		fnCalled = true
		return testStruct{Name: "fresh", Value: 200}, nil
	})

	assert.NoError(t, err)
	assert.False(t, fnCalled) // Function should not be called
	assert.Equal(t, "cached", result.Name)
	assert.Equal(t, 100, result.Value)
}

func TestCache_GetOrSet_CacheMiss(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// GetOrSet on non-existent key
	var result testStruct
	fnCalled := false
	err := cache.GetOrSet(ctx, "newkey", &result, 5*time.Minute, func() (interface{}, error) {
		fnCalled = true
		return testStruct{Name: "fresh", Value: 200}, nil
	})

	assert.NoError(t, err)
	assert.True(t, fnCalled) // Function should be called
	assert.Equal(t, "fresh", result.Name)
	assert.Equal(t, 200, result.Value)

	// Verify it was cached
	var cached testStruct
	err = cache.Get(ctx, "newkey", &cached)
	assert.NoError(t, err)
	assert.Equal(t, "fresh", cached.Name)
}

func TestCache_TTL(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Set with TTL
	cache.Set(ctx, "ttlkey", "value", 10*time.Second)

	// Check TTL
	ttl, err := cache.TTL(ctx, "ttlkey")
	assert.NoError(t, err)
	assert.True(t, ttl > 0 && ttl <= 10*time.Second)
}

func TestCache_Expire(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Set without specific TTL
	cache.Set(ctx, "expirekey", "value", 1*time.Hour)

	// Update TTL
	err := cache.Expire(ctx, "expirekey", 5*time.Second)
	assert.NoError(t, err)

	// Check new TTL
	ttl, err := cache.TTL(ctx, "expirekey")
	assert.NoError(t, err)
	assert.True(t, ttl <= 5*time.Second)
}

func TestCache_DeleteByPattern(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Set multiple keys with pattern
	cache.Set(ctx, "user:1:profile", "data1", 5*time.Minute)
	cache.Set(ctx, "user:1:settings", "data2", 5*time.Minute)
	cache.Set(ctx, "user:2:profile", "data3", 5*time.Minute)

	// Delete by pattern for user:1
	err := cache.DeleteByPattern(ctx, "user:1:*")
	assert.NoError(t, err)

	// Verify user:1 keys are gone, user:2 key remains
	exists, _ := cache.Exists(ctx, "user:1:profile")
	assert.False(t, exists)
	exists, _ = cache.Exists(ctx, "user:1:settings")
	assert.False(t, exists)
	exists, _ = cache.Exists(ctx, "user:2:profile")
	assert.True(t, exists)
}

func TestCache_Client(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	client := cache.Client()
	assert.NotNil(t, client)
}

func TestCache_SetMarshalError(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()
	ctx := context.Background()

	// Channels cannot be marshaled to JSON
	ch := make(chan int)
	err := cache.Set(ctx, "invalid", ch, 5*time.Minute)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal")
}
