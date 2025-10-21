package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bharathbbg/inventory-service/internal/config"
	"github.com/bharathbbg/inventory-service/internal/model"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(config config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) CacheProduct(ctx context.Context, product *model.Product) error {
	key := fmt.Sprintf("product:%s", product.ID)
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, 1*time.Hour).Err()
}

func (c *RedisCache) GetCachedProduct(ctx context.Context, productID string) (*model.Product, error) {
	key := fmt.Sprintf("product:%s", productID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var product model.Product
	if err := json.Unmarshal(data, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (c *RedisCache) DeleteCachedProduct(ctx context.Context, productID string) error {
	key := fmt.Sprintf("product:%s", productID)
	return c.client.Del(ctx, key).Err()
}