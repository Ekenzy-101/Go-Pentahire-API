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

func CreatePostgresConnectionPool() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	connectionPool, err = pgxpool.Connect(ctx, config.DatabaseURL)
	helpers.ExitIfError(err)

	err = connectionPool.Ping(ctx)
	helpers.ExitIfError(err)

	log.Println("Successfully connected to PostgreSQL database")
	return connectionPool
}

func GetPostgresConnectionPool() *pgxpool.Pool {
	return connectionPool
}
