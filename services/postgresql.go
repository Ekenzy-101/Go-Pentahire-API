package services

import (
	"context"
	"log"
	"time"

	"github.com/Ekenzy-101/Pentahire-API/config"
	"github.com/Ekenzy-101/Pentahire-API/helpers"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	connectionPool *pgxpool.Pool
)

func CreatePostgresConnectionPool(ctx context.Context) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	connectionPool, err := pgxpool.Connect(ctx, config.DatabaseURL)
	helpers.ExitIfError(err)

	helpers.ExitIfError(connectionPool.Ping(ctx))
	if !config.IsTesting {
		log.Println("Successfully connected to PostgreSQL database")
	}

	return connectionPool
}

func GetPostgresConnectionPool() *pgxpool.Pool {
	return connectionPool
}
