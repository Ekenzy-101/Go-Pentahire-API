package tests

import (
	"context"
	"testing"

	"github.com/Ekenzy-101/Pentahire-API/services"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAuthRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth")
}

var (
	pool        *pgxpool.Pool
	redisClient *redis.Client
	ctx         = context.Background()

	_ = BeforeSuite(func() {
		pool = services.CreatePostgresConnectionPool(ctx)
		redisClient = services.CreateRedisClient(ctx)
	})

	_ = AfterSuite(func() {
		pool.Close()
	})
)
