package services

import (
	"context"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
)

func CreateRedisClient() *redis.Client {
	options, err := redis.ParseURL(config.RedisURL)
	helpers.ExitIfError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	redisClient = redis.NewClient(options)
	helpers.ExitIfError(redisClient.Ping(ctx).Err())
	return redisClient
}

func GetRedisClient() *redis.Client {
	return redisClient
}
