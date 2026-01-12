package services

import (
	"context"
	"log"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
)

func CreateRedisClient(ctx context.Context) *redis.Client {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	options, err := redis.ParseURL(config.RedisURL)
	helpers.ExitIfError(err)

	redisClient = redis.NewClient(options)
	helpers.ExitIfError(redisClient.Ping(ctx).Err())
	if !config.IsTesting {
		log.Println("Successfully connected to Redis database")
	}

	return redisClient
}

func GetRedisClient() *redis.Client {
	return redisClient
}
