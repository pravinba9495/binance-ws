package main

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

type PubSub interface {
	Publish(ch string, msg []byte) error
}

type redisClient struct {
	client *redis.Client
}

func (rc *redisClient) Publish(ch string, msg []byte) error {
	return rc.client.Publish(context.Background(), ch, msg).Err()
}

func NewRedisPubSub() PubSub {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &redisClient{client: rdb}
}
