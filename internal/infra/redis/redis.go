package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/IamOnah/storefronthq/internal/config"
	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, key string) error
}

type rediscache struct {
	client *redis.Client
}

func NewCache(cfg *config.RedisConfig) *rediscache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,
	})
	return &rediscache{client: rdb}
}

func (rc *rediscache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return fmt.Errorf("encode payload: %w", err)
	}
	if err := rc.client.SetEx(ctx, key, buf.Bytes(), ttl).Err(); err != nil {
		return fmt.Errorf("redis setex: %w", err)
	}
	return nil
}

func (rc *rediscache) Get(ctx context.Context, key string, dest any) error {
	data, err := rc.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return ErrCacheMiss
	} else if err != nil {
		return fmt.Errorf("redis get: %w", err)
	}

	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(dest); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}
	return nil
}

func (rc *rediscache) Delete(ctx context.Context, key string) error {
	if err := rc.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}
