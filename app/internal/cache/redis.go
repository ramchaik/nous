package cache

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	GetUncompressed(ctx context.Context, key string) ([]byte, error)
	SetCompressed(ctx context.Context, key string, value []byte, expiration time.Duration) error
	HashKey(key string) string
	GetAllValues(ctx context.Context, pattern string) ([][]byte, error)
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: client}
}

func (rc *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return rc.client.Get(ctx, key).Bytes()
}

func (rc *RedisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

func (rc *RedisCache) GetUncompressed(ctx context.Context, key string) ([]byte, error) {
	compressed, err := rc.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func (rc *RedisCache) SetCompressed(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	var compressedBuf bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBuf)
	_, err := gzipWriter.Write(value)
	if err != nil {
		return err
	}
	gzipWriter.Close()

	return rc.Set(ctx, key, compressedBuf.Bytes(), expiration)
}

func (rc *RedisCache) HashKey(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (rc *RedisCache) GetAllValues(ctx context.Context, pattern string) ([][]byte, error) {
	keys, err := rc.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var values [][]byte
	for _, key := range keys {
		value, err := rc.Get(ctx, key)
		if err == nil {
			values = append(values, value)
		}
	}

	return values, nil
}
